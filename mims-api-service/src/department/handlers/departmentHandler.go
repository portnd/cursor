package handlers

import (
	"gitlab.com/mims-api-service/src/department/domains"

	"github.com/gin-gonic/gin"
)

// init handler
type DepartmentHandler struct {
	departmentUseCase domains.DepartmentUseCase
}

// init handler
func NewDepartmentHandler(usecase domains.DepartmentUseCase) *DepartmentHandler {
	return &DepartmentHandler{
		departmentUseCase: usecase,
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

func (t *DepartmentHandler) Role(c *gin.Context) {
	// request form
	// var param Param
	// offset, _ := strconv.Atoi(c.Query("start"))
	// limit, _ := strconv.Atoi(c.Query("length"))
	// param.Offset = offset
	// param.Limit = limit
	// param.Search = c.QueryMap("search")
	// param.ColumnOrder = map[string]string{
	// 	"0": "id",
	// 	"1": "first_name",
	// 	"2": "last_name",
	// 	"3": "role_name",
	// }
	// param.Column = c.DefaultQuery("order[0][column]", "0")
	// param.Dir = c.DefaultQuery("order[0][dir]", "desc")
	// user, _ := t.userUseCase.User(param)

	// page, size := helpers.CheckPageSize(c)

	// limitPerPage := helpers.IntToInt64(page * size)
	// limitSize := helpers.IntToInt64(size)
	// pagesSize := helpers.IntToInt64(page)
	// c.JSON(200, helpers.Response(&user, 1, limitSize, limitPerPage, pagesSize))
}
