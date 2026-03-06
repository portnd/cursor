<script setup lang="ts">
import { useForm } from "vee-validate"
import { useMaintenanceHistorySearchTableStore } from "../store/MaintenanceHistoryProjectListTableStore"
import type { THeader } from "~~/core/shared/types/Datatable"
import ServerSideDataTable from "~/core/modules/common/datatable/ui"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()
const route = useRoute()

const store = useMaintenanceHistorySearchTableStore()
const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 75 },
	{ text: "ชื่อโครงการ", value: "name", width: 200 },
	{ text: "หน่วยงาน", value: "division", width: 200 },
	{ text: "สายทาง", value: "road", width: 200 },
	{ text: "ปีงบประมาณ", value: "year", width: 100 },
	{ text: "ราคาจัดซื้อจัดจ้าง", value: "price", width: 150 },
	{ text: "ประเภทงบประมาณ", value: "type_budget", width: 150 },
	{ text: "วิธีการซ่อมบำรุง", value: "intervention_criteria", width: 150 },
	// { text: "ประเภทการซ่อมบำรุง", value: "type_maintenance", width: 350 },
	// { text: "งบประมาณการซ่อมบำรุง", value: "budget", width: 175 },
	{ text: "ระยะทางรวม (กม.)", value: "total_distance", width: 150 },
	{ text: "วันที่ตรวจรับงานงวดสุดท้าย", value: "due_date" },
	{ text: "วันที่หมดการค้ำประกัน", value: "expiration" },
	// { text: "จัดการ", value: "oparetion" },
]

onMounted(async () => {
	// โหลด dropdown แบบ background ไม่ block - ตารางโหลดทันที (ServerSideDataTable is-init)
	Promise.all([
		store.getYearList(),
		store.getBudgetCriteria(),
		store.getDevision(),
		store.getRoadDropdownList(),
	])
	if (Object.keys(route.query).length) {
		store.setQuriesParams(route.query)
		await nextTick()
		await dataTable.value?.searchData(store.params)
	}
})

useForm()

// ค้นหาข้อมูล
const dataTable: Ref = ref()

// const search: IMaintenanceHistorySearch = reactive(store.params)

const onSearch = async () => {
	// store.fetchMaintenanceHistoryList()
	await dataTable.value.searchData(store.params)
}

const resetSearch = async () => {
	store.resetSearch()
	// store.fetchMaintenanceHistoryList()
	await dataTable.value.searchData(store.params)
}

const getLink = (id: number) => {
	return `/maintenances/history/${id}/info`
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5">
				<div class="row mb-3">
					<div class="col-12 col-md-4 mb-2">
						<VSelect
							v-model="store.params.budget_year"
							:options="store.getYearsOptions"
							label="ปีงบประมาณ"
							name="budget_year"
							:close-on-select="true"
						/>
					</div>
					<div class="col-12 col-md-4 mb-2">
						<VTree
							v-model="store.params.budget_method_id"
							label="ประเภทงบประมาณ"
							:multiple="true"
							:searchable="true"
							:options="store.getBudgetCriteriaOptions"
							placeholder="ทั้งหมด"
							:name="`budget_method_id`"
							:limit="0"
							mode="LEAF_PRIORITY"
						/>
					</div>
					<div class="col-12 col-md-4 mb-2">
						<VTree
							v-model="store.params.road_group_id"
							label="สายทาง"
							:multiple="true"
							:searchable="true"
							:options="store.getRoadDropdownOption"
							placeholder="ทั้งหมด"
							:name="`road_group_id`"
							:limit="0"
							mode="LEAF_PRIORITY"
						/>
					</div>
					<div class="col-12 col-md-4 mb-2">
						<VTree
							v-model="store.params.owner_code"
							label="หน่วยงาน"
							:multiple="false"
							:searchable="true"
							:options="store.getDivisionOption()"
							placeholder="ทั้งหมด"
							:name="`owner_code`"
							:limit="0"
							mode="LEAF_PRIORITY"
						/>
					</div>
					<div class="col-12 col-md-4 mb-2">
						<VTextInput v-model="store.params.name" label="ค้นหาด้วยชื่อ" name="name" />
					</div>
					<div class="col-12 col-md-4 col-lg-4 mb-2 align-self-end">
						<BtnSearch class="mt-md-3" @click="onSearch" />
						<button
							type="button"
							class="btn btn-outline-primary rounded-4 ms-5 fw-semibold text-gray-700 align-self-end mt-md-3"
							@click="resetSearch()"
						>
							รีเซ็ต
						</button>
					</div>
					<div class="col-12 mt-8 text-end">
						<NuxtLink
							v-if="
								initUserStore.accessPermissions[IUserRolesAccess.manage_all_maint_history] ||
								initUserStore.accessPermissions[IUserRolesAccess.manage_owner_maint_history]
							"
							type="button"
							class="btn btn-outline btn-outline-primary align-items-center btn-create"
							style="padding-top: 0.8em"
							@click="navigateTo('history/create')"
						>
							+ เพิ่มประวัติ
						</NuxtLink>
						<!-- <BtnCreate class="btn-create" @click="createItem" /> -->
					</div>
				</div>
				<!-- <VDatatable ref="dataTable" :items="data" :headers="headers"> -->
				<ServerSideDataTable
					ref="dataTable"
					url="/maintenance"
					:is-init="Object.keys(route.query).length === 0"
					:headers="headers"
					:rows-per-page="25"
					@get-data="store.setDatas"
				>
					<template #item-no="{ item }">
						<div class="text-center">{{ item.no }}</div>
					</template>
					<template #item-name="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-start">{{ item.name }}</div>
						</NuxtLink>
					</template>
					<template #item-division="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-start">{{ item.owner_name !== "" && item.owner_name ? item.owner_name : "-" }}</div>
						</NuxtLink>
					</template>
					<template #item-road="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-start" v-html="store.getRoadGroupNames(item.road_group_names)"></div>
						</NuxtLink>
					</template>
					<template #item-year="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-center">{{ item.budget_year + 543 }}</div>
						</NuxtLink>
					</template>
					<template #item-price="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-center">{{ toNumber(item.budget_procurement) }}</div>
						</NuxtLink>
					</template>
					<template #item-type_budget="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<!-- <div class="text-center">{{ item.type_budget }}</div> -->
							<div class="text-start">{{ item.budget.name }}</div>
						</NuxtLink>
					</template>
					<template #item-intervention_criteria="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-start">{{ item.budget_method.method_name }}</div>
						</NuxtLink>
					</template>
					<template #item-total_distance="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-end">{{ item.km_total.toFixed(2) }}</div>
						</NuxtLink>
					</template>
					<template #item-due_date="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<div class="text-center">{{ buddhistFormatDate(item.project_end_date, "dd mmm yyyy") }}</div>
						</NuxtLink>
					</template>
					<template #item-expiration="{ item }">
						<NuxtLink :to="getLink(item.id_parent)" class="text-black">
							<!-- <div class="text-center" :style="{ color: item.id === 1 ? '#1F70F3' : '#F1416C' }"> -->
							<div class="text-center" :style="{ color: item.color }">
								{{ buddhistFormatDate(item.guarantee_expiration_date, "dd mmm yyyy") }}
							</div>
						</NuxtLink>
					</template>
				</ServerSideDataTable>
				<!-- ใช้ connect API ใช้ server side นะ -->
				<!-- <ServerSideDataTable ref="dataTable" :headers="headers" url="/maintenance/history">
					<template #item-no="{ item }">
						<div class="text-center">{{ item.no }}</div>
					</template>
					<template #item-name="{ item }">
						<NuxtLink :to="getLink(item.id)" class="text-black">
							<div class="text-start">{{ item.name }}</div>
						</NuxtLink>
					</template>
					<template #item-road="{ item }">
						<NuxtLink :to="getLink(item.id)" class="text-black">
							<div class="text-start">{{ item.road_group.name }}</div>
						</NuxtLink>
					</template>
					<template #item-year="{ item }">
						<NuxtLink :to="getLink(item.id)" class="text-black">
							<div class="text-center">{{ item.budget_year + 543 }}</div>
						</NuxtLink>
					</template>
					<template #item-type_budget="{ item }">
						<NuxtLink :to="getLink(item.id)" class="text-black">
							<div class="text-center">{{ item.budget.name }}</div>
						</NuxtLink>
					</template>
					<template #item-total_distance="{ item }">
						<NuxtLink :to="getLink(item.id)" class="text-black">
							<div class="text-end">{{ (item.km_total / 1000).toFixed(2) }}</div>
						</NuxtLink>
					</template>
					<template #item-end_date="{ item }">
						<NuxtLink :to="getLink(item.id)" class="text-black">
							<div class="text-center">{{ buddhistFormatDate(item.last_inspection_date, "dd mmm yyyy") }}</div>
						</NuxtLink>
					</template>
					<template #item-expiration="{ item }">
						<NuxtLink
							:to="getLink(item.id)"
							class="text-black"
							:class="store.setDateColor(item.guarantee_expiration_date)"
						>
							<div class="text-center">{{ buddhistFormatDate(item.guarantee_expiration_date, "dd mmm yyyy") }}</div>
						</NuxtLink>
					</template>
				</ServerSideDataTable> -->
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.btn-create {
	margin-bottom: -25px !important;

	@media (max-width: 767px) {
		margin-bottom: 1rem !important;
	}
}

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
