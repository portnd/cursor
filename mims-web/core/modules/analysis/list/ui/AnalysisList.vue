<script setup lang="ts">
import { useAnalysisListStore } from "../store/AnalysisListStore"
import { IAnalysisParams } from "../infrastructure"
import ServerSideDataTable from "~/core/modules/common/datatable/ui/ServerSideDataTable.vue"
import type { THeader } from "~~/core/shared/types/Datatable"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const store = useAnalysisListStore()
const canEdit = ref<boolean>()

onMounted(async () => {
	await store.getStrategics()
	canEdit.value = initUserStore.accessPermissions[IUserRolesAccess.manage_myself_maintenance_analysis]
})

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 50 },
	{ text: "ชื่อรายการ", value: "name", width: 250 },
	{ text: "ประเภท", value: "type", width: 200 },
	{ text: "ความเห็น", value: "comment", width: 200 },
	{ text: "วิเคราะห์เมื่อ", value: "analysis_date", width: 200 },
	{ text: "สถานะ", value: "status", width: 200 },
]

const dataTable = ref()
const search = reactive<IAnalysisParams>(store.params)

const onSearch = () => {
	dataTable.value.searchData(search)
}

const onResetSearch = () => {
	store.resetParams()
	onSearch()
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<div class="row mb-3">
					<div class="col-xl-2 col-lg-3 col-md-3 col-12 mb-2">
						<VSelect
							v-model="store.params.type"
							:options="store.getStrategicsOptions"
							label="ประเภท"
							name=""
							:close-on-select="true"
							:can-clear="true"
							:can-deselect="true"
						/>
					</div>
					<div class="col-xl-2 col-lg-3 col-md-3 col-12 mb-2">
						<VTextInput v-model="store.params.condition" label="ชื่อรายการ" name="condition" @keyup.enter="onSearch" />
					</div>
					<div class="col-xl-3 col-md-6 col-12 mb-2 align-self-end">
						<BtnSearch @click="onSearch" />
						<button
							type="button"
							class="btn btn-outline-primary rounded-4 ms-5 mt-md-0 mt-3 fw-semibold text-gray-700"
							@click="onResetSearch()"
						>
							รีเซ็ต
						</button>
					</div>
					<div v-if="canEdit" class="col-xl-5 col-md-12 col-12 mb-2 text-end align-self-end">
						<NuxtLink
							to="/analyses/strategic"
							class="btn btn-outline btn-outline-primary rounded-4 mb-2 mb-md-0 px-5 fw-semibold"
						>
							<i class="fi fi-rr-plus align-middle fs-8"></i>
							บำรุงรักษาเชิงกลยุทธ์
						</NuxtLink>
						<NuxtLink to="/analyses/annual" class="btn btn-outline btn-outline-primary rounded-4 ms-5 px-5 fw-semibold">
							<i class="fi fi-rr-plus align-middle fs-8"></i>
							บำรุงรักษาประจำปี
						</NuxtLink>
					</div>
				</div>
				<div class="col-12">
					<ServerSideDataTable ref="dataTable" :headers="headers" :url="'/analyze'">
						<template #item-no="{ item }">
							<NuxtLink class="text-black" :class="store.togglerCursor(item.status)" @click="store.toDetailsPage(item)">
								<div class="text-center">{{ item.no }}</div>
							</NuxtLink>
						</template>
						<template #item-name="{ item }">
							<NuxtLink class="text-black" :class="store.togglerCursor(item.status)" @click="store.toDetailsPage(item)">
								<div class="text-center">{{ item.name === "" ? "-" : item.name }}</div>
							</NuxtLink>
						</template>
						<template #item-type="{ item }">
							<NuxtLink class="text-black" :class="store.togglerCursor(item.status)" @click="store.toDetailsPage(item)">
								<div class="text-center">{{ item.type_analysis }}</div>
							</NuxtLink>
						</template>
						<template #item-condition="{ item }">
							<NuxtLink class="text-black" :class="store.togglerCursor(item.status)" @click="store.toDetailsPage(item)">
								<div class="text-start text-preline my-2">
									<span class="fw-semibold">สายทาง:</span> {{ item.condition.road_name?.split(",")?.join(", ") }}
									<i v-show="item.is_favorite" class="fi fi-ss-star align-middle fs-5 text-warning"></i>
									<div v-html="store.generateCondition(item.condition)"></div>
								</div>
							</NuxtLink>
						</template>
						<template #item-comment="{ item }">
							<NuxtLink
								class="text-black data-table"
								:class="store.togglerCursor(item.status)"
								@click="store.toDetailsPage(item)"
							>
								<div
									class="text-wrap text-preline"
									:class="item.comment === '' ? 'text-center' : 'text-start'"
									v-html="store.generateComment(item.comment)"
								></div>
							</NuxtLink>
						</template>
						<template #item-analysis_date="{ item }">
							<NuxtLink class="text-black" :class="store.togglerCursor(item.status)" @click="store.toDetailsPage(item)">
								<div class="text-center">
									{{
										item.analysis_date === ""
											? "-"
											: buddhistFormatDate(item.analysis_date, "dd mmm yyyy เวลา HH:ii น.")
									}}
								</div>
							</NuxtLink>
						</template>
						<template #item-status="{ item }">
							<NuxtLink
								class="cursor-pointer badge"
								:class="store.generateStatusColors(item.status)"
								@click="store.toDetailsPage(item)"
							>
								<div class="text-center">{{ `${item.status} (${item.percentage}%)` }}</div>
							</NuxtLink>
						</template>
					</ServerSideDataTable>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
@mixin width-size($max_width: false) {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	display: block;
	@if $max_width {
		max-width: 300px;
	} @else {
		max-width: 450px;
	}
}
.truncate-condition {
	@include width-size($max_width: false);
}

.truncate-comment {
	@include width-size($max_width: true);
}

i.position-absolute {
	top: 5px;
}

.data-table {
	align-items: start !important;
}

@media (max-width: 9999px) and (min-width: 1400px) {
	.col-tool {
		width: calc(100% / 4);
	}
}
@media (min-width: 1300px) {
	.col-tool {
		width: 100%;
	}
}
</style>
