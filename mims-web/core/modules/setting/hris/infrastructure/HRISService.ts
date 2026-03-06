import { IHRISItem, IHRISPreview, IRequestHRIS } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class HRISService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestHRIS): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/hris`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IHRISItem>> {
		return this.http.get(`/settings/hris/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestHRIS): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/hris/${id}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public preview(): Promise<IResponse<IHRISPreview>> {
		return this.http
			.get(`/settings/hris_preview`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public import(): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/hris_import`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
