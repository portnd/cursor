<script setup lang="ts">
import { useForm } from "vee-validate"
import { useReportProjectMaintenanceReportStore } from "../../store"
import { IValidate } from "~/core/shared/types/Validate"

const submitNum = ref()
const store = useReportProjectMaintenanceReportStore()

onMounted(async () => {
	handleReset()
	await store.getReportProjectMaintenanceFilter()
})

const handleValidate = computed(() => {
	let validations: IValidate = {}
	if (submitNum.value > 0) {
		validations = {
			road_group_id: "required",
			road_section_id: "required",
			year_start: "required",
			year_end: "required",
		}
	}
	return validations
})

const { handleSubmit, handleReset, submitCount } = useForm({ validationSchema: handleValidate })

watch(
	() => submitCount.value,
	() => {
		submitNum.value = submitCount.value
	}
)

const onSubmit = handleSubmit((_, actions) => {
	useAction(actions)
	useDownloadFile(
		"รายงานประวัติการซ่อมบำรุง",
		`/report/type13?road_section_id=${store.data.road_section_id}&year_start=${store.data.year_start}&year_end=${store.data.year_end}&type=${store.data.type}`
	)
})

const onExport = (type: string) => {
	store.data.type = type
	onSubmit()
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<!-- Begin:: รายงานสินทรัพย์ -->
	<div class="row">
		<div class="col-12">
			<VSelect
				v-model="store.data.road_group_id"
				:options="store.getRoadGroupOptions"
				label="หมายเลขทางหลวง"
				name="road_group_id"
				placeholder="เลือก"
				@update:model-value="
					() => {
						store.data.road_section_id = null
					}
				"
			/>
		</div>
		<div class="col-12">
			<VSelect
				v-model="store.data.road_section_id"
				:options="store.getRoadSectionOptions"
				label="ตอนควบคุม"
				name="road_section_id"
				placeholder="เลือก"
			/>
		</div>

		<div class="col-12">
			<VSelect
				v-model="store.data.year_start"
				:options="store.getYearStartOptions"
				label="ปีเริ่มต้น"
				name="year_start"
				placeholder="เลือก"
			/>
		</div>

		<div class="col-12">
			<VSelect
				v-model="store.data.year_end"
				:options="store.getYearEndOptions"
				label="ปีสิ้นสุด"
				name="year_end"
				placeholder="เลือก"
			/>
		</div>
	</div>
	<div class="d-flex flex-xxl-row flex-column gap-3 mt-5">
		<button class="btn btn-code px-10 mb-xxl-0 mb-2" @click="onExport('html')">
			<div class="d-flex align-items-center justify-content-center">
				<i class="fi fi-ss-file-code fs-1 lh-0"></i>
				<span>HTML</span>
			</div>
		</button>
		<button class="btn btn-outline btn-outline-danger px-10 mb-xxl-0 mb-2" @click="onExport('pdf')">
			<div class="d-flex align-items-center justify-content-center">
				<i class="fi fi-ss-file-pdf fs-1 lh-0"></i>
				<span>PDF</span>
			</div>
		</button>
		<button class="btn btn-outline btn-outline-success px-10 mb-xxl-0 mb-2" @click="onExport('excel')">
			<div class="d-flex align-items-center justify-content-center">
				<i class="fi fi-ss-file-excel fs-1 lh-0"></i>
				<span>EXCEL</span>
			</div>
		</button>
	</div>
</template>

<style lang="scss" scoped>
.btn-code {
	border-radius: 16px;
	background-color: var(--kt-white);
	border: 1px solid #1f70f3 !important;
	color: #1f70f3;

	&:hover {
		color: var(--kt-white);
		background-color: #1f70f3;
	}
}
</style>
