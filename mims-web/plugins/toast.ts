import Toast from "vue-toastification"
import "vue-toastification/dist/index.css"

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.use(Toast, {
		transition: "Vue-Toastification__fade",
		maxToasts: 1,
		newestOnTop: false,
		filterBeforeCreate: (toast: any, toasts: any) => {
			if (toasts.filter((t: any) => t.type === toast.type).length !== 0) {
				return false
			}
			return toast
		},
	})
})
