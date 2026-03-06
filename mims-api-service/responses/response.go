package responses

import (
	"net/http"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/models"
)

// When any routes use function SuccessResponse but there is no data
// for response, Please use this struct instead.
type NoData struct{}

type DataId struct {
	Id int `json:"id"`
}

type Pagination struct {
	CurrentPage  int64       `json:"current_page" extensions:"x-order=1"`
	NextPage     int64       `json:"next_page" extensions:"x-order=2"`
	PreviousPage int64       `json:"previous_page" extensions:"x-order=3"`
	SizePerPage  int64       `json:"size_per_page" extensions:"x-order=4"`
	TotalPages   int64       `json:"total_pages" extensions:"x-order=5"`
	TotalItems   int64       `json:"total_items" extensions:"x-order=6"`
	Items        interface{} `json:"items" extensions:"x-order=7"`
}

type ErrorFormat struct {
	Message string      `json:"message" extensions:"x-order=0"`
	Field   interface{} `json:"field" extensions:"x-order=1"`
}

// 400, 500
type Fail struct {
	Status bool        `json:"status" extensions:"x-order=0"`
	Code   int         `json:"code" extensions:"x-order=1"`
	Err    ErrorFormat `json:"error" extensions:"x-order=2"`
}

// 200, 202, 201
type Success struct {
	Status bool        `json:"status" example:"true" extensions:"x-order=0"`
	Code   int         `json:"code" extensions:"x-order=1"`
	Data   interface{} `json:"data" extensions:"x-order=2"`
}

func ValidateResponse(msg map[string]string) Fail {
	var request Fail
	request.Code = 422
	request.Err.Message = "validate error"
	request.Err.Field = msg
	request.Status = false
	return request
}

func FailRespone(err error) Fail {
	var request Fail
	appErr, ok := err.(*AppErr)
	if ok {
		request.Code = appErr.StatusCode
	}
	request.Err.Message = err.Error()
	request.Err.Field = NoData{}
	request.Status = false
	return request
}

func SuccessResponse(data interface{}, code int) Success {
	return Success{
		Status: true,
		Code:   code,
		Data:   data,
	}
}

type AppErr struct {
	StatusCode int
	Message    string
}

func (r *AppErr) Error() string {
	return r.Message
}

func NewAppErr(statusCode int, msg string) *AppErr {
	return &AppErr{
		StatusCode: statusCode,
		Message:    msg,
	}
}

func NewInternalServerError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}
}

func NewLoginFailError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    constants.INVALID_ID_PASSWORD,
	}
}

func NewInvalidTokenError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid token",
	}
}

func NewNotFoundError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "data was not found",
	}
}

func NewEmptyIDError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "id is required",
	}
}

func NewLimitNumberOfColumnsError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "number of columns exceed the limit, cannot add columns",
	}
}

func NewAccessDenyError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusForbidden,
		Message:    "access deny",
	}
}

func NewColumnNameDuplicateError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "column name is duplicate",
	}
}

func NewImageTypeError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "image file type must be only either png or jpeg",
	}
}

func NewFileSizeExceedLimitError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusBadRequest,
		Message:    "file size exceed the limit, cannot upload",
	}
}

func NewPermissionError() *AppErr {
	return &AppErr{
		StatusCode: http.StatusUnauthorized,
		Message:    constants.INVALID_USER_PERMISSION,
	}
}

func NewDuplicatedNameError(msg string) *AppErr {
	return &AppErr{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    msg,
	}
}

type Permission struct {
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
	CanSend    bool `json:"can_send"`
	CanReject  bool `json:"can_reject"`
}

type IntRes struct {
}

type TheGeom struct {
	TheGeom string `json:"the_geom"`
}

type Year struct {
	Year int `json:"year"`
}

type Sum struct {
	Sum float64 `json:"sum"`
}

type ReportTrafficVolumeHeader struct {
	RoadID          int     `json:"road_id"`
	RoadGroupName   string  `json:"road_group_name"`
	RoadSectionName string  `json:"road_section_name"`
	RoadName        string  `json:"road_name"`
	KmStart         float64 `json:"km_start"`
	KmEnd           float64 `json:"km_end"`
	TotalKm         float64 `json:"total_km"`
}

type UserInfoRespond struct {
	models.UserDepartment
	Roles         []RoleRespond          `json:"roles"`
	AccessControl []models.AccessControl `json:"access_control"`
}

type RoleRespond struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedBy int    `json:"-"`
	UpdatedBy int    `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	IsActive  bool   `json:"-"`
}
