export type TFileStatus = "upload" | "not_edit" | "delete" | "no_file"

export interface IFile {
	data: data | null
	status: TFileStatus
	isUpload: boolean
}

interface data {
	id: string
	_relativePath: string
	lastModified: number
	lastModifiedDate: Object
	name: string
	size: number
	type: string
	webkitRelativePath: string | undefined
	file: File | undefined
	base64: string // Base64 encoded file
}

export interface IMultiFile {
	id: number | null
	file: string | undefined // Base64 encoded file
	file_name: string | undefined
	status: TFileStatus
}

export interface IFileMetadata {
	lastModifiedDate: Date
	name: string
	size: number
	type: string
}
