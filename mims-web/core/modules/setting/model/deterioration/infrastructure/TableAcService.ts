import { ITableAC, ITableACRequest } from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class TableACService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: ITableACRequest): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/deterioration/asphalt`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<ITableAC>> {
		return this.http
			.get(`/settings/deterioration/asphalt/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
