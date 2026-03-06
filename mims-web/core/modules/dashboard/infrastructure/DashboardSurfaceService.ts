import {
	IDashboardAsset,
	IDashboardAssetRequest,
	IDashboardRoad,
	IDashboardSurface,
	IDataMart,
	IDataMartCheck,
	ISurfaceMap,
} from "./index"
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

	public getSurfaceMap(display: number, roadIds: string): Promise<IResponse<ISurfaceMap[]>> {
		let roadParam = ""
		if (roadIds) {
			roadParam = `&road_id=${roadIds}`
		}
		return this.http
			.get(`/summary/surface_map?display=${display}${roadParam}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDataMart(): Promise<IResponse<IDataMart[]>> {
		return this.http.get(`/summary/data_mart`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getDataMartCheck(): Promise<IResponse<IDataMartCheck>> {
		return this.http
			.get(`/dashboard/surface/data_mart_check`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboardAssets(params: IDashboardAssetRequest): Promise<IResponse<IDashboardAsset[]>> {
		const keyword = this.createParams(params)
		return this.http
			.get(`/dashboard/asset_map?${keyword}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboardRoads(): Promise<IResponse<IDashboardRoad>> {
		return this.http.get(`dashboard/road`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
