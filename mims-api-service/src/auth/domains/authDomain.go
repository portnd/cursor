package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
)

// business logic
type AuthUseCase interface {
	Login(email, password string) (string, string, error)
	// GenerateToken(email string, isUser bool) (string, string, error)
	RefreshToken(accessToken, refreshToken string) (string, string, error)
	Logout(authHeader string) error
	ForgotPassword(request requests.ForgotPasswordRequest) error
	ResetPassword(requet requests.ResetPasswordRequest) error
	CheckResetPasswordToken(token string) (string, error)
	ResendVerifyEmail(request requests.ResendVerifyEmail) error
	VerifyEmail(token string) error
}

// อะไรเชื่อมต่อกับ DB
type AuthRepository interface {
	GetUserByUserName(username string) (models.Users, error)
	GetUserByID(userId uint) (models.Users, error)
	GetUserByEmail(email string) (models.Users, error)
	GetUserByResetPasswordToken(token string) (models.Users, error)
	UpdateUser(user models.Users) error
	CreateAuth(models.Auth) error
	GetAuthByUserId(userId uint) (models.Auth, error)
	DeleteAuthByUserId(userId uint) error
	UpdateAuth(authDetail models.Auth) error
	GetRoleByUser(uint) ([]models.UserRole, error)
	GetAccessControlByRole(int) ([]models.RoleAccessJoinControl, error)
	GetUserByVerifyEmailToken(token string) (models.Users, error)
}
