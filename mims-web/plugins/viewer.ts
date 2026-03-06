import VueViewer from "v-viewer"
import "viewerjs/dist/viewer.css"

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.use(VueViewer, {
		defaultOptions: {
			zIndex: 1055,
			inline: false,
			button: false,
			navbar: false,
			title: true,
			toolbar: false,
			tooltip: false,
			movable: false,
			zoomable: true,
			rotatable: true,
			scalable: true,
			transition: true,
			fullscreen: false,
			keyboard: true,
		},
	})
})
