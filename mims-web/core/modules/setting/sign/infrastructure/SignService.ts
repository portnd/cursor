import { ISign, IRequestSign } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class SignService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestSign): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("name", params.name)
		param.append("abbr", params.abbr)
		param.append("image", params.image)
		param.append("image_status", params.image_status)

		return this.http
			.post(`/settings/signs`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<ISign>> {
		return this.http
			.get(`/settings/signs/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestSign): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("name", params.name)
		param.append("abbr", params.abbr)
		param.append("image", params.image)
		param.append("image_status", params.image_status)

		return this.http
			.put(`/settings/signs/${id}`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
