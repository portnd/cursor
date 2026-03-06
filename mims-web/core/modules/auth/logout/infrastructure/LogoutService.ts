import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class LogoutService extends HttpService {
	public constructor() {
		super("/")
	}

	public logout(): Promise<IResponse<{}>> {
		return this.http.get(`/auth/logout`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
