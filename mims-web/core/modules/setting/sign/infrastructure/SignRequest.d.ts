import type { TFileStatus } from "~~/core/shared/types/File"

export interface IRequestSign {
	name: string
	abbr: string
	image: file
	image_status: TFileStatus
}
