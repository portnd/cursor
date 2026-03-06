import { IAccount, IRequestAccount } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class AccountService extends HttpService {
	public constructor() {
		super("/")
	}

	public getAccountData(): Promise<IResponse<IAccount>> {
		return this.http.get(`/user_info`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}

	public updateAccount(params: IRequestAccount): Promise<IResponse<IAccount>> {
		return this.http.put(`/user_info`, params).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
