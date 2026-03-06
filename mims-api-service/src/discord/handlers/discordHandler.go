package handlers

import (
	"net/http"

	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/discord/domains"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase domains.UseCase
}

func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

// @summary
// @description
// @tags discord
// @id get_discord
// @response 200 {object} responses.Success{data=[]string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/discord/log [get]
func (h *Handler) GetHrisAndHsmsLog(c *gin.Context) {
	err := h.Usecase.GetHrisAndHsmsLog()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
}
