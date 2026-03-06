import { IDashboardSurface } from "./index"
import { ITree } from "~/core/shared/types/Tree"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DashboardSurfaceService extends HttpService {
	public constructor() {
		super("/")
	}

	public get(params: any): Promise<IResponse<IDashboardSurface>> {
		return this.http
			.get(`/summary/surface${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadsTree(): Promise<IResponse<ITree[]>> {
		return this.http.get(`/roads/tree`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
