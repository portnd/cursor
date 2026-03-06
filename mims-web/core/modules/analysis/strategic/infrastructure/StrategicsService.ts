import { IStrategicsList } from "../../list/infrastructure"
import {
	IStrategicDashboard,
	ICopy,
	IStrategicAnalyzeData,
	IStrategicStep1,
	IStrategicStep2,
	IStrategicRoadGroup,
	IStrategicInterventionCriteria,
	IStrtegicMapFilter,
	IStrategicMapData,
} from "./StrategicModel"
import {
	IMapDataReq,
	IReportParams,
	IStrategicCreateAnalyzeParams,
	IStrategicCreatePrepareDataReq,
	IStrategicUpdatePrepareData,
	IUpdateModelReq,
} from "./StrategicRequest"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class StrategicsService extends HttpService {
	public constructor() {
		super("/")
	}

	public getRoadTree(): Promise<IResponse<IStrategicRoadGroup[]>> {
		return this.http
			.get("maintenance/road_dropdown_list_analyze")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getStrategicList(): Promise<IResponse<IStrategicsList[]>> {
		return this.http
			.get("ref/maintenance_analysis_strategic")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createPreparingData(params: IStrategicCreatePrepareDataReq): Promise<IResponse<IStrategicStep1>> {
		return this.http
			.post("analyze/prepare_data", params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public creaeteAnalyseStep2(id: number, prepareDataId: number[]): Promise<IResponse<IStrategicStep2>> {
		const param = { prepare_data_id: prepareDataId }
		return this.http
			.post(`analyze/${id}/condition`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createAnalyzing(id: number, params: IStrategicCreateAnalyzeParams): Promise<IResponse<{}>> {
		return this.http
			.post(`analyze/${id}/analyzing`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAnalyzeDefaultDetails(id: number): Promise<IResponse<IStrategicAnalyzeData>> {
		return this.http.get(`analyze/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public createFavorite(id: number): Promise<IResponse<{ is_favorite: boolean }>> {
		return this.http
			.post(`analyze/${id}/favorite`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public copy(id: number): Promise<IResponse<ICopy>> {
		return this.http.post(`analyze/${id}/copy`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public deleteAnalyze(id: number): Promise<IResponse<{}>> {
		return this.http.delete(`analyze/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	// public updatePrepareData(id: number, params: IStrategicUpdatePrepareData): Promise<IResponse<IStrategicAnalyzeData>> {
	// 	return this.http
	// 		.put(`analyze/${id}/search`, params)
	// 		.then(this.handleResponse.bind(this))
	// 		.catch(this.handleError.bind(this))
	// }

	public getMaintenancePaymentReport(id: number): Promise<IResponse<{}>> {
		return this.http
			.get(`analyze/${id}/report/report1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditonMaintenanceReport(id: number): Promise<IResponse<{}>> {
		return this.http
			.get(`analyze/${id}/report/report2`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenancePaymentIRIReport(id: number): Promise<IResponse<{}>> {
		return this.http
			.get(`analyze/${id}/report/report3`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenancePlanReport(id: number): Promise<IResponse<{}>> {
		return this.http
			.get(`analyze/${id}/report/report4`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceSurfacePlanReport(id: number): Promise<IResponse<{}>> {
		return this.http
			.get(`analyze/${id}/report/report5`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDashboard(id: number): Promise<IResponse<IStrategicDashboard>> {
		return this.http
			.get(`analyze/dashboard/strategic/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAnalyzeReport(id: number, typeId: number, params?: IReportParams): Promise<IResponse<{ data: string }>> {
		const keyword = params ? `?${this.createParams(params)}` : ""
		return this.http
			.get(`analyze/${id}/report/report${typeId}${keyword}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getPrepareDataId(id: number): Promise<IResponse<number[]>> {
		return this.http
			.get(`analyze/${id}/prepare_data_id`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getCheckPrepareData(id: number): Promise<IResponse<{ status: boolean }>> {
		return this.http
			.get(`analyze/${id}/check_prepare_data`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updatePrepareData(id: number, params: IStrategicUpdatePrepareData): Promise<IResponse<{ id: number }>> {
		return this.http
			.put(`analyze/${id}/prepare_data`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getPrepareDataSelectedId(id: number): Promise<IResponse<number[]>> {
		return this.http
			.get(`analyze/${id}/prepare_data_selected`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getInterventionCriteria(): Promise<IResponse<IStrategicInterventionCriteria[]>> {
		return this.http
			.get(`analyze/intervention_criterias`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMapFilter(id: number): Promise<IResponse<IStrtegicMapFilter>> {
		return this.http
			.get(`analyze/dashboard-map/${id}/filter`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMap(id: number, params: IMapDataReq): Promise<IResponse<IStrategicMapData>> {
		return this.http
			.get(`analyze/dashboard-map/${id}?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateModel(id: number, params: IUpdateModelReq[]): Promise<IResponse<{}>> {
		return this.http
			.put(`analyze/${id}/model`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
