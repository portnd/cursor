import type { TFileStatus } from "~~/core/shared/types/File"

export interface IParams {
	lane_no: string
	surveyed_date: string
	damage_filename: file
	damage_filename_status: TFileStatus
	image_filename: file
	image_filename_status: TFileStatus
}
