package requests

type ForgotPasswordRequest struct {
	Email       string `json:"email" validate:"nonzero,regexp=^[0-9a-z]+(\\.[0-9a-z]+)*@[a-z]+(\\.[a-z]{2\\,3})+(\\.[a-z]{2\\,3})*$" extensions:"x-order=0"`
	CallbackUrl string `json:"callback_url" validate:"nonzero" extensions:"x-order=1"`
}

type ResetPasswordRequest struct {
	ResetPasswordToken string `json:"reset_password_token" validate:"nonzero" extensions:"x-order=0"`
	NewPassword        string `json:"new_password" validate:"min=8" extensions:"x-order=1"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"min=8" extensions:"x-order=2"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token" extensions:"x-order=0"`
	RefreshToken string `json:"refresh_token" extensions:"x-order=1"`
}

type ResendVerifyEmail struct {
	UserId      uint   `json:"user_id" validate:"nonzero" extensions:"x-order=0"`
	CallbackUrl string `json:"callback_url" validate:"nonzero" extensions:"x-order=1"`
}
