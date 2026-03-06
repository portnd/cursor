import { IRucListData, IRucDataTable } from "./RoadUserCostRucModel.d"
import { IRucParentParams } from "./RoadUserCostRucRequest"
import { IResponse } from "~/core/shared/http/Response"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RoadUserCostRucService extends HttpService {
	public constructor() {
		super("/")
	}

	public getRucList(): Promise<IResponse<IRucListData[]>> {
		return this.http
			.get("ref/road_user_cost/ruc")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRucData(type: string): Promise<IResponse<IRucDataTable>> {
		return this.http
			.get(`settings/road_user_cost/ruc/${type}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public postRucParams(type: string, params: IRucParentParams): Promise<IResponse<IRucParentParams>> {
		return this.http
			.post(`settings/road_user_cost/ruc/${type}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
