<script setup lang="ts">
import { useDashboardStore } from "../store"

const data = [
	{
		id: 0,
		km_start: "0+775",
		km_end: "0+800",
		value: 0,
		child: [
			{
				id: 0,
				km_start: "0+775",
				km_end: "0+800",
				value: 0,
				image_id: 0,
				image: "",
				geom_cl: "LINESTRING(100.55069783333073 13.74187357989189,100.55071697908836 13.741646696901423)",
			},
		],
	},
	{
		id: 1,
		km_start: "0+800",
		km_end: "0+995",
		value: 0,
		child: [
			{
				id: 0,
				km_start: "0+800",
				km_end: "0+995",
				value: 0,
				image_id: 0,
				image: "",
				geom_cl: "LINESTRING(100.55069783333073 13.74187357989189,100.55071697908836 13.741646696901423)",
			},
		],
	},
]
const expanded = ref(false)
const store = useDashboardStore()

const expandAll = () => {
	expanded.value = true
}

const collapseAll = () => {
	expanded.value = false
}

const conditionOptions = [
	{
		label: "IRI",
		value: 1,
	},
	{
		label: "MPD",
		value: 2,
	},
	{
		label: "RUT",
		value: 3,
	},
	{
		label: "IFI",
		value: 4,
	},
	{
		label: "แถบสะท้อนแสง",
		value: 5,
	},
]
</script>
<template>
	<div class="row align-items-center">
		<div class="col-3">
			<h6 class="fw-normal fs-5 mb-0">
				ค่า {{ conditionOptions.find((item) => item.value === store.conditionType)?.label }}
			</h6>
		</div>
		<div class="col-9 text-end">
			<button class="btn btn-outline btn-outline-primary rounded-4 btn-sm me-3" @click="expandAll()">Expand All</button>
			<button class="btn btn-outline btn-outline-primary rounded-4 btn-sm" @click="collapseAll()">Collapse All</button>
		</div>
	</div>
	<div v-if="data.length > 0" class="mt-4">
		<div class="iri-list">
			<div class="d-flex title-wrapper">
				<div class="title"></div>
				<div class="title">กม.เริ่มต้น - กม.สิ้นสุด</div>
				<div class="title">{{ conditionOptions.find((item) => item.value === store.conditionType)?.label }}</div>
			</div>
			<div class="item-wrapper">
				<template v-for="item in data" :key="item.id">
					<div
						class="d-flex cursor-pointer expand"
						data-bs-toggle="collapse"
						:data-bs-target="`#collapse${item.id}`"
						:aria-expanded="expanded"
					>
						<div class="item">
							<i class="fi fi-br-caret-down fs-4 lh-0 text-gray-600"></i>
						</div>
						<div class="item">{{ item.km_start }} - {{ item.km_end }}</div>
						<div class="item">{{ item.value }}</div>
					</div>
					<template v-for="child of item.child" :key="child.id">
						<div :id="`collapse${item.id}`" class="collapse hover" :class="{ show: expanded }">
							<div class="d-flex">
								<div ref="activeEl" class="item" :class="{ active: child.id == id }" :data-id="child.id">
									<span class="fs-5 text-gray-600">|</span>
								</div>
								<div class="item cursor-pointer" :class="{ active: child.id == id }" :data-id="child.id">
									{{ child.km_start }} - {{ child.km_end }}
								</div>
								<div class="item cursor-pointer" :class="{ active: child.id == id }" :data-id="child.id">
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
</template>

<style scoped lang="scss">
.active {
	background-color: #d9d9d9;
	color: #000;
	font-weight: 500;
}

.hover:hover {
	background-color: #e7e7e7;
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
} */

/* .btn.active {
	background-color: var(--kt-primary) !important;
	color: #fff !important;
	i {
		color: #fff !important;
	}
} */
</style>
