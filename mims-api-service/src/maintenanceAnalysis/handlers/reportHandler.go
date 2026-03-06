package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/responses"
)

// @Summary รายงานเงื่อนไขค่าใช้จ่ายการซ่อมบำรุง
// @Description
// @tags Analyze Report
// @id GetReport1
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Maintenance Analysis ID"
// @Param type query string false "ประเภทไฟล์ของรายงาน (pdf, excel, html)"
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/analyze/{id}/report/report1 [get]
func (h *Handler) GetReport1(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	userID := helpers.GetUserInfo(c).UserId

	typeFile, ok := c.Request.URL.Query()["type"]
	reportType := ""
	if ok {
		reportType = typeFile[0]
	} else {
		errResponse := responses.FailRespone(errors.New("กรุณาเลือกประเภทรายงาน"))
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	pathFile := fmt.Sprintf("storages/analyze_report/%d", 1)
	if err := os.MkdirAll(pathFile, 0775); err != nil {
		log.Fatal(err)
	}

	switch reportType {
	case "pdf":
		resp, err := h.Usecase.GetReportType1(ID, userID, "pdf")
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "html":
		resp, err := h.Usecase.GetReportType1(ID, userID, "html")
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "excel":
		resp, err := h.Usecase.Report1Excel(ID, userID)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(200, responses.SuccessResponse(resp, 200))
	}
}

// @Summary รายงานเงื่อนไขการซ่อมบำรุง
// @Description
// @tags Analyze Report
// @id GetReportTypeTwo
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Maintenance Analysis ID"
// @Param type query string false "ประเภทไฟล์ของรายงาน (pdf, excel, html)"
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/analyze/{id}/report/report2 [get]
func (h *Handler) GetReport2(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId
	typeFile, ok := c.Request.URL.Query()["type"]
	reportType := ""
	if ok {
		reportType = typeFile[0]
	} else {
		errResponse := responses.FailRespone(errors.New("กรุณาเลือกประเภทรายงาน"))
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	pathFile := fmt.Sprintf("storages/analyze_report/%d", 1)
	if err := os.MkdirAll(pathFile, 0775); err != nil {
		log.Fatal(err)
	}

	switch reportType {
	case "pdf":
		resp, err := h.Usecase.GetReportType2(ID, userID, "pdf")
		fmt.Println(111)
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "html":
		resp, err := h.Usecase.GetReportType2(ID, userID, "html")
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "excel":
		resp, err := h.Usecase.Report2Excel(ID, userID)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(200, responses.SuccessResponse(resp, 200))
	}
}

// @Summary รายงานสรุปค่าซ่อมบำรุงและค่า IRI
// @Description
// @tags Analyze Report
// @id GetReport3
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Maintenance Analysis ID"
// @Param type query string false "ประเภทไฟล์ของรายงาน (pdf, excel, html)"
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/analyze/{id}/report/report3 [get]
func (h *Handler) GetReport3(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId

	typeFile, ok := c.Request.URL.Query()["type"]
	reportType := ""
	if ok {
		reportType = typeFile[0]
	} else {
		errResponse := responses.FailRespone(errors.New("กรุณาเลือกประเภทรายงาน"))
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	pathFile := fmt.Sprintf("storages/analyze_report/%d", 1)
	if err := os.MkdirAll(pathFile, 0775); err != nil {
		log.Fatal(err)
	}

	switch reportType {
	case "pdf":
		resp, err := h.Usecase.GetReportType3(ID, userID, "pdf")
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(200, responses.SuccessResponse(resp, 200))
	case "html":
		resp, err := h.Usecase.GetReportType3(ID, userID, "html")
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(200, responses.SuccessResponse(resp, 200))
	case "excel":
		resp, err := h.Usecase.GetReportType3Excel(ID, userID)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(200, responses.SuccessResponse(resp, 200))
	}
}

// @Summary รายงานรายละเอียดแผนงานซ่อมบำรุง ตามสายทาง
// @Description
// @tags Analyze Report
// @id GetReport4
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Maintenance Analysis ID"
// @Param type query string false "ประเภทไฟล์ของรายงาน (pdf, excel, html)"
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/analyze/{id}/report/report4 [get]
func (h *Handler) GetReport4(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId

	typeFile, ok := c.Request.URL.Query()["type"]
	reportType := ""
	if ok {
		reportType = typeFile[0]
	} else {
		errResponse := responses.FailRespone(errors.New("กรุณาเลือกประเภทรายงาน"))
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	planID, ok := c.Request.URL.Query()["plan"]
	plan := 0
	if ok {
		plan, _ = strconv.Atoi(string(planID[0]))
	}

	pathFile := fmt.Sprintf("storages/analyze_report/%d", 1)
	if err := os.MkdirAll(pathFile, 0775); err != nil {
		log.Fatal(err)
	}

	switch reportType {
	case "pdf":
		resp, err := h.Usecase.GetReportType4(ID, userID, plan, "pdf")
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "html":
		resp, err := h.Usecase.GetReportType4(ID, userID, plan, "html")
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "excel":
		resp, err := h.Usecase.Report4Excel(ID, userID, plan)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		c.JSON(200, responses.SuccessResponse(resp, 200))
	}
}

// @Summary รายงานแผนการดำเนินงานการปรับปรุงผิวทาง
// @Description
// @tags Analyze Report
// @id GetReport5
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Maintenance Analysis ID"
// @Param type query string false "ประเภทไฟล์ของรายงาน (pdf, excel, html)"
// @Param plan query string false "แผน (1,2,3)"
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/analyze/{id}/report/report5 [get]
func (h *Handler) GetReport5(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId
	plan := 0
	planID, ok := c.Request.URL.Query()["plan"]
	if ok {
		plan, _ = strconv.Atoi(planID[0])
	} else {
		plan = 0
	}

	typeFile, ok := c.Request.URL.Query()["type"]
	reportType := ""
	if ok {
		reportType = typeFile[0]
	} else {
		errResponse := responses.FailRespone(errors.New("กรุณาเลือกประเภทรายงาน"))
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	pathFile := fmt.Sprintf("storages/analyze_report/%d", 1)
	if err := os.MkdirAll(pathFile, 0775); err != nil {
		log.Fatal(err)
	}

	switch reportType {
	case "pdf":
		resp, err := h.Usecase.GetReportType5(ID, userID, plan, "pdf")
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "html":
		resp, err := h.Usecase.GetReportType5(ID, userID, plan, "html")
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
	case "excel":
		resp, err := h.Usecase.Report5Excel(ID, userID, plan, "excel")
		if err != nil {
			appErr, ok := err.(*responses.AppErr)
			if ok {
				errResponse := responses.FailRespone(appErr)
				c.JSON(appErr.StatusCode, errResponse)
			}
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
		// c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
		// c.JSON(200, responses.SuccessResponse(os.Getenv("STORAGE_IP")+"/public/excel/report5.xlsx", 200))
	}
}
