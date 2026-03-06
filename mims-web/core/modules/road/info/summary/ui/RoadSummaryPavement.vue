<script setup lang="ts">
import { IRoadSummaryItem } from "../infrastructure"
import { useRoadSummaryStore } from "../store"
import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"
import type { THeader } from "~~/core/shared/types/Datatable"

const initUserStore = useInitUserStore()

const store = useRoadSummaryStore()
const route = useRoute()

const headers: THeader[] = [
	{ text: "เริ่มต้น", value: "km_start" },
	{ text: "สิ้นสุด", value: "km_end" },
	{ text: "กว้าง (ม.)", value: "width_surface" },
	{ text: "ความหนาผิวทางลาดยาง (ซม.)", value: "thickness_surface" },
	{ text: "ความหนาผิวคอนกรีต (ซม.)", value: "thickness_surface_concrete" },
	{ text: "ไหล่ทาง ซ้าย (ม.)", value: "width_shoulder_left" },
	{ text: "ไหล่ทาง ขวา (ม.)", value: "width_shoulder_right" },
	{ text: "ความหนา Slab (ซม.)", value: "thickness_concrete_slab" },
	{ text: "ชั้น Base (ซม.)", value: "thickness_base" },
	{ text: "ชั้น Subbase (ซม.)", value: "thickness_subbase" },
	{ text: "ชั้น Subgrade (ซม.)", value: "thickness_subgrade" },
]

const activeRowNumber = ref(1)

const selectRow = (item: IRoadSummaryItem) => {
	store.dataInPicture = item
	store.surfaceSectionCode = item?.surface_cross_section_code
	store.setLocation(item)

	// เพิ่ม class ใส่ใน row นั้น ๆ
	activeRowNumber.value = item.no ?? 0
}

const editItem = () => {
	navigateTo(`/roads/${route.params.roadId}/summary/edit?tab=pavement`)
}

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div id="row" class="row mb-3">
			<div class="col-12 col-md-8 text-start align-self-center ps-4 order-last order-md-first">
				<template v-if="store.update_date !== ''">
					<label class="text-gray-900 fs-6 me-1">ปรับปรุงข้อมูลโดย</label>
					<VUser :label="store.update_by.full_name" :name="store.update_by.full_name" :role="'-'" />
					<label class="text-gray-900 fs-6 ms-1">
						เมื่อวันที่
						{{ buddhistFormatDate(new Date(store.update_date), "dd mmm yyyy เวลา HH:ii น.") }}
					</label>
				</template>
			</div>
			<div class="col-12 col-md-4 text-end">
				<NuxtLink
					v-if="
						initUserStore.accessPermissions[IUserRolesAccess.manage_road_summary] ||
						initUserStore.getIsOwnerManagePermission(
							initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_summary],
							store.road.ref_depot.id
						)
					"
					class="btn btn-outline btn-outline-primary rounded-2 px-5 fw-semibold fs-6"
					@click="editItem()"
				>
					ปรับปรุงข้อมูล
				</NuxtLink>
			</div>
		</div>
		<VDatatable
			:headers="headers"
			:items="store.roadSurface"
			:no-border="true"
			active-item-class-name="active"
			:active-row-number="activeRowNumber"
		>
			<template #customize-headers>
				<thead>
					<tr>
						<th colspan="2" class="text-center border-bottom border-1">กม.</th>
						<th colspan="3" class="text-center border-bottom border-1">ผิวทาง</th>
						<th rowspan="2" class="text-center">ไหล่ทาง<br />ซ้าย<br />(ม.)</th>
						<th rowspan="2" class="text-center">ไหล่ทาง<br />ขวา<br />(ม.)</th>
						<th rowspan="2" class="text-center">ความหนา<br />Slab<br />(ซม.)</th>
						<th rowspan="2" class="text-center">ความหนา<br />Base<br />(ซม.)</th>
						<th rowspan="2" class="text-center">ความหนา<br />Subbase<br />(ซม.)</th>
						<th rowspan="2" class="text-center">ความหนา<br />Subgrade<br />(ซม.)</th>
					</tr>
					<tr>
						<th class="text-center border-0">เริ่มต้น</th>
						<th class="text-center">สิ้นสุด</th>
						<th class="text-center">กว้าง (ม.)</th>
						<th class="text-center">
							ความหนา<br />ผิวทางลาดยาง<br />
							(ซม.)
						</th>
						<th class="text-center">ความหนา<br />ผิวคอนกรีต <br />(ซม.)</th>
					</tr>
				</thead>
			</template>

			<!-- begin::Items -->
			<template #item-km_start="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ item.km_start }}</div>
				</div>
			</template>
			<template #item-km_end="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ item.km_end }}</div>
				</div>
			</template>
			<template #item-width_surface="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.width_surface) }}</div>
				</div>
			</template>
			<template #item-thickness_surface="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.thickness_surface) }}</div>
				</div>
			</template>
			<template #item-thickness_surface_concrete="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.thickness_surface_concrete) }}</div>
				</div>
			</template>
			<template #item-width_shoulder_left="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.width_shoulder_left) }}</div>
				</div>
			</template>
			<template #item-width_shoulder_right="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.width_shoulder_right) }}</div>
				</div>
			</template>
			<template #item-thickness_concrete_slab="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.thickness_concrete_slab) ?? "-" }}</div>
				</div>
			</template>
			<template #item-thickness_base="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.thickness_base) ?? "-" }}</div>
				</div>
			</template>
			<template #item-thickness_subbase="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.thickness_subbase) ?? "-" }}</div>
				</div>
			</template>
			<template #item-thickness_subgrade="{ item }">
				<div class="cursor-pointer" @click="selectRow(item)">
					<div class="text-center">{{ toNumber(item.thickness_subgrade) ?? "-" }}</div>
				</div>
			</template>
			<!-- end::Items -->
		</VDatatable>
		<!-- begin::รูปประกอบ -->
		<div class="row mt-5 surfaceSectionCode-container">
			<div v-show="store.surfaceSectionCode === 1" key="surfaceSectionCode1" class="col-12 p-0">
				<div class="surfaceSectionCode-1">
					<div class="image-container">
						<div class="surface-symbol-top">
							<span class="symbol-1">{{ store.dataInPicture.width_shoulder_left }} ม.</span>
							<span class="symbol-2">{{ store.dataInPicture.width_surface }} ม.</span>
							<span class="symbol-3">{{ store.dataInPicture.width_shoulder_right }} ม.</span>
						</div>
						<img src="/images/roads/1_asphalt_on_concrete_deck.svg" class="w-100" />
						<span class="symbol-end-1">{{ store.dataInPicture.thickness_surface }} ซม.</span>
						<span class="symbol-end-2">{{ store.dataInPicture.thickness_base }} ซม.</span>
					</div>
				</div>
			</div>
			<div v-show="store.surfaceSectionCode === 2" key="surfaceSectionCode2" class="col-12 p-0">
				<div class="surfaceSectionCode-2">
					<div class="image-container">
						<div class="surface-symbol-top">
							<span class="symbol-1">{{ store.dataInPicture.width_shoulder_left }} ม.</span>
							<span class="symbol-2">{{ store.dataInPicture.width_surface }} ม.</span>
							<span class="symbol-3">{{ store.dataInPicture.width_shoulder_right }} ม.</span>
						</div>
						<img src="/images/roads/2_asphalt_on_steel_deck.svg" class="w-100" />
						<span class="symbol-end-1">{{ store.dataInPicture.thickness_surface }} ซม.</span>
						<span class="symbol-end-2">{{ store.dataInPicture.thickness_base }} ซม.</span>
					</div>
				</div>
			</div>
			<div v-show="store.surfaceSectionCode === 3" key="surfaceSectionCode3" class="col-12 p-0">
				<div class="surfaceSectionCode-3">
					<div class="image-container">
						<div class="surface-symbol-top">
							<span class="symbol-1">{{ store.dataInPicture.width_shoulder_left }} ม.</span>
							<span class="symbol-2">{{ store.dataInPicture.width_surface }} ม.</span>
							<span class="symbol-3">{{ store.dataInPicture.width_shoulder_right }} ม.</span>
						</div>
						<img src="/images/roads/3_asphalt_on_ground.svg" class="w-100" />
						<span class="symbol-end-1">{{ store.dataInPicture.thickness_surface }} ซม.</span>
						<span class="symbol-end-2">{{ store.dataInPicture.thickness_base }} ซม.</span>
						<span class="symbol-end-3">{{ store.dataInPicture.thickness_subbase }} ซม.</span>
						<span class="symbol-end-4">{{ store.dataInPicture.thickness_subgrade }} ซม.</span>
					</div>
				</div>
			</div>
			<div v-show="store.surfaceSectionCode === 4" key="surfaceSectionCode4" class="col-12 p-0">
				<div class="surfaceSectionCode-4">
					<div class="image-container">
						<div class="surface-symbol-top">
							<span class="symbol-1">{{ store.dataInPicture.width_shoulder_left }} ม.</span>
							<span class="symbol-2">{{ store.dataInPicture.width_surface }} ม.</span>
							<span class="symbol-3">{{ store.dataInPicture.width_shoulder_right }} ม.</span>
						</div>
						<img src="/images/roads/4_composite_pavement.svg" class="w-100" />
						<span class="symbol-end-1">{{ store.dataInPicture.thickness_surface }} ซม.</span>
						<span class="symbol-end-2">{{ store.dataInPicture.thickness_concrete_slab }} ซม.</span>
						<span class="symbol-end-3">{{ store.dataInPicture.thickness_base }} ซม.</span>
						<span class="symbol-end-4">{{ store.dataInPicture.thickness_subbase }} ซม.</span>
						<span class="symbol-end-5">{{ store.dataInPicture.thickness_subgrade }} ซม.</span>
					</div>
				</div>
			</div>
			<div v-show="store.surfaceSectionCode === 5" key="surfaceSectionCode5" class="col-12 p-0">
				<div class="surfaceSectionCode-5">
					<div class="image-container">
						<div class="surface-symbol-top">
							<span class="symbol-1">{{ store.dataInPicture.width_shoulder_left }} ม.</span>
							<span class="symbol-2">{{ store.dataInPicture.width_surface }} ม.</span>
							<span class="symbol-3">{{ store.dataInPicture.width_shoulder_right }} ม.</span>
						</div>
						<img src="/images/roads/5_concrete_on_ground.svg" class="w-100" />
						<span class="symbol-end-1">{{ store.dataInPicture.thickness_surface }} ซม.</span>
						<span class="symbol-end-2">{{ store.dataInPicture.thickness_base }} ซม.</span>
						<span class="symbol-end-3">{{ store.dataInPicture.thickness_subbase }} ซม.</span>
						<span class="symbol-end-4">{{ store.dataInPicture.thickness_subgrade }} ซม.</span>
					</div>
				</div>
			</div>
		</div>
		<!-- end::รูปประกอบ -->
	</VSkeletonLoader>
</template>

<style scoped lang="scss">
// begin::surfaceSectionCode === 1,2
@mixin symbol-1-end-position-1() {
	top: 46%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 45%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 45%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 45%;
		right: -9%;
	}
}

@mixin symbol-1-end-position-2() {
	top: 61%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 60%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 60%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 60%;
		right: -9%;
	}
}

.image-container {
	position: relative;
	text-align: center;
	font-size: 0.9rem;
	max-width: 720px;
	@media only screen and (max-width: 475px) {
		font-size: 0.75rem;
	}
}

.surfaceSectionCode-1,
.surfaceSectionCode-2 {
	text-align: -webkit-center;
	.symbol-end-1 {
		@include symbol-1-end-position-1;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-2 {
		position: absolute;
		@include symbol-1-end-position-2;
		width: 75px;
		text-align: start;
	}

	.surface-symbol-top {
		width: 100%;
		display: flex;
		margin-top: 6%;
		position: absolute;
		.symbol-1 {
			width: 35%;
		}
		.symbol-2 {
			width: 18%;
		}
		.symbol-3 {
			width: 27%;
		}
	}
}
// end::surfaceSectionCode === 1,2

// begin::surfaceSectionCode === 3
@mixin symbol-3-end-position-1() {
	top: 43%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 43%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 43%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 43%;
		right: -9%;
	}
}

@mixin symbol-3-end-position-2() {
	top: 57%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 56%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 56%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 56%;
		right: -9%;
	}
}

@mixin symbol-3-end-position-3() {
	top: 71%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 69%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 69%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 69%;
		right: -9%;
	}
}

@mixin symbol-3-end-position-4() {
	top: 84%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 83%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 83%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 83%;
		right: -9%;
	}
}

.surfaceSectionCode-3 {
	text-align: -webkit-center;
	.symbol-end-1 {
		@include symbol-3-end-position-1;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-2 {
		position: absolute;
		@include symbol-3-end-position-2;
		width: 75px;
		text-align: start;
	}
	.symbol-end-3 {
		@include symbol-3-end-position-3;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-4 {
		position: absolute;
		@include symbol-3-end-position-4;
		width: 75px;
		text-align: start;
	}

	.surface-symbol-top {
		width: 100%;
		display: flex;
		margin-top: 6%;
		position: absolute;
		.symbol-1 {
			width: 35%;
		}
		.symbol-2 {
			width: 18%;
		}
		.symbol-3 {
			width: 27%;
		}
	}
}
// end::surfaceSectionCode === 3

// begin::surfaceSectionCode === 4
@mixin symbol-4-end-position-1() {
	top: 30%;
	right: -1%;
	@media only screen and (max-width: 900px) {
		top: 29%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 29%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 29%;
		right: -9%;
	}
}

@mixin symbol-4-end-position-2() {
	top: 43%;
	right: -1%;
	@media only screen and (max-width: 900px) {
		top: 42%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 42%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 42%;
		right: -9%;
	}
}

@mixin symbol-4-end-position-3() {
	top: 56%;
	right: -1%;
	@media only screen and (max-width: 900px) {
		top: 54%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 54%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 54%;
		right: -9%;
	}
}

@mixin symbol-4-end-position-4() {
	top: 69.5%;
	right: -1%;
	@media only screen and (max-width: 900px) {
		top: 67.5%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 67.5%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 67.5%;
		right: -9%;
	}
}

@mixin symbol-4-end-position-5() {
	top: 82%;
	right: -1%;
	@media only screen and (max-width: 900px) {
		top: 80%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 80%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 80%;
		right: -9%;
	}
}

.surfaceSectionCode-4 {
	text-align: -webkit-center;
	.symbol-end-1 {
		@include symbol-4-end-position-1;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-2 {
		position: absolute;
		@include symbol-4-end-position-2;
		width: 75px;
		text-align: start;
	}
	.symbol-end-3 {
		@include symbol-4-end-position-3;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-4 {
		position: absolute;
		@include symbol-4-end-position-4;
		width: 75px;
		text-align: start;
	}
	.symbol-end-5 {
		position: absolute;
		@include symbol-4-end-position-5;
		width: 75px;
		text-align: start;
	}

	.surface-symbol-top {
		width: 100%;
		display: flex;
		margin-top: 4.5%;
		position: absolute;
		.symbol-1 {
			width: 34.5%;
			padding-left: 13%;
		}
		.symbol-2 {
			width: 21%;
		}
		.symbol-3 {
			width: 18%;
		}
	}
}
// end::surfaceSectionCode === 4

// begin::surfaceSectionCode === 5
@mixin symbol-5-end-position-1() {
	top: 43%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 43%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 43%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 43%;
		right: -9%;
	}
}

@mixin symbol-5-end-position-2() {
	top: 57.5%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 56.5%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 56.5%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 56.5%;
		right: -9%;
	}
}

@mixin symbol-5-end-position-3() {
	top: 70.5%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 69.5%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 69.5%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 69.5%;
		right: -9%;
	}
}

@mixin symbol-5-end-position-4() {
	top: 84.5%;
	right: -2%;
	@media only screen and (max-width: 900px) {
		top: 83.5%;
		right: -2%;
	}
	@media only screen and (max-width: 600px) {
		top: 83.5%;
		right: -6%;
	}
	@media only screen and (max-width: 475px) {
		top: 83.5%;
		right: -9%;
	}
}

.surfaceSectionCode-5 {
	text-align: -webkit-center;
	.symbol-end-1 {
		@include symbol-5-end-position-1;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-2 {
		position: absolute;
		@include symbol-5-end-position-2;
		width: 75px;
		text-align: start;
	}
	.symbol-end-3 {
		@include symbol-5-end-position-3;
		position: absolute;
		width: 75px;
		text-align: start;
	}
	.symbol-end-4 {
		position: absolute;
		@include symbol-5-end-position-4;
		width: 75px;
		text-align: start;
	}

	.surface-symbol-top {
		width: 100%;
		display: flex;
		margin-top: 6%;
		position: absolute;
		.symbol-1 {
			width: 35%;
		}
		.symbol-2 {
			width: 18%;
		}
		.symbol-3 {
			width: 27%;
		}
	}
}
// end::surfaceSectionCode === 5

.surfaceSectionCode-container {
	min-height: 200px;
	@media only screen and (max-width: 600px) {
		min-height: 150px;
	}
}
</style>
