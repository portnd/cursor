<script setup lang="ts">
import { IReportParams } from "../../infrastructure"
import {
	useStrategicAnalyzeSummaryStore,
	useStrategicAnalysisDashboardStore,
	useStrategicEditStore,
	DISPLAY_IRI,
	DISPLAY_METHOD,
} from "../../store"
import StrategicSummaryGraph from "./StrategicSummaryGraph.vue"
import StrategicSummaryTab from "./StrategicSummaryTab.vue"

import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const route = useRoute()
const id = Number(route.params.id)
const path = route.path.split("/")[route.path.split("/").length - 1]

const isModalShow = ref(false)
const modal = ref()
const canEdit = ref<boolean>()
// const favorite = ref<boolean>(false)

const store = useStrategicAnalysisDashboardStore()
const editStore = useStrategicEditStore()
const summary = useStrategicAnalyzeSummaryStore()
useStoreLifecycle([store, editStore, summary])

const isModalMap = ref(false)
const device = useDevice()
const windowWidth = ref()

onMounted(async () => {
	await store.getdata(id)
	await store.getMapfilter(id)
	await store.getMapData(id)
})

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

const reportMenu = computed(() => {
	const planCount = store.data?.number_plan

	interface IReportMenu {
		name: string
		id: number
		report_no: string
		plan?: number
	}
	let basedMenu: IReportMenu[] = [
		{ name: "เงื่อนไขการซ่อมบำรุง", id: 1, report_no: "1" },
		{ name: "เงื่อนไขค่าใช้จ่ายการซ่อมบำรุง", id: 2, report_no: "2" },
		{ name: "สรุปค่าซ่อมบำรุงและค่า IRI ของแต่ละปี", id: 3, report_no: "3" },
		{ name: "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง", id: 4, report_no: "4" },
		{ name: "แผนการดำเนินงานการปรับปรุงผิวทาง", id: 5, report_no: "5" },
	]

	basedMenu = basedMenu.map((item) => {
		const newItem = { ...item }

		switch (path) {
			case "no-budget-limit":
				if (item.id === 4) {
					newItem.name = "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง"
				}
				if (item.id === 5) {
					newItem.name = "แผนการดำเนินงานการปรับปรุงผิวทาง "
				}
				break
			case "budget-limit":
				if (item.id === 4) {
					newItem.name = "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง (ไม่จำกัดงบประมาณ)"
				}
				if (item.id === 5) {
					newItem.name = "แผนการดำเนินงานการปรับปรุงผิวทาง (ไม่จำกัดงบประมาณ)"
				}
				break
		}

		return newItem
	})

	if (path === "iri-target") {
		basedMenu = basedMenu.filter((item) => item.report_no !== "4" && item.report_no !== "5")
	}

	if (!planCount) {
		return basedMenu
	}

	let idCounter = basedMenu[basedMenu.length - 1].id

	if (path !== "iri-target") {
		basedMenu[basedMenu.length - 1].id += planCount

		const lastIndex = basedMenu[basedMenu.length - 1]

		for (let i = 1; i <= planCount; i++) {
			basedMenu.push(
				{
					name: `รายละเอียดแผนงานซ่อมบำรุงตามสายทาง (แผนที่ ${i})`,
					id: idCounter++,
					report_no: "4",
					plan: i,
				},
				{
					name: `แผนการดำเนินงานการปรับปรุงผิวทาง (แผนที่ ${i})`,
					id: lastIndex.id + i,
					report_no: "5",
					plan: i,
				}
			)
		}
	} else {
		for (let i = 1; i <= planCount; i++) {
			basedMenu.push({
				name: `รายละเอียดแผนงานซ่อมบำรุงตามสายทาง (แผนที่ ${i})`,
				id: ++idCounter,
				report_no: "4",
				plan: i,
			})
		}

		for (let i = 1; i <= planCount; i++) {
			basedMenu.push({
				name: `แผนการดำเนินงานการปรับปรุงผิวทาง (แผนที่ ${i})`,
				id: ++idCounter,
				report_no: "5",
				plan: i,
			})
		}
	}

	return basedMenu.sort((a, b) => a.id - b.id)
})

// const bookmark = async (id: number) => {
// 	const res = await store.setFavorite(id)

// 	if (res?.status) {
// 		favorite.value = res.data.is_favorite
// 	}
// }

onMounted(() => {
	canEdit.value = initUserStore.accessPermissions[IUserRolesAccess.manage_myself_maintenance_analysis]
	window.addEventListener("resize", updateWidth)
})

const onCopy = async () => {
	const res = await editStore.copy(id)

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				editStore.isCopy = true
				navigateTo(`/analyses/strategic/select-road/${res.data.id}`)
			},
		})
	}
}

const downloadFile = (fileType: string) => {
	const params = {} as IReportParams
	params.type = fileType

	if (store.exportReport.reportPlan) {
		params.plan = store.exportReport.reportPlan
	} else {
		delete params.plan
	}

	const keyword = params ? store.encodeQuery(params) : ""

	useDownloadFile(
		`ดาวน์โหลด: ${store.exportReport.reportName}`,
		`analyze/${id}/report/report${store.exportReport.reportTypeNo}${keyword !== "" ? "?" + keyword : ""}`,
		fileType as "html" | "pdf" | "excel"
	)
}

watch(
	() => store.exportReport.reportId,
	() => {
		if (!store.exportReport.reportId) {
			isModalShow.value = false
		} else {
			isModalShow.value = true
		}
	}
)

// const onDelete = () => {
// 	useDeleteItem({
// 		name: "บำรุงรักษาเชิงกลยุทธ์",
// 		url: `analyze/${id}`,
// 		callBack: () => {
// 			navigateTo("/analyses")
// 			store.$reset()
// 		},
// 	})
// }

onUnmounted(() => {
	store.$reset()
	window.removeEventListener("resize", updateWidth)
})

const openModal = () => {
	isModalMap.value = !isModalMap.value
}

const conditionRule = ref([
	{
		name: "0 <= เรียบมาก < 2.5",
		color: "#50CD89",
	},
	{
		name: "2.5 <= เรียบ < 3.5",
		color: "#87C442",
	},
	{
		name: "2.5 <= ขรุขระ < 4.5",
		color: "#FDB833",
	},
	{
		name: "4.5 <= ขรุขระมาก < 15",
		color: "#DC3545",
	},
])
// const maintenance = [
// 	{
// 		name: "OL-Overlay",
// 		color: "#FF7A8D",
// 	},
// 	{
// 		name: "M&OL-Mill&Overlay",
// 		color: "#FFB800",
// 	},
// 	{
// 		name: "RCL-Recycling",
// 		color: "#AF85FF",
// 	},
// 	{
// 		name: "Rc-Reconstruction",
// 		color: "#B22727",
// 	},
// 	{
// 		name: "SS-SlurrySeal",
// 		color: "#82E0AA",
// 	},
// 	{
// 		name: "FDR",
// 		color: "#418FFF",
// 	},
// 	{
// 		name: "BCO",
// 		color: "#FF69B4",
// 	},
// 	{
// 		name: "M&OL",
// 		color: "#7FFFD4",
// 	},
// 	{
// 		name: "Seal",
// 		color: "#FF7F33",
// 	},
// ]

watch(
	() => summary.condition,
	() => {
		if (summary.condition === 1) {
			conditionRule.value = [
				{
					name: "0 <= เรียบมาก < 2.5",
					color: "#50CD89",
				},
				{
					name: "2.5 <= เรียบ < 3.5",
					color: "#87C442",
				},
				{
					name: "2.5 <= ขรุขระ < 4.5",
					color: "#FDB833",
				},
				{
					name: "4.5 <= ขรุขระมาก < 15",
					color: "#DC3545",
				},
			]
		} else {
			conditionRule.value = [
				{
					name: "0 <= ผ่าน < 2.5",
					color: "#50CD89",
				},
				{
					name: "2.5 <= ไม่ผ่าน < 3.5",
					color: "#DC3545",
				},
			]
		}
	}
)

watch(
	() => summary.fullScreen,
	() => {
		const element = document.querySelector("#content") as HTMLElement
		if (summary.fullScreen) {
			element.style.display = "none"
		} else {
			element.style.display = "block"
		}
	}
)

const updateWidth = () => {
	windowWidth.value = window.innerWidth
	const element = document.querySelector("#content") as HTMLElement | null
	const topEnd = document.querySelector(".top-end") as HTMLElement | null
	const isMobile = ref()
	if (summary.fullScreen === true) {
		if (device.isMobileOrTablet || windowWidth.value < 992) {
			if (element) element.style.display = "block"
			if (topEnd) topEnd.style.top = "20px"
			summary.map?.Ui?.LayerSelector.visible(false)
			summary.map?.Ui?.Fullscreen.visible(false)
			isMobile.value = true
		} else {
			if (element) element.style.display = "none"
			if (topEnd) topEnd.style.top = "40px"
			summary.map?.Ui?.LayerSelector.visible(true)
			summary.map?.Ui?.Fullscreen.visible(true)
			isMobile.value = false
		}
	} else if (summary.fullScreen === false) {
		if (device.isMobileOrTablet || windowWidth.value < 992) {
			summary.map?.Ui?.LayerSelector.visible(false)
			summary.map?.Ui?.Fullscreen.visible(false)
		} else {
			summary.map?.Ui?.LayerSelector.visible(true)
			summary.map?.Ui?.Fullscreen.visible(true)
		}
		const newSrc = "https://api.longdo.com/map/images/ui/fullscreen-up.png"
		const buttonFullScreen = document.querySelector(".ldmap_fullscreen") as HTMLImageElement | null
		if (buttonFullScreen) {
			buttonFullScreen.src = newSrc
		}
		if (topEnd) topEnd.style.top = "40px"
	}
	if (isMobile.value) {
		const elementFullScreen = document.querySelector(".ldmap_placeholder_fullscreen")
		if (elementFullScreen) {
			elementFullScreen.classList.remove("ldmap_placeholder_fullscreen")
		}
		summary.fullScreen = false
	}
}

const onDelete = () => {
	useDeleteItem({
		name: "บำรุงรักษาเชิงกลยุทธ์",
		url: `analyze/${id}`,
		callBack: () => {
			navigateTo("/analyses")
			store.$reset()
		},
	})
}
</script>

<template>
	<div class="row">
		<div id="content" class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<div class="card p-5 mb-5 position-relative">
				<!-- Loading overlay -->
				<div v-if="store.loading" class="loading-overlay">
					<div class="loading-overlay__inner">
						<span class="spinner-border spinner-border text-primary" role="status"></span>
						<span class="loading-overlay__text ms-3 fs-6 text-gray-600">โปรดรอสักครู่...</span>
					</div>
				</div>
				<div class="row mb-md-5">
					<div class="col-lg-5 col-md-4 col-12 align-self-center">
						<!-- <button
							type="button"
							class="btn btn-primary rounded-4 mb-3 mb-md-0 px-5 py-2 fw-semibold"
							@click="bookmark(id)"
						>
							<i v-if="favorite === true" class="fi fi-ss-star align-middle fs-5"></i>
							<i v-else class="fi fi-rs-star align-middle fs-5"></i>
							เพิ่มไปยังรายการโปรด
						</button> -->
					</div>
					<div class="col-lg-7 col-md-8 col-12 d-flex align-items-center justify-content-end gap-4 text-md-end">
						<span class="dropdown">
							<button
								id="dropdownMenuButton1"
								class="btn btn-primary rounded-4 mb-3 mb-md-0 py-3 fw-semibold dropdown-toggle"
								type="button"
								data-bs-toggle="dropdown"
								aria-expanded="false"
								data-bs-auto-close="outside"
							>
								รายงาน <i class="fi fi-rr-angle-small-down fs-5 mt-1 ms-1 pe-0"></i>
							</button>
							<ul
								ref="modal"
								class="dropdown-menu lh-2 mb-0 pt-0"
								aria-labelledby="dropdownMenuButton1"
								style="overflow: hidden"
							>
								<li v-for="(item, index) of reportMenu" :key="index">
									<NuxtLink
										class="dropdown-report fs-6 py-2 cursor-pointer"
										:class="store.exportReport.reportId === item.id ? 'active' : ''"
										@click="store.handleReport(item.id, item.report_no, item.name, item.plan ? item.plan : undefined)"
										>{{ item.name }}</NuxtLink
									>
								</li>
								<li v-if="isModalShow" class="d-flex my-1">
									<NuxtLink
										class="dropdown-item text-center rounded-1 py-3 mx-4 cursor-pointer"
										@click="downloadFile('html')"
									>
										<img src="/images/files/html.png" width="36" alt="logo" />
									</NuxtLink>
									<NuxtLink
										class="dropdown-item text-center rounded-1 py-3 mx-4 cursor-pointer"
										@click="downloadFile('pdf')"
									>
										<img src="/images/files/pdf.png" width="36" alt="logo" />
									</NuxtLink>
									<NuxtLink
										class="dropdown-item text-center rounded-1 py-3 mx-4 cursor-pointer"
										@click="downloadFile('excel')"
									>
										<img src="/images/files/xlsx.png" width="36" alt="logo" />
									</NuxtLink>
								</li>
							</ul>
						</span>
						<div
							v-show="canEdit"
							id="list-menu-1"
							class="cursor-pointer mt-md-3 fs-3"
							data-bs-toggle="dropdown"
							aria-expanded="false"
							data-bs-auto-close="outside"
						>
							<i class="fi fi-br-menu-dots-vertical"></i>
						</div>
						<ul class="dropdown-menu lh-2 mb-0 py-1 p-4 list-menu" aria-labelledby="list-menu-1">
							<li
								v-show="canEdit"
								class="px-4 py-2"
								@click="navigateTo(`/analyses/strategic/summary/${id}/${path}/detail`)"
							>
								รายละเอียด
							</li>
							<li v-show="canEdit" class="px-4 py-2" @click="navigateTo(`/analyses/strategic/select-road/${id}/edit`)">
								แก้ไข
							</li>
							<li v-show="canEdit" class="px-4 py-2" @click="onCopy">คัดลอก</li>
							<li v-show="canEdit" class="px-4 py-2 text-danger" @click="onDelete">
								<i class="fi fi-sr-trash align-middle fs-6" />
								ลบ
							</li>
						</ul>
						<!-- <NuxtLink
							:to="`/analyses/strategic/summary/${id}/budget-limit/detail`"
							class="btn btn-outline btn-outline-primary rounded-4 ms-3 mb-3 mb-md-0 px-8 py-3 fw-semibold"
						>
							รายละเอียด
						</NuxtLink>
						<NuxtLink
							v-show="canEdit"
							:to="`/analyses/strategic/select-road/${id}/edit`"
							class="btn btn-outline btn-outline-primary rounded-4 ms-3 mb-3 mb-md-0 px-8 py-3 fw-semibold"
						>
							แก้ไข
						</NuxtLink>
						<NuxtLink
							v-show="canEdit"
							class="btn btn-outline btn-outline-primary rounded-4 ms-3 mb-3 mb-md-0 px-8 py-3 fw-semibold"
							@click="onCopy"
						>
							คัดลอก
						</NuxtLink> -->
						<!-- <button
							v-show="canEdit"
							type="button"
							class="btn btn-outline btn-outline-danger rounded-4 ms-md-3 mb-3 mb-md-0 px-5 py-2 fw-semibold"
							@click="onDelete"
						>
							<i class="fi fi-sr-trash align-middle fs-6"></i>
							ลบ
						</button> -->
					</div>
				</div>
				<div class="row mt-5">
					<!-- กราฟ -->
					<StrategicSummaryGraph />
					<!-- กราฟ -->
				</div>
				<div class="row mb-3">
					<div class="col-12">
						<StrategicSummaryTab />
					</div>
				</div>
			</div>
		</div>
		<div id="map" class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap v-model="mapShow" :loading="store.loading" height="97vh" :is-sticky="true" @map="store.setMap">
						<template #top-end>
							<button type="button" class="btn btn-layer" @click="openModal()">
								<i v-if="!isModalMap" class="fi fi-sr-layers fs-4 p-0 text-gray-600"></i>
								<i v-else class="fi fi-br-cross fs-4 p-0 text-gray-600"></i>
							</button>
						</template>
						<template #top-end-modal>
							<div v-if="isModalMap" class="map-layer rounded-2 shadow">
								<div class="map-title px-4 py-2">การแสดงผลชั้นข้อมูล</div>
								<div class="map-detail px-4 pb-8 pt-2">
									<VSelect
										v-model="store.map_params.plan"
										:options="store.getMapPlanOptions"
										label="แผน"
										name="plan"
										placeholder="เลือก"
										:can-clear="false"
										:can-deselect="false"
										@update:model-value="store.onMapSelected(id)"
									/>
									<VSelect
										v-model="store.map_params.year"
										:options="store.getMapYearOptions"
										label="ปีที่"
										name="year"
										placeholder="เลือก"
										:can-clear="false"
										:can-deselect="false"
										@update:model-value="store.onMapSelected(id)"
									/>
									<VSelect
										v-model="store.map_params.display"
										:options="store.getMapDisplayOptions"
										label="การแสดงผลลัพธ์"
										name="display"
										placeholder="เลือก"
										:can-clear="false"
										:can-deselect="false"
										@update:model-value="store.onMapSelected(id)"
									/>
									<VSelect
										v-if="store.map_params.display === DISPLAY_IRI"
										v-model="store.map_params.cirteria"
										:options="store.getCriteriaOptions"
										label="เกณฑ์จำแนกสภาพทาง"
										name="cirteria"
										placeholder="เลือก"
										:can-clear="false"
										:can-deselect="false"
										@update:model-value="store.onMapSelected(id)"
									/>
								</div>
							</div>
						</template>
						<template #bottom-start>
							<div class="position-relative">
								<div
									:style="{
										'margin-top': `-${store.getSymbolBoxPositionTop}px`,
									}"
								>
									<div
										style="
											left: 0px;
											top: -40px;
											position: absolute;
											background: rgba(255, 255, 255);
											box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
											border-radius: 16px;
											border: 1px #fdb833 solid;
										"
										:style="{
											width: store.map_params.display !== DISPLAY_METHOD ? '185px' : `${store.getSymbolBoxWidth}px`,
											height: `${store.getSymbolBoxHeight}px`,
										}"
									></div>
									<div
										style="height: 0px; left: 0; top: 15px; position: absolute; border: 0.5px #fdb833 solid"
										:style="{
											width:
												store.map_params.display === DISPLAY_IRI
													? '184px'
													: store.getCriteriaLegend.length > 5
													? '280px'
													: '230px',
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
										:style="{
											left:
												store.map_params.display === DISPLAY_IRI
													? '13px'
													: store.getCriteriaLegend.length > 5
													? '100px'
													: '68px',
											width: store.map_params.display === DISPLAY_IRI ? '160px' : '80px',
										}"
									>
										<div class="text-gray-800" style="font-size: 12px; font-weight: 500">
											<span v-if="store.map_params.display === DISPLAY_IRI">เกณฑ์ IRI ผิวทางลาดยาง</span>
											<span v-else>วิธีการซ่อม</span>
										</div>
									</div>
									<div
										style="
											width: 75px;
											height: 33px;
											padding: 4px 8px;
											top: -39px;
											position: absolute;
											background: #fdb833;
											border-radius: 0px 16px 0px 16px;
											overflow: hidden;
											justify-content: center;
											align-items: center;
											display: inline-flex;
										"
										:style="{
											left: store.map_params.display === DISPLAY_IRI ? '110px' : `${store.getSymbolPositionLeft}px`,
										}"
									>
										<div class="text-gray-800" style="font-size: 12px; line-height: 12px; word-wrap: break-word">
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
											v-if="store.map_params.display === DISPLAY_IRI"
											class="row ms-0 d-flex align-items-center"
											style="width: 205px"
										>
											<template v-for="(item, index) in store.getGradeAC" :key="index">
												<div class="col-2 square mb-3" :style="{ 'background-color': item.color }"></div>
												<div class="col-10 mb-2 text-gray-700" style="font-size: 10px">{{ item.name }}</div>
											</template>
										</div>
										<div v-else>
											<div class="parent">
												<template v-for="(item, index) in store.getCriteriaLegend" :key="index">
													<div class="row d-flex">
														<div class="col-2 square mb-3" :style="{ 'background-color': item.color }"></div>
														<div class="col-10 mb-3 text-gray-700">{{ item.name }}</div>
													</div>
												</template>
											</div>
										</div>
									</div>

									<!-- iri grad_cc -->
									<div
										v-if="store.map_params.display === DISPLAY_IRI"
										style="position: absolute"
										:style="{
											top: `${store.getGradeCCPositionTop}px`,
										}"
									>
										<div
											style="height: 0px; width: 184px; border: 0.5px #fdb833 solid"
											:style="{
												width: '184px',
											}"
										></div>
										<div
											style="
												height: 16px;
												width: 160px;
												top: -7px;
												position: absolute;
												background: white;
												justify-content: center;
												align-items: center;
												gap: 10px;
												display: inline-flex;
											"
											:style="{
												left: '13px',
											}"
										>
											<div class="text-gray-800 text-center" style="font-size: 12px; font-weight: 500">
												<span> เกณฑ์ iri ผิวทางคอนกรีต </span>
											</div>
										</div>

										<div
											style="
												left: 10px;
												top: 25px;
												position: absolute;
												justify-content: center;
												align-items: center;
												gap: 10px;
												display: inline-flex;
											"
										>
											<div class="row ms-0 d-flex align-items-center" style="width: 205px">
												<template v-for="(item, index) in store.getGradeCC" :key="index">
													<div class="row">
														<div class="col-2 square mb-3" :style="{ 'background-color': item.color }"></div>
														<div class="col-10 mb-2 text-gray-700 criteria-ellipsis-one" style="font-size: 10px">
															{{ item.name }}
														</div>
													</div>
												</template>
											</div>
										</div>
									</div>
									<!--  -->
								</div>
							</div>
						</template>
					</VMap>
				</KeepAlive>
			</div>
		</div>
	</div>
	<!-- </VSkeletonLoader> -->
</template>

<style scoped lang="scss">
.lh-2 {
	line-height: 2;
}

.btn.btn-primary.dropdown-toggle:after {
	display: none;
}
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
	width: 18em;
}
.square {
	width: 2em;
	height: 2em;
}
.parent {
	display: flex;
	height: 150px;
	width: 215px;
	flex-flow: column wrap;
}
.parent > div {
	width: 130px;
	margin-right: 5px;
	margin-left: 2px;
	line-height: 20px;
	font-size: 10px;
	text-wrap: nowrap;
}
.dropdown-report {
	display: block;
	width: 100%;
	font-weight: 400;
	white-space: nowrap;
	border: 0;
	border-radius: 0;
	padding: 0.25rem 1rem;
	color: var(--kt-gray-800);
	&:hover {
		background-color: #f9f9f9;
	}
}
.dropdown-report.active {
	background-color: #fff0d9;
}

.list-menu {
	padding: 1rem 0 !important;

	li {
		cursor: pointer !important;
	}

	li:hover {
		background-color: #f9f9f9;
	}
}

.loading-overlay {
	position: absolute;
	inset: 0;
	background: rgba(255, 255, 255, 0.75);
	z-index: 10;
	border-radius: inherit;
	display: flex;
	align-items: center;
	justify-content: center;
}

.loading-overlay__inner {
	display: flex;
	align-items: center;
}

.loading-overlay__text {
	font-weight: 500;
}
</style>
