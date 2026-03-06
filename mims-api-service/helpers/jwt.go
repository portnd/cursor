package helpers

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// )

// type Jwt struct {
// }

// func (j Jwt) CreateToken() (models.Token, error) {
// 	var err error
// 	fmt.Println(user)
// 	claims := jwt.MapClaims{}
// 	claims["user_id"] = user.Id
// 	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	jwt := models.Token{}

// 	jwt.AccessToken, err = token.SignedString([]byte(os.Getenv("SECRET_KEY")))
// 	if err != nil {
// 		return jwt, err
// 	}

// 	return j.createRefreshToken(jwt)
// }

// func (Jwt) ValidateToken(accessToken string) (models.User, error) {
// 	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("SECRET_KEY")), nil
// 	})

// 	user := models.User{}
// 	if err != nil {
// 		return user, err
// 	}

// 	_, ok := token.Claims.(jwt.MapClaims)
// 	if ok && token.Valid {
// 		return user, nil
// 	}

// 	return user, errors.New("invalid token")
// }

// func (j Jwt) ValidateRefreshToken(model models.Token) (models.User, error) {
// 	token, err := jwt.Parse(model.RefreshToken, func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(os.Getenv("SECRET_KEY")), nil
// 	})

// 	var user models.User
// 	if err != nil {
// 		return user, err
// 	}

// 	payload, ok := token.Claims.(jwt.MapClaims)
// 	if !(ok && token.Valid) {
// 		return user, errors.New("invalid token")
// 	}

// 	claims := jwt.MapClaims{}
// 	parser := jwt.Parser{}
// 	token, _, err = parser.ParseUnverified(payload["token"].(string), claims)
// 	if err != nil {
// 		return user, err
// 	}

// 	payload, ok = token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return user, errors.New("invalid token")
// 	}

// 	user.Id = int(payload["user_id"].(float64))

// 	return user, nil
// }

// func (Jwt) createRefreshToken(token models.Token) (models.Token, error) {
// 	var err error

// 	claims := jwt.MapClaims{}
// 	claims["token"] = token.AccessToken
// 	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

// 	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	token.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
// 	if err != nil {
// 		return token, err
// 	}

// 	return token, nil
// }

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}
		return []byte(GetSecretKey()), nil
	})
}
