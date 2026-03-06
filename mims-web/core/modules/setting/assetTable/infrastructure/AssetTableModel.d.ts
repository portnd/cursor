export interface IAssetTable {
	asset_id: number
	table_name: string
	table_label: string
	asset_type: string
	asset_group: IDataOtherList
	geom_type: string
	line_color_code: string
	icon_filepath: string | null
	approver: string[]
	viewer: string[]
	columns: IColumnAssetTable[]
	delete_columns: string[]
}

export interface IDataAssetTable {
	asset_id: number
	table_name: string
	table_label: string
	asset_type: string
	asset_group: IDataOtherList
	geom_type: string | number
	line_color_code: string
	approver: string[]
	viewer: string[]
	columns: IColumnAssetTable[]
	delete_columns: string[]
}

export interface IDataOtherList {
	id: string
	name: string
	can_delete: boolean
}

export interface IColumnAssetTable {
	column_id: number
	column_name: string
	table_name_ref: string
	component_title: string
	component_type: string
	is_required: boolean
	is_visible_view: boolean
	is_visible_edit: boolean
	is_visible_report: boolean
	is_ban: boolean
	is_new: boolean
	is_edit: boolean
}
