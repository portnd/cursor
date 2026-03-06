export interface IDashboardAsset {
	asset_group: IAssetGroup
	asset_list: IAssetList[]
}

export interface IAssetGroup {
	id: number
	name: string
}

export interface IAssetList {
	asset: IAsset
	is_range: boolean
	value: number
}

export interface IAsset {
	default_color: string
	default_icon_url: string
	id: number
	name: string
	thumbnail_icon_filepath: string
}

export interface IAssetLocation {
	id: number
	name: string
	type: string
	icon_filepath: string
	thumbnail_icon_filepath: string
	line_color_code: string
	wkt: string
}
