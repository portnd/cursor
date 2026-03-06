import {
	IReportAccidentData,
	IReportAssetData,
	IReportAssetImproveData,
	IReportAssetSummaryData,
	IReportRoadConditionSummaryData,
	IReportRoadsTree,
	IReportSurfaceSummary,
	IReportTrackingMaintenance,
	IReportTrafficData,
	IReportsHistMaintenance,
	IReportsYears,
} from "./ReportModels"
import { IReportConditionParams } from "./ReportsRequest"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class ReportSerivce extends HttpService {
	public constructor() {
		super("")
	}

	public getConditionsYears(): Promise<IResponse<IReportsYears>> {
		return this.http
			.get("report/condition/year")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionsReport(): Promise<IResponse<{}>> {
		return this.http
			.get("report/condition/report")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDamageYears(): Promise<IResponse<IReportsYears>> {
		return this.http.get("report/damage/year").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getDamageReport(params: IReportConditionParams): Promise<IResponse<{ data: string }>> {
		const keyword = this.createParams(params)
		return this.http
			.get(`report/damage/report?${keyword}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceHistDataOptions(): Promise<IResponse<IReportsHistMaintenance[]>> {
		return this.http
			.get(`/report/maintenance_history/data`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenaceHistReport(): Promise<IResponse<{}>> {
		return this.http
			.get("report/maintenance_history/report")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceTrackingDataOptions(): Promise<IResponse<IReportTrackingMaintenance[]>> {
		return this.http
			.get(`report/maintenance_tracking/data`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceTrackingReport(): Promise<IResponse<{}>> {
		return this.http
			.get(`report/maintenance_tracking/report`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadTrees(): Promise<IResponse<IReportRoadsTree[]>> {
		return this.http.get("roads/tree").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getTrafficData(): Promise<IResponse<IReportTrafficData>> {
		return this.http
			.get("report/traffic_volume/data")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAccidentData(): Promise<IResponse<IReportAccidentData>> {
		return this.http
			.get("report/accident_volume/data")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAssetOptions(): Promise<IResponse<IReportAssetData>> {
		return this.http.get("report/asset/data").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getAssetMapOptions(): Promise<IResponse<IReportAssetData>> {
		return this.http.get("report/map/data").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getAssetSummaryOptions(): Promise<IResponse<IReportAssetSummaryData[]>> {
		return this.http
			.get("report/summary_asset/data")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAssetImproveOptions(): Promise<IResponse<IReportAssetImproveData>> {
		return this.http
			.get("report/asset_adjustment/data")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionSummaryOptions(): Promise<IResponse<IReportRoadConditionSummaryData>> {
		return this.http
			.get("report/summary_condition/data")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getSurfaceOptions(): Promise<IResponse<IReportSurfaceSummary>> {
		return this.http.get("report/surface/data").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
