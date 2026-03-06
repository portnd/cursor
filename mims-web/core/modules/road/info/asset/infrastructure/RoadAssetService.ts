import {
	IRoadAssetRevision,
	IRoadAssetTableTemplate,
	IRoadsAssetItem,
	Km,
	IRequestRoadsAsset,
	IRequestAssetDetail,
	IRoadAssetDetail,
	IAssetDetail,
} from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class RoadsAssetService extends HttpService {
	public constructor() {
		super("/")
	}

	public getMenu(assetType: string): Promise<IResponse<IRoadsAssetItem[]>> {
		return this.http
			.get(`/roads/menu?assetType=${assetType}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getRevision(id: number, refAssetTableId: number): Promise<IResponse<IRoadAssetRevision[]>> {
		return this.http
			.get(`/roads/${id}/asset_revision_list?ref_asset_table_id=${refAssetTableId}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getTemplate(
		refAssetTableId: number,
		action: string,
		assetObjectId: number
	): Promise<IResponse<IRoadAssetTableTemplate[]>> {
		return this.http
			.get(
				`/roads/asset_edit_template?ref_asset_table_id=${refAssetTableId}&action=${action}&asset_object_id=${assetObjectId}`
			)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getKm(id: number, geom: string): Promise<IResponse<Km>> {
		return this.http
			.get(`/roads/${id}/km?geom=${geom}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public postAssetRoad(id: number, params: IRequestRoadsAsset): Promise<IResponse<Km>> {
		return this.http
			.post(`/roads/${id}/asset`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public putAssetRoad(id: number, idParent: number, params: IRequestRoadsAsset): Promise<IResponse<Km>> {
		return this.http
			.put(`/roads/${id}/asset/${idParent}`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAssetList(
		RoadId: number,
		roadAssetId: number,
		params: IRequestAssetDetail
	): Promise<IResponse<IRoadAssetDetail>> {
		return this.http
			.get(
				`/roads/${RoadId}/asset_details/${roadAssetId}?ref_asset_table_id=${params.ref_asset_table_id}&page=${params.page}&limit=${params.limit}
				`
			)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public putAssetCancel(idParent: number): Promise<IResponse<Km>> {
		return this.http
			.put(`/roads/asset_cancel/${idParent}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public putAssetConfirm(idParent: number): Promise<IResponse<Km>> {
		return this.http
			.put(`/roads/asset_confirm/${idParent}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAssetDetail(idAsset: number): Promise<IResponse<IAssetDetail>> {
		return this.http
			.get(`/roads/asset_table/${idAsset}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
