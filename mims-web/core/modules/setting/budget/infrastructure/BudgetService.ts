import { IBudget, IRequestBudget } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class BudgetService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestBudget): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/budget`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IBudget>> {
		return this.http
			.get(`/settings/budget/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(params: IRequestBudget): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/budget`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
