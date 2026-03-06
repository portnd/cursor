package responses

import "gitlab.com/mims-api-service/models"

type Role struct {
	Items        []models.Role `json:"items"`
	CurrentPage  int64         `json:"current_page"`
	NextPage     int64         `json:"next_page"`
	PreviousPage int64         `json:"previous_page"`
	SizePerPage  int64         `json:"size_per_page"`
	TotalPages   int64         `json:"total_pages"`
	TotalItems   int64         `json:"total_items"`
}

type RoleById struct {
	Id          int           `json:"id"`
	Role        string        `json:"role"`
	AccessGroup []AccessGroup `json:"access_group"`
}

type AccessGroup struct {
	Name          string          `json:"name"`
	AccessControl []AccessControl `json:"access_control"`
}

type AccessControl struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	IsCheck bool   `json:"is_check"`
}
