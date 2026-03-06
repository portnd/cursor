<script setup lang="ts">
import { IFile } from "~~/core/shared/types/File"

const props = defineProps({
	modelValue: {
		type: String,
		request: true,
	},
	assetType: {
		type: String,
		defualt: "in",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	colors: {
		type: String,
		default: "#fdb833",
	},
	files: {
		type: String,
		default: null,
	},
	images: {
		type: Object as PropType<IFile>,
	},
})
const color = computed(() => {
	return props.colors
})

const image = ref<IFile>(
	props.images
		? props.images
		: {
				data: null,
				status: "no_file",
				isUpload: false,
		  }
)
const uploadFileKm: Ref = ref()
const uploadFileLatLon: Ref = ref()
const geomType = ref()

const emit = defineEmits(["update:modelValue", "update:lineColorCode", "update:image", "update:deleteImageDefault"])
const onSelect = (event: string) => {
	if (event === "km") {
		// uploadFileLatLon.value.clearFile()
		emit("update:lineColorCode", "")
	} else if (event === "km_range") {
		// uploadFileKm.value.clearFile()
		// uploadFileLatLon.value.clearFile()
		emit("update:lineColorCode", "#fdb833")
	} else {
		// uploadFileKm.value.clearFile()
		emit("update:lineColorCode", "")
	}
	emit("update:modelValue", event)
}

const updateColor = (color: string) => {
	emit("update:lineColorCode", color)
}

// watch(color, (newValue) => {
// 	emit("update:lineColorCode", newValue)
// })

watch(image, (newValue) => {
	emit("update:image", newValue)
})

onMounted(() => {
	geomType.value = props.modelValue
})

onUpdated(() => {
	geomType.value = props.modelValue
})
</script>

<template>
	<VLabel label="ประเภทพิกัด" :required="true" />
	<div class="nav nav-tabs geom-type-tabs col-12" role="tablist" :class="{ disabled: disabled }">
		<!-- <button
			:class="[
				assetType !== 'in' ? 'disabled' : '',
				disabled ? 'disabled' : '',
				modelValue === 'km' && 'active',
				modelValue === '' && 'active',
			]"
			class="col-lg col-md col-4 btn fs-6 px-0 fw-semibold rounded-0"
			data-bs-toggle="tab"
			data-bs-target="#geom-type1"
			type="button"
			role="tab"
			@click="onSelect('km')"
		>
			กม.
			<span class="line"></span>
		</button>
		<button
			:class="[assetType !== 'in' ? 'disabled' : '', disabled ? 'disabled' : '', modelValue === 'km_range' && 'active']"
			class="col-lg col-md col-4 btn fs-6 px-0 fw-semibold rounded-0"
			data-bs-toggle="tab"
			data-bs-target="#geom-type2"
			type="button"
			role="tab"
			@click="onSelect('km_range')"
		>
			ช่วงกม.
			<span class="line"></span>
		</button> -->
		<!-- <button
			:class="[assetType === 'in' ? 'disabled' : '', disabled ? 'disabled' : '', modelValue === 'point' && 'active']"
			class="col-lg col-md col-4 btn fs-6 px-0 fw-semibold rounded-0"
			data-bs-toggle="tab"
			data-bs-target="#geom-type3"
			type="button"
			role="tab"
			@click="onSelect('point')"
		>
			LAT, LON
			<span class="line"></span>
		</button> -->
		<div class="form-check form-check-custom mb-2 me-8 form-check-inline mt-1">
			<input
				id="geom-type1"
				v-model="geomType"
				class="form-check-input cursor-pointer"
				type="radio"
				name="geom-type"
				value="km"
				:disabled="disabled ? true : fasle"
				@click="onSelect('km')"
			/>
			<label for="geom-type1" class="form-check-label ms-3 cursor-pointer">กม.</label>
		</div>
		<div class="form-check form-check-custom mb-2 me-6 form-check-inline mt-1">
			<input
				id="geom-type2"
				v-model="geomType"
				class="form-check-input cursor-pointer"
				type="radio"
				name="geom-type"
				value="km_range"
				:disabled="disabled ? true : fasle"
				@click="onSelect('km_range')"
			/>
			<label for="geom-type2" class="form-check-label ms-3 cursor-pointer">ช่วงกม.</label>
		</div>
		<div class="form-check form-check-custom mb-2 me-6 form-check-inline mt-1">
			<input
				id="geom-type3"
				v-model="geomType"
				class="form-check-input cursor-pointer"
				type="radio"
				name="geom-type"
				value="point"
				:disabled="disabled ? true : fasle"
				@click="onSelect('point')"
			/>
			<label for="geom-type3" class="form-check-label ms-3 cursor-pointer">LAT, LON</label>
		</div>
	</div>
	<div class="col-12">
		<div
			v-if="modelValue === 'km' || modelValue === ''"
			id="geom-type1"
			:class="[disabled ? 'disabled' : '', modelValue === 'km' && 'active show', modelValue === '' && 'active show']"
			class="tab-pane geom-type-tab fade pt-2"
			role="tabpanel"
		>
			<VUploadFile
				ref="uploadFileKm"
				v-model="image"
				label="รูปสัญลักษณ์"
				:files="files"
				total-file-size="1MB"
				:image-size="300"
				name="file"
				:multiple="false"
				aspect-ratio="0.225"
				:accepted-file-types="['image/png', 'image/jpg', 'image/jpeg']"
				@update:delete-value="($file:any) => {emit('update:deleteImageDefault', $file)}"
			/>
		</div>

		<div
			v-if="modelValue === 'km_range'"
			id="geom-type2"
			:class="modelValue === 'km_range' && 'active show'"
			class="tab-pane geom-type-tab fade"
			role="tabpanel"
		>
			<VColorPicker
				v-model="color"
				name="color"
				label="เลือกสี"
				:required="true"
				@update:model-value="(e:string)=>updateColor(e)"
			/>
		</div>

		<div
			v-if="modelValue === 'point'"
			id="geom-type3"
			:class="[disabled ? 'disabled' : '', modelValue === 'point' && 'active show']"
			class="tab-pane geom-type-tab fade pt-2"
			role="tabpanel"
		>
			<VUploadFile
				ref="uploadFileLatLon"
				v-model="image"
				label="รูปสัญลักษณ์"
				:files="files"
				total-file-size="1MB"
				:image-size="300"
				name="file"
				:multiple="false"
				aspect-ratio="0.225"
				:accepted-file-types="['image/png', 'image/jpg', 'image/jpeg']"
				@update:delete-value="($file:any) => {emit('update:deleteImageDefault', $file)}"
			/>
		</div>
	</div>
</template>

<style scoped>
.geom-type-tabs {
	border-bottom: 0px solid #f4f4f4 !important;
}
.geom-type-tabs > button {
	display: grid;
}
.geom-type-tabs .active > .line {
	border: 2px solid var(--kt-primary) !important;
	border-radius: 10px;
	margin-top: 11px;
}
@media only screen and (max-width: 767px) {
	.geom-type-tabs .active > .line {
		margin-top: 13px;
	}
}
.geom-type-tab:not(.show) {
	display: none;
}

/* disabled */
.geom-type-tabs.disabled {
	cursor: no-drop;
}
.geom-type-tabs.disabled .active > .line {
	border: 2px solid #ddd !important;
}
</style>
