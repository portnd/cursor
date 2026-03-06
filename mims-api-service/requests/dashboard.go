package requests

type MaintenanceDashboard struct {
	Page    int     `form:"page"`
	Limit   int     `form:"limit"`
	KmStart float64 `form:"km_start"`
	KmEnd   float64 `form:"km_end"`
	Year    string  `form:"year"`
}
