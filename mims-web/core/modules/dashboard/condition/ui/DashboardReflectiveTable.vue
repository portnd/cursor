<script setup lang="ts">
import { useDashboardReflectiveStore } from "../store"

const store = useDashboardReflectiveStore()

const expanded = ref(false)

const expandAll = () => {
	store.reflectivityDataTable.forEach((item) => {
		if (store.params.ref_reflectivity_range_id === 2) {
			item.expanded = true
		} else {
			item.expanded = true
			item.child.forEach((child) => {
				child.expanded = true
			})
		}
	})
	expanded.value = true
}

const collapseAll = () => {
	store.reflectivityDataTable.forEach((item) => {
		if (store.params.ref_reflectivity_range_id === 2) {
			item.expanded = false
		} else {
			item.expanded = false
			item.child.forEach((child) => {
				child.expanded = false
			})
		}
	})
	expanded.value = false
}

onMounted(() => {
	// onSearch()
})

const activeEl = ref()
</script>

<template>
	<div class="card card-rounded p-5 my-5">
		<div class="row align-items-center">
			<div class="col-sm-6 col-12 mb-5 mb-md-0">
				<h6 class="fw-normal fs-5 mb-0">ค่า Retro Reflectivity Average</h6>
			</div>
			<div
				class="col-sm-6 col-12 mb-5 mb-md-0 text-end"
				:style="{ visibility: store.params.ref_reflectivity_range_id === 2 ? 'hidden' : 'visible' }"
			>
				<button class="btn btn-outline btn-outline-primary rounded-4 btn-sm me-3" @click="expandAll()">
					Expand All
				</button>
				<button
					class="btn btn-outline btn-outline-primary rounded-4 btn-sm"
					:style="{ visibility: store.params.ref_reflectivity_range_id === 2 ? 'hidden' : 'visible' }"
					@click="collapseAll()"
				>
					Collapse All
				</button>
			</div>
		</div>
		<div v-if="store.reflectivityDataTable.length > 0" class="mt-4">
			<div class="iri-list">
				<div class="d-flex title-wrapper">
					<div
						class="title"
						:style="{ borderBottom: store.params.ref_reflectivity_range_id === 2 ? '1px solid #E7E9F1' : '' }"
					></div>
					<div
						class="title"
						:style="{ borderBottom: store.params.ref_reflectivity_range_id === 2 ? '1px solid #E7E9F1' : '' }"
					>
						กม.เริ่มต้น - กม.สิ้นสุด
					</div>
					<div
						class="title"
						:style="{ borderBottom: store.params.ref_reflectivity_range_id === 2 ? '1px solid #E7E9F1' : '' }"
					>
						Retro Reflectivity Average
					</div>
				</div>
				<div class="item-wrapper">
					<template v-for="(item, index) in store.reflectivityDataTable" :key="item.id">
						<div
							class="d-flex cursor-pointer"
							data-bs-toggle="collapse"
							:data-bs-target="`#collapse${item.id}`"
							:aria-expanded="item.expanded"
							:class="{
								collapsed: item.expanded && store.params.ref_reflectivity_range_id !== 2,
								expand: store.params.ref_reflectivity_range_id !== 2,
								collapse: store.params.ref_reflectivity_range_id === 2,
								hover: store.params.ref_reflectivity_range_id === 2,
								active: Number(item.id) === store.params.toggle_id,
							}"
							@click="store.toggleParentDataTable(item, index)"
						>
							<div
								class="item"
								:style="{
									backgroundColor: store.params.ref_reflectivity_range_id === 2 ? '' : '#f4f4f4',
									transform: item.expanded ? `rotate(180deg)` : ``,
								}"
							>
								<i
									class="fi fi-br-caret-down fs-4 lh-0 text-gray-600"
									:style="
										item.child?.length === 0 || store.params.ref_reflectivity_range_id === 2
											? 'visibility: hidden;'
											: ''
									"
								></i>
							</div>
							<div
								class="item"
								:style="{
									backgroundColor: store.params.ref_reflectivity_range_id === 2 ? '' : '#f4f4f4',
								}"
							>
								{{ item.km_start }} - {{ item.km_end }}
							</div>
							<div
								class="item"
								:style="{
									backgroundColor: store.params.ref_reflectivity_range_id === 2 ? '' : '#f4f4f4',
								}"
							>
								{{ item.value }}
							</div>
						</div>
						<template v-for="child of item.child" :key="child.id">
							<div :id="`collapse${item.id}`" class="collapse hover" :class="{ show: item.expanded }">
								<div class="d-flex" @click="store.toggleDataTable(child)">
									<div
										ref="activeEl"
										class="item"
										:style="{ backgroundColor: store.params.ref_reflectivity_range_id === 2 ? '#fff' : '' }"
										:class="{ active: child.id == store.params.toggle_id }"
										:data-id="child.id"
									>
										<span class="fs-5 text-gray-600">|</span>
									</div>
									<div
										class="item cursor-pointer"
										:class="{ active: child.id == store.params.toggle_id }"
										:style="{ backgroundColor: store.params.ref_reflectivity_range_id === 2 ? '#fff' : '' }"
										:data-id="child.id"
									>
										{{ child.km_start }} - {{ child.km_end }}
									</div>
									<div
										class="item cursor-pointer"
										:class="{ active: child.id == store.params.toggle_id }"
										:style="{ backgroundColor: store.params.ref_reflectivity_range_id === 2 ? '#fff' : '' }"
										:data-id="child.id"
									>
										{{ child.value }}
									</div>
								</div>
							</div>
						</template>
					</template>
				</div>
			</div>
		</div>
		<span v-else class="text-center m-10 p-3">ไม่พบข้อมูล</span>
	</div>
</template>

<style scoped lang="scss">
.active {
	background-color: #d9d9d9;
	color: #000;
	font-weight: 500;
}

.hover:hover {
	background-color: #f4f4f4f4 !important;
}

.iri-list {
	border: 1px solid var(--kt-gray-300);
	border-radius: 8px;
	width: 100%;

	.collapse {
		width: 100%;
	}

	.title {
		padding: 0.75rem 1.5rem;
		border-left: 1px solid var(--kt-gray-300);
		text-align: center;
		vertical-align: middle;
		font-size: 1.1rem;
		font-weight: 500;
	}

	.item {
		padding: 0.75rem 1.5rem;
		border-left: 1px solid var(--kt-gray-300);
		text-align: center;
		vertical-align: middle;
		font-size: 1rem;
	}

	.title:nth-child(1),
	.item:nth-child(1) {
		width: 16%;
	}

	.title:nth-child(2),
	.item:nth-child(2) {
		width: 42%;
	}

	.title:nth-child(3),
	.item:nth-child(3) {
		width: 42%;
	}

	.title:first-of-type,
	.item:first-of-type {
		border-left: none !important;
	}

	.item-wrapper {
		overflow-y: overlay;
		max-height: 385px;
	}

	.expand {
		background-color: #fff0d980 !important;
	}
}
</style>
