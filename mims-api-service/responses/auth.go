package responses

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	Status bool          `json:"status" example:"true"`
	Data   LoginResponse `json:"data"`
}

type CheckResetPasswordTokenResponse struct {
	Email string `json:"email"`
}
