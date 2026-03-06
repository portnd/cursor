// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	ssr: false,
	app: {
		baseURL: process.env.BASE_URL || "/",
		rootId: "nuxt",
		head: {
			htmlAttrs: {
				lang: "en",
			},
		bodyAttrs: {
			style: "background-color: #f5f8fa",
		},
			charset: "utf-8",
			titleTemplate: "%s | ระบบบริหารทรัพย์สินและโครงสร้างพื้นฐานในเขตทาง",
			meta: [
				{
					name: "viewport",
					content: "width=device-width, initial-scale=1, maximum-scale=5",
				},
				{
					name: "description",
					content: "ระบบบริหารงานบำรุงทางพิเศษของการทางพิเศษแห่งประเทศไทย",
				},
				{
					name: "keywords",
					content: "ระบบบริหารงานบำรุงทางพิเศษของการทางพิเศษแห่งประเทศไทย",
				},
				{
					name: "author",
					content: "ระบบบริหารงานบำรุงทางพิเศษของการทางพิเศษแห่งประเทศไทย",
				},
			],
			link: [{ rel: "icon", type: "image/ico", href: "/images/favicon.png" }],
		},
	},
	components: [
		{
			path: "components",
			pathPrefix: false,
		},
	],
	modules: [
		[
			"@pinia/nuxt",
			{
				autoImports: ["defineStore", "definePiniaStore"],
			},
		],
	"@pinia-plugin-persistedstate/nuxt",
	"@vueuse/nuxt",
	"nuxt-icon",
	"@nuxtjs/device",
],
	plugins: [
		"plugins/toast.ts",
		"plugins/vee-validate.ts",
		"plugins/datatable.ts",
		"plugins/highcharts.ts",
		"plugins/apexcharts.client.ts",
		"plugins/datepicker.ts",
		"plugins/draggable.ts",
		"plugins/multiselect.ts",
		"plugins/bootstrap.client.ts",
		"plugins/sweetalert2.ts",
		"plugins/longdo-map-vue.client.ts",
		"plugins/colorpicker.ts",
		"plugins/file-pond.ts",
		"plugins/viewer.ts",
		"plugins/mask.ts",
	],
	runtimeConfig: {
		public: {
			siteUrl: process.env.SITE_URL || "http://localhost:3000",
			baseUrl: process.env.BASE_URL || "/",
			apiUrl: process.env.API_URL || "http://localhost:8080/api/v1/",
			autoRefreshJwtInterval: process.env.AUTO_REFRESH_JWT_INTERVAL || "30", // 30 Minutes
			jwtExpirationTime: process.env.JWT_EXPIRATION_TIME || "20", // 20 Minutes
			longdoMapApiKey: process.env.LONGDO_MAP_API_KEY,
		},
	},
	piniaPersistedstate: {
		cookieOptions: {
			sameSite: "strict",
		},
		storage: "localStorage",
	},
	devServer: {
		host: "0.0.0.0",
		port: 3000,
	},
	build: {
		transpile: ["vue-toastification", "sweetalert2"],
	},
	vite: {
		optimizeDeps: {
			include: ["highcharts"],
		},
	},
})
