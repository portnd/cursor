import { IInitData } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class InitDataService extends HttpService {
	public constructor() {
		super("/")
	}

	public initData(): Promise<IResponse<IInitData>> {
		return this.http.get(`/initdata`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
