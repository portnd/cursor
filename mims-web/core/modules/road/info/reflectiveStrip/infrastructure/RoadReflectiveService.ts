import {
	IReflectiveListData,
	ILine,
	ICompareYearData,
	ICompareLaneData,
	IReflectiveData,
	IReflectivityDetails,
} from "./RoadReflectiveModel"
import { IReflectivePostParams, ICompareParams, IReflectiveUpdateParams } from "./RoadReflectiveRequest"
import { HttpService } from "~/core/shared/http/HttpService"
import { IResponse } from "~/core/shared/http/Response"

export class RoadReflectiveService extends HttpService {
	public constructor() {
		super("/")
	}

	public getReflectivityList(id: number): Promise<IResponse<IReflectiveListData[]>> {
		return this.http
			.get(`roads/${id}/retro_reflectivity/list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createReflectiveData(roadId: number, params: IReflectivePostParams): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("line_no", params.line_no)
		param.append("surveyed_date", params.surveyed_date)
		param.append("remarks", params.remarks)
		param.append("csv_file", params.csv_file)
		param.append("csv_file_status", params.csv_filename_status)

		return this.http
			.post(`roads/${roadId}/retro_reflectivity`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getLaneList(id: number): Promise<IResponse<ILine[]>> {
		return this.http
			.get(`roads/${id}/retro_reflectivity/line_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateReflectiveData(
		roadId: number,
		params: IReflectiveUpdateParams
	): Promise<IResponse<{ id: number; id_parent: number }>> {
		const param = new FormData()
		param.append("line_no", params.line_no.toString())
		param.append("surveyed_date", params.surveyed_date)
		param.append("remarks", params.remarks)
		param.append("id_parent", params.id_parent.toString())
		param.append("csv_file", params.csv_file)
		param.append("csv_file_status", params.csv_filename_status)

		return this.http
			.put(`roads/${roadId}/retro_reflectivity/${params.id_parent}`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReflectivityData(idParent: number): Promise<IResponse<IReflectiveData>> {
		return this.http
			.get(`roads/retro_reflectivity/${idParent}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getReflectivityDetails(idParent: number, typeId: number | null): Promise<IResponse<IReflectivityDetails>> {
		const rangeParam = typeId ? `?range_type=${typeId}` : ""
		return this.http
			.get(`roads/retro_reflectivity/details/${idParent}${rangeParam}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public compareLane(roadID: number, params: ICompareParams): Promise<IResponse<ICompareLaneData[]>> {
		return this.http
			.get(`roads/${roadID}/retro_reflectivity/compare_line?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public compareYear(roadID: number, params: ICompareParams): Promise<IResponse<ICompareYearData[]>> {
		return this.http
			.get(`roads/${roadID}/retro_reflectivity/compare_year?${this.createParams(params)}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
