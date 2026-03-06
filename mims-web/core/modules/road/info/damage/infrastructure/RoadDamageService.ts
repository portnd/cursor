import { IRoadDamageData, IDatum, IDataImport, ILane } from "./RoadDamageModel"
import { IParams } from "./RoadDamageRequest"
import { HttpService } from "~/core/shared/http/HttpService"
import { IResponse } from "~/core/shared/http/Response"

export class RoadDamageService extends HttpService {
	public constructor() {
		super("/")
	}

	public getLaneList(id: number): Promise<IResponse<ILane[]>> {
		return this.http
			.get(`roads/${id}/lane_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadDamageList(id: number): Promise<IResponse<IDatum[]>> {
		return this.http
			.get(`roads/${id}/damage_list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRoadDamageDetails(id: number, parentId: number): Promise<IResponse<IRoadDamageData>> {
		return this.http
			.get(`roads/${id}/damage_detail/${parentId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public createDamageData(id: number, params: IParams): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("lane_no", params.lane_no)
		param.append("surveyed_date", params.surveyed_date)
		param.append("damage_filename", params.damage_filename)
		param.append("damage_filename_status", params.damage_filename_status)
		param.append("image_filename", params.image_filename)
		param.append("image_filename_status", params.image_filename_status)

		return this.http
			.post(`roads/${id}/damage_import`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getDamageDefault(id: number, parentId: number): Promise<IResponse<IDataImport>> {
		return this.http
			.get(`roads/${id}/damage_import/${parentId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public updateDamageData(
		id: number,
		parentId: number,
		params: IParams
	): Promise<IResponse<{ id: number; id_parent: number }>> {
		const param = new FormData()
		param.append("lane_no", params.lane_no)
		param.append("surveyed_date", params.surveyed_date)
		param.append("damage_filename", params.damage_filename)
		param.append("damage_filename_status", params.damage_filename_status)
		param.append("image_filename", params.image_filename)
		param.append("image_filename_status", params.image_filename_status)

		return this.http
			.put(`roads/${id}/damage_import/${parentId}`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
