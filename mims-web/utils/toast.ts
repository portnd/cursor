import { useToast, TYPE } from "vue-toastification"
import VToast from "~~/components/extras/VToast.vue"

interface IToast {
	title?: string
	message?: string
	type?: string
	callBack?(): void
}

export const showToast = (options: IToast): void => {
	const toast = useToast()

	const content = {
		component: VToast,
		props: {
			title: options.title,
			message: options.message,
		},
	}

	const getType = () => {
		switch (options.type) {
			case TYPE.ERROR: {
				return TYPE.ERROR
			}
			case TYPE.INFO: {
				return TYPE.INFO
			}
			case TYPE.SUCCESS: {
				return TYPE.SUCCESS
			}
			case TYPE.WARNING: {
				return TYPE.WARNING
			}
			default: {
				return TYPE.DEFAULT
			}
		}
	}

	toast(content, {
		type: getType(),
		pauseOnHover: false,
		timeout: 3000,
		onClose() {
			if (options.callBack !== undefined) {
				options.callBack()
			}
		},
	})
}
