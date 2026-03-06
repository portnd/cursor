import { IUser } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class InitUserService extends HttpService {
	public constructor() {
		super("/")
	}

	public initUser(): Promise<IResponse<IUser>> {
		return this.http.get(`/user_info`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
