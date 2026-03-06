import { ISwalDelete, deleteItem } from "~~/core/modules/common/deleteItem/ui"
import { downloadFile } from "~~/core/modules/common/downloadFile/ui"

export type TIcon = "success" | "error" | "warning" | "info" | "question"

export interface ISwal {
	title?: string
	message?: string
	html?: string
	type?: TIcon
	callBack?(): void
}

export const showAlert = (options: ISwal): void => {
	const { $swal }: any = useNuxtApp()

	$swal
		.fire({
			title: options.title,
			icon: options.type,
			text: options.message,
			html: options.html,
			confirmButtonText: "ตกลง",
			confirmButtonColor: "#FDB833",
			showCancelButton: options.type === "question",
			cancelButtonText: "ยกเลิก",
			reverseButtons: true,
		})
		.then((result: any) => {
			if (options.type === "question") {
				if (result.isConfirmed) {
					if (options.callBack !== undefined) {
						options.callBack()
					}
				}
			} else if (options.callBack !== undefined) {
				options.callBack()
			}
		})
}

export const useDeleteItem = (options: ISwalDelete): void => {
	const { name, url, callBack, showAlert = true } = options
	deleteItem({ name, url, callBack, showAlert })
}

/** fileType: 'html' | 'pdf' = open in new tab; 'excel' = download; undefined = default download */
export const useDownloadFile = (title: string, url: string, fileType?: "html" | "pdf" | "excel"): void => {
	downloadFile(title, url, fileType)
}
