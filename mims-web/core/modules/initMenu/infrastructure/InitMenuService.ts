import { IInitMenu } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class InitMenuService extends HttpService {
	public constructor() {
		super("/")
	}

	public initMenu(): Promise<IResponse<IInitMenu[]>> {
		return this.http.get(`/menu`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
