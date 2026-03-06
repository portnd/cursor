import {
	IReportChildRoadSummaryRequest,
	IRoadsAssetFilter,
	IRoadsAssetMapFilter,
	IRoadsVolumeAadtFilter,
	IRoadsRoadConditionAssetFilter,
	IRoadsRoadConditionFilter,
	IRoadsRoadConditionSummaryFilter,
	IRoadsRoadReflectLightFilter,
	IRoadsRoadReflectLightSummaryFilter,
	IRoadsRoadSurfaceFilter,
	IReportProjectMaintenanceFilter,
	IReportMaintenanceFilter,
	IReportChildRoadDamageFilter,
	IRoadsRoadDamageSummaryFilter,
} from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class ChildReportSerivce extends HttpService {
	public constructor() {
		super("")
	}

	public getChildAssetFilter(): Promise<IResponse<IRoadsAssetFilter>> {
		return this.http
			.get(`/report/asset/filter/type1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildAssetMapFilter(): Promise<IResponse<IRoadsAssetMapFilter>> {
		return this.http
			.get(`/report/asset/filter/type2`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadConditionAssetFilter(): Promise<IResponse<IRoadsRoadConditionAssetFilter>> {
		return this.http
			.get(`/report/asset/filter/type3`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadSummaryFilter(): Promise<IResponse<IRoadsRoadConditionAssetFilter>> {
		return this.http
			.get(`/report/road/filter/type1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getExportRoadSummary(params: IReportChildRoadSummaryRequest): Promise<IResponse<{ url: string }>> {
		return this.http
			.get(`/report/road/report` + this.createParams(params))
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadSurfaceFilter(): Promise<IResponse<IRoadsRoadSurfaceFilter>> {
		return this.http
			.get(`/report/road/filter/type2`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadConditionFilter(): Promise<IResponse<IRoadsRoadConditionFilter>> {
		return this.http
			.get(`/report/road/filter/type3`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadConditionSummaryFilter(): Promise<IResponse<IRoadsRoadConditionSummaryFilter>> {
		return this.http
			.get(`/report/road/filter/type4`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadReflectLightFilter(): Promise<IResponse<IRoadsRoadReflectLightFilter>> {
		return this.http
			.get(`/report/road/filter/type5`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChildRoadReflectLightSummaryFilter(): Promise<IResponse<IRoadsRoadReflectLightSummaryFilter>> {
		return this.http
			.get(`/report/road/filter/type6`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReportRoadDamageFilter(): Promise<IResponse<IReportChildRoadDamageFilter>> {
		return this.http
			.get(`/report/road_damage/filter/type1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReportRoadDamageSummaryFilter(): Promise<IResponse<IRoadsRoadDamageSummaryFilter>> {
		return this.http
			.get(`/report/road_damage/filter/type2`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReportMaintenanceFilter(): Promise<IResponse<IReportMaintenanceFilter>> {
		return this.http
			.get(`/report/maintenance_kpi/filter/type1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReportProjectMaintenanceFilter(): Promise<IResponse<IReportProjectMaintenanceFilter>> {
		return this.http
			.get(`/report/maintenance/filter/type1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReportVolumeAadtFilter(): Promise<IResponse<IRoadsVolumeAadtFilter>> {
		return this.http
			.get(`/report/aadt/filter/type1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
