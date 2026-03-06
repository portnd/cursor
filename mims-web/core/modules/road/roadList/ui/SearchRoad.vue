<script setup lang="ts">
import * as bootstrap from "bootstrap"
import { useForm } from "vee-validate"
import { useRoadListStore } from "../store/index"
import FilterModal from "./FilterModal.vue"
import { IValidate } from "~/core/shared/types/Validate"
import { ITree } from "~/core/shared/types/Tree"

const store = useRoadListStore()
const validate = computed(() => {
	const validations: IValidate = {}
	validations.km_start = store.params.km_start === "" ? "" : "km"
	validations.km_end = store.params.km_end === "" ? "" : "km"

	return validations
})

const { resetField } = useForm({ validationSchema: validate })

const closeDropdown = () => {
	const dropdownMenu = document.querySelector(".dropdownSearchAdvanced")
	const dropdownInstance = new bootstrap.Dropdown(dropdownMenu as HTMLElement)
	dropdownInstance.hide()
}

// const mock = ref("1000 รายการ")

const conditionModal: Ref = ref()

const filterItem = () => {
	conditionModal.value.showModal()
}

const roadGroupOptions = computed(() => {
	const roadGroupData = useInitData()?.refRoadGroup()
	const options = roadGroupData?.map((item) => ({
		id: item.id.toString(),
		label: "ทางหลวงพิเศษหมายเลข " + item.number,
	}))
	return options ?? []
})

const districtOptions = computed(() => {
	const initData = useInitData()
	const districts = initData?.refDistrict()
	const roadGroupData = initData?.refRoadGroup()
	const hasRoadGroup = (store.params.road_group_id?.length ?? 0) > 0 && store.params.road_group_id
	const roadGroupIdSet = new Set(store.params.road_group_id?.map(Number))

	const districtCodes = roadGroupData
		?.reduce((acc: string[], item) => {
			if (roadGroupIdSet.has(item.id)) {
				acc.push(...item.ref_district_codes)
			}
			return acc
		}, [])
		.map(Number)

	const districtCodesSet = hasRoadGroup ? new Set(districtCodes?.map(Number)) : null

	const options = districts
		?.filter((district) => !hasRoadGroup || districtCodesSet?.has(Number(district.district_code)))
		.map((district) => ({
			id: `parennt${district.id.toString()}`,
			label: district.name,
			children: district.depots.map((depot) => ({
				id: depot.depot_code,
				label: depot.name,
			})),
		}))

	return options ?? []
})

const roadSectionOptions = computed(() => {
	const initData = useInitData()
	const roadSections = initData?.refRoadSection()
	const districts = initData?.refDistrict()
	const hasDepotCodes = (store.params.depot_code?.length ?? 0) > 0 && store.params.depot_code
	const depotCodeSet = new Set(store.params.depot_code?.map(Number))
	const hasRoadGroup = (store.params.road_group_id?.length ?? 0) > 0 && store.params.road_group_id

	const roadSectionId =
		districts?.reduce((acc: number[], district) => {
			district.depots.forEach((depot) => {
				if (depotCodeSet.has(Number(depot.depot_code))) {
					acc.push(...depot.road_section_id)
				}
			})
			return acc
		}, []) ?? []

	const roadSectionIdSet = new Set(roadSectionId)

	const roadGroupData = useInitData()?.refRoadGroup() || []

	roadGroupData.forEach((group) => {
		group.roadGroups = (roadSections || [])
			.filter((obj) => Number(obj.road_group_id) === group.id)
			.sort((a, b) => Number(a.number) - Number(b.number))
	})

	const sortRoadSections = roadGroupData.flatMap((group) => group.roadGroups)

	const options =
		// roadSections
		sortRoadSections.reduce((acc: ITree[], road) => {
			if (!hasDepotCodes || roadSectionIdSet.has(road.id)) {
				if (hasRoadGroup) {
					if (hasRoadGroup.includes(road.road_group_id.toString())) {
						acc.push({
							id: road.id.toString(),
							label: `${road.number} ${road.name_origin} - ${road.name_destination}`,
						})
					}
				} else {
					acc.push({
						id: road.id.toString(),
						label: `${road.number} ${road.name_origin} - ${road.name_destination}`,
					})
				}
			}
			return acc
		}, []) ?? []

	return options ?? []
})

const surfaceOptions = computed(() => {
	const hasSection = (store.params.road_section_id?.length ?? 0) > 0 && store.params.road_section_id
	const sectionSet = new Set(store.params.road_section_id?.map(Number))
	const initData = useInitData()
	const roadSections = initData?.refRoadSection()
	const roadSurface = initData?.refSurface()
	const surfaceId: number[] = []
	roadSections?.forEach((section) => {
		if (sectionSet.has(Number(section.id))) {
			surfaceId.push(...section.ref_surface_id)
		}
	})
	const roadSurfaceIdSet = new Set(surfaceId)
	const options =
		roadSurface?.reduce((acc: ITree[], surface) => {
			if (!hasSection || roadSurfaceIdSet.has(surface.id)) {
				acc.push({
					id: surface.id.toString(),
					label: `${surface.name}`,
				})
			}
			return acc
		}, []) ?? []
	if (hasSection) {
		store.params.ref_surface_id = []
	}

	return options || []
})

const resetSearch = () => {
	// store.params.depot_code = []
	// store.params.keyword = ""
	// store.params.road_group_id = []
	// store.params.road_section_id = []
	// store.params.km_start = ""
	// store.params.km_end = ""
	navigateTo("roads")
	store.$reset()
	store.getData()
}

const searchData = () => {
	store.getData()
	closeDropdown()
}

// const handleKeyDown = (event: any) => {
// 	if (event.key === "Enter") {
// 		console.log("enter")
// 		store.getData()
// 	}
// }

// onMounted(() => {
// 	window.addEventListener("keydown", handleKeyDown)
// })

// onUnmounted(() => {
// 	window.removeEventListener("keydown", handleKeyDown)
// })
</script>

<template>
	<div class="card card-rounded mb-5">
		<div class="card-body p-4">
			<div class="d-flex">
				<i class="fi-rr-search fs-3 text-gray-600 position-absolute translate-middle-y ms-5 icon-search"></i>
				<input
					v-model="store.params.keyword"
					type="text"
					class="form-control text-gray-800 ps-14 pe-14"
					placeholder="ค้นหาด้วยชื่อหรือรหัสสายทาง เช่น ศรีนครินทร์, 00070101"
					@keyup.enter="store.getData"
				/>
				<div class="btn btn-sm btn-icon position-absolute z-index-3 icon-settings">
					<i
						id="dropdownSearchAdvanced"
						class="fi-rr-settings-sliders fs-3 text-gray-600"
						data-bs-toggle="dropdown"
						aria-expanded="false"
					></i>
					<div
						class="dropdown-menu p-7 rounded-5 w-lg-400px dropdownSearchAdvanced"
						aria-labelledby="dropdownSearchAdvanced"
						data-bs-auto-close="outside"
						@click.stop
					>
						<h4 class="fw-semibold text-gray-800 mb-5">ค้นหาขั้นสูง</h4>
						<div class="row">
							<div class="col-12 mb-1">
								<VTree
									v-model="store.params.road_group_id"
									:options="roadGroupOptions"
									label="สายทาง"
									:multiple="true"
									mode="LEAF_PRIORITY"
									name="road_group_id"
									placeholder="ทั้งหมด"
									:searchable="true"
									:limit="0"
									@update:model-value="resetField('depot_code', { value: [] })"
								/>
							</div>
							<div class="col-12 mb-1">
								<VTree
									v-model="store.params.depot_code"
									:options="districtOptions"
									label="หน่วยงานที่รับผิดชอบ"
									:multiple="true"
									mode="LEAF_PRIORITY"
									name="depot_code"
									placeholder="ทั้งหมด"
									:searchable="true"
									:limit="0"
									@update:model-value="() => resetField('road_section_id', { value: [] })"
								/>
							</div>
							<div class="col-12 mb-1">
								<VTree
									v-model="store.params.road_section_id"
									:options="roadSectionOptions"
									label="ตอนควบคุม"
									name="road_section_id"
									:multiple="true"
									mode="BRANCH_PRIORITY"
									placeholder="ทั้งหมด"
									:searchable="true"
									:limit="0"
								/>
							</div>
							<div class="col-12 mb-1">
								<VLabel label="ช่วงกม.เริ่มต้น - สิ้นสุด" />
								<div class="row">
									<div class="col-6">
										<VTextInput v-model="store.params.km_start" label="" name="km_start" placeholder="เริ่มต้น" />
									</div>
									<div class="col-6">
										<VTextInput v-model="store.params.km_end" label="" name="km_end" placeholder="สิ้นสุด" />
									</div>
								</div>
							</div>

							<div class="col-12 mb-1">
								<VTree
									v-model="store.params.ref_surface_id"
									placeholder="ทั้งหมด"
									label="ชนิดผิว"
									:multiple="true"
									:limit="0"
									:options="surfaceOptions"
									name="surface"
									mode="LEAF_PRIORITY"
								/>
							</div>
							<div class="col-12 mb-1">
								<VLabel label="สภาพทาง" />
								<div class="position-relative">
									<VTextInput
										v-model="store.condition_value"
										class="text-input"
										placeholder="ทั้งหมด"
										:disabled="true"
										name="road-condition"
									/>
									<button type="button" class="btn btn-primary position-absolute select-btn" @click="filterItem">
										เลือก
									</button>
								</div>
							</div>
						</div>

						<div class="d-flex justify-content-end mt-8">
							<BtnCancel @click="closeDropdown" />
							<BtnSubmit label="ตกลง" :loading="store.loading" @click="searchData" />
						</div>
					</div>
				</div>
				<BtnSubmit label="ค้นหา" :disabled="store.loading" class="ms-5 mt-0 btn-search" @click.stop="store.getData" />
				<button type="button" class="btn rounded-4 ms-5 fw-semibold text-gray-700" @click="resetSearch()">
					รีเซ็ต
				</button>
			</div>
		</div>
	</div>
	<FilterModal ref="conditionModal" />
</template>

<style scoped lang="scss">
.icon-search {
	top: 55%;
}

.icon-settings {
	top: 25%;
	right: 205px;
}

.icon-settings i {
	width: 50px;
	height: 25px;
	padding: 0.25rem 0 0.25rem;
	justify-content: center;
}

.dropdownSearchAdvanced {
	min-width: 320px;
	cursor: auto;
}

.btn-search {
	margin-top: 0px !important;
}

.select-btn {
	top: -2px !important;
	right: 0;
	border-top-left-radius: 0;
	border-bottom-left-radius: 0;
	border-left: 1px solid #e4e6ef;
}
</style>
