<script setup lang="ts">
import { useMaintenanceHistorySearchTableStore, useMaintenanceHistoryDetailsStore } from "../store"
import { MaintenanceHistoryTable } from "./index"

const store = useMaintenanceHistorySearchTableStore()
const device = useDevice()
const windowWidth = ref()

// Event delegation: คลิกป๊อปอัปแผนที่ (หน้ารายการ) → นำทางไปหน้า info
const handleMapClick = (e: MouseEvent) => {
	const link = (e.target as HTMLElement).closest("[data-maintenance-info-id]") as HTMLElement | null
	if (link) {
		e.preventDefault()
		e.stopPropagation()
		const id = link.getAttribute("data-maintenance-info-id")
		if (id) {
			navigateTo(`/maintenances/history/${id}/info`)
		}
	}
}

const mapShow = ref({
	collapsed: false,
})

onMounted(() => {
	// Clear details store when entering list so next info visit doesn’t show stale data
	useMaintenanceHistoryDetailsStore().$reset()
	addEventListener("resize", updateWidth)
})

const updateWidth = () => {
	windowWidth.value = window.innerWidth
	const element = document.querySelector("#content") as HTMLElement | null
	const topEnd = document.querySelector(".bottom-start") as HTMLElement | null
	const map = store.map
	if (!map || !element || !topEnd) return
	const isMobile = ref()
	if (store.fullScreen === true) {
		if (device.isMobileOrTablet || windowWidth.value < 992) {
			element.style.display = "block"
			topEnd.style.top = "20px"
			map.Ui.LayerSelector.visible(false)
			map.Ui.Fullscreen.visible(false)
			isMobile.value = true
		} else {
			element.style.display = "none"
			topEnd.style.top = "40px"
			map.Ui.LayerSelector.visible(true)
			map.Ui.Fullscreen.visible(true)
			isMobile.value = false
		}
	} else if (store.fullScreen === false) {
		if (device.isMobileOrTablet || windowWidth.value < 992) {
			map.Ui.LayerSelector.visible(false)
			map.Ui.Fullscreen.visible(false)
		} else {
			map.Ui.LayerSelector.visible(true)
			map.Ui.Fullscreen.visible(true)
		}
		const newSrc = "https://api.longdo.com/map/images/ui/fullscreen-up.png"
		const buttonFullScreen = document.querySelector(".ldmap_fullscreen") as HTMLImageElement
		if (buttonFullScreen) {
			buttonFullScreen.src = newSrc
		}
		// topEnd.style.top = "40px"
	}
	if (isMobile.value) {
		const elementFullScreen = document.querySelector(".ldmap_placeholder_fullscreen")
		if (elementFullScreen) {
			elementFullScreen.classList.remove("ldmap_placeholder_fullscreen")
		}
		store.fullScreen = false
	}
}

watch(
	() => store.fullScreen,
	() => {
		const element = document.querySelector("#content") as HTMLElement
		if (store.fullScreen) {
			element.style.display = "none"
		} else {
			element.style.display = "block"
		}
	}
)

onUnmounted(() => {
	removeEventListener("resize", updateWidth)
})
</script>

<template>
	<div class="row">
		<div id="content" class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<MaintenanceHistoryTable />
		</div>
		<div id="map" class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'" @click="handleMapClick">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" height="97vh" @map="store.setMap">
						<template #bottom-start>
							<div class="position-relative">
								<div style="margin-top: -190px">
									<div
										style="
											width: 205px;
											height: 150px;
											left: 0px;
											top: 12px;
											position: absolute;
											background: #fff;
											box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
											border-radius: 16px;
											border: 1px #fdb833 solid;
										"
									></div>
									<div
										style="
											width: 184px;
											height: 0px;
											left: 10px;
											top: 66px;
											position: absolute;
											border: 0.5px #fdb833 solid;
										"
									></div>
									<div
										style="
											width: 120px;
											height: 28px;
											left: 42px;
											top: 50px;
											position: absolute;
											background: #fff;
											justify-content: center;
											align-items: center;
											gap: 10px;
											display: inline-flex;
										"
									>
										<div
											style="
												width: 300px;
												height: 27px;
												text-align: center;
												color: #727272;
												font-size: 12px;
												font-weight: 500;
												word-wrap: break-word;
											"
										>
											ระยะเวลา<br />โครงการติดค้ำประกัน
										</div>
									</div>

									<div
										style="
											width: 75px;
											height: 33px;
											padding: 4px 8px;
											left: 130px;
											top: 12px;
											position: absolute;
											background: #fdb833;
											border-radius: 0px 16px 0px 16px;
											overflow: hidden;
											justify-content: center;
											align-items: center;
											display: inline-flex;
										"
									>
										<div
											style="color: white; font-size: 12px; font-weight: 500; line-height: 12px; word-wrap: break-word"
										>
											สัญลักษณ์
										</div>
									</div>
								</div>

								<div style="width: 100px; display: inline-flex; top: 110px; position: absolute; left: 12px">
									<div
										style="
											width: 15px;
											height: 29px;
											transform: rotate(-90deg);
											transform-origin: 0 0;
											background: #1f70f3;
											border: 2px #1f70f3 solid;
											position: absolute;
											top: 5px;
										"
									></div>
									<div
										style="
											color: rgb(114, 114, 114);
											font-size: 11px;
											font-weight: 300;
											margin-left: 23px;
											position: absolute;
											top: -12px;
											width: 150px;
											left: 12px;
										"
									>
										ระยะเวลาติดค้ำประกัน &gt; 20%
									</div>

									<div
										style="
											width: 15px;
											height: 29px;
											transform: rotate(-90deg);
											transform-origin: 0 0;
											background: #f1416c;
											border: 2px #f1416c solid;
											position: absolute;
											top: 30px;
										"
									></div>
									<div
										style="
											color: rgb(114, 114, 114);
											font-size: 11px;
											font-weight: 300;
											margin-left: 23px;
											position: absolute;
											top: 14px;
											width: 150px;
											left: 12px;
										"
									>
										ระยะเวลาติดค้ำประกัน &lt;= 20%
									</div>
								</div>
							</div>
						</template>
					</VMap>
				</KeepAlive>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
<!-- <img class="position-absolute bottom-0" src="/images/dashboards/symbol.png" width="180" alt="" />
								<img
									class="position-absolute"
									style="bottom: 30px; left: 12px"
									src="/images/dashboards/period.png"
									width="160"
									alt=""
								/> -->
