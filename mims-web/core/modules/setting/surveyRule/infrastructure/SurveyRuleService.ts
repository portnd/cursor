import { IResSurveyRule, IRequestSurveyRule } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class SurveyRuleService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestSurveyRule): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/owners`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IResSurveyRule>> {
		return this.http
			.get(`/settings/owners/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestSurveyRule): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/owners/${id}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
