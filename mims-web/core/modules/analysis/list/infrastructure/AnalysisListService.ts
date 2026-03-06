import { IStrategicsList } from "./AnalysisListModel"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class AnalysisListService extends HttpService {
	public constructor() {
		super("/")
	}

	public getStrategics(): Promise<IResponse<IStrategicsList[]>> {
		return this.http
			.get("ref/maintenance_analysis_strategic")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
