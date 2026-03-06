// @ts-ignore
import LongdoMap from "longdo-map-vue"

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.use(LongdoMap, {
		load: {
			apiKey: useRuntimeConfig().longdoMapApiKey,
		},
	})
})
