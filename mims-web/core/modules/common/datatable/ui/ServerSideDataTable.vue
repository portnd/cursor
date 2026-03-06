<script setup lang="ts">
import { DatatableService } from "../infrastructure"
import type { THeader, TItem, TServerOptions, TClickRowArgument } from "~~/core/shared/types/Datatable"

interface ISeach {
	[key: string]: any
}

const props = defineProps({
	headers: {
		type: Array<THeader>,
		required: true,
	},
	noBorder: {
		type: Boolean,
		default: true,
	},
	headTextDirection: {
		type: String,
		default: "center",
	},
	bodyTextDirection: {
		type: String,
		default: "center",
	},
	rowsPerPage: {
		type: Number,
		default: 10,
	},
	url: {
		type: String,
		required: true,
	},
	additionalParams: {
		type: String as PropType<ISeach>,
		default: "",
	},
	itemsSelected: {
		type: Boolean,
		default: false,
	},
	selected: {
		type: Array<number>, // ให้ใส่ id ของ row นั้นมา
		default: [],
	},
	isInit: {
		type: Boolean,
		default: true,
	},
	allIds: {
		type: Array<number>, // ให้ใส่ id ทั้งหมด
		default: [],
	},
	loading: {
		type: Boolean,
		default: false,
	},
	loadingTop: {
		type: Boolean,
		default: false,
	},
})

// เพิ่มเติม Class
const className: string =
	"customize-table serverside-datatable " +
	(props.noBorder ? "border-none" : "") +
	" " +
	(props.loadingTop ? "loading-top" : "")

const dataTable = ref()

// แสดงจำนวนหน้า, จำนวนข้อมูล
const currentPageFirstIndex = computed(() => dataTable.value?.currentPageFirstIndex)
const currentPageLastIndex = computed(() => dataTable.value?.currentPageLastIndex)
const clientItemsLength = computed(() => dataTable.value?.clientItemsLength)

const serverItemsLength = ref(0)
const ALLOWED_ROWS_PER_PAGE = [10, 25, 50, 100]
const serverOptions = ref<TServerOptions>({
	page: 1,
	rowsPerPage: props.rowsPerPage,
})

// ป้องกัน vue3-easy-data-table overwrite rowsPerPage เป็นค่าไม่ถูกต้อง (เช่น 15) ทำให้ API ได้ limit ผิด
watch(
	() => serverOptions.value.rowsPerPage,
	(val) => {
		if (val != null && !ALLOWED_ROWS_PER_PAGE.includes(val)) {
			serverOptions.value.rowsPerPage = props.rowsPerPage
		}
	},
	{ immediate: true }
)
const loadingStatus = ref(props.loading)
const rows = ref<TItem[]>([])
const search = ref<ISeach[]>([])

const computedLoading = computed(() => {
	return loadingStatus.value
})

watch(
	() => props.loading,
	() => {
		loadingStatus.value = props.loading
		rows.value = []
	}
)

// Server side
const loadData = async () => {
	loadingStatus.value = true

	const { page, rowsPerPage } = serverOptions.value

	// เชื่อมต่อ Api
	const datatableService = new DatatableService()

	// รวม parameter
	const params = {
		page,
		limit: rowsPerPage,
	}
	let mergedParams = { ...params, ...search.value }

	// กรณีมี Params เพิ่มเติม
	if (Object.keys(props.additionalParams).length !== 0) {
		mergedParams = { ...mergedParams, ...props.additionalParams }
	}

	// เชื่อมต่อ Api
	const res = await datatableService.get(props.url, mergedParams)
	// หน้าปัจจุบัน
	serverOptions.value.page = res.data?.current_page

	const data = res.data?.items
	serverItemsLength.value = res.data?.total_items || 0
	if (data) {
		const index = (page - 1) * rowsPerPage + 1

		rows.value = data.map((item, key) => {
			return { ...item, no: toNumber(index + key) }
		})
	} else {
		rows.value = []
	}
	loadingStatus.value = false
	isSearch.value = false
	emit("getData", rows.value)
}

// กรณีใช้ค้นหาข้อมูล
const isSearch: Ref<Boolean> = ref(false)
const searchData = (params: ISeach[]) => {
	isSearch.value = true
	serverOptions.value.page = 1
	search.value = params
	loadData()
}

// initial load
if (props.isInit) {
	loadData()
}

watch(
	serverOptions,
	() => {
		// กรณี ใช้คำสั่งค้นหาผ่าน expose
		if (!isSearch.value) {
			loadData()
		}
	},
	{ deep: true }
)

// Items selected
const selectedIds: Ref<TItem[] | null> = ref(null)
const resultIds: Ref<Number[]> = ref([])
const container = ref()

const emit = defineEmits(["selected", "update:itemsSelected", "getData"])

const updateCheckboxAll = () => {
	// checkbox เลือกทั้งหมด
	const header = container.value.querySelector(".vue3-easy-data-table__header")
	if (header) {
		// allSelected เลือกทั้งหมด, noneSelected ไม่ได้เลือก, partSelected เลือกบางอัน
		header.classList.remove("noneSelected", "allSelected", "partSelected")

		setTimeout(() => {
			if (resultIds.value.length === 0) {
				header.classList.add("noneSelected")
			} else if (serverItemsLength.value === resultIds.value.length) {
				header.classList.add("allSelected")
			} else {
				header.classList.add("partSelected")
			}
		}, 200)
	}
}

const setItemsSelect = () => {
	if (props.itemsSelected) {
		// Reset ค่า
		selectedIds.value = []
		rows.value.forEach((row) => {
			// Update is checked
			if (selectedIds.value && resultIds.value.includes(row.id)) {
				selectedIds.value.push(row)
			}
		})

		updateCheckboxAll()
	}
}

watch(
	() => rows.value,
	() => {
		setItemsSelect()
	}
)

// initial
onMounted(() => {
	setDefaultSelected()
})

watch(
	() => props.selected,
	() => {
		setDefaultSelected()
	}
)

const setDefaultSelected = () => {
	if (props.itemsSelected) {
		resultIds.value = props.selected
		setItemsSelect()

		emit("selected", resultIds.value)
	}
}

const updateValue = (items: TItem) => {
	const partSelected = container.value?.querySelector(".vue3-easy-data-table__header.partSelected")

	// Uncheck all items
	if (items.length === 0 && !partSelected) {
		selectedIds.value = []
		resultIds.value = []

		emit("selected", [])
	}

	updateCheckboxAll()
}

const selectRow = (item: TClickRowArgument) => {
	pushOrReplace(resultIds.value, item.id)
	emit("selected", resultIds.value)
}

const deselectRow = (item: TClickRowArgument) => {
	if (selectedIds.value) {
		selectedIds.value = selectedIds.value.filter((selectedId) => selectedId.id !== item.id)
	}
	resultIds.value = resultIds.value.filter((id) => id !== item.id)

	emit("selected", resultIds.value)
}

const selectAll = () => {
	selectedIds.value = rows.value
	resultIds.value = props.allIds

	emit("selected", props.allIds)
}

const resetSelected = () => {
	selectedIds.value = []
	resultIds.value = []

	emit("selected", resultIds.value)
}

const resetRow = () => {
	rows.value = []
}

const pushOrReplace = (array: Number[], value: Number) => {
	const index = array.indexOf(value)

	if (index === -1) {
		array.push(value)
	} else {
		array.splice(index, 1, value)
	}
}

const getData = () => {
	return rows.value
}
defineExpose({
	getData,
	loadData,
	searchData,
	resetSelected,
	resetRow,
})
</script>

<template>
	<div ref="container" class="customize-index text-gray-700 fs-7 mb-2">
		<i>
			จำนวนข้อมูลทั้งหมด {{ toNumber(clientItemsLength) }} รายการ แสดงรายการที่
			{{ toNumber(clientItemsLength) === "0" ? 0 : toNumber(currentPageFirstIndex) }} ถึง
			{{ toNumber(currentPageLastIndex) }}
		</i>
		<VueDatatable
			ref="dataTable"
			v-model:items-selected="selectedIds"
			v-model:server-options="serverOptions"
			:server-items-length="serverItemsLength"
			:loading="computedLoading"
			:headers="headers"
			:items="rows"
			border-cell
			buttons-pagination
			theme-color="#FDB833"
			:table-class-name="className"
			:header-text-direction="headTextDirection"
			:header-class-name="itemsSelected ? 'customize-items-selected' : ''"
			:body-text-direction="bodyTextDirection"
			rows-per-page-message="จำนวนรายการ"
			rows-of-page-separator-message="จาก"
			empty-message="ไม่พบข้อมูล"
			:rows-per-page="rowsPerPage"
			:rows-items="[10, 25, 50, 100]"
			@select-row="selectRow"
			@deselect-row="deselectRow"
			@select-all="selectAll"
			@update:items-selected="updateValue"
		>
			<template v-if="$slots['customize-headers']" #customize-headers>
				<slot name="customize-headers" />
			</template>

			<!-- begin::Slot for dynamic table items -->
			<template v-for="(header, key) in headers" #[`header-${header.value}`]="item" :key="key">
				<div v-if="!item.align" class="text-center w-100">
					{{ item.text }}
				</div>
				<div v-else :class="`text-${item.align}`" class="w-100">
					{{ item.text }}
				</div>
			</template>

			<template v-for="(header, key) in headers" #[`item-${header.value}`]="item" :key="key">
				<slot :name="`item-${header.value}`" :item="item" />
			</template>
			<!-- end::Slot for dynamic table items -->
		</VueDatatable>
	</div>
</template>

<style>
/* ของเก่า */
.easy-checkbox input[type="checkbox"].allSelected + label:after {
	border: 0px !important;
}
.easy-checkbox input[type="checkbox"].allSelected + label:before,
.easy-checkbox input[type="checkbox"].partSelected + label:before {
	background: #fff;
}

/* ของใหม่ */
.partSelected .easy-checkbox input[type="checkbox"] + label:after {
	transform: translate(0.2em, 0.5875em) !important;
	width: 0.75em !important;
	height: 0.375em !important;
	border: 0.125em solid #fff !important;
	border-bottom-style: none !important;
	border-right-style: none !important;
	border-left-style: none !important;
	margin: 2px !important;
}
.partSelected .easy-checkbox input[type="checkbox"] + label:before,
.allSelected .easy-checkbox input[type="checkbox"] + label:before {
	background: var(--kt-primary) !important;
}
.noneSelected .easy-checkbox input[type="checkbox"] + label:before {
	background: #fff !important;
}
.allSelected .easy-checkbox input[type="checkbox"] + label:after {
	transform: translate(0.2em, 0.3038461538em) rotate(-45deg) !important;
	width: 0.75em !important;
	height: 0.375em !important;
	border: 0.13em solid #fff !important;
	border-top-style: none !important;
	border-right-style: none !important;
	margin: 2px !important;
}

.loading-top .loading-entity {
	position: absolute !important;
	top: 18px !important;
}
</style>
