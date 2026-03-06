import { DeleteItemService } from "../infrastructure"
import httpStatusCode from "~~/core/shared/http/HttpStatusCode"

export interface ISwalDelete {
	name: string
	url?: string
	callBack?(): void
	showAlert?: boolean
}

export const deleteItem = (options: ISwalDelete): void => {
	const { $swal }: any = useNuxtApp()
	const itemName =
		options.name === ""
			? " "
			: ` "<span class="text-truncate" style="max-width: 250px; display: table-cell;">${options.name}</span>" `

	$swal
		.fire({
			html: `<div class="text-center mt-2">
        <img src="/images/icons/gif/trash.gif" alt="ลบข้อมูล" width="65" />
        <h3 class="swal2-title mb-4">ลบข้อมูล ?</h3>
        <span>คุณต้องการลบข้อมูล${itemName}หรือไม่</span>
      </div>`,
			showCancelButton: !0,
			confirmButtonColor: "#FDB833",
			confirmButtonText: "ตกลง",
			cancelButtonText: "ยกเลิก",
			reverseButtons: true,
		})
		.then(async (result: any) => {
			if (result.isConfirmed) {
				if (options.url) {
					const deleteItemService = new DeleteItemService()
					const res = await deleteItemService.delete(options.url)

					if (res.status === false) {
						useHandlerError(res.code, res.error, { showAlert: true })
					} else if (options.callBack !== undefined) {
						if (options.showAlert) {
							useHandlerSuccess(httpStatusCode.DELETED, {
								showAlert: true,
								fn: options.callBack,
							})
						} else {
							options.callBack()
						}
					} else {
						// กรณีไม่ใช้ callBack
						useHandlerSuccess(res.code, { showAlert: true })
					}
				} else if (options.callBack !== undefined) {
					options.callBack()
				}
			}
		})
}
