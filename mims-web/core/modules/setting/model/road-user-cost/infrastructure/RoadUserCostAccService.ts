import { IRucAccChanceParams, IRucAccLossParams } from "./RoadUserCostAccRequest"
import { IRucAccChanceData, IRucAccMasterData, IRucAccLossData } from "./RoadUserCostAccModel"
import { IResponse } from "~/core/shared/http/Response"
import { HttpService } from "~/core/shared/http/HttpService"

export class RoadUserCostAccService extends HttpService {
	public constructor() {
		super("/")
	}

	public getMasterAccident(): Promise<IResponse<IRucAccMasterData[]>> {
		return this.http
			.get("ref/road_user_cost/acc")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getChanceOfAccident(roadGroupId: number): Promise<IResponse<IRucAccChanceData>> {
		return this.http
			.get(`settings/road_user_cost/acc/chance_of_accident/${roadGroupId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public postChanceOfAccident(params: IRucAccChanceParams): Promise<IResponse<IRucAccChanceParams>> {
		return this.http
			.post("settings/road_user_cost/acc/chance_of_accident", params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getLossValueAccident(): Promise<IResponse<IRucAccLossData>> {
		return this.http
			.get("settings/road_user_cost/acc/loss_value")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public postLossValueAccident(params: IRucAccLossParams): Promise<IResponse<IRucAccLossParams>> {
		return this.http
			.post("settings/road_user_cost/acc/loss_value", params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
