import { IMaintenanceSequence } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class MaintenanceSequenceService extends HttpService {
	public constructor() {
		super("/")
	}

	public get(): Promise<IResponse<IMaintenanceSequence>> {
		return this.http
			.get(`/settings/intervention_criteria/sequence`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public post(params: IMaintenanceSequence): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/intervention_criteria/sequence`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
