import { IDashboardAsset, IAssetLocation } from "./index"
import { ITree } from "~/core/shared/types/Tree"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DashboardAssetService extends HttpService {
	public constructor() {
		super("/")
	}

	public get(params: any): Promise<IResponse<IDashboardAsset[]>> {
		return this.http
			.get(`/summary/asset?${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getLocation(params: any): Promise<IResponse<IAssetLocation[]>> {
		return this.http
			.get(`/summary/asset_location?${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadsTree(): Promise<IResponse<ITree[]>> {
		return this.http.get(`/roads/tree`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
