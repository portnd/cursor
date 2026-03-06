import type { TFileStatus } from "~~/core/shared/types/File"

export interface IConditionPostParams {
	lane_no: string
	surveyed_date: string
	remarks: string
	iri_filename: file
	iri_filename_status: TFileStatus
	image_filename: file
	image_filename_status: TFileStatus
}

export interface IConditionPutParams {
	lane_no: number
	surveyed_date: string
	remarks: string
	id_parent: number
	iri_filename: file
	iri_filename_status: TFileStatus
	image_filename: file
	image_filename_status: TFileStatus
}

export interface ICompareParams {
	years: Array
	lanes: Array
	condition_type: string
}
