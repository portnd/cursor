package usecases

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/auth/domains"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	authRepo domains.AuthRepository
}

// init usecase
func NewAuthUseCase(repo domains.AuthRepository) domains.AuthUseCase {
	return &authUseCase{
		authRepo: repo,
	}
}

// =========================================================
// type authCustomClaims struct {
// 	Name   string `json:"name"`
// 	User   bool   `json:"user"`
// 	UserId int    `json:"user_id"`
// 	jwt.StandardClaims
// }

func (t *authUseCase) Login(username, password string) (string, string, error) {
	user, err := t.authRepo.GetUserByUserName(username)
	if err != nil {
		log.Println(err)
		return "", "", responses.NewLoginFailError()
	}

	if user.Username != username {
		log.Println(err)
		return "", "", responses.NewLoginFailError()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err)
		return "", "", responses.NewLoginFailError()
	}

	newAuthDetail := GenerateAuthDetail(user.Id)
	authDetail, _ := t.authRepo.GetAuthByUserId(user.Id)

	if authDetail.ID == 0 {
		err := t.authRepo.CreateAuth(newAuthDetail)
		if err != nil {
			return "", "", responses.NewInternalServerError()
		}
		authDetail = newAuthDetail
	}
	// get access control
	roles, _ := t.authRepo.GetRoleByUser(user.Id)
	accCtrl := make(map[int]string)
	for _, item := range roles {
		roleId := item.RoleID
		roleAsscessCtrls, _ := t.authRepo.GetAccessControlByRole(roleId)
		for _, item2 := range roleAsscessCtrls {
			accCtrl[item2.AccessControlId] = item2.AccessKey
		}
	}
	// departmentId := user.DepartmentId
	accCtrls := converstMapToArray(accCtrl)

	accessToken, refreshToken, err := CreateToken(authDetail, accCtrls)
	if err != nil {
		log.Println(err)
		return "", "", responses.NewInternalServerError()
	}
	return accessToken, refreshToken, nil
}

func CreateToken(authDetail models.Auth, accCtrl []string) (string, string, error) {
	claims := jwt.MapClaims{}
	if len(accCtrl) == 0 {
		claims["access_control"] = []string{}
	} else {
		claims["access_control"] = accCtrl
	}

	claims["user_id"] = authDetail.UserID
	claims["access_uuid"] = authDetail.AccessUUID
	claims["exp"] = time.Now().Add(time.Minute * 3600).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(helpers.GetSecretKey()))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := createRefreshToken(accessToken, authDetail)
	if err != nil {
		// return jwt, err
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func createRefreshToken(accessToken string, authDetail models.Auth) (string, error) {
	claims := jwt.MapClaims{}
	claims["token"] = accessToken
	claims["user_id"] = authDetail.UserID
	claims["refresh_uuid"] = authDetail.RefreshUUID
	claims["exp"] = time.Now().Add(time.Minute * 3900).Unix()
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(helpers.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
func converstMapToArray(m map[int]string) []string {
	var data []string
	for _, value := range m {
		data = append(data, value)
	}
	return data
}

func (t *authUseCase) RefreshToken(accessToken, refreshToken string) (string, string, error) {
	token, err := helpers.ValidateToken(refreshToken)
	if err != nil {
		log.Println(err)
		return "", "", responses.NewInvalidTokenError()
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return "", "", responses.NewInvalidTokenError()
	}

	authData, err := t.authRepo.GetAuthByUserId(uint(payload["user_id"].(float64)))
	if err != nil {
		return "", "", responses.NewInvalidTokenError()
	}

	if payload["refresh_uuid"] != authData.RefreshUUID {
		return "", "", responses.NewInvalidTokenError()
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err = parser.ParseUnverified(payload["token"].(string), claims)
	if err != nil {
		log.Println(err)
		return "", "", responses.NewInvalidTokenError()
	}

	payload, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", responses.NewInvalidTokenError()
	}

	updatedAuth := GenerateAuthDetail(uint(payload["user_id"].(float64)))

	err = t.authRepo.UpdateAuth(updatedAuth)
	if err != nil {
		return "", "", responses.NewInternalServerError()
	}

	// get access control
	userId := uint(payload["user_id"].(float64))
	roles, err := t.authRepo.GetRoleByUser(userId)
	if err != nil {
		return "", "", responses.NewInternalServerError()
	}

	accCtrl := make(map[int]string)
	for _, item := range roles {
		roleId := item.RoleID
		roleAsscessCtrls, _ := t.authRepo.GetAccessControlByRole(roleId)
		for _, item2 := range roleAsscessCtrls {
			accCtrl[item2.AccessControlId] = item2.AccessKey
		}
	}
	accCtrls := converstMapToArray(accCtrl)

	newAccessToken, newRefreshToken, err := CreateToken(updatedAuth, accCtrls)
	if err != nil {
		return "", "", responses.NewInternalServerError()
	}
	return newAccessToken, newRefreshToken, nil
}

func (t *authUseCase) Logout(authHeader string) error {
	if authHeader == "" {
		return responses.NewAppErr(http.StatusUnauthorized, "unauthorized")
	}

	extractToken := strings.Split(authHeader, " ")[1]

	token, err := helpers.ValidateToken(extractToken)
	if err != nil {
		log.Println(err)
		return responses.NewInvalidTokenError()
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return responses.NewInvalidTokenError()
	}

	authDetail, err := t.authRepo.GetAuthByUserId(uint(payload["user_id"].(float64)))
	if err != nil {
		log.Println(err)
		return responses.NewInvalidTokenError()
	}

	if payload["access_uuid"] != authDetail.AccessUUID {
		return responses.NewInvalidTokenError()
	}

	err = t.authRepo.DeleteAuthByUserId(uint(payload["user_id"].(float64)))
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func GenerateAuthDetail(userId uint) models.Auth {
	return models.Auth{
		UserID:      userId,
		AccessUUID:  GenerateUUIDV4(),
		RefreshUUID: GenerateUUIDV4(),
	}
}

func GenerateUUIDV4() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func (t *authUseCase) ForgotPassword(request requests.ForgotPasswordRequest) error {
	user, err := t.authRepo.GetUserByEmail(request.Email)
	if err != nil {
		return responses.NewAppErr(http.StatusBadRequest, "email is incorrect")
	}

	resetPasswordToken := createResetPasswordToken(request.Email, 3)
	resetPasswordUrl := request.CallbackUrl + "/" + resetPasswordToken

	content := helpers.EmailContent{
		SendFromEmail: fmt.Sprintf("<%s>", os.Getenv("SMTP_USER")),
		SendToEmail:   request.Email,
		EmailSubject:  "Reset Your Password",
		EmailMessage: fmt.Sprintf(`
		<html>
		<head>
		<style>
			* {
				box-sizing: boder-box;
				padding: 0;
				margin: 0;
			
		</style>
		</head>
			<body>
				<p>Reset password's time will expire in : %s</p>
				<a href="%s">reset email link</a>
			</body>
		</html>
	`, time.Now().Add(time.Duration(3)*time.Minute).Format("2006-01-02 15:04:05"), resetPasswordUrl),
	}

	// user.ResetPasswordToken = resetPasswordToken

	err = t.authRepo.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	err = helpers.CreateEmailAndSend(content)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func (t *authUseCase) ResetPassword(requet requests.ResetPasswordRequest) error {
	err := validateResetPasswordToken(requet.ResetPasswordToken)
	if err != nil {
		log.Println(err)
		return responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	if requet.NewPassword != requet.ConfirmNewPassword {
		return responses.NewAppErr(http.StatusUnprocessableEntity, "ConfirmNewPassword:mismatch")
	}

	user, err := t.authRepo.GetUserByResetPasswordToken(requet.ResetPasswordToken)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(requet.ConfirmNewPassword), 14)
	user.Password = string(newPassword)
	// user.ResetPasswordToken = " "

	err = t.authRepo.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func (au *authUseCase) CheckResetPasswordToken(token string) (string, error) {
	err := validateResetPasswordToken(token)
	if err != nil {
		return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	user, err := au.authRepo.GetUserByResetPasswordToken(token)
	if err != nil {
		return "", responses.NewInternalServerError()
	}

	return user.Email, nil
}

func (au *authUseCase) VerifyEmail(token string) error {
	user, err := au.authRepo.GetUserByVerifyEmailToken(token)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return responses.NewNotFoundError()
		}
		return responses.NewInternalServerError()
	}

	// user.IsVerifyEmail = true
	// user.VerifyEmailToken = " "

	if err := au.authRepo.UpdateUser(user); err != nil {
		return responses.NewInternalServerError()
	}

	return nil
}

func (au *authUseCase) ResendVerifyEmail(request requests.ResendVerifyEmail) error {
	// user, err := au.authRepo.GetUserByID(request.UserId)
	// if err != nil {
	// 	log.Println(err)
	// 	if err.Error() == "record not found" {
	// 		return responses.NewNotFoundError()
	// 	}
	// 	return responses.NewInternalServerError()
	// }

	// if user.IsVerifyEmail {
	// 	return responses.NewAppErr(http.StatusBadRequest, "use was already verify email")
	// }

	// verifyEmailURL := request.CallbackUrl + "/" + user.VerifyEmailToken

	// content := helpers.EmailContent{
	// 	SendFromEmail: fmt.Sprintf("<%s>", os.Getenv("SMTP_USER")),
	// 	SendToEmail:   user.Email,
	// 	EmailSubject:  "Resend: verify your email",
	// 	EmailMessage: fmt.Sprintf(`
	// 	<html>
	// 	<head>
	// 	<style>
	// 		* {
	// 			box-sizing: boder-box;
	// 			padding: 0;
	// 			margin: 0;

	// 	</style>
	// 	</head>
	// 		<body>
	// 			<a href="%s">verify email link</a>
	// 		</body>
	// 	</html>
	// `, verifyEmailURL),
	// }

	// if err := helpers.CreateEmailAndSend(content); err != nil {
	// 	log.Println(err)
	// 	return responses.NewInternalServerError()
	// }

	return nil
}

func createResetPasswordToken(email string, expiredTime int) string {
	expiredTokenTime := time.Now().Add(time.Duration(expiredTime) * time.Minute).Format(time.RFC3339)
	hashEmail := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))

	return base64.StdEncoding.EncodeToString([]byte(expiredTokenTime + " " + hashEmail))
}

func validateResetPasswordToken(token string) error {
	decodeToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return errors.New("token is invalid")
	}

	expiredTimeInString := strings.Split(string(decodeToken), " ")[0]

	expiredTime, err := time.Parse(time.RFC3339, expiredTimeInString)
	if err != nil {
		return errors.New("token is invalid")
	}

	if time.Now().After(expiredTime) {
		return errors.New("token's time is expired")
	}

	return nil
}
