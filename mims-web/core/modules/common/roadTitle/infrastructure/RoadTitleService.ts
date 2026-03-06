import { IRoad } from "../../../road/info/summary/infrastructure/RoadDetailModel.d"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RoadTitleService extends HttpService {
	public constructor() {
		super("/")
	}

	public getRoad(id: number): Promise<IResponse<IRoad>> {
		return this.http.get(`roads/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
