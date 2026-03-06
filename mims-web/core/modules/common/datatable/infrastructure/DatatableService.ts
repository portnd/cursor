import { IDatatable, IRequestDatatable } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DatatableService extends HttpService {
	public constructor() {
		super("/")
	}

	public get(url: string, params: IRequestDatatable): Promise<IResponse<IDatatable>> {
		return this.http
			.get(`${url}?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
