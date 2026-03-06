package responses

type RefSurfaceParam struct {
	ID           int                    `json:"id"`
	RefSurfaceID int                    `json:"ref_surface_id" gorm:"column:ref_surface_id"`
	JsonParams   map[string]interface{} `json:"params" gorm:"column:params"`
	CreateBy     int                    `json:"create_by" gorm:"column:create_by"`
	CreateDate   string                 `json:"create_date" gorm:"column:create_date"`
	IsLatest     bool                   `json:"is_latest" gorm:"column:is_latest"`
}
