import { IReflectivityRule, IRequestReflectivityRule } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class ReflectivityRuleService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestReflectivityRule): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/owners_road_line`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IReflectivityRule>> {
		return this.http
			.get(`/settings/owners_road_line/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestReflectivityRule): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/owners_road_line/${id}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
