import { IRequestChangePassword } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class ChangePasswordService extends HttpService {
	public constructor() {
		super("/")
	}

	public updatePassword(params: IRequestChangePassword): Promise<IResponse<{}>> {
		return this.http
			.put(`/user_info/change_password`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
