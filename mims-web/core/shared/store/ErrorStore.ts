import HttpStatusCode from "~~/core/shared/http/HttpStatusCode"

export interface IError {
	title: string
	statusCode: string
	messageEn: string
	messageTh: string
	imagePath: string
}

export const useErrorStore = defineStore("error", {
	state: (): IError => ({
		title: "",
		statusCode: "",
		messageEn: "",
		messageTh: "",
		imagePath: "",
	}),
	actions: {
		showError(statusCode: string, message: string) {
			this.statusCode = statusCode

			switch (Number(statusCode)) {
				case HttpStatusCode.ACCESS_DENIED: {
					this.title = `${statusCode} คุณไม่ได้รับอนุญาตเข้าถึงหน้านี้`
					this.messageEn = ""
					this.messageTh = "ขออภัย ! คุณไม่ได้รับอนุญาตเข้าถึงหน้านี้"
					this.imagePath = "/images/errors/403-error.png"
					break
				}
				case HttpStatusCode.NOT_FOUND: {
					this.title = `${statusCode} ไม่พบหน้าที่คุณกำลังค้นหา`
					this.messageEn = ""
					this.messageTh = "ขออภัย ! ไม่พบหน้าที่คุณกำลังค้นหา"
					this.imagePath = "/images/errors/404-error.png"
					break
				}
				case HttpStatusCode.INTERNAL_SERVER_ERROR: {
					this.title = `${statusCode} พบข้อผิดพลาดเกี่ยวกับเซิร์ฟเวอร์`
					this.messageEn = ""
					this.messageTh = "ขออภัย ! พบข้อผิดพลาดเกี่ยวกับเซิร์ฟเวอร์"
					this.imagePath = "/images/errors/500-error.png"
					break
				}
				default: {
					this.title = `${statusCode} พบข้อผิดพลาด`
					this.messageEn = `Status Code : ${statusCode}`
					this.messageTh = `พบข้อผิดพลาด : ${message || ""}`
					this.imagePath = "/images/errors/default-error.png"
					break
				}
			}
		},
	},
	getters: {},
})
