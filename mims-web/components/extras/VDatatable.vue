<script setup lang="ts">
import { v4 as uuidv4 } from "uuid"
import { BodyRowClassNameFunction } from "vue3-easy-data-table"
import type { THeader, TItem } from "~~/core/shared/types/Datatable"

const props = defineProps({
	headers: {
		type: Array<THeader>,
		required: true,
	},
	items: {
		type: Array<TItem>,
		default: [],
	},
	noBorder: {
		type: Boolean,
		default: true,
	},
	hideRowsPerPage: {
		type: Boolean,
		default: false,
	},
	headTextDirection: {
		type: String,
		default: "",
	},
	bodyTextDirection: {
		type: String,
		default: "center",
	},
	rowsPerPage: {
		type: Number,
		default: 25,
	},
	xAxis: {
		type: Number,
		default: 0,
	},
	currentPage: {
		type: Number,
		default: 1,
	},
	itemsSelected: {
		type: Boolean,
		default: false,
	},
	selected: {
		type: Array<number>, // ให้ใส่ id ของ row นั้นมา
		default: [],
	},
	activeItemClassName: {
		type: String,
		default: "",
	},
	activeRowNumber: {
		type: Number,
	},
	whiteSpace: {
		type: Boolean,
		default: true,
	},
	loading: {
		type: Boolean,
		default: false,
	},
})

const dataTable = ref()
const uuid = ref("uuid-" + uuidv4())

// เพิ่มเติม Class
const className = ref(`customize-table ${uuid.value} ` + (props.noBorder ? "border-none" : ""))
className.value = props.whiteSpace === true ? className.value : className.value + ` no-white-space`

// แสดงจำนวนหน้า, จำนวนข้อมูล
const currentPageFirstIndex = computed(() => dataTable.value?.currentPageFirstIndex)
const currentPageLastIndex = computed(() => dataTable.value?.currentPageLastIndex)
const clientItemsLength = computed(() => dataTable.value?.clientItemsLength)
const currentPaginationNumber = computed(() => dataTable.value?.currentPaginationNumber)

// เพิ่ม Key เข้าไปใน Items
const rows = computed(() =>
	props.items.map((item, key) => {
		return { ...item, no: toNumber(key + 1) }
	})
)

const emit = defineEmits(["scrollbarX", "currentPage", "selected"])

// Scrollbar X
onMounted(() => {
	const tableWrapper = document.querySelector(`.${uuid.value}`)
	const scrollbar = tableWrapper?.querySelector(".vue3-easy-data-table__main")

	// xAxis
	if (scrollbar) {
		scrollbar.scrollLeft = props.xAxis
	}

	scrollbar?.addEventListener("scroll", () => {
		const xAxis = scrollbar.scrollLeft
		emit("scrollbarX", xAxis)
	})

	watch(
		() => props.xAxis,
		(newValue) => {
			// xAxis
			if (scrollbar) {
				scrollbar.scrollLeft = newValue
			}
		}
	)

	// ข้อมูลหน้าปัจจุบัน
	watch(
		() => props.currentPage,
		(_) => {
			if (tableWrapper) {
				const paginationbuttons = tableWrapper.querySelectorAll(".item.button")
				const filteredPaginationbuttons = Array.from(paginationbuttons).filter(
					(item) => item.textContent === String(props.currentPage)
				)
				filteredPaginationbuttons.forEach((paginationbutton: Element) => {
					;(paginationbutton as HTMLElement).click()
				})
			}
		}
	)
})

// ข้อมูลหน้าปัจจุบัน
const handlePageChange = () => {
	emit("currentPage", currentPaginationNumber.value)
}

// Items selected
const setItemsSelect = () => {
	if (props.itemsSelected) {
		const selecteds: TItem[] = []
		props.selected.forEach((id: number) => {
			selecteds.push(rows.value.find((obj: any) => obj.id === id) as TItem)
		})
		selectedIds.value = selecteds
	}
}

const selectedIds: Ref<TItem[] | null> = ref(null)
setItemsSelect()

watch(
	() => props.selected,
	() => {
		setItemsSelect()
	}
)

watch(
	() => props.items,
	() => {
		if (props.itemsSelected) {
			selectedIds.value = []
		}
	}
)

const updateValue = () => {
	emit("selected", selectedIds.value)
}

const bodyRowClassName: BodyRowClassNameFunction = (_, rowNumber: number): string => {
	if (typeof props.activeRowNumber !== "undefined") {
		if (rowNumber === Number(props.activeRowNumber)) {
			return props.activeItemClassName
		}
	}
	return ""
}

updateValue()
</script>

<template>
	<div class="customize-index text-gray-700 fs-7 mb-2">
		<i>
			จำนวนข้อมูลทั้งหมด {{ toNumber(clientItemsLength) }} รายการ แสดงรายการที่
			{{ clientItemsLength === 0 ? 0 : toNumber(currentPageFirstIndex) }} ถึง
			{{ toNumber(currentPageLastIndex) }}
		</i>
	</div>
	<VueDatatable
		ref="dataTable"
		v-model:items-selected="selectedIds"
		:headers="headers"
		:items="rows"
		border-cell
		buttons-pagination
		theme-color="#FDB833"
		:table-class-name="className"
		:header-text-direction="headTextDirection"
		:body-text-direction="bodyTextDirection"
		:loading="loading"
		rows-per-page-message="จำนวนรายการ"
		rows-of-page-separator-message="จาก"
		empty-message="ไม่พบข้อมูล"
		:hide-rows-per-page="hideRowsPerPage"
		:rows-per-page="rowsPerPage"
		:rows-items="[25, 50, 100]"
		:current-page="currentPage"
		:header-class-name="itemsSelected ? 'customize-items-selected' : ''"
		:body-row-class-name="bodyRowClassName"
		@update-page-items="handlePageChange"
		@update:items-selected="updateValue"
	>
		<template v-if="$slots['customize-headers']" #customize-headers>
			<slot name="customize-headers" />
		</template>

		<!-- begin::Slot for dynamic table items -->
		<template v-for="(header, key) in headers" #[`header-${header.value}`]="item" :key="key">
			<div
				class="w-100"
				:class="[
					!item.align ? 'text-center' : `text-${item.align}`,
					{ 'fixed-start': item.fixed, 'fixed-end': item.fixedEnd },
				]"
			>
				{{ item.text }}
			</div>
		</template>

		<template v-for="(header, key) in headers" #[`item-${header.value}`]="item" :key="key">
			<div :class="{ 'fixed-start': header.fixed, 'fixed-end': header.fixedEnd }">
				<slot :name="`item-${header.value}`" :item="item" />
			</div>
		</template>
		<!-- end::Slot for dynamic table items -->
	</VueDatatable>
</template>

<style scoped></style>
