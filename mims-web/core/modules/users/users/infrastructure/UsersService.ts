import { IRequestUsers, IUsersRolesData, IUsersDepartmentsData, IDefaultUsersData } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class UsersService extends HttpService {
	public constructor() {
		super("/")
	}

	public createUser(params: IRequestUsers): Promise<IResponse<{}>> {
		return this.http.post(`/users`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getDefaultUser(id: number): Promise<IResponse<IDefaultUsersData>> {
		return this.http.get(`/users/${id}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public updateUser(id: number, params: IRequestUsers): Promise<IResponse<{}>> {
		return this.http.put(`/users/${id}`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getRolesList(): Promise<IResponse<IUsersRolesData>> {
		return this.http.get(`roles`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public getDepartmentsList(): Promise<IResponse<IUsersDepartmentsData>> {
		return this.http.get("settings/departments").then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
