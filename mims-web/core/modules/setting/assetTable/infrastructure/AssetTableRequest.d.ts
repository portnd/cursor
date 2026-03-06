import type { TFileStatus } from "~~/core/shared/types/File"

export interface IRequestAssetTable {
	icon_filepath: file
	icon_filepath_status: TFileStatus
	data: IRequestDataAssetTable
}

export interface IRequestDataAssetTable {
	table_name: string
	table_label: string
	delete_columns: string[]
	asset_type: string
	asset_group: string
	geom_type: string
	line_color_code: string
	approver_id: string[]
	viewer_id: string[]
	columns: IRequestColumnAssetTable[]
}

export interface IRequestColumnAssetTable {
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
}
