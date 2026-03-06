import {
	IAnnualRoadsTree,
	IAnnualStrategicsList,
	IAnnualAnalyzeData,
	IAnnualAnalyzeDataDefault,
	IAnnualCopy,
	IAnnualDefaultDataStep2,
	IAnnualDashboard,
	IAnnualMapFilter,
	IAnnualMapData,
	IAnnualRoadGroup,
	IAnnualInterventionCriteria,
} from "./AnnualModel"
import {
	IAnnualParams,
	IAnnualStepParams2,
	IAnnualUpdateModelReq,
	IAnnualUpdatePrepareDataParams,
	IMapDataReq,
} from "./AnnualRequest"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class AnnualService extends HttpService {
	public constructor() {
		super("/")
	}

	public getRoadTrees(): Promise<IResponse<IAnnualRoadsTree[]>> {
		return this.http
			.get("maintenance/road_dropdown_list_analyze")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getStrategicList(): Promise<IResponse<IAnnualStrategicsList[]>> {
		return this.http
			.get("ref/maintenance_analysis_strategic")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createPreparingData(params: IAnnualParams): Promise<IResponse<IAnnualAnalyzeData>> {
		return this.http
			.post("analyze/prepare_data", params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createAnnualStep2(id: number, prepareDataId: number[]): Promise<IResponse<IAnnualDefaultDataStep2>> {
		const param = { prepare_data_id: prepareDataId }
		return this.http
			.post(`analyze/${id}/condition`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createAnalye(id: number, params: IAnnualStepParams2): Promise<IResponse<{}>> {
		return this.http
			.post(`analyze/${id}/analyzing`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDefaultData(id: number): Promise<IResponse<IAnnualAnalyzeDataDefault>> {
		return this.http.get(`analyze/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public updatePrepareData(id: number, params: IAnnualUpdatePrepareDataParams): Promise<IResponse<IAnnualAnalyzeData>> {
		return this.http
			.put(`analyze/${id}/prepare_data`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createFavorite(id: number): Promise<IResponse<{ is_favorite: boolean }>> {
		return this.http
			.post(`analyze/${id}/favorite`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createCopy(id: number): Promise<IResponse<IAnnualCopy>> {
		return this.http.post(`analyze/${id}/copy`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getMaintenancePaymentReport(id: number): Promise<IResponse<string>> {
		return this.http
			.get(`analyze/${id}/report/report1`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditonMaintenanceReport(id: number): Promise<IResponse<string>> {
		return this.http
			.get(`analyze/${id}/report/report2`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenancePaymentIRIReport(id: number): Promise<IResponse<string>> {
		return this.http
			.get(`analyze/${id}/report/report3`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenancePlanReport(id: number): Promise<IResponse<string>> {
		return this.http
			.get(`analyze/${id}/report/report4`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceSurfacePlanReport(id: number): Promise<IResponse<string>> {
		return this.http
			.get(`analyze/${id}/report/report5`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAnnualDashboard(id: number): Promise<IResponse<IAnnualDashboard>> {
		return this.http
			.get(`analyze/dashboard/annual/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getCheckPrepareData(id: number): Promise<IResponse<{ status: boolean }>> {
		return this.http
			.get(`analyze/${id}/check_prepare_data`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getPrepareDataId(id: number): Promise<IResponse<number[]>> {
		return this.http
			.get(`analyze/${id}/prepare_data_id`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMapFilterOptions(id: number): Promise<IResponse<IAnnualMapFilter>> {
		return this.http
			.get(`analyze/dashboard-map/${id}/filter`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMapData(id: number, params: Partial<IMapDataReq>): Promise<IResponse<IAnnualMapData>> {
		const queries = this.createParams(params)
		return this.http
			.get(`analyze/dashboard-map/${id}?${queries}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadTree(): Promise<IResponse<IAnnualRoadGroup[]>> {
		return this.http
			.get("maintenance/road_dropdown_list")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getInterventionCriteria(): Promise<IResponse<IAnnualInterventionCriteria[]>> {
		return this.http
			.get(`maintenance/intervention_criteria`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAnalyzeDefaultDetails(id: number): Promise<IResponse<IAnnualAnalyzeData>> {
		return this.http.get(`analyze/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public updateModel(id: number, params: IAnnualUpdateModelReq[]): Promise<IResponse<{}>> {
		return this.http
			.put(`analyze/${id}/model`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getPrepareDataSelectedId(id: number): Promise<IResponse<number[]>> {
		return this.http
			.get(`analyze/${id}/prepare_data_selected`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
