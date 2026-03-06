import { IRequestCheckPasswordToken, ICheckResetPasswordToken } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class CheckResetPasswordTokenService extends HttpService {
	public constructor() {
		super("/")
	}

	public checkResetPasswordToken(params: IRequestCheckPasswordToken): Promise<IResponse<ICheckResetPasswordToken>> {
		return this.http
			.get(`/auth/check_reset_password_token/` + params.reset_password_token)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
