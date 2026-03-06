import {
	IMaintenanceHistoryBudgetsData,
	IMaintenanceHistoryDetailData,
	IMaintenanceHistoryRoadGroupListData,
	IRoadChildListData,
	IMaintenanceMethodListData,
	IMaintenanceHistoryPlanStatusData,
	IMaintenancePlanListData,
	IMaintenanceHistoryAttrachments,
	IMaintenanceHistoryPlanGraph,
	IMaintenanceBudgetCriteria,
	IPlanProgressGraphReportHistTableData,
	IMaintenanceDefaultData,
	IMaintenanceDivision,
	IMaintenanceRoadGroup,
	IMaintenanceHistoryListData,
	IMaintenanceRoadData,
	IMaintenanceWarrantyData,
} from "./MaintenanceHistoryModel"
import {
	IMaintenanceHistoryCreateRequest,
	IMaintenanceHistoryEditParams,
	IMaintenanceHistoryFileParams,
	IMaintenanceHistoryGuaranteeCreateParams,
	IMaintenanceHistoryUpdateRequest,
	IMaintenanceWarrantyCreateRequest,
	IMaintenanceHistorySearch,
	IMaintenanceHistoryRoadsUpdateRequest,
} from "./MaintenanceHistoryRequest"
import { HttpService } from "~~/core/shared/http/HttpService"
import { IResponse } from "~/core/shared/http/Response"

export class MaintenanceHistoryService extends HttpService {
	public constructor() {
		super("/")
	}

	public getHistoryRoadGroup(): Promise<IResponse<IMaintenanceHistoryRoadGroupListData[]>> {
		return this.http.get("road_group").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getMaintenanceHistoryDetails(id: number): Promise<IResponse<IMaintenanceHistoryDetailData>> {
		return (
			this.http
				// .get(`maintenance/history/${id}`)
				.get(`maintenance/${id}`)
				.then(this.handleResponse.bind(this))
				.catch(this.handleError.bind(this))
		)
	}

	public createMaintenanceHistoryGuarantee(
		id: number,
		params: IMaintenanceHistoryGuaranteeCreateParams
	): Promise<IResponse<IMaintenanceHistoryGuaranteeCreateParams>> {
		return this.http
			.post(`maintenance/history/${id}/road`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceBudgetMethodList(): Promise<IResponse<IMaintenanceHistoryBudgetsData[]>> {
		return this.http.get("maintenance/budgets").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getRoadList(roadId: number): Promise<IResponse<IRoadChildListData>> {
		return this.http.get(`road_group/${roadId}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getMaintenanceMethodList(): Promise<IResponse<IMaintenanceMethodListData[]>> {
		return this.http
			.get("maintenance/intervention_criteria")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenancePlanStatusList(id: number): Promise<IResponse<IMaintenanceHistoryPlanStatusData[]>> {
		return this.http
			.get(`maintenance/${id}/plan_stauts`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceStatusOptions(id: number): Promise<IResponse<IMaintenancePlanListData[]>> {
		return this.http
			.get(`maintenance/${id}/plan`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateMaintenanceHistoryGuarantee(
		id: number,
		historyId: number,
		params: IMaintenanceHistoryEditParams
	): Promise<IResponse<IMaintenanceHistoryEditParams>> {
		return this.http
			.put(`maintenance/history/${id}/road/${historyId}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceAttrachment(
		id: number,
		params: IMaintenanceHistoryFileParams
	): Promise<IResponse<IMaintenanceHistoryAttrachments[]>> {
		return this.http
			.get(`maintenance/${id}/attachments?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceHistoryPlanGraphRerport(
		id: number,
		planId: number[]
	): Promise<IResponse<IMaintenanceHistoryPlanGraph[]>> {
		const params = planId.length > 0 ? `?plan_id=${encodeURIComponent(planId.join(","))}` : ""
		return this.http
			.get(`maintenance/${id}/plan_progress_report${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceHistoryPlanTableReport(
		id: number,
		planId: number[]
	): Promise<IResponse<IPlanProgressGraphReportHistTableData>> {
		const params = planId.length > 0 ? `?plan_id=${encodeURIComponent(planId.join(","))}` : ""
		return this.http
			.get(`maintenance/${id}/plan_progress_table_report${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceDefault(id: number): Promise<IResponse<IMaintenanceDefaultData>> {
		return this.http.get(`maintenance/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getMaintenanceBudgetCriteria(): Promise<IResponse<IMaintenanceBudgetCriteria[]>> {
		return this.http.get("maintenance/budgets").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getMaintenanceYearList(): Promise<IResponse<number[]>> {
		return this.http.get(`maintenance/years`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public createMaintenance(
		params: IMaintenanceHistoryCreateRequest
	): Promise<IResponse<{ id: number; id_parent: number }>> {
		return this.http.post(`maintenance`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public updateMaintenanceData(id: number, params: IMaintenanceHistoryUpdateRequest): Promise<IResponse<{}>> {
		return this.http
			.put(`maintenance/${id}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createMaintenanceRoad(id: number, params: any): Promise<IResponse<{}>> {
		return this.http
			.post(`maintenance/${id}/road`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceRoadInfo(id: number, mRoadId: number): Promise<IResponse<IMaintenanceRoadData>> {
		return this.http
			.get(`maintenance/${id}/road/${mRoadId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateMaintenanceRoadInfo(
		id: number,
		mRoadId: number,
		params: IMaintenanceHistoryRoadsUpdateRequest
	): Promise<IResponse<{}>> {
		return this.http
			.put(`maintenance/${id}/road/${mRoadId}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createMaintenanceWarranty(id: number, params: IMaintenanceWarrantyCreateRequest): Promise<IResponse<{}>> {
		return this.http
			.post(`maintenance/${id}/road_history`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceWarranty(id: number, mRoadId: number): Promise<IResponse<IMaintenanceWarrantyData>> {
		return this.http
			.get(`maintenance/${id}/road_history/${mRoadId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateMaintenanceWarranty(id: number, mRoadId: number, params: any): Promise<IResponse<{}>> {
		return this.http
			.put(`maintenance/${id}/road_history/${mRoadId}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceDivision(): Promise<IResponse<IMaintenanceDivision[]>> {
		return this.http.get("maintenance/division").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getMaintenanceRoadGroup(): Promise<IResponse<IMaintenanceRoadGroup[]>> {
		return this.http
			.get("maintenance/road_dropdown_list")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceHistoryList(params: IMaintenanceHistorySearch): Promise<IResponse<IMaintenanceHistoryListData>> {
		return (
			this.http
				// .get(`/maintenance/history?${this.createParams(params)}`)
				.get(`/maintenance?${this.createParams(params)}`)
				.then(this.handleResponse.bind(this))
				.catch(this.handleError.bind(this))
		)
	}
}
