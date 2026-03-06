// @ts-ignore
import Treeselect from "@komgrip/vue3-treeselect"
import "@komgrip/vue3-treeselect/dist/vue3-treeselect.css"

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component("VueTreeselect", Treeselect)
})
