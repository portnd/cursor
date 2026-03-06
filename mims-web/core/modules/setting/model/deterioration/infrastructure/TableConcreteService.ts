import { ITableConcrete, ITableConcreteRequest } from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class TableConcreteService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: ITableConcreteRequest): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/deterioration/concrete`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<ITableConcrete>> {
		return this.http
			.get(`/settings/deterioration/concrete/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
