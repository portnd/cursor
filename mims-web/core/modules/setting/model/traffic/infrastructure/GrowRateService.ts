import { IGrowRate, IGrowRateRequest } from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class GrowthRateService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IGrowRateRequest[]): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/aadt/growth_rate`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(): Promise<IResponse<IGrowRate[]>> {
		return this.http
			.get(`/settings/aadt/growth_rate`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
