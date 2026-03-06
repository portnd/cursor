import { ISurface, IRequestSurface } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class SurfaceService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestSurface): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/ref/surface`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<ISurface>> {
		return this.http
			.get(`/settings/ref/surface/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestSurface): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/ref/surface/${id}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
