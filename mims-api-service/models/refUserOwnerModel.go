package models

// Todo ...
type RefUserOwner struct {
	Id   int    `json:"id"`
	Name string `json:"email"`
}

func (b *RefUserOwner) TableName() string {
	return "ref_user_owner"
}
