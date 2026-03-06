package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/mims-api-service/databases"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		errUnAuth := responses.AppErr{
			StatusCode: http.StatusUnauthorized,
			Message:    "unauthorized to access the resource",
		}
		errResponse := responses.FailRespone(&errUnAuth)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, err := ValidateToken(tokenString)
		if err != nil {
			logs.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		if !token.Valid {
			logs.Error(token.Valid)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logs.Error(payload)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		authDetail, err := getAuthDetailByUserID(uint64(payload["user_id"].(float64)))
		if err != nil {
			logs.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		if payload["access_uuid"] != authDetail.AccessUUID {
			logs.Error("access_uuid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		c.Set("accessControl", payload["access_control"])
		c.Set("departmentID", payload["department_id"])
		c.Set("userID", payload["user_id"])

		userPermission := []string{}
		if payload["access_control"] != nil {
			userPermission = helpers.GetUserPermission(payload["access_control"].([]interface{}))
		}
		c.Set("userInfo", helpers.UserInfo{
			UserId: int(payload["user_id"].(float64)),
			// UserDept:       int(payload["department_id"].(float64)),
			UserPermission: userPermission,
		})

		c.Next()
	}
}

func getAuthDetailByUserID(userId uint64) (models.Auth, error) {
	authDetail := models.Auth{}

	err := databases.DB.Where("user_id = ?", userId).First(&authDetail).Error
	if err != nil {
		return models.Auth{}, err
	}

	return authDetail, nil
}

func ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(helpers.GetSecretKey()), nil
	})
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(r *http.Request) (int, error) {
	tokenString := ExtractToken(r)
	token, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if !ok {
			return 0, err
		}

		if err != nil {
			return 0, err
		}
		// return claims["user_id"], nil
	}
	userId, _ := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	return userId, err
}
