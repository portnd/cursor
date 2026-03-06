import { Modal, Popover } from "bootstrap"

export default defineNuxtPlugin(() => ({
	provide: {
		bootstrap: {
			Modal,
			Popover,
		},
	},
}))
