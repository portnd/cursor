<script setup lang="ts">
import { useInitUserStore } from "../../initUser/store/InitUserStore"
import { DashboardAsset, DashboardSurface, DashboardMaintenance } from "./index"
import { DashboardCondition } from "~/core/modules/dashboard/condition/ui/index"

defineProps({
	mapCollapsed: {
		type: Boolean,
		default: false,
	},
	searchCall: {
		type: Boolean,
		default: false,
	},
	loading: {
		type: Boolean,
		default: false,
	},
	currentMenu: {
		type: String,
		default: "asset",
	},
})

interface DataTableFunction {
	getData: Function
	loadData: Function
	searchData: Function
	resetSelected: Function
	resetRow: Function
}

const initUserStore = useInitUserStore()

const emit = defineEmits(["tabMenu", "dataTable"])
const changeTab = (tab: string) => {
	emit("tabMenu", tab)
}

const handleDataTable = (dataTable: DataTableFunction) => {
	emit("dataTable", dataTable)
}
</script>
<template>
	<div class="row">
		<div class="col-12">
			<div class="card shadow">
				<div class="card-body p-5">
					<div class="mt-0">
						<ul class="nav nav-tabs nav-line-tabs mb-5">
							<li
								v-show="
									initUserStore.accessPermissions[IUserRolesAccess.view_all_asset_dashboard] ||
									initUserStore.accessPermissions[IUserRolesAccess.view_owner_asset_dashboard]
								"
								class="nav-item"
								:class="currentMenu === 'asset' ? 'active' : ''"
								data-bs-toggle="tab"
								data-bs-target="#detail-asset"
								role="tab"
								aria-selected="true"
								@click="changeTab('asset')"
							>
								<span class="nav-link cursor-pointer">ข้อมูลสรุปทรัพย์สิน</span>
								<span class="line"></span>
							</li>
							<li
								v-show="
									initUserStore.accessPermissions[IUserRolesAccess.view_all_road_condition_dashboard] ||
									initUserStore.accessPermissions[IUserRolesAccess.view_owner_road_condition_dashboard]
								"
								class="nav-item"
								:class="currentMenu === 'condition' ? 'active' : ''"
								data-bs-toggle="tab"
								data-bs-target="#detail-condition"
								role="tab"
								aria-selected="false"
								@click="changeTab('condition')"
							>
								<span class="nav-link cursor-pointer">ข้อมูลสรุปสภาพทาง</span>
								<span class="line"></span>
							</li>
							<li
								v-show="
									initUserStore.accessPermissions[IUserRolesAccess.view_all_surface_dashboard] ||
									initUserStore.accessPermissions[IUserRolesAccess.view_owner_surface_dashboard]
								"
								class="nav-item"
								:class="currentMenu === 'surface' ? 'active' : ''"
								data-bs-toggle="tab"
								data-bs-target="#detail-surface"
								role="tab"
								aria-selected="false"
								@click="changeTab('surface')"
							>
								<span class="nav-link cursor-pointer">ข้อมูลสรุปผิวทาง</span>
								<span class="line"></span>
							</li>
							<li
								v-show="
									initUserStore.accessPermissions[IUserRolesAccess.view_all_maint_history_dashboard] ||
									initUserStore.accessPermissions[IUserRolesAccess.view_owner_maint_history_dashboard]
								"
								class="nav-item"
								:class="currentMenu === 'maintenance' ? 'active' : ''"
								data-bs-toggle="tab"
								data-bs-target="#detail-maintenance"
								role="tab"
								aria-selected="false"
								@click="changeTab('maintenance')"
							>
								<span class="nav-link cursor-pointer">ข้อมูลสรุปซ่อมบำรุง</span>
								<span class="line"></span>
							</li>
						</ul>
					</div>

					<div class="row">
						<div v-show="loading === false" class="col-12 order-2">
							<!-- begin::Content -->
							<div class="tab-content p-2">
								<div
									id="detail-asset"
									class="tab-pane fade"
									:class="currentMenu === 'asset' ? 'active show' : ''"
									role="tabpanel"
								>
									<DashboardAsset :map-collapsed="mapCollapsed" />
								</div>
								<div
									id="detail-condition"
									class="tab-pane fade"
									:class="currentMenu === 'condition' ? 'active show' : ''"
									role="tabpanel"
								>
									<!-- <DashboardCondition :map-collapsed="mapCollapsed" /> -->
									<DashboardCondition :map-collapsed="mapCollapsed" :search-call="searchCall" />
								</div>
								<div
									id="detail-surface"
									class="tab-pane fade"
									:class="currentMenu === 'surface' ? 'active show' : ''"
									role="tabpanel"
								>
									<DashboardSurface />
								</div>
								<div
									id="detail-maintenance"
									class="tab-pane fade"
									:class="currentMenu === 'maintenance' ? 'active show' : ''"
									role="tabpanel"
								>
									<DashboardMaintenance @data-table="handleDataTable" />
								</div>
							</div>
							<!-- end::Content -->
						</div>
						<div
							v-if="loading === true"
							class="d-flex justify-content-center align-items-center my-20 h-25"
							style="transform: scale(2)"
						>
							<div class="spinner-border text-primary" role="status">
								<span class="visually-hidden">Loading...</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 22px !important;
}
</style>
