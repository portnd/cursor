import {
	IDashboardAsset,
	IDashboardAssetData,
	IDashboardAssetDetail,
	IDashboardAssetDetailsRequest,
	IDashboardAssetRequest,
	IDashboardMaintenance,
	IDashboardMaintenanceMap,
	IDashboardMaintenanceMapRequest,
	IDashboardMaintenanceRequest,
	IDashboardMapAssetRequest,
	IDashboardRoad,
	IDashboardRoadOptions,
	IDashboardSurface,
	IDashboardSurfaceMap,
	IDashboardSurfaceMapRequest,
	IDashboardSurfaceRequest,
	IDataMart,
	IDataMartCheck,
	IDashboardConditionMap,
} from "./index"
import { ITree } from "~/core/shared/types/Tree"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DashboardService extends HttpService {
	public constructor() {
		super("/")
	}

	public getRoadsTree(): Promise<IResponse<ITree[]>> {
		return this.http.get(`/roads/tree`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getSurfaceMap(params: IDashboardSurfaceMapRequest): Promise<IResponse<IDashboardSurfaceMap[]>> {
		const query = this.createParams(params).size === 0 ? "" : `?${this.createParams(params)}`
		return this.http
			.get(`/dashboard/surface_map${query}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDataMart(): Promise<IResponse<IDataMart[]>> {
		return this.http
			.get(`/dashboard/surface/data_mart`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDataMartCheck(): Promise<IResponse<IDataMartCheck>> {
		return this.http
			.get(`/dashboard/surface/data_mart_check`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboardMapAssets(params: IDashboardMapAssetRequest): Promise<IResponse<IDashboardAsset[]>> {
		const keyword = this.createParams(params)
		return this.http
			.get(`/dashboard/asset_map${keyword.size > 0 ? `?${keyword}` : ""}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboardAsset(params: IDashboardAssetRequest): Promise<IResponse<IDashboardAssetData>> {
		const keyword = this.createParams(params)
		return this.http
			.get(`dashboard/asset${keyword.size > 0 ? `?${keyword}` : ""}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboardAssetDetails(params: IDashboardAssetDetailsRequest): Promise<IResponse<IDashboardAssetDetail>> {
		const keyword = this.createParams(params)
		return this.http
			.get(`dashboard/asset_detail${keyword.size > 0 ? `?${keyword}` : ""}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAssetsMapDetails(id: number, assetId: number): Promise<IResponse<string>> {
		return this.http
			.get(`dashboard/asset_map/${id}/detail/${assetId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadDropdown(type: string): Promise<IResponse<IDashboardRoadOptions[]>> {
		return this.http
			.get(`maintenance/road_dropdown_list_dashboard?type=${type}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboardYear(): Promise<IResponse<number[]>> {
		return this.http.get("dashboard/years").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getDashboardRoads(): Promise<IResponse<IDashboardRoad>> {
		return this.http.get(`dashboard/road`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getDashboardSurface(params: IDashboardSurfaceRequest): Promise<IResponse<IDashboardSurface>> {
		const query = this.createParams(params).size === 0 ? "" : `?${this.createParams(params)}`
		return this.http
			.get(`dashboard/surface${query}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenance(params: IDashboardMaintenanceRequest): Promise<IResponse<IDashboardMaintenance>> {
		const queries = this.createParams(params).size === 0 ? "" : `?${this.createParams(params)}`
		return this.http
			.get(`dashboard/maintenance${queries}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceMap(params: IDashboardMaintenanceMapRequest): Promise<IResponse<IDashboardMaintenanceMap[]>> {
		const queries = this.createParams(params).size === 0 ? "" : `?${this.createParams(params)}`
		return this.http
			.get(`dashboard/maintenance_map${queries}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getCondition_map(params: any): Promise<IResponse<IDashboardConditionMap>> {
		const queries = this.createParams(params).size === 0 ? "" : `?${this.createParams(params)}`
		return this.http
			.get(`/dashboard/condition_map${queries}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
