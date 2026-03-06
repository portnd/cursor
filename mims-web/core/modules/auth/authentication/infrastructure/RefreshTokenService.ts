import { IRefreshToken, IRequestRefreshToken } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RefreshTokenService extends HttpService {
	public constructor() {
		super("/")
	}

	public refreshToken(params: IRequestRefreshToken): Promise<IResponse<IRefreshToken>> {
		return this.http
			.post(`/auth/refresh_token`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
