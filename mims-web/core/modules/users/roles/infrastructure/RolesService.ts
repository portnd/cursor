import { IRolesDetail, IRequestUpdateRoles, IAccessRoles } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RolesService extends HttpService {
	public constructor() {
		super("/")
	}

	public createRoles(params: IRequestUpdateRoles): Promise<IResponse<{}>> {
		return this.http.post(`/roles`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IRolesDetail>> {
		return this.http.get(`/roles/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestUpdateRoles): Promise<IResponse<{}>> {
		return this.http.put(`/roles/${id}`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getAccess(): Promise<IResponse<IAccessRoles[]>> {
		return this.http
			.get(`/roles/access_control`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
