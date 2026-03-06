import { TIcon } from "~/utils/alert"
import { IError } from "~~/core/shared/http/Error"
import { ICallback } from "~~/core/shared/types/Callback"
import httpStatusCode from "~~/core/shared/http/HttpStatusCode"
import { useActionStore } from "~~/core/shared/store/ActionStore"

const useHandlerError = (
	statusCode: number,
	error: IError,
	callback: ICallback = { showToast: false, showAlert: false }
) => {
	// กรณี ไม่พบข้อมูล
	if (statusCode === httpStatusCode.NOT_FOUND) {
		return showError({ statusCode: httpStatusCode.NOT_FOUND })
	}

	// action จาก vee-validate
	const actionStore = useActionStore()

	// กรณีเซสชั่นหมดอายุ หรือเชื่อมต่อ Api แล้วเซิฟเวอร์ตอบกลับมาเป็น 401
	if (statusCode === httpStatusCode.UNAUTHORIZED) {
		return navigateTo("/auth/expired")
	}

	// กรณีติด Validate
	if (statusCode === httpStatusCode.VALIDATE) {
		const actions = actionStore.actions

		const setErrors = () => {
			return new Promise((resolve) => {
				resolve(actions.setErrors(error.field))
			})
		}

		setErrors().then(() => {
			// ให้ scrollbar เลื่อนไปหา input ที่มีการ validate อันแรกของหน้าจอ
			scrollIntoInvalidField()
		})

		return
	}

	// กรณีไม่มีสิทธิ์เข้าถึงหน้าเว็บไซต์ 403
	if (statusCode === httpStatusCode.ACCESS_DENIED) {
		return showError({ statusCode: httpStatusCode.ACCESS_DENIED })
	}

	// กรณี Error อื่น ๆ
	const options = reactive({
		title: "",
		message: "",
		type: "error" as TIcon,
	})

	switch (statusCode) {
		case httpStatusCode.INTERNAL_SERVER_ERROR: {
			options.title = "เกิดข้อผิดพลาด"
			options.message = "พบบางอย่างผิดพลาด"
			options.type = "error"
			break
		}

		case httpStatusCode.LOCKED: {
			options.title = "แจ้งเตือน"
			options.message = "ใช้งานนี้ถูกปิดการใช้งาน กรุณาติดต่อผู้ดูแลระบบ"
			options.type = "error"
			break
		}

		default: {
			options.title = "แจ้งเตือน"
			options.message = error?.message
			options.type = "error"
			break
		}
	}

	// กรณีแสดง Swal
	if (callback.showAlert) {
		showAlert({
			title: options.title,
			message: options.message,
			type: options.type,
			callBack: callback.fn,
		})
	}

	// กรณีแสดง Toast
	if (callback.showToast) {
		showToast({
			title: options.title,
			message: options.message,
			type: options.type,
			callBack: () => {
				if (callback.fn) {
					return callback.fn
				} else {
					return navigateTo(callback.to)
				}
			},
		})
	}

	/* กรณีไม่แสดง Swal, Toast */
	// กรณีมี function
	if (!callback.showAlert && !callback.showToast && callback.fn) {
		return callback.fn()
	}

	// กรณีให้ Redirect
	if (!callback.showAlert && !callback.showToast && callback.to) {
		return navigateTo(callback.to)
	}
}

export default useHandlerError
