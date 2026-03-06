export interface IRoadsAssetItem {
	id: number
	name: string
	items: IitemsRoadsAssetItem[]
}

interface IitemsRoadsAssetItem {
	id: number
	name: string
	geom_type: string
	icon_filepath: string
}

export interface IRoadAssetRevision {
	id: number
	id_parent: number
	updated_date: string
	revision_no: number
	is_exclusive_lock: boolean
	status: string
}

export interface IRoadAssetTableTemplate {
	column_data_type: string
	column_name: string
	table_name_ref: string
	table_name: string
	component_title: string
	component_type: string
	is_required: boolean
	value: any
}

export interface Km {
	km: string
	km_end: string
	km_start: string
	lane_no: string
	the_geom: string
	road_id: string
	type: string
}

export interface IRoadAssetDetail {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IRoadAssetDetailItem
}

export interface IRoadAssetDetailItem {
	id: number
	id_parent: number
	updated_date: string
	revision: number
	status: string
	reject_reason: string
	updated_by: updated
	permissions: permission
	is_exclusive_lock: boolean
	road_assets: IRoadAssetItem[]
	data_columns: IRoadColumn[]
	icon_filepath: string
	thumbnail_icon_filepath: string
	line_color_code: string
	status_code: string
}

interface permission {
	can_edit: boolean
	can_delete: boolean
	can_approve: boolean
	can_send: boolean
	can_reject: boolean
}

interface updated {
	uid: string
	user_name: string
	full_name: string
	department: department
	profile_picture: string
	// id: number
	// email: string
	// department_id: number
	// firstname: string
	// lastname: string
	// profile_img_path: string
	// status: string
	// tel: number
	// created_by: number
	// updated_by: number
}

interface department {
	id: number
	name: string
	can_delete: boolean
}

export interface IRoadAssetItem {
	clearance: number
	geom_2d: location
	geom_3d: string
	id_parent: number
	km: number
	no: string
	note: string
	ref_asset_location: IrefAsset
	ref_asset_position: IrefAsset
	ref_asset_sign: IrefAsset
	ref_asset_sign_image: IrefAsset
	sign_area: number
	sign_caption: string
	sign_filepath: string
	surveyed_date: string
	thumbnail_sign_filepath: string
	[key: string]: any
}

interface IRefAsset {
	id: number
	name: string
	filepath: string
	thumbnail_filepath: string
	abbr: string
}

export interface IRoadColumn {
	key: string
	value: IRoadAssetColumn
}

export interface IRoadAssetColumn {
	component_title: string
	component_type: string
	column_data_type: string
}

export interface location {
	lat: string
	lon: string
}

export interface IAssetDetail {
	geom_type_id: number
	geom_type_name: string
	department_manage: number[]
	color: string
}
