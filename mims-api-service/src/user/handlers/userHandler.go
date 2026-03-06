package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/middlewares"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/user/domains"

	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
)

// init handler
type UserHandler struct {
	userUseCase domains.UserUseCase
}

// init handler
func NewUserHandler(usecase domains.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// ================================== start function  ==================================

// request form
type LoginCredentials struct {
	Email    string `form:"email" validate:"min=1"`
	Password string `form:"password" validate:"min=1"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Param struct {
	Offset      int               `json:"offset"`
	Limit       int               `json:"limit"`
	Search      map[string]string `json:"search"`
	Column      string            `json:"column"`
	Dir         string            `json:"dir"`
	ColumnOrder map[string]string `json:"column_order"`
}

// @summary
// @description
// @tags user info
// @id user_info
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.UserInfo "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/user_info [get]
func (t *UserHandler) GetUserInfo(c *gin.Context) {
	userId, errToken := middlewares.ExtractTokenMetadata(c.Request)
	if errToken != nil {
		fmt.Println(errToken)
		return
	}
	data, err := t.userUseCase.UserInfo(userId)
	if err != nil {
		fmt.Println(errToken)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags user info
// @id user_info_update
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param User body requests.UserInfoUpdateReq true "Insert your user"
// @response 200 {object} responses.UpdateResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/user_info [put]
func (t *UserHandler) UpdateUserInfo(c *gin.Context) {
	userId, errToken := middlewares.ExtractTokenMetadata(c.Request)
	if errToken != nil {
		errResponse := responses.FailRespone(errToken)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	var req requests.UserInfoUpdateReq
	errBind := c.BindJSON(&req)
	if errBind != nil {
		fmt.Println(errBind)
		errResponse := responses.FailRespone(errBind)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err := t.userUseCase.UpdateUserInfo(userId, req)
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
	// userId, errToken := middlewares.ExtractTokenMetadata(c.Request)
	if errToken != nil {
		fmt.Println(errToken)
		return
	}
	data, err := t.userUseCase.UserInfo(userId)
	if err != nil {
		fmt.Println(errToken)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 202))
}

// @summary
// @description
// @tags user info
// @id user_info_update_password
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param User body requests.UpdatePasswordUserInfoReq true "Insert your password data"
// @response 200 {object} responses.UpdateResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/user_info/change_password [put]
func (t *UserHandler) UpdatePasswordUserInfo(c *gin.Context) {
	userId, errToken := middlewares.ExtractTokenMetadata(c.Request)
	if errToken != nil {
		errResponse := responses.FailRespone(errToken)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var req requests.UpdatePasswordUserInfoReq
	errBind := c.BindJSON(&req)
	if errBind != nil {
		fmt.Println(errBind)
		errResponse := responses.FailRespone(errBind)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	isPassword, err := t.userUseCase.CheckPassword(userId, req.CurrentPassword)
	if err != nil {
		errData := responses.NewAppErr(400, "รหัสผ่านปัจุบันไม่ถูกต้อง")
		errResponse := responses.FailRespone(errData)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if !isPassword {
		errData := responses.NewAppErr(400, "รหัสผ่านปัจุบันไม่ถูกต้อง")
		errResponse := responses.FailRespone(errData)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	err = t.userUseCase.UserInfoChangePassword(userId, req)
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
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
	}

	//
	c.JSON(202, responses.SuccessResponse(responses.NoData{}, 202))
}

// @summary
// @description
// @tags users
// @id users
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param email query string false "root@root.com" example(root@root.com)
// @Param username query string false "root" example(root)
// @Param department_id query string false "1" example(1)
// @response 200 {object} responses.Users "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/users [get]
func (t *UserHandler) User(c *gin.Context) {
	var queryParams requests.UserQueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	data, _ := t.userUseCase.GetUser(queryParams)
	c.JSON(200, responses.SuccessResponse(data, 200))

}

// @summary
// @description
// @tags users
// @id users_by_id
// @param id path int true "id of User to be gotten"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.User "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/users/{id} [get]
func (t *UserHandler) GetUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	userId, _ := strconv.Atoi(id)
	user, err := t.userUseCase.GetUserById(userId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(user, 200))
}

// @summary
// @description
// @tags users
// @id users_create
// @param User body responses.User true "Insert your user"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 201 {object} responses.User "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/users [post]
func (t *UserHandler) CreateUser(c *gin.Context) {
	userId, errToken := middlewares.ExtractTokenMetadata(c.Request)
	if errToken != nil {
		fmt.Println(errToken)
		// return errToken
	}
	var req requests.UserReq
	errBind := c.BindJSON(&req)
	if errBind != nil {
		errResponse := responses.FailRespone(errBind)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	req.CreatedBy = userId
	data, err := t.userUseCase.CreateUser(req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags users
// @id users_update
// @param id path int true "id of User to be gotten"
// @param User body responses.User true "Update your user"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 201 {object} responses.User "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/users/{id} [put]
func (t *UserHandler) UpdateUserById(c *gin.Context) {
	// user id token
	userId, errToken := middlewares.ExtractTokenMetadata(c.Request)
	if errToken != nil {
		errResponse := responses.FailRespone(errToken)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	id, _ := strconv.Atoi(c.Params.ByName("id"))

	var req requests.UserReq
	errBind := c.BindJSON(&req)
	if errBind != nil {
		c.JSON(http.StatusUnprocessableEntity, errBind.Error())
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	req.UpdatedBy = userId
	data, err := t.userUseCase.UpdateUser(id, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags users
// @id users_update
// @param id path int true "id of User to be gotten"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/users/{id} [delete]
func (t *UserHandler) DeleteUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	userId, _ := strconv.Atoi(id)
	err := t.userUseCase.DeleteUserById(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Status(204)
}
