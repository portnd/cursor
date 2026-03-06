import { IRoadWorkEffectACData, IRoadWorkEffectCCData } from "./RoadWorkEffectModel"
import { IRoadWorkEffectACParams, IRoadWorkEffectCCParams } from "./RoadWorkEffectRequest"
import { IResponse } from "~/core/shared/http/Response"
import { HttpService } from "~/core/shared/http/HttpService"

export class RoadWorkEffectService extends HttpService {
	public constructor() {
		super("/")
	}

	public getAsphalt(): Promise<IResponse<IRoadWorkEffectACData>> {
		return this.http
			.get(`settings/road_work_effect/asphalt`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public postAsphalt(params: IRoadWorkEffectACParams): Promise<IResponse<IRoadWorkEffectACParams>> {
		return this.http
			.post(`settings/road_work_effect/asphalt`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getConcrete(): Promise<IResponse<IRoadWorkEffectCCData>> {
		return this.http
			.get("settings/road_work_effect/concrete")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public postConcrete(params: IRoadWorkEffectCCParams): Promise<IResponse<IRoadWorkEffectCCParams>> {
		return this.http
			.post("settings/road_work_effect/concrete", params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
