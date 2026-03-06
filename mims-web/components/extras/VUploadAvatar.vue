<script setup lang="ts">
import { useField } from "vee-validate"
import { FilePond as FilePondFile } from "filepond"
import { IFile } from "~~/core/shared/types/File"

interface FileProperties {
	webkitRelativePath?: string
	lastModified?: number
	_relativePath: string
	lastModifiedDate: string
}

interface FilePondError {
	main: string
	sub: string
}

const props = defineProps({
	modelValue: {
		type: Object as PropType<IFile>,
		required: true,
	},
	file: {
		type: [String],
		default: null,
	},
	name: {
		type: String,
		required: true,
	},
	totalFileSize: {
		type: String,
		required: true,
	},
	imageSize: {
		type: Number,
		default: null,
	},
	label: {
		type: String,
		default: "",
	},
	required: {
		type: Boolean,
		default: false,
	},
	disabled: {
		type: Boolean,
		default: false,
	},
})

// Ref
const filePondRef = ref()

const result = ref<IFile>()

const emit = defineEmits(["input", "update:modelValue", "update:deleteImage"])

const handleFile = async (_err: FilePondError | null, fileItem: FilePondFile) => {
	const file = new File([fileItem.file], fileItem.file.name, { type: fileItem.file.type })
	const properties = fileItem.file as unknown as FileProperties
	const isUpload = !(
		typeof properties.webkitRelativePath === "undefined" && typeof fileItem.file.lastModified === "undefined"
	)

	if (_err === null) {
		result.value = {
			data: {
				id: fileItem.id,
				_relativePath: properties._relativePath,
				lastModified: fileItem.file.lastModified,
				lastModifiedDate: properties.lastModifiedDate,
				name: fileItem.file.name,
				size: fileItem.file.size,
				type: fileItem.file.type,
				webkitRelativePath: properties.webkitRelativePath,
				file: isUpload ? file : undefined,
				base64: (await fileToBase64(file)) as string,
			},
			isUpload,
			status: isUpload ? "upload" : "not_edit",
		}

		emit("update:modelValue", result.value)
	} else {
		showToast({
			title: "แจ้งเตือน",
			message: _err.sub,
			type: "primary",
		})

		result.value = {
			data: null,
			isUpload: false,
			status: "no_file",
		}
		emit("update:modelValue", result.value)
	}
}

const handleRemoveFile = () => {
	result.value = {
		data: null,
		isUpload: false,
		status: "delete",
	}
	emit("update:modelValue", result.value)
	emit("update:deleteImage", result.value)
}

// clearFile
const showFilePond = ref(true)
const clearFile = () => {
	const filePond = filePondRef.value

	if (!filePond) {
		return
	}

	showFilePond.value = false
	nextTick(() => {
		showFilePond.value = true

		result.value = {
			data: null,
			isUpload: false,
			status: "no_file",
		}
		emit("update:modelValue", result.value)
	})
}

// Files
const toFile = (file: String) => {
	const result: any[] = []
	if (file) {
		result.push({
			source: file,
			option: {
				type: "limbo",
			},
		})
	}
	return result
}

// Initialize
const initFiles: any = ref()
const computedInitFiles = computed(() => {
	if (props.file !== null) {
		initFiles.value = toFile(props.file)
	}
	return initFiles.value
})

// Validate
const { meta } = useField(props.name)

defineExpose({
	clearFile,
})
</script>
<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />
	<div class="py-3">
		<VField style="display: none" :name="name" :model-value="modelValue" />
		<file-pond
			ref="filePondRef"
			:accepted-file-types="['image/png', 'image/jpeg']"
			class="filepond-avatar"
			:class="{ 'is-invalid': !meta.valid && meta.validated, disabled: disabled }"
			:disabled="disabled"
			:name="name"
			label-button-process-item="Upload"
			:model-value="modelValue"
			:files="computedInitFiles"
			:max-total-file-size="totalFileSize"
			label-max-total-file-size-exceeded="ไฟล์มีขนาดใหญ่เกิน"
			:label-max-total-file-size="`ขนาดไฟล์ต้องไม่เกิน ${totalFileSize}`"
			:image-validate-size-max-width="imageSize ? imageSize : ''"
			:image-validate-size-max-height="imageSize ? imageSize : ''"
			image-validate-size-label-format-error="ไม่รองรับประเภทรูปภาพนี้"
			:image-validate-size-label-expected-max-size="`ขนาดรูปภาพต้องไม่เกิน ${imageSize}px × ${imageSize}px`"
			label-button-image-overlay="custom label"
			label-file-loading="กำลังโหลดไฟล์"
			label-file-processing="กำลังอัพโหลดไฟล์"
			label-file-processing-complete="อัพโหลดไฟล์สำเร็จ"
			label-file-processing-aborted="ยกเลิกการอัพโหลด"
			label-file-processing-error="อัพโหลดไม่สำเร็จ"
			label-tap-to-cancel="กดเพื่อยกเลิก"
			label-tap-to-retry="กดเพื่อลองอีกครั้ง"
			label-tap-to-undo="กดเพื่อเลิกทำ"
			label-idle="ลากและวางรูปภาพของคุณ <br/>หรือกด <span class='filepond--label-action'>อัปโหลด</span>"
			:image-preview-height="170"
			image-crop-aspect-ratio="1:1"
			:image-resize-target-width="150"
			:image-resize-target-height="150"
			style-panel-layout="compact circle"
			style-load-indicator-position="center bottom"
			style-progress-indicator-position="right bottom"
			style-button-remove-item-position="left bottom"
			style-button-process-item-position="right bottom"
			:allow-file-encode="true"
			@addfile="handleFile"
			@removefile="handleRemoveFile"
			@update:delete-value="($event: any) => (initFiles = $event)"
		/>
		<div class="text-center text-gray-600 mt-1">
			รองรับไฟล์ png, jpg, jpeg <span v-show="imageSize">ขนาดรูปภาพไม่เกิน {{ imageSize }} x {{ imageSize }} px</span
			><br />
			และมีขนาดไฟล์ไม่เกิน {{ totalFileSize }}
			<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
		</div>
	</div>
</template>

<style lang="scss">
.filepond-avatar {
	.filepond--credits {
		display: none;
	}
	.filepond--drop-label {
		padding: 2rem;
		color: #4c4e53;
		height: 85%;
		z-index: 10;
		transform: unset;
	}
	.filepond--label-action {
		text-decoration-color: #babdc0;
	}
	.filepond--panel-root {
		background-color: #edf0f4;
		border-radius: 100%;
	}
	.filepond--root {
		max-width: 175px;
		margin: 0 auto;
		border: 1px solid #ddd;
		height: 170px !important;
		border-radius: 100%;
	}
	.filepond--image-preview {
		background-color: #fff;
		margin: 0 auto;
		max-width: 175px;
		height: 175px;
		border-radius: 100%;
	}
	.filepond--list.filepond--list {
		position: unset;
	}
	.filepond--file-info,
	.filepond--file-status {
		display: none !important;
	}
	.filepond--item {
		margin: 0 auto;
		max-width: 175px;
		height: 175px !important;
		border-radius: 100%;
		margin-top: -1px !important;
	}
	.filepond--image-canvas-wrapper {
		max-width: 175px !important;
		height: 175px !important;
		border-radius: 100%;
		padding: 0 !important;
	}
	.filepond--list-scroller {
		margin: 0 auto;
		overflow: hidden;
		height: 175px !important;
		border-radius: 100%;
	}
	.filepond--image-preview-overlay {
		max-width: 175px;
		height: 175px;
		border-radius: 100%;
	}
	.filepond--image-canvas-wrapper {
		transform: unset !important;

		.filepond--image-bitmap {
			canvas {
				width: 175px !important;
				height: 175px !important;
				padding: 0px !important;
			}
		}
	}
	.filepond--image-clip {
		width: 175px !important;
		height: 175px !important;
	}

	.filepond--list-scroller[data-state="overflow"] {
		mask: unset !important;
	}

	// Preview
	.filepond--fullsize-overlay {
		z-index: 1055;
		background-color: rgb(0 0 0 / 80%);
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		overflow: hidden;
		outline: 0;
		display: flex;
	}
	.filepond--fullsize-overlay .image-container {
		height: 50%;
		margin: auto;
	}

	// Disabled
	.filepond-avatar.disabled .filepond--root {
		background-color: var(--kt-input-disabled-bg);
		cursor: no-drop;

		.filepond--drop-label.filepond--drop-label label {
			color: var(--kt-text-gray-400);
		}

		.filepond--label-action {
			cursor: no-drop;
		}
	}

	// Validate
	.filepond-avatar.is-invalid .filepond--root {
		border: 1px var(--kt-danger) solid;
	}

	// Drop label
	.filepond--drop-label:has(+ .filepond--list-scroller .filepond--list .filepond--item) {
		display: none !important;
	}
}
</style>
