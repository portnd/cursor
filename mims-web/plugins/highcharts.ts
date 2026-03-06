export default defineNuxtPlugin(async () => {
	console.log("[HighchartsPlugin] starting, process.client:", process.client)
	if (!process.client) return

	const mod: any = await import("highcharts")
	console.log("[HighchartsPlugin] mod loaded, mod type:", typeof mod, "mod.chart:", typeof mod?.chart, "mod.default:", typeof mod?.default, "mod.default?.chart:", typeof mod?.default?.chart)

	// esbuild pre-bundles highcharts in Node.js context where window is absent,
	// causing module.exports = K (the UMD factory fn, not the Highcharts namespace).
	// We must detect and resolve the real namespace before using it.
	let HC: any
	if (typeof mod?.chart === "function") {
		// Named export available directly (Vite CJS-interop best case)
		HC = mod
	} else if (typeof mod?.default?.chart === "function") {
		// Default export IS the namespace
		HC = mod.default
	} else if (typeof mod?.default === "function" && typeof window !== "undefined") {
		// mod.default is the UMD factory K — call it with window to get the namespace
		HC = mod.default(window)
	} else {
		HC = mod?.default ?? mod
	}

	console.log("[HighchartsPlugin] HC resolved, type:", typeof HC, "HC.chart:", typeof HC?.chart, "HC.setOptions:", typeof HC?.setOptions)

	if (typeof HC?.setOptions === "function") {
		HC.setOptions({
			lang: {
				resetZoom: "รีเซ็ต",
				resetZoomTitle: "รีเซ็ต",
			},
		})

		const BoostModule: any = await import("highcharts/modules/boost")
		const initBoost = BoostModule?.default ?? BoostModule
		if (typeof initBoost === "function") {
			initBoost(HC)
		}
	}

	console.log("[HighchartsPlugin] providing $highcharts, HC:", HC, "HC.chart:", typeof HC?.chart)

	// Provide the resolved Highcharts instance so components skip repeated imports
	return {
		provide: {
			highcharts: HC,
		},
	}
})
