package handlers

import (
	"fmt"
	"net/http"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/auth/domains"

	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
)

// init handler
type AuthHandler struct {
	authUseCase domains.AuthUseCase
}

// init handler
func NewAuthHandler(usecase domains.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: usecase,
	}
}

// ================================== start function  ==================================
// @summary
// @description
// @tags authentication
// @id login
// @param Login body requests.LoginRequest true "Login data"
// @response 200 {object} responses.TokenResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/login [post]
func (t *AuthHandler) Login(c *gin.Context) {
	// request form
	var request requests.LoginRequest
	err := c.ShouldBind(&request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// validate
	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// call usecase
	accessToken, refreshToken, err := t.authUseCase.Login(request.Username, request.Password)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
		}
		return
	}

	var loginResponse responses.LoginResponse
	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = refreshToken
	c.JSON(http.StatusOK, responses.SuccessResponse(loginResponse, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id refresh_token
// @param Login body requests.RefreshTokenRequest true "Login data"
// @response 200 {object} responses.TokenResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/refresh_token [post]
func (t *AuthHandler) RefreshToken(c *gin.Context) {

	var request requests.RefreshTokenRequest
	err := c.ShouldBind(&request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	accessToken, refreshToken, err := t.authUseCase.RefreshToken(request.AccessToken, request.RefreshToken)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	var loginResponse responses.LoginResponse
	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = refreshToken
	c.JSON(http.StatusOK, responses.SuccessResponse(loginResponse, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id logout
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/logout [get]
func (t *AuthHandler) Logout(c *gin.Context) {
	err := t.authUseCase.Logout(c.GetHeader("Authorization"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))

}

// @summary
// @description
// @tags authentication
// @id forgot_password
// @param Login body requests.ForgotPasswordRequest true "Login data"
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entit
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/forgot_password [post]
func (t *AuthHandler) ForgotPassword(c *gin.Context) {
	var forgotRequest requests.ForgotPasswordRequest
	err := c.ShouldBind(&forgotRequest)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(forgotRequest); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	if err := t.authUseCase.ForgotPassword(forgotRequest); err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id resset_password
// @param Login body requests.ResetPasswordRequest true "Login data"
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/reset_password [post]
func (t *AuthHandler) ResetPassword(c *gin.Context) {
	var resetRequest requests.ResetPasswordRequest
	err := c.ShouldBind(&resetRequest)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(resetRequest); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err = t.authUseCase.ResetPassword(resetRequest)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}

		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id check_resset_password_token
// @param reset_password_token path string false "Insert your reset password token"
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/check_reset_password_token/{reset_password_token} [get]
func (t *AuthHandler) CheckResetPasswordToken(c *gin.Context) {
	email, err := t.authUseCase.CheckResetPasswordToken(c.Param("reset_password_token"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.CheckResetPasswordTokenResponse{Email: email}, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id verify_email
// @param verify_email_token path string false "Insert your verify email token"
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/verify_email/{verify_email_token} [get]
func (t *AuthHandler) VerifyEmail(c *gin.Context) {
	err := t.authUseCase.VerifyEmail(c.Param("verify_email_token"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id resend_verify_email
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param resend_verify_email body requests.ResendVerifyEmail false "Insert your resend veify email body"
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/resend_verify_email [post]
func (t *AuthHandler) ResendVerifyEmail(c *gin.Context) {
	var request requests.ResendVerifyEmail
	err := c.ShouldBind(&request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)

	}

	if err := t.authUseCase.ResendVerifyEmail(request); err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
}

// @summary
// @description
// @tags authentication
// @id test
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.NoDataResponse "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/auth/test [get]
func (t *AuthHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hi therer!!"})
	userId, _ := c.Get("userID")
	fmt.Println(userId)
	fmt.Printf("%T\n", uint(userId.(float64)))
}
