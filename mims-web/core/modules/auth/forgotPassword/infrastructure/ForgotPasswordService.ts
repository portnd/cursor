import { IRequestForgotPassword } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class ForgotPasswordService extends HttpService {
	public constructor() {
		super("/")
	}

	public forgotPassword(params: IRequestForgotPassword): Promise<IResponse<{}>> {
		return this.http
			.post(`/auth/forgot_password`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
