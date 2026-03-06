import { ILogin, IRequestLogin } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class LoginService extends HttpService {
	public constructor() {
		super("/")
	}

	public login(params: IRequestLogin): Promise<IResponse<ILogin>> {
		return this.http.post(`/auth/login`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
