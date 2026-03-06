import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DeleteItemService extends HttpService {
	public constructor() {
		super("/")
	}

	public delete(url: string): Promise<IResponse<{}>> {
		return this.http.delete(`${url}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
