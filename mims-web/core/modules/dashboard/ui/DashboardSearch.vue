<script setup lang="ts">
import { defineRule, useForm } from "vee-validate"
import { IDashboardMaintenanceTableRequest } from "../infrastructure"
import { useDashboardStore } from "../store"
import DashboardTab from "./DashboardTab.vue"
import { IValidate } from "~/core/shared/types/Validate"

defineProps({
	collapsed: {
		type: Boolean,
		default: false,
	},
	searchCall: {
		type: Boolean,
		default: false,
	},
})

defineRule("kmOverlapStart", (value: any) => {
	const start = convertStringToKm(value)
	const end = convertStringToKm(store.params.km_end)

	if (start > end) {
		return "กม. เริ่มต้น ต้องมีค่าน้อยกว่า กม. สิ้นสุด"
	}
	return true
})

defineRule("kmOverlapEnd", (value: any) => {
	const end = convertStringToKm(value)
	const start = convertStringToKm(store.params.km_start)

	if (end < start) {
		return "กม. สิ้นสุด ต้องมีค่ามากกว่า กม. เริ่มต้น"
	}
	return true
})

const localStorageUser = window.localStorage.getItem("init-user")
const initUserStore = localStorageUser ? JSON.parse(localStorageUser) : { accessPermissions: {} }
const store = useDashboardStore()
useStoreLifecycle(store)
const dataTable = ref()

interface DataTableFunction {
	getData: Function
	loadData: Function
	searchData: Function
	resetSelected: Function
	resetRow: Function
}

const validates = computed(() => {
	const validate: IValidate = {}
	if (store.params.km_start) {
		validate.km_start = "km|kmOverlapStart"
	}

	if (store.params.km_end) {
		validate.km_end = "km|kmOverlapEnd"
	}

	return validate
})

useForm({
	validationSchema: validates,
})

const handleDataTable = (table: DataTableFunction) => {
	dataTable.value = table
}

const onReset = async () => {
	store.resetParams()
	await onSearch()
}

const onSearch = async () => {
	await store.onSearch()

	store.loading = true

	switch (store.menu) {
		case "maintenance":
			const params: IDashboardMaintenanceTableRequest = {
				road_id: store.params.road_id,
				depot_code: store.params.depot_code,
				km_start: Number.isNaN(convertStringToKm(store.params.km_start))
					? null
					: convertStringToKm(store.params.km_start),
				km_end: Number.isNaN(convertStringToKm(store.params.km_end)) ? null : convertStringToKm(store.params.km_end),
				year: store.params.year,
			}

			await dataTable.value?.searchData(params)
			break
		case "condition":
			store.isSingleRoad = store.params.road_id.length === 1
			break
		default:
			break
	}

	store.loading = false

	// store.searchCall = !store.searchCall
}

const refDonutChart = ref()

onMounted(() => {
	store.setMenuByPermission()
	console.log("menu =", store.menu)
	store.initial()
	console.log("initUserStore", initUserStore)
	console.log("initUserStore 🇹🇭", initUserStore.accessPermissions)

	setTimeout(async () => {
		if (initUserStore.accessPermissions[IUserRolesAccess.view_total_dashboard]) {
			await store.getRoad()
			store.onUpdateData(refDonutChart)
		}
	}, 1000)
})

const conditionArray = [1, 2]
const donutChart = reactive({
	chartOptions: {
		responsive: [
			{
				breakpoint: 1190,
				options: {
					chart: {
						height: 150,
					},
					plotOptions: {
						pie: {
							donut: {
								size: "75%",
								labels: {
									total: {
										fontSize: "10px",
									},
								},
							},
						},
					},
				},
			},
			{
				breakpoint: 991,
				options: {
					chart: {
						height: 168,
					},
					plotOptions: {
						pie: {
							donut: {
								size: "70%",
								labels: {
									total: {
										fontSize: "11px",
									},
								},
							},
						},
					},
				},
			},
			{
				breakpoint: 950,
				options: {
					chart: {
						height: 150,
					},
					plotOptions: {
						pie: {
							donut: {
								size: "70%",
								labels: {
									total: {
										fontSize: "10px",
									},
								},
							},
						},
					},
				},
			},
		],
		plotOptions: {
			pie: {
				donut: {
					size: "70%",
					labels: {
						show: true,
						value: {
							show: true,
							fontSize: "20px",
							fontWeight: 500,
							color: "#FDB833",
						},
						total: {
							show: true,
							showAlways: true,
							label: "จำนวนถนน (สาย)",
							fontSize: "11px",
						},
					},
				},
			},
		},
		stroke: {
			show: false,
		},
		labels: ["หมายเลข 7", "หมายเลข 9"],
		dataLabels: {
			enabled: false,
		},
		colors: ["#FDB833", "#FFDA6A"],
		legend: {
			show: false,
		},
		tooltip: {
			Html: true,
			enabled: true,
			theme: "light",
			custom: function ({ series, seriesIndex, w }: any) {
				let html = `<div class="apexcharts-result">`
				html += `<label class="fs-7">ทางหลวงพิเศษ${w.config.labels[seriesIndex]}</label>
				<label class="fs-7 text-center">${series[seriesIndex]} สาย</label>`
				html += `</div>`
				return html
			},
		},
	},
})

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row mb-5">
		<div class="col-12">
			<div class="card p-5">
				<div class="row">
					<div class="col-md-6 col-12">
						<VTree
							v-model="store.params.depot_code"
							label="หน่วยงานที่รับผิดชอบ"
							:multiple="true"
							:searchable="true"
							:options="store.getOwnerOptions"
							placeholder="ทั้งหมด"
							:name="`depot_code`"
							:limit="0"
							mode="LEAF_PRIORITY"
						/>
					</div>
					<div class="col-md-6 col-12">
						<VTree
							v-model="store.params.road_id"
							label="สายทาง"
							:multiple="true"
							:searchable="true"
							:options="store.roadsOptions"
							placeholder="ทั้งหมด"
							:name="`road_id`"
							:limit="0"
							mode="LEAF_PRIORITY"
						/>
					</div>
					<div class="km">
						<VTextInput v-model="store.params.km_start" label="ช่วง กม." :name="`km_start`" placeholder="00+000" />
					</div>
					<div class="dash d-sm-flex d-none align-items-center p-0" style="margin-top: 2em">
						<span>-</span>
					</div>
					<div class="km">
						<VTextInput
							v-model="store.params.km_end"
							style="margin-top: 2.3em"
							label=""
							:name="`km_end`"
							placeholder="00+000"
						/>
					</div>
					<div class="col-sm-3">
						<VSelect
							v-model="store.params.year"
							:options="store.yearOptions"
							label="ปี"
							name="budget_year"
							:close-on-select="true"
							:can-deselect="false"
						/>
					</div>
					<div class="col-sm-2 align-self-end text-sm-start text-end d-flex">
						<BtnSearch :disabled="store.loading" @click="onSearch" />
						<button
							type="button"
							class="btn rounded-4 ms-5 mt-md-0 mt-sm-2 mt-2 fw-semibold text-gray-700"
							:disabled="store.loading"
							@click="onReset"
						>
							รีเซ็ต
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>

	<div v-show="initUserStore.accessPermissions[IUserRolesAccess.view_total_dashboard]" class="row">
		<div :class="!collapsed ? 'card-graph' : 'card-graph-collapsed'">
			<div class="card">
				<div class="graph-container row justify-content-center">
					<div class="col-graph d-flex align-items-center p-0">
						<ClientOnly>
							<apexchart
								ref="refDonutChart"
								type="donut"
								height="168"
								:options="donutChart.chartOptions"
								:series="store.getQuantityRoadSeries.series"
							/>
						</ClientOnly>
					</div>
					<div class="col-condition align-self-center">
						<div class="row d-flex justify-content-center">
							<template v-for="(item, key) in store.getRoadsLegend.legend" :key="key">
								<div
									class="col-auto d-flex align-items-center cursor-pointer condition-item mx-1"
									:class="conditionArray.includes(key + 1) ? '' : 'selected'"
								>
									<div class="square my-2 me-1" :style="`background: ${item.color}`"></div>
									<span>{{ item.name }}</span>
								</div>
							</template>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div
			v-for="(road, index) of store.getRoadsLenght"
			:key="index"
			:class="!collapsed ? 'card-summary' : 'card-summary-collapsed'"
		>
			<div
				id="slide-1"
				class="card px-3 py-5 row me-1"
				:style="index === 0 ? 'border-left: 16px solid #fdb833' : 'border-left: 16px solid #FFDA6A'"
			>
				<span class="text-gray-800 owner col-12">{{ road.name }}</span>
				<div class="col-12 d-flex py-2">
					<span class="text-primary fw-semibold text-total d-flex align-items-center me-3">{{ road.total }}</span>
					<span class="align-self-center text-gray-600 fw-semibold">กม.</span>
				</div>
				<div class="col-12 text-gray-600 fs-7 pe-0">
					<span>ความยาวผิวทางคอนกรีต </span>
					<span class="text-center text-primary fw-semibold ps-0">{{ road.concrete }}</span>
					<span class="col-auto text-gray-600 fs-7 pe-0 ps-1">%</span>
				</div>
				<div class="col-12 text-gray-600 fs-7 pe-0">
					<span class="col-auto text-gray-600 fs-7 pe-0">ความยาวผิวทาง</span>
					<span class="col-auto text-center text-primary fw-semibold pe-0 ps-1">{{ road.asphalt }}</span>
					<span class="col-auto text-gray-600 fs-7 pe-0 ps-1">%</span>
				</div>
			</div>
		</div>
		<div
			v-for="(road, index) of store.getRoadsAadt"
			:key="index"
			:class="!collapsed ? 'card-traffic' : 'card-traffic-collapsed'"
		>
			<div class="card px-7 py-5 w-100 d-flex justify-content-around">
				<h5 class="fw-semibold">{{ road.name }}</h5>
				<div class="d-flex flex-wrap">
					<span class="text-primary d-flex align-items-center fw-semibold me-2" style="font-size: 2rem">
						{{ toNumber(road.aadt) }}
					</span>
					<span
						class="d-flex align-items-center"
						:class="road.growth_rate === 'up' ? 'text-success' : 'text-danger'"
						style="font-size: 1rem; margin-top: 0.5em"
					>
						{{ toNumber(road.percent, 2) }}%
					</span>
				</div>
				<span class="text-gray-600">เปรียบเทียบระหว่างปี {{ road.year1 }} - {{ road.year2 }}</span>
			</div>
		</div>
	</div>

	<DashboardTab
		v-if="store.hasPermissionAccressTab"
		:current-menu="store.menu"
		:loading="store.loading"
		:map-collapsed="collapsed"
		:search-call="store.searchCall"
		@data-table="handleDataTable"
		@tab-menu="store.handleMenu"
	/>
</template>

<style scoped lang="scss">
.card-graph {
	width: 35%;
	margin-bottom: 1.25rem;
	.card {
		height: 100%;
		.graph-container {
			height: 100%;
		}
	}
	.col-graph {
		width: 55%;
	}
	.col-condition {
		width: 35%;
		.condition-item {
			padding: 0;
		}
	}
	@media (max-width: 9999px) and (min-width: 1600px) {
		width: 20%;
		.col-graph {
			width: 100%;
		}
		.col-condition {
			width: 90%;
			.condition-item {
				padding: 0;
			}
		}
	}
	@media (max-width: 1200px) and (min-width: 992px) {
		width: 35%;
		height: auto;
		.col-graph {
			width: 100%;
		}
		.col-condition {
			width: 60%;
		}
	}
	@media (max-width: 991px) {
		.col-graph {
			width: 60%;
		}
		.col-condition {
			width: 40%;
			.condition-item {
				padding: 0;
			}
		}
	}
	@media (max-width: 900px) and (min-width: 576px) {
		.col-graph {
			width: 100%;
		}
		.col-condition {
			width: 80%;
			.condition-item {
				padding: 0;
			}
		}
	}
	@media (max-width: 575px) {
		width: 100%;
		.col-graph {
			width: 50%;
		}
		.col-condition {
			width: 40%;
			.condition-item {
				padding: 0;
			}
		}
	}
}
.card-graph-collapsed {
	width: 35%;
	margin-bottom: 1.25rem;
	.card {
		height: 100%;
		.graph-container {
			height: 100%;
		}
	}
	.col-graph {
		width: 50%;
	}
	.col-condition {
		width: 50%;
	}
	@media (max-width: 9999px) and (min-width: 1188px) {
		width: 18%;
		height: auto;
		.col-graph {
			width: 100%;
		}
		.col-condition {
			width: 90%;
		}
	}
	@media (max-width: 1189px) and (min-width: 576px) {
		width: 35%;
		height: auto;
		margin-bottom: 1em;
		.col-graph {
			width: auto;
		}
		.col-condition {
			width: 100%;
		}
	}
	@media (max-width: 575px) {
		width: 100%;
		margin-bottom: 1em;
		.col-graph {
			width: 50%;
		}
		.col-condition {
			width: 50%;
			.condition-item {
				padding: 10px;
			}
		}
	}
}
.card-summary {
	width: 32.5%;
	margin-bottom: 1.25rem;
	.card {
		height: 100%;
	}
	.owner {
		width: 80%;
	}
	@media (max-width: 9999px) and (min-width: 1600px) {
		width: 20%;
	}
	@media (max-width: 1189px) and (min-width: 992px) {
		width: 32.5%;
		.owner {
			width: 80%;
		}
	}
	@media (max-width: 991px) {
		.owner {
			width: 80%;
		}
	}
	@media (max-width: 575px) {
		width: 50%;
		.owner {
			width: 100%;
		}
	}
}
.card-summary-collapsed {
	width: 32.5%;
	height: "auto";
	margin-bottom: 1.25rem;
	@media (max-width: 9999px) and (min-width: 1188px) {
		width: 19.8%;
		.card {
			height: 100%;
		}
	}
	@media (max-width: 1189px) and (min-width: 576px) {
		margin-bottom: 1em;
		.card {
			height: 100%;
		}
	}
	@media (max-width: 575px) {
		width: 50%;
		margin-bottom: 1em;
	}
}
.card-traffic {
	height: auto;
	width: 50%;
	margin-bottom: 1.25rem;
	@media (max-width: 9999px) and (min-width: 1600px) {
		width: 20%;
		.card {
			height: 100%;
		}
	}
}
.card-traffic-collapsed {
	width: 50%;
	margin-bottom: 1.25rem;
	@media (max-width: 9999px) and (min-width: 1188px) {
		width: 21.2%;
		.card {
			height: 100%;
		}
	}
}
.text-total {
	font-size: 2rem;
}
.square {
	width: 15px;
	height: 15px;
	border-radius: 5px;
	display: inline-block;
}
.selected {
	opacity: 0.2;
}
.km {
	width: 24.7%;
	@media (max-width: 575px) {
		width: 100%;
	}
}
.dash {
	width: fit-content;
}
</style>
