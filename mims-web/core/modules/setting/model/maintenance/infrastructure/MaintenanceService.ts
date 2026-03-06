import { IAnalysiSurfacesRule, IAnalysisRule, IMaintenanceItemList, IMaintenanceRequest } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class MaintenanceService extends HttpService {
	public constructor() {
		super("/")
	}

	public get(): Promise<IResponse<IAnalysisRule>> {
		return this.http
			.get(`/settings/intervention_criteria`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(params: IMaintenanceItemList): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/intervention_criteria`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public post(params: IMaintenanceRequest): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/intervention_criteria`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public delete(id: number): Promise<IResponse<{}>> {
		return this.http
			.delete(`/settings/intervention_criteria/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRule(): Promise<IResponse<IAnalysiSurfacesRule>> {
		return this.http
			.get(`/settings/intervention_criteria/list`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public create(params: IAnalysisRule): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/intervention_criteria`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
