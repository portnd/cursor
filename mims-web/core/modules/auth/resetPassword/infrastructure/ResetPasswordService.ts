import { IRequestResetPassword } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class ResetPasswordService extends HttpService {
	public constructor() {
		super("/")
	}

	public resetPassword(params: IRequestResetPassword): Promise<IResponse<{}>> {
		return this.http
			.post(`/auth/reset_password`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
