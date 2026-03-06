<script setup lang="ts">
import { useConditionStore } from "../store"

const store = useConditionStore()

const expanded = ref(false)

const expandAll = () => {
	store.conditionDataTable.forEach((item) => (item.expanded = true))
	expanded.value = true
}

const collapseAll = () => {
	store.conditionDataTable.forEach((item) => (item.expanded = false))
	expanded.value = false
	store.image.expanded = expanded.value
}

watch(
	() => store.image.expanded,
	() => {
		expanded.value = store.image.expanded
		store.conditionDataTable.forEach((item) => (item.expanded = store.image.expanded))
	}
)

watch(
	() => store.image.play,
	(newPlay) => {
		if (newPlay) {
			expanded.value = true
			store.conditionDataTable.forEach((item) => (item.expanded = true))
		} else {
			expanded.value = false
			store.conditionDataTable.forEach((item) => (item.expanded = false))
		}
	}
)

const toggleType = (name: string) => {
	const typeName = name.toUpperCase()

	if (typeName === "IRI") {
		return "ค่า " + typeName + " (ม./กม.)"
	} else if (typeName === "RUT" || typeName === "MPD") {
		return "ค่า " + typeName + " (มม.)"
	} else {
		return "ค่า " + typeName
	}
}

const activeEl = ref()

watch(
	() => store.image.imageID,
	async (newId) => {
		if (newId) {
			store.conditionDataTable?.forEach((parent) => {
				parent.child?.forEach((child) => {
					if (child.id === newId) {
						const isGeomNotEmpty = !child.geom_cl?.toLowerCase().includes("empty")
						if (isGeomNotEmpty) {
							const geom = child.geom_cl?.split(",")[0]?.split("(")[1]?.split(" ")
							store.map?.location({
								lon: geom[0],
								lat: geom[1],
							})
						}
					}
				})
			})
		}

		await nextTick()
		const activeElement = document.querySelector(`.item.active[data-id="${newId}"]`)
		const container = document.querySelector(".item-wrapper")

		if (activeElement && container) {
			const containerRect = container.getBoundingClientRect()
			const activeElementRect = activeElement.getBoundingClientRect()
			container.scrollTop += activeElementRect.top - containerRect.top
		}
	}
)
</script>

<template>
	<div class="card card-rounded p-5 mt-5">
		<div class="row align-items-center">
			<div class="col-3">
				<h6 class="fw-normal fs-5 mb-0">ค่า {{ store.params.condition_type }}</h6>
			</div>
			<div class="col-9 text-end">
				<button
					:style="{ visibility: store.params.ref_condition_range_id === 3 ? 'hidden' : 'visible' }"
					class="btn btn-outline btn-outline-primary rounded-4 btn-sm me-3"
					@click="expandAll()"
				>
					Expand All
				</button>
				<button
					:style="{ visibility: store.params.ref_condition_range_id === 3 ? 'hidden' : 'visible' }"
					class="btn btn-outline btn-outline-primary rounded-4 btn-sm"
					@click="collapseAll()"
				>
					Collapse All
				</button>
			</div>
		</div>
		<div v-if="store.conditionDataTable.length > 0" class="mt-4">
			<div class="iri-list">
				<div class="d-flex title-wrapper">
					<div
						class="title"
						:style="{ borderBottom: store.params.ref_condition_range_id === 3 ? '1px solid #E7E9F1' : '' }"
					></div>
					<div
						class="title"
						:style="{ borderBottom: store.params.ref_condition_range_id === 3 ? '1px solid #E7E9F1' : '' }"
					>
						กม.เริ่มต้น - กม.สิ้นสุด
					</div>
					<div
						class="title"
						:style="{ borderBottom: store.params.ref_condition_range_id === 3 ? '1px solid #E7E9F1' : '' }"
					>
						{{ toggleType(store.params.condition_type) }}
					</div>
				</div>

				<div class="item-wrapper">
					<template v-for="(item, index) in store.conditionDataTable" :key="item.id">
						<div
							class="d-flex cursor-pointer"
							data-bs-toggle="collapse"
							:data-bs-target="`#collapse${item.id}`"
							:aria-expanded="item.expanded"
							:class="{
								collapsed: !store.image.play || (item.expanded && store.params.ref_condition_range_id !== 3),
								expand: store.params.ref_condition_range_id !== 3,
								collapse: store.params.ref_condition_range_id === 3,
								hover: store.params.ref_condition_range_id === 3,
								active: Number(item.id) == store.image.imageID,
							}"
							@click="store.toggleParentDataTable(item, index)"
						>
							<div
								class="item"
								:style="{
									backgroundColor: store.params.ref_condition_range_id === 3 ? '' : '#f4f4f4',
									transform: item.expanded ? `rotate(180deg)` : ``,
								}"
							>
								<i
									class="fi fi-br-caret-down fs-4 lh-0 text-gray-600"
									:style="item.child?.length === 0 ? 'visibility: hidden;' : ''"
								></i>
							</div>
							<div
								class="item"
								:style="{ backgroundColor: store.params.ref_condition_range_id === 3 ? '' : '#f4f4f4' }"
							>
								{{ item.km_start }} - {{ item.km_end }}
							</div>
							<div
								class="item"
								:style="{ backgroundColor: store.params.ref_condition_range_id === 3 ? '' : '#f4f4f4' }"
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
										:style="{ backgroundColor: store.params.ref_condition_range_id === 3 ? '#fff' : '' }"
										:class="{ active: child.id == store.image.imageID }"
										:data-id="child.id"
									>
										<span class="fs-5 text-gray-600">|</span>
									</div>
									<div
										class="item cursor-pointer"
										:style="{ backgroundColor: store.params.ref_condition_range_id === 3 ? '#fff' : '' }"
										:class="{ active: child.id == store.image.imageID }"
										:data-id="child.id"
									>
										{{ child.km_start }} - {{ child.km_end }}
									</div>
									<div
										class="item cursor-pointer"
										:style="{ backgroundColor: store.params.ref_condition_range_id === 3 ? '#fff' : '' }"
										:class="{ active: child.id == store.image.imageID }"
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

/* .btn {
	&:hover,
	&:focus {
		background-color: var(--kt-primary) !important;
		color: #fff !important;
		i {
			color: #fff !important;
		}
	}
}

.btn.active {
	background-color: var(--kt-primary) !important;
	color: #fff !important;
	i {
		color: #fff !important;
	}
} */
</style>
