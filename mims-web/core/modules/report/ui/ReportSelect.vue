<script setup lang="ts">
import { useInitUserStore } from "../../initUser/store"

// import { useReportStore } from "../store/ReportStore"

// interface IReportList {
// 	name: string
// 	access_key: string
// 	show: boolean
// }

const initUserStore = useInitUserStore()

const emit = defineEmits(["select"])

const activeId = ref(0)
const clickActive = (id: number, accessKey: string) => {
	activeId.value = id
	emit("select", accessKey)
}
const reportList = ref([
	{
		id: 0,
		name: "รายงานสินทรัพย์",
		access_key: "asset",
		show:
			initUserStore.accessPermissions[IUserRolesAccess.download_all_asset_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_asset_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_map_asset_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_map_asset_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_result_asset_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_result_asset_report],
		children: [
			{
				id: 0,
				name: "รายงานสินทรัพย์",
				access_key: "child_asset",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_asset_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_asset_report],
			},
			{
				id: 12,
				name: "รายงานแผนที่สินทรัพย์",
				access_key: "child_asset_map",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_map_asset_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_map_asset_report],
			},
			{
				id: 22,
				name: "รายงานสรุปสินทรัพย์",
				access_key: "child_road_condition_asset",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_result_asset_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_result_asset_report],
			},
		],
	},

	{
		id: 1,
		name: "รายงานข้อมูลสภาพทาง",
		access_key: "road_condition",
		show:
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_surface_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_surface_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_condition_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_condition_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_condition_result_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_condition_result_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_retro_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_retro_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_retro_result_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_retro_result_report],

		children: [
			{
				id: 45,
				name: "รายงานสรุปข้อมูลสายทาง",
				access_key: "child_road_summary",
				show: initUserStore.accessPermissions[IUserRolesAccess.download_all_road_report],
			},
			{
				id: 46,
				name: "รายงานสรุปรายละเอียดชนิดผิวทาง",
				access_key: "child_road_surface",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_surface_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_surface_report],
			},
			{
				id: 47,
				name: "รายงานข้อมูลสภาพทาง",
				access_key: "child_road_condition",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_road_condition_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_condition_report],
			},
			{
				id: 48,
				name: "รายงานสรุปข้อมูลสภาพทาง",
				access_key: "child_road_condition_summary",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_road_condition_result_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_condition_result_report],
			},
			{
				id: 49,
				name: "รายงานค่าการสะท้อนแสงของเส้นจราจร",
				access_key: "child_reflect_light",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_road_retro_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_retro_report],
			},
			{
				id: 50,
				name: "รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร",
				access_key: "child_reflect_light_summary",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_road_retro_result_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_retro_result_report],
			},
		],
	},

	{
		id: 2,
		name: "รายงานข้อมูลความเสียหาย",
		access_key: "road_damage",
		show:
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_damage_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_damage_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_all_road_damage_result_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_damage_result_report],
		children: [
			{
				id: 51,
				name: "รายงานข้อมูลความเสียหาย",
				access_key: "child_road_damage",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_road_damage_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_damage_report],
			},
			{
				id: 52,
				name: "รายงานสรุปข้อมูลความเสียหาย",
				access_key: "child_road_damage_summary",
				show:
					initUserStore.accessPermissions[IUserRolesAccess.download_all_road_damage_result_report] ||
					initUserStore.accessPermissions[IUserRolesAccess.download_owner_road_damage_result_report],
			},
		],
	},

	{
		id: 3,
		name: "รายงานการซ่อมบำรุงตามเกณฑ์ KPI",
		access_key: "maintenance",
		show:
			initUserStore.accessPermissions[IUserRolesAccess.download_all_maint_kpi_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_maint_kpi_report],
	},

	{
		id: 4,
		name: "รายงานประวัติการซ่อมบำรุง",
		access_key: "project_maintenance",
		show:
			initUserStore.accessPermissions[IUserRolesAccess.download_all_maint_history_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_maint_history_report],
	},

	{
		id: 5,
		name: "รายงานปริมาณจราจร",
		access_key: "volume_aadt",
		show:
			initUserStore.accessPermissions[IUserRolesAccess.download_all_aadt_report] ||
			initUserStore.accessPermissions[IUserRolesAccess.download_owner_aadt_report],
	},
])
</script>

<template>
	<div class="card p-6 mb-8">
		<h4 class="fw-semibold mb-5">
			<div class="d-flex align-items-center gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none">
					<g clip-path="url(#clip0_742_47655)">
						<path
							d="M21.414 5H17V0.586L21.414 5ZM22 7V24H2V3C2 2.20435 2.31607 1.44129 2.87868 0.87868C3.44129 0.31607 4.20435 0 5 0L15 0V7H22ZM7 16H14V14H7V16ZM17 18H7V20H17V18ZM17 10H7V12H17V10Z"
							fill="#FDB833"
						/>
					</g>
					<defs>
						<clipPath id="clip0_742_47655">
							<rect width="24" height="24" fill="white" />
						</clipPath>
					</defs>
				</svg>
				<span>ขั้นตอนที่ 1 : เลือกรายงาน</span>
			</div>
		</h4>

		<div id="report" class="accordion">
			<template v-for="(report, index) of reportList.filter((e) => e.show)" :key="report.id">
				<div v-if="report.children" class="accordion-item card border border-1 border-gray-300 shadow-none mb-5">
					<h2 :id="`#report_header${index}`" class="accordion-header">
						<button
							class="accordion-button fs-5 fw-semibold py-5 text-gray-800 bg-white shadow-none"
							:class="index !== 0 ? 'collapsed' : ''"
							type="button"
							data-bs-toggle="collapse"
							:data-bs-target="`#report_body${index}`"
							:aria-controls="`report_body${index}`"
						>
							{{ report.name }}
						</button>
					</h2>
					<div
						:id="`report_body${index}`"
						data-bs-parent="#report"
						:aria-labelledby="`#report_header${index}`"
						class="accordion-collapse collapse"
						:class="{ show: index === 0 }"
					>
						<div class="p-5 pt-0">
							<div id="nav-tab" class="row" role="tablist">
								<div v-for="child of report.children.filter((e) => e.show)" :key="child.id" class="col-12 mb-4">
									<div
										class="border border-1 border-gray-300 fs-6 fw-normal rounded-3 py-4 px-5 bg-hover-gray-100 cursor-pointer"
										:class="activeId === child.id ? 'active' : ''"
										data-bs-toggle="tab"
										role="tab"
										aria-selected="true"
										@click="clickActive(child.id, child.access_key)"
									>
										{{ child.name }}
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div v-else id="nav-tab" class="row" role="tablist">
					<div class="col-12 mb-4">
						<div
							class="border border-1 border-gray-300 fs-6 fw-normal rounded-3 py-6 px-5 bg-hover-gray-100 cursor-pointer"
							:class="activeId === report.id ? 'active' : ''"
							data-bs-toggle="tab"
							role="tab"
							aria-selected="true"
							@click="clickActive(report.id, report.access_key)"
						>
							<span class="fs-5 fw-semibold text-gray-800">{{ report.name }}</span>
						</div>
					</div>
				</div>
			</template>
		</div>
	</div>
</template>

<style lang="scss" scoped>
#nav-tab .active {
	background-color: var(--kt-primary-light) !important;
}

.accordion,
.accordion-button {
	border-radius: 16px !important;
	z-index: 0 !important;
	box-shadow: none !important;
}
.accordion-item,
.accordion-item .accordion-button.collapsed {
	border-top-right-radius: 16px !important;
	border-top-left-radius: 16px !important;
	border-bottom-left-radius: 16px !important;
	border-bottom-right-radius: 16px !important;
}

.accordion-item .accordion-item {
	border: 1px solid var(--kt-gray-300);
	box-shadow: none;
}
</style>
