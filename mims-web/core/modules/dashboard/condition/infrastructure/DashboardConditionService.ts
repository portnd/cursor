import {
	ILine,
	ILane,
	IReflectiveListData,
	IReflectivityDetails,
	IRoadConditionList,
	IRoadConditionDetails,
	ICompareYearData,
	ICompareLaneData,
	IDashboardCondition,
	IDashboardConditionMap,
} from "./DashboardConditionModel"
import { ICompareParams } from "./DashboardConditionResponse"

import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DashboardConditionService extends HttpService {
	public constructor() {
		super("/")
	}

	public getLaneList(id: number): Promise<IResponse<ILine[]>> {
		return this.http
			.get(`roads/${id}/retro_reflectivity/line_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReflectivityList(id: number): Promise<IResponse<IReflectiveListData[]>> {
		return this.http
			.get(`roads/${id}/retro_reflectivity/list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReflectivityDetails(idParent: number, typeId: number): Promise<IResponse<IReflectivityDetails>> {
		return this.http
			.get(`roads/retro_reflectivity/details/${idParent}?range_type=${typeId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionList(id: number): Promise<IResponse<IRoadConditionList[]>> {
		return this.http
			.get(`roads/${id}/condition_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionDetails(
		idParent: number,
		params: { condition_range_type: string }
	): Promise<IResponse<IRoadConditionDetails>> {
		const query = this.createParams(params)
		return this.http
			.get(`roads/condition_details/${idParent}?${query}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public compareYear(roadID: number, params: ICompareParams): Promise<IResponse<ICompareYearData[]>> {
		return this.http
			.get(`roads/${roadID}/condition_compare_year?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public compareLane(roadID: number, params: ICompareParams): Promise<IResponse<ICompareLaneData[]>> {
		return this.http
			.get(`roads/${roadID}/condition_compare_lane?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getLaneLists(id: number): Promise<IResponse<ILane[]>> {
		return this.http
			.get(`roads/${id}/lane_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getCondition(params: any): Promise<IResponse<IDashboardCondition>> {
		return this.http
			.get(`/dashboard/condition?${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getCondition_map(params: any): Promise<IResponse<IDashboardConditionMap>> {
		return this.http
			.get(`/dashboard/condition_map?${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
