import { IMaintenanceProjectsData } from "./RoadSummaryModel"
import {
	IRoadSummary,
	IRequestRoadSummary,
	IRoad,
	IRequestRoad,
	ISurfaceIcon,
	ILaneList,
	IConditionCompareAverage,
	ITrafficModel,
	ITrafficDetail,
	IRequestTraffic,
} from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RoadSummaryService extends HttpService {
	public constructor() {
		super("/")
	}

	public getSurface(id: number): Promise<IResponse<IRoadSummary[]>> {
		const data = this.http
			.get(`road/surface?road_id=${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
		return data
	}

	public getSurfaceIcon(id: number): Promise<IResponse<ISurfaceIcon[]>> {
		const data = this.http
			.get(`road/surface/icon/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
		return data
	}

	public post(params: IRequestRoadSummary): Promise<IResponse<{}>> {
		return this.http
			.post(`road/surface`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoads(id: number): Promise<IResponse<IRoad>> {
		return this.http.get(`roads/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public deleteRoads(id: number): Promise<IResponse<{}>> {
		return this.http.delete(`roads/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public updateRoads(id: number, params: IRequestRoad): Promise<IResponse<{}>> {
		const formData = new FormData()
		formData.append("road_code", params.road_code)
		formData.append("name", params.name)
		formData.append("road_group_id", params.road_group_id?.toString())
		formData.append("origin", params.origin)
		formData.append("destination", params.destination)
		formData.append("km_start", params.km_start?.toString())
		formData.append("km_end", params.km_end?.toString())
		formData.append("road_color_code", params.road_color_code)
		formData.append("ref_road_type_id", params.ref_road_type_id?.toString())
		formData.append("register_date", params.register_date)
		formData.append(
			"center_line_shape_file",
			params.center_line_shape_file?.status === "not_edit"
				? "undefined"
				: params.center_line_shape_file?.data?.file ?? ""
		)
		formData.append("center_line_shape_file_status", params.center_line_shape_file_status)
		formData.append(
			"center_lane_shape_file",
			params.center_lane_shape_file?.status === "not_edit"
				? "undefined"
				: params.center_lane_shape_file?.data?.file ?? ""
		)
		formData.append("center_lane_shape_file_status", params.center_lane_shape_file_status)
		formData.append("remark", params.remark)
		formData.append("ramp_id", params.ramp_id.toString())
		formData.append("year_construction_completed", params.year_construction_completed.toString())

		return this.http
			.put(`roads/${id}`, formData)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getLaneList(id: number): Promise<IResponse<ILaneList[]>> {
		return this.http
			.get(`roads/${id}/lane_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConditionCompareAverage(id: number, laneId: number): Promise<IResponse<IConditionCompareAverage>> {
		return this.http
			.get(`roads/${id}/condition_compare_average/${laneId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceYears(id: number): Promise<IResponse<number[]>> {
		return this.http
			.get(`roads/${id}/maintenance_year`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getMaintenanceProjects(id: number, year: number): Promise<IResponse<IMaintenanceProjectsData[]>> {
		const params = !year ? "" : `?year=${year}`
		return this.http
			.get(`roads/${id}/maintenance${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getTrafficRevision(id: number): Promise<IResponse<ITrafficModel[]>> {
		return this.http
			.get(`roads/${id}/volume_aadt/revision`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getTrafficDetail(id: number, revid: number): Promise<IResponse<ITrafficDetail>> {
		return this.http
			.get(`roads/${id}/volume_aadt/${revid}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createTraffic(id: number, params: IRequestTraffic): Promise<IResponse<{ id: number; id_parent: number }>> {
		return this.http
			.post(`roads/${id}/volume_aadt`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateTraffic(
		id: number,
		idParent: number,
		aadtId: number,
		params: IRequestTraffic
	): Promise<IResponse<{ id: number; id_parent: number }>> {
		return this.http
			.put(`roads/${id}/volume_aadt/${idParent}/aadt/${aadtId}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public deleteTraffic(id: number, revid: number): Promise<IResponse<{}>> {
		return this.http
			.delete(`roads/${id}/volume_aadt/${revid}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
