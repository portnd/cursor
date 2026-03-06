import { IFile, IMultiFile } from "~/core/shared/types/File"

export const fileToBase64 = async (file: File): Promise<string> => {
	const reader = new FileReader()
	const result = await new Promise<string>((resolve, reject) => {
		reader.onload = () => {
			resolve(reader.result as string)
		}
		reader.onerror = (error) => reject(error)
		reader.readAsDataURL(file)
	})
	return result
}

export const getMimeTypes = (): Record<string, string> => {
	return {
		"image/jpeg": "jpeg",
		"image/jpg": "jpg",
		"image/png": "png",
		"image/gif": "gif",
		"image/webp": "webp",
		"image/svg+xml": "svg",
		"application/pdf": "pdf",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
		"application/msword": "doc",
		"application/vnd.ms-excel": "xls",
		"application/vnd.ms-powerpoint": "ppt",
		"text/plain": "txt",
		"text/csv": "csv",
		"text/html": "html",
		"text/css": "css",
		"text/javascript": "js",
		"application/json": "json",
		"application/zip": "zip",
		"application/x-rar-compressed": "rar",
		"application/x-rar": "rar",
		"application/x-zip-compressed": "zip",
		".dwg": "dwg",
	}
}

export const getExtensionFromMimeType = (mimeType: string): string => {
	const mimeToExtMap = getMimeTypes()

	const extension = mimeToExtMap[mimeType.toLowerCase()]

	return extension || ""
}

export const checkFileExist = async (url: string): Promise<boolean> => {
	try {
		const response = await fetch(url)
		if (!response.ok) {
			return false
		}
		return true
	} catch (error) {
		return false
	}
}

export const getFileExtension = (filename: string): string => {
	return filename.substring(filename.lastIndexOf(".") + 1)
}

// files: [
// 	{
// 		// create
// 		id: null,
// 		file: "base64",
// 		file_name: "file_name.jpg",
// 		status: "upload",
// 	},
// 	{
// 		// delete
// 		id: null,
// 		file: undefined,
// 		file_name: "",
// 		status: "delete",
// 	},
// 	{
// 		// not edit
// 		id: 0,
// 		file: undefined,
// 		file_name: "",
// 		status: "not_edit",
// 	},
// ]

interface IBeforeFile {
	id: number
	path: string
	file_name: string
	file_type?: string
}

// beforeFiles คือ ไฟล์ทั้งหมดก่อนที่จะโดนกระทำ
// afterfiles คือ ไฟล์ทั้งหมดที่รับค่าจาก VUpload แบบ multiple
export const toFiles = (beforeFiles: IBeforeFile[], afterfiles: IFile[]): Array<IMultiFile> => {
	const files: Array<IMultiFile> = []
	// กรณีสร้างครั้งแรก
	if (beforeFiles.length === 0) {
		afterfiles.forEach((afterfile) => {
			files.push({
				id: null,
				file: afterfile.data?.base64,
				file_name: afterfile.data?.name,
				status: "upload",
			})
		})
	} else {
		// กรณี not_edit, upload
		afterfiles.forEach((afterfile) => {
			const oldValue = beforeFiles.find((obj) => obj.path.includes(String(afterfile.data?.name)))
			files.push({
				id: oldValue?.id ? oldValue.id : null,
				file: afterfile.data?.base64,
				file_name: oldValue?.file_name ? oldValue.file_name : afterfile.data?.name,
				status: afterfile.status,
			})
		})

		// กรณี delete
		beforeFiles.forEach((beforeFile) => {
			const findFile = files.find((obj) => obj.id === beforeFile.id)
			if (!findFile) {
				files.push({
					id: beforeFile.id,
					file: "",
					file_name: beforeFile.file_name,
					status: "delete",
				})
			}
		})
	}

	return files
}
