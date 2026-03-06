package models

type Auth struct {
	ID          int
	UserID      uint
	AccessUUID  string
	RefreshUUID string
}

type AuthAccCtrl struct {
	ID            int
	UserID        uint
	AccessUUID    string
	RefreshUUID   string
	AccessControl []string
}

func (a *Auth) TableName() string {
	return "auths"
}

func (a *AuthAccCtrl) TableName() string {
	return "auths"
}
