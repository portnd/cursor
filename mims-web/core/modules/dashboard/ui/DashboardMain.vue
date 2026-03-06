<script setup lang="ts">
import { useDashboardStore } from "../store"
import { useDashboardConditionStore } from "../condition/store"
import RoadConditionPhotoViewer from "../../road/info/condition/ui/RoadConditionPhotoViewer.vue"
import { DashboardSearch } from "./index"

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})
const store = useDashboardStore()
const conditionStore = useDashboardConditionStore()
useStoreLifecycle(store)

watch(
	() => store.menu,
	() => {
		handleShowImage()
	}
)

// Fallback for IRI when API labels not yet loaded (e.g. before getCondition)
const conditionSymbolFallback = [
	{ name: "เรียบมาก", color: "rgb(66, 210, 53)" },
	{ name: "เรียบ", color: "rgb(164, 252, 165)" },
	{ name: "ขรุขระ", color: "rgb(247, 122, 20)" },
	{ name: "ขรุขระมาก", color: "rgb(255, 41, 10)" },
]

// Legend for condition map: use API labels/colors per condition type (IRI → เรียบมาก/เรียบ/…, MPD → หยาบ (ดีมาก)/ปานกลาง/ละเอียด (แย่มาก), etc.)
const conditionSymbol = computed(() => {
	const labels = conditionStore.conditionLabel ?? []
	const colors = conditionStore.conditionColors ?? []
	if (labels.length === 0 || colors.length === 0) return conditionSymbolFallback
	return labels.map((name: string, i: number) => ({
		name,
		color: colors[i] ? (colors[i].startsWith("#") ? colors[i] : `#${colors[i]}`) : "#cccccc",
	}))
})

const handleShowImage = () => {
	if (store.menu === "condition") {
		if (store.roads.length === 1 && store.conditionType !== 5) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

onUnmounted(() => {
	store.$reset()
})

const isModalMap = ref(false)
const openModal = () => {
	isModalMap.value = !isModalMap.value
}
</script>

<template>
	<!-- {{ store.menu }} -->
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<DashboardSearch :collapsed="mapShow.collapsed" />
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div :class="store.loadingMap ? 'position-relative bg-tran' : 'widget'">
				<div v-show="store.loadingMap" class="loading-icon text-white">
					<div class="spinner-border" role="status"></div>
					<span class="fs-4 text-loading" style="text-shadow: 1px 1px var(--kt-gray-800)">กำลังโหลดข้อมูล</span>
				</div>
				<KeepAlive>
					<VMap
						v-model="mapShow"
						:height="'97vh'"
						:is-sticky="true"
						:opacity-load="store.loadingMap || store.loading ? true : false"
						@map="store.setMap"
					>
						<template v-if="store.menu === 'surface'" #top-end>
							<button type="button" class="btn btn-layer" @click="openModal()">
								<i v-if="!isModalMap" class="fi fi-sr-layers fs-4 p-0 text-gray-600"></i>
								<i v-else class="fi fi-br-cross fs-4 p-0 text-gray-600"></i>
							</button>
						</template>
						<template v-if="store.menu === 'surface'" #top-end-modal>
							<div v-if="isModalMap" class="map-layer rounded-2 shadow">
								<div class="map-title px-4 py-2">การแสดงผลชั้นข้อมูล {{ store.params.display }}</div>
								<div class="map-detail px-4 pb-8 pt-2">
									<VRadio
										v-model="store.params.display"
										:options="[
											{ label: 'อายุผิวทาง', value: 1 },
											{ label: 'สภาพผิวทาง', value: 2 },
										]"
										name="status"
										:inline="false"
										:required="true"
										@update:model-value="store.getSurfaceMap"
									/>
									<VLoading :loading="store.loading" margin-top="24" />
								</div>
							</div>
						</template>
						<template v-if="store.menu === 'surface'" #bottom-start>
							<div class="position-relative">
								<div
									class="detail-modal"
									:style="{
										'margin-top': store.params.display !== 2 ? '-170px' : '-218px',
									}"
								>
									<div
										style="
											left: 0px;
											top: -40px;
											position: absolute;
											background: rgba(255, 255, 255, 0.95);
											box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
											border-radius: 16px;
											border: 1px var(--kt-primary) solid;
										"
										:style="{
											height: store.params.display !== 2 ? '200px' : '250px',
											width: store.params.display !== 2 ? '185px' : '285px',
										}"
									></div>
									<div
										style="height: 0px; left: 5px; top: 15px; position: absolute; border: 0.5px var(--kt-primary) solid"
										:style="{
											width: store.params.display !== 2 ? '174px' : '275px',
										}"
									></div>
									<div
										style="
											height: 16px;
											width: 100px;
											top: 8px;
											position: absolute;
											background: white;
											justify-content: center;
											align-items: center;
											gap: 10px;
											display: inline-flex;
										"
										:style="{ left: store.params.display !== 2 ? '42px' : '90px' }"
									>
										<div class="text-gray-800" style="font-size: 12px; font-weight: 500">
											<span v-if="store.params.display !== 2">อายุผิวทาง</span>
											<span v-else>สภาพผิวทาง</span>
										</div>
									</div>
									<div
										style="
											width: 75px;
											height: 33px;
											padding: 4px 8px;
											top: -39px;
											position: absolute;
											background: var(--kt-primary);
											border-radius: 0px 16px 0px 16px;
											overflow: hidden;
											justify-content: center;
											align-items: center;
											display: inline-flex;
										"
										:style="{
											left: store.params.display !== 2 ? '110px' : '210px',
										}"
									>
										<div class="text-white" style="font-size: 12px; line-height: 12px; word-wrap: break-word">
											สัญลักษณ์
										</div>
									</div>
									<div
										style="
											left: 10px;
											top: 30px;
											position: absolute;
											justify-content: center;
											align-items: center;
											gap: 10px;
											display: inline-flex;
										"
									>
										<div
											v-if="store.params.display === 2"
											class="row ms-0 d-flex align-items-center"
											style="width: 205px"
										>
											<div class="parent">
												<template v-for="(item, index) in store.data?.surface?.summary" :key="index">
													<div class="row d-flex">
														<div
															class="col-2 square mb-3"
															:style="{ 'background-color': item?.surface?.color_code }"
														></div>
														<div class="col-10 mb-2 text-gray-700" style="font-size: 10px">
															{{ item?.surface?.name }}
														</div>
													</div>
												</template>
											</div>
										</div>
										<div v-else>
											<div class="parent">
												<template v-for="(item, index) in store.pavementAge" :key="index">
													<div class="row d-flex">
														<div class="col-2 square mb-3" :style="{ 'background-color': item.color }"></div>
														<div class="col-10 mb-3 text-gray-700">{{ item.name }}</div>
													</div>
												</template>
											</div>
										</div>
									</div>
								</div>
							</div>
						</template>
						<template v-else-if="store.menu === 'maintenance'" #bottom-start>
							<div class="position-relative">
								<div style="margin-top: -190px">
									<div
										style="
											width: 205px;
											height: 150px;
											left: 0px;
											top: 12px;
											position: absolute;
											background: rgba(255, 255, 255, 0.95);
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
											background: #fcfcfb;
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
						<template v-else-if="store.menu === 'condition'" #bottom-start>
							<div class="position-relative">
								<div
									class="detail-modal"
									:style="{
										'margin-top': store.params.display !== 2 ? '-170px' : '-218px',
									}"
								>
									<div
										style="
											left: 0px;
											top: -40px;
											position: absolute;
											background: rgba(255, 255, 255, 0.95);
											box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
											border-radius: 16px;
											border: 1px var(--kt-primary) solid;
										"
										:style="{
											height: store.params.display !== 2 ? '200px' : '250px',
											width: store.params.display !== 2 ? '185px' : '285px',
										}"
									></div>
									<div
										style="height: 0px; left: 5px; top: 15px; position: absolute; border: 0.5px var(--kt-primary) solid"
										:style="{
											width: store.params.display !== 2 ? '174px' : '275px',
										}"
									></div>
									<div
										style="
											height: 16px;
											width: 100px;
											top: 8px;
											position: absolute;
											background: white;
											justify-content: center;
											align-items: center;
											gap: 10px;
											display: inline-flex;
										"
										:style="{ left: store.params.display !== 2 ? '42px' : '90px' }"
									>
										<div class="text-gray-800" style="font-size: 12px; font-weight: 500">
											{{ conditionStore.conditionTypeString || "ข้อมูลสรุปสภาพทาง" }}
										</div>
									</div>
									<div
										style="
											width: 75px;
											height: 33px;
											padding: 4px 8px;
											top: -39px;
											position: absolute;
											background: var(--kt-primary);
											border-radius: 0px 16px 0px 16px;
											overflow: hidden;
											justify-content: center;
											align-items: center;
											display: inline-flex;
										"
										:style="{
											left: store.params.display !== 2 ? '110px' : '210px',
										}"
									>
										<div class="text-white" style="font-size: 12px; line-height: 12px; word-wrap: break-word">
											สัญลักษณ์
										</div>
									</div>
									<div
										style="
											left: 10px;
											top: 30px;
											position: absolute;
											justify-content: center;
											align-items: center;
											gap: 10px;
											display: inline-flex;
										"
									>
										<div>
											<div class="parent">
												<template v-for="(item, index) in conditionSymbol" :key="index">
													<div class="row d-flex">
														<div class="col-2 square mb-3" :style="{ 'background-color': item.color }"></div>
														<div class="col-10 mb-3 text-gray-700">{{ item.name }}</div>
													</div>
												</template>
											</div>
										</div>
									</div>
								</div>
							</div>
						</template>
					</VMap>
				</KeepAlive>
				<RoadConditionPhotoViewer
					v-if="handleShowImage()"
					:map-show="mapShow.collapsed"
					:image="'/images/icons/png/picture-asset.png'"
				/>
			</div>
		</div>
	</div>
</template>

<style scoped>
.btn-layer {
	background-color: white;
	width: 2em;
	height: auto;
	border-radius: 0.5em;
	display: flex;
	align-items: center;
	justify-content: center;
}
.map-layer {
	background-color: white;
}
.map-title {
	border-top-left-radius: 1em;
	border-top-right-radius: 1em;
	background-color: var(--kt-primary);
	width: 12em;
	color: var(--kt-text-white);
}
.map-loading {
	position: absolute;
	z-index: 999;
	width: 100%;
}

.loading-icon {
	position: absolute;
	left: 47%;
	top: 50%;
	z-index: 999999 !important;
}
.text-loading {
	width: 300px;
	position: absolute;
	top: 100%;
	left: -220%;
	z-index: 999999 !important;
}
.bg-tran {
	background: linear-gradient(360deg, rgba(0, 0, 0, 0.64) -22%, rgba(0, 0, 0, 0) 100%);
	border-radius: 16px;
	z-index: 999999 !important;
	position: sticky;
	top: 1.2%;
}
.parent {
	display: flex;
	height: 200px;
	width: 215px;
	flex-flow: column wrap;
}
.parent > div {
	width: 140px;
	margin-right: 5px;
	margin-left: 2px;
	line-height: 20px;
	font-size: 10px;
}
</style>
