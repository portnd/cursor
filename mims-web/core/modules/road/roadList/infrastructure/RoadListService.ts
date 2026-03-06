import {
	// IRoadData,
	// IRequestRoadList,
	IRoadList,
	IRequestRoadInit,
	ICreateRoad,
} from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RoadListService extends HttpService {
	public constructor() {
		super("/")
	}

	// public getRoads(params: IRequestRoadList): Promise<IResponse<IRoadData[]>> {
	// 	return this.http
	// 		.get(`roads?${this.createParams(params)}`)
	// 		.then(this.handleResponse.bind(this))
	// 		.catch(this.handleError.bind(this))
	// }

	public getRoads(quries: any): Promise<IResponse<IRoadList[]>> {
		const params = this.createParams(quries)
		return this.http.get(`roads?${params}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getRoadinit(id: string, level: string, roadType: string): Promise<IResponse<IRequestRoadInit>> {
		const params = !roadType ? "" : `?ref_road_type_id=${roadType}`
		return this.http
			.get(`roads/init/${id}/${level}${params}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createRoad(params: ICreateRoad): Promise<IResponse<{ id: number }>> {
		const formData = new FormData()
		if (params.road_level === 1) {
			formData.append("road_section_id", params.road_section_id?.toString())
		} else {
			formData.append("name", params.name?.toString())
			formData.append("road_id", params.road_id?.toString())
			formData.append("ramp_id", params.ramp_id?.toString())
		}
		formData.append("road_level", params.road_level?.toString())
		formData.append("km_start", params.km_start?.toString())
		formData.append("km_end", params.km_end?.toString())
		formData.append("year_construction_completed", params.year_construction_completed?.toString())
		formData.append("road_color_code", params.road_color_code)
		formData.append("ref_road_type_id", params.ref_road_type_id?.toString())
		formData.append("register_date", params.register_date)
		formData.append("center_lane_shape_file", params.center_lane_shape_file)
		formData.append("center_line_shape_file", params.center_line_shape_file)
		formData.append("remark", params.remark)

		console.log("formData :", formData)

		return this.http.post("roads", formData).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
