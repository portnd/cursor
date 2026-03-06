import {
	IConditionListData,
	ILane,
	IGraphData,
	ICompareYearData,
	ICompareLaneData,
	IConditionData,
	IRoadConditionList,
	IRoadConditionDetails,
} from "./RoadConditionModel"
import { IConditionPostParams, IConditionPutParams, ICompareParams } from "./RoadConditionRequest"
import { HttpService } from "~/core/shared/http/HttpService"
import { IResponse } from "~/core/shared/http/Response"

export class RoadConditionService extends HttpService {
	public constructor() {
		super("/")
	}

	public conditionList(id: number): Promise<IResponse<IConditionListData[]>> {
		return this.http
			.get(`roads/${id}/condition_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public conditionDetails(idParent: number, type: string): Promise<IResponse<IGraphData>> {
		return this.http
			.get(`roads/condition_details/${idParent}?condition_type=${type}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public compareLane(roadID: number, params: ICompareParams): Promise<IResponse<ICompareLaneData[]>> {
		return this.http
			.get(`roads/${roadID}/condition_compare_lane?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public compareYear(roadID: number, params: ICompareParams): Promise<IResponse<ICompareYearData[]>> {
		return this.http
			.get(`roads/${roadID}/condition_compare_year?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionList(id: number): Promise<IResponse<IRoadConditionList[]>> {
		return this.http
			.get(`roads/${id}/condition_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getLaneList(id: number): Promise<IResponse<ILane[]>> {
		return this.http
			.get(`roads/${id}/lane_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateConditionData(
		roadId: number,
		params: IConditionPutParams
	): Promise<IResponse<{ id: number; id_parent: number }>> {
		const param = new FormData()
		param.append("lane_no", params.lane_no.toString())
		param.append("surveyed_date", params.surveyed_date)
		param.append("remarks", params.remarks)
		param.append("id_parent", params.id_parent.toString())
		param.append("iri_filename", params.iri_filename)
		param.append("iri_filename_status", params.iri_filename_status)
		param.append("image_filename", params.image_filename)
		param.append("image_filename_status", params.image_filename_status)

		return this.http
			.put(`roads/${roadId}/condition/${params.id_parent}`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createConditionData(roadId: number, params: IConditionPostParams): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("lane_no", params.lane_no)
		param.append("surveyed_date", params.surveyed_date)
		param.append("remarks", params.remarks)
		param.append("iri_filename", params.iri_filename)
		param.append("iri_filename_status", params.iri_filename_status)
		param.append("image_filename", params.image_filename)
		param.append("image_filename_status", params.image_filename_status)

		return this.http
			.post(`roads/${roadId}/condition`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionData(idParent: number): Promise<IResponse<IConditionData>> {
		return this.http
			.get(`roads/condition/${idParent}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionDetails(idParent: number, rangeId: number): Promise<IResponse<IRoadConditionDetails>> {
		return this.http
			.get(`roads/condition_details/${idParent}?condition_range_type=${rangeId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
