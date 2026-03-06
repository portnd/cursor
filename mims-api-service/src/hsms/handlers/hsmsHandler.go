package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/hsms/domains"
)

// init handler
type Handler struct {
	UseCase domains.UseCase
}

// init handler
func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		UseCase: usecase,
	}
}

// ================================== start function  ==================================

// @summary
// @description
// @tags hsms01
// @id hsms01_bridge
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/bridge [get]
func (h *Handler) GetHsmsBridge(c *gin.Context) {
	data, err := h.UseCase.GetHsmsBridge()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_guard
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/guard [get]
func (h *Handler) GetHsmsGuard(c *gin.Context) {
	data, err := h.UseCase.GetHsmsGuard()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_interchange
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/interchange [get]
func (h *Handler) GetHsmsInterchange(c *gin.Context) {
	data, err := h.UseCase.GetHsmsInterchange()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_intersection
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/intersection [get]
func (h *Handler) GetHsmsIntersection(c *gin.Context) {
	data, err := h.UseCase.GetHsmsIntersection()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_streetlight
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/streetlight [get]
func (h *Handler) GetHsmsStreetlight(c *gin.Context) {
	data, err := h.UseCase.GetHsmsStreetlight()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_railwaycrossing
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/railwaycrossing [get]
func (h *Handler) GetHsmsRailwaycrossing(c *gin.Context) {
	data, err := h.UseCase.GetHsmsRailwaycrossing()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_trafficlight
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/trafficlight [get]
func (h *Handler) GetHsmsTrafficlight(c *gin.Context) {
	data, err := h.UseCase.GetHsmsTrafficlight()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags hsms01
// @id hsms01_uturnbridge
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hsms/uturnbridge [get]
func (h *Handler) GetHsmsUturnbridge(c *gin.Context) {
	data, err := h.UseCase.GetHsmsUturnbridge()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}
