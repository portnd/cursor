import { IOptimization, IOptimizationRequest } from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class OptizationService extends HttpService {
	public constructor() {
		super("/")
	}

	public get(): Promise<IResponse<IOptimization>> {
		return this.http
			.get(`/settings/optimization`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public post(params: IOptimizationRequest): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/optimization`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
