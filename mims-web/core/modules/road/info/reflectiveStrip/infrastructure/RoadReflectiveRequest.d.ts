import type { TFileStatus } from "~~/core/shared/types/File"

export interface IReflectivePostParams {
	line_no: number | null
	surveyed_date: string
	remarks: string
	csv_file: file
	csv_filename_status: TFileStatus
}

export interface IReflectiveUpdateParams {
	line_no: number | null
	surveyed_date: string
	remarks: string
	id_parent: number
	csv_file: file
	csv_filename_status: TFileStatus
}

export interface ICompareParams {
	years: number[]
	lines: number[]
}
