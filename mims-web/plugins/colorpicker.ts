import PickColors from "vue-pick-colors"

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component("VueColorPicker", PickColors)
})
