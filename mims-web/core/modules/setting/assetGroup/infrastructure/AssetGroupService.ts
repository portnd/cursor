import { IAssetGroup, IRequestAssetGroup } from "./index"
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class AssetGroupService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IRequestAssetGroup): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/asset_groups`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IAssetGroup>> {
		return this.http
			.get(`/settings/asset_groups/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public put(id: number, params: IRequestAssetGroup): Promise<IResponse<{}>> {
		return this.http
			.put(`/settings/asset_groups/${id}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
