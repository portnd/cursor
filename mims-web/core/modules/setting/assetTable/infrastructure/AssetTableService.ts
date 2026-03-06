import { IDataAssetTable, IRequestAssetTable } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class AssetTableService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestAssetTable): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("icon_filepath", params.icon_filepath)
		param.append("icon_filepath_status", params.icon_filepath_status)
		param.append("data", JSON.stringify(params.data))

		return this.http
			.post(`/settings/asset_tables`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: string): Promise<IResponse<IDataAssetTable>> {
		return this.http
			.get(`/settings/asset_tables/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(id: string, params: IRequestAssetTable): Promise<IResponse<{}>> {
		const param = new FormData()
		param.append("icon_filepath", params.icon_filepath)
		param.append("icon_filepath_status", params.icon_filepath_status)
		param.append("data", JSON.stringify(params.data))

		return this.http
			.put(`/settings/asset_tables/${id}`, param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
