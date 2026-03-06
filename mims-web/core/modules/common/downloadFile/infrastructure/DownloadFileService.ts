import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class DownloadFileService extends HttpService {
	public constructor() {
		super("/")
	}

	public download(url: string): Promise<IResponse<{}>> {
		return this.http.get(`${url}`).then(this.handleResponse.bind(this)).catch(this.handleError.bind(this))
	}
}
