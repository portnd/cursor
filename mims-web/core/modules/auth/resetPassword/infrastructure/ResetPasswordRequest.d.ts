export interface IRequestResetPassword {
	reset_password_token: string
	new_password: string
	confirm_new_password: string
}
