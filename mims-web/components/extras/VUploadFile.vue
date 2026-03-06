<script lang="ts" setup>
import { useField } from "vee-validate"
import { FilePond as IFilePond, FilePondFile, FilePondErrorDescription } from "filepond"
import { IFile } from "~~/core/shared/types/File"

interface FileProperties {
	webkitRelativePath?: string
	lastModified?: number
	_relativePath: string
	lastModifiedDate: string
}

const props = defineProps({
	modelValue: {
		type: Object as PropType<IFile | Array<IFile>>,
		required: true,
	},
	files: {
		type: [String, Array<String>],
		default: null,
	},
	name: {
		type: String,
		required: true,
	},
	totalFileSize: {
		type: String,
	},
	maxFileSize: {
		type: String,
	},
	multiple: {
		type: Boolean,
		required: false,
	},
	imageSize: {
		type: Number,
		default: null,
	},
	acceptedFileTypes: {
		type: Array,
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
	aspectRatio: {
		type: String,
		default: "0.025",
	},
	allowFileTypeValidation: {
		type: Boolean,
		default: true,
	},
	fileValidateTypeDetectType: {
		type: Function,
	},
	maxFiles: {
		type: Number,
		default: null,
	},
})

// Ref
const filePondRef = ref<IFilePond | null>(null)
const multiples = ref<Array<IFile>>([])
const single = ref<IFile>()

const emit = defineEmits(["input", "update:modelValue", "update:deleteImage"])

const handleFile = async (_err: FilePondErrorDescription | null, fileItem: FilePondFile) => {
	switch (fileItem.fileType) {
		case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		case "application/vnd.ms-excel":
			fileItem.setMetadata("poster", "/images/files/xlsx.png")
			break
		case "application/msword":
		case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
			fileItem.setMetadata("poster", "/images/files/docx.png")
			break
		case "application/vnd.ms-powerpoint":
		case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
			fileItem.setMetadata("poster", "/images/files/ppt.png")
			break
		case "text/plain":
			fileItem.setMetadata("poster", "/images/files/txt.png")
			break
		case "text/csv":
			fileItem.setMetadata("poster", "/images/files/csv.png")
			break
		case "application/zip":
			fileItem.setMetadata("poster", "/images/files/zip.png")
			break
		case "application/x-rar":
			fileItem.setMetadata("poster", "/images/files/rar.png")
			break
		case "application/x-zip-compressed":
			fileItem.setMetadata("poster", "/images/files/rar.png")
			break
		case "application/pdf":
			// fileItem.setMetadata("poster", "/images/files/pdf.png")
			break
		default:
			// กรณี filepond ไม่สามารถดึง mimetype มาได้
			const ext = fileItem.file.name.split(".").pop()
			if (ext === "dwg") {
				fileItem.setMetadata("poster", "/images/files/dwg.png")
			} else {
				fileItem.setMetadata("poster", "/images/files/unknown.png")
			}
			break
	}

	const file = new File([fileItem.file], fileItem.file.name, { type: fileItem.file.type })
	const properties = fileItem.file as unknown as FileProperties
	const isUpload = !(
		typeof properties.webkitRelativePath === "undefined" && typeof fileItem.file.lastModified === "undefined"
	)

	if (_err === null) {
		if (props.multiple) {
			multiples.value.push({
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
					base64: isUpload ? ((await fileToBase64(file)) as string) : "",
				},
				isUpload,
				status: isUpload ? "upload" : "not_edit",
			})
			emit("update:modelValue", multiples.value)
		} else {
			single.value = {
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
			emit("update:modelValue", single.value)
		}
	} else {
		single.value = {
			data: null,
			isUpload: false,
			status: "no_file",
		}
		emit("update:modelValue", single.value)
	}
}

const handleRemoveFile = (_arr: FilePondErrorDescription | null, fileItem: FilePondFile) => {
	if (props.multiple) {
		multiples.value.map((file: IFile) => file.data?.id === fileItem.id)
		for (let i = 0; i < multiples.value.length; i++) {
			if (multiples.value[i].data?.id === fileItem.id) {
				// multiples.value.splice(i--, 1)
				multiples.value[i].status = "delete"
			}
		}

		emit("update:modelValue", multiples.value)
		emit("update:deleteImage", multiples.value)
	} else {
		single.value = {
			data: null,
			isUpload: false,
			status: "delete",
		}
		emit("update:modelValue", single.value)
		emit("update:deleteImage", single.value)
	}
}

const handleWarning = (error: FilePondErrorDescription, file: FilePondFile[]) => {
	if (props.multiple) {
		if (error.body === "Max files") {
			// ซ่อน div
			// if (multiples.value.length === 0 || multiples.value.length > props.maxFiles) {
			// 	const element = document.querySelector(".filepond--drop-label") as HTMLElement
			// 	if (element) {
			// 		element.style.visibility = "hidden"
			// 	}
			// }

			// กรณีเลือกไฟล์เกิน
			if (typeof file === "object") {
				if (multiples.value.length + file.length > props.maxFiles) {
					showAlert({
						title: "แจ้งเตือน",
						message: `โปรดอัปโหลดไฟล์ไม่เกิน ${props.maxFiles} ไฟล์`,
						type: "warning",
						callBack() {
							emit("update:modelValue", multiples.value)
						},
					})
				}
			}
		}
	}
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

		if (props.multiple) {
			multiples.value = []
			emit("update:modelValue", multiples.value)
		} else {
			single.value = {
				data: null,
				isUpload: false,
				status: "no_file",
			}
			emit("update:modelValue", single.value)
		}
	})
}

// นามสกุลไฟล์
const extensions = ref("")
if (props.acceptedFileTypes!.length > 0) {
	extensions.value = [
		...new Set(props.acceptedFileTypes!.map((mimeType: any) => "." + getExtensionFromMimeType(mimeType as string))),
	]
		.filter((extension) => extension !== "")
		.join(", ")
}

// Files
const toFiles = (files: String | Array<String>) => {
	const result: any[] = []
	if (Array.isArray(files)) {
		for (let i = 0; i < files.length; i++) {
			result.push({
				source: files[i],
				option: {
					type: "limbo",
				},
			})
		}
	} else if (files) {
		result.push({
			source: files,
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
	if (props.files !== null) {
		initFiles.value = toFiles(props.files)
	}
	return initFiles.value
})

const onParentActive = () => {
	setTimeout(() => {
		updateWidth()
	}, 300)
}

const updateWidth = () => {
	const filepondElement = document.querySelectorAll(".filepond--wrapper")
	filepondElement.forEach((el) => {
		const item = el as HTMLElement
		const filepondWidth = item.offsetWidth
		const contentElement = item.querySelector(".content") as HTMLElement
		const iconElement = item.querySelector(".icon") as HTMLElement
		const buttonElement = item.querySelector(".button") as HTMLElement
		const uploadText = ref()
		if (buttonElement) {
			uploadText.value = buttonElement.querySelectorAll("span")
		}
		if (filepondWidth && contentElement && iconElement && buttonElement) {
			if (filepondWidth < 395 && filepondWidth > 324) {
				contentElement.style.fontSize = "11.3px"
				uploadText.value[0].style.fontSize = "11.3px"
				iconElement.style.width = "35px"
				iconElement.style.height = "35px"
				buttonElement.style.height = "35px"
			} else if (filepondWidth <= 324) {
				contentElement.style.fontSize = "9px"
				uploadText.value[0].style.fontSize = "8px"
				iconElement.style.width = "30px"
				iconElement.style.height = "30px"
				buttonElement.style.height = "30px"
			} else {
				contentElement.style.fontSize = "1rem"
				uploadText.value[0].style.fontSize = "1rem"
				iconElement.style.width = "45px"
				iconElement.style.height = "45px"
				buttonElement.style.height = "40px"
			}
		}
	})
}

onMounted(() => {
	setTimeout(() => {
		updateWidth()
	}, 300)
	window.addEventListener("resize", updateWidth)
})

onUnmounted(() => {
	window.removeEventListener("resize", updateWidth)
})

// Validate
const { meta } = useField(props.name)

defineExpose({
	clearFile,
	onParentActive,
})
</script>
<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />
	<div v-if="showFilePond" class="py-3 pt-0">
		<VField style="display: none" :name="name" :model-value="modelValue" />
		<file-pond
			ref="filePondRef"
			class="filepond"
			:class="{
				multiple: multiple,
				single: !multiple,
				disabled: disabled,
				'is-invalid': !meta.valid && meta.validated,
			}"
			:disabled="disabled"
			:name="name"
			:allow-multiple="multiple"
			label-button-process-item="Upload"
			:model-value="modelValue"
			:files="computedInitFiles"
			:max-file-size="maxFileSize"
			label-max-file-size-exceeded="ไฟล์มีขนาดใหญ่เกิน"
			:label-max-file-size="`ขนาดไฟล์ต้องไม่เกิน ${maxFileSize}`"
			label-max-total-file-size-exceeded="ไฟล์มีขนาดใหญ่เกิน"
			:label-max-total-file-size="(multiple ? 'ไฟล์มีขนาดรวมเกิน' : 'ขนาดไฟล์ต้องไม่เกิน') + ` ${totalFileSize}`"
			:image-validate-size-max-width="imageSize ? imageSize : ''"
			:image-validate-size-max-height="imageSize ? imageSize : ''"
			image-validate-size-label-format-error="ไม่รองรับประเภทรูปภาพนี้"
			image-validate-size-label-image-size-too-big="ขนาดรูปภาพไม่ถูกต้อง"
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
			label-file-load-error="ไม่สามารถโหลดข้อมูลได้"
			label-file-waiting-for-size="รอโหลดข้อมูลขนาด"
			label-file-type-not-allowed="ประเภทไฟล์ไม่ถูกต้อง"
			file-validate-type-label-expected-types="รองรับ {allButLastType} หรือ {lastType}"
			:allow-pdf-preview="true"
			pdf-preview-height="320"
			pdf-component-extra-params="toolbar=0&scrollbar=0&zoom=1&page=1"
			:allow-download-by-url="false"
			:label-idle="`
			<div class='row align-items-center'>
				<div class='col-2 p-2'>
					<svg class='icon' width='45' height='45' viewBox='0 0 78 78' fill='none' xmlns='http://www.w3.org/2000/svg'>
						<path d='M58.0725 25.5787C57.5335 25.427 57.0416 25.1448 56.6414 24.7577C56.2412 24.3706 55.9455 23.8909 55.7812 23.3622C54.8687 20.311 53.3383 17.4742 51.2824 15.0234C49.2266 12.5726 46.6881 10.5588 43.8207 9.10383C40.9532 7.64886 37.8163 6.78297 34.5998 6.55858C31.3833 6.33419 28.1541 6.75595 25.1076 7.79834C22.0611 8.84074 19.2607 10.4821 16.8757 12.6231C14.4908 14.7641 12.5709 17.3602 11.2323 20.2543C9.89376 23.1485 9.16427 26.2805 9.08802 29.4608C9.01177 32.6411 9.59035 35.8036 10.7888 38.757C11.0641 39.3785 11.1233 40.0724 10.9573 40.7306C10.7914 41.3888 10.4094 41.9744 9.87103 42.3962C7.49843 44.1335 5.64605 46.4713 4.50993 49.1622C3.37381 51.853 2.99623 54.7969 3.41715 57.6822C4.07402 61.5827 6.1266 65.1216 9.20234 67.6566C12.2781 70.1916 16.1731 71.5546 20.1805 71.4983H36.0232C36.8135 71.4983 37.5714 71.1886 38.1301 70.6373C38.6889 70.0861 39.0029 69.3383 39.0029 68.5587C39.0029 67.7791 38.6889 67.0314 38.1301 66.4801C37.5714 65.9288 36.8135 65.6191 36.0232 65.6191H20.1805C17.607 65.6807 15.0962 64.8307 13.1032 63.2231C11.1102 61.6155 9.76695 59.3567 9.31681 56.8561C9.03136 55.0181 9.26271 53.138 9.98566 51.4207C10.7086 49.7035 11.8954 48.2149 13.4168 47.1172C15.0148 45.9389 16.1708 44.2702 16.7031 42.3734C17.2354 40.4765 17.1139 38.4589 16.3577 36.6375C14.8624 32.7417 14.7855 28.4543 16.1402 24.5087C17.2224 21.4179 19.1603 18.6879 21.7344 16.6279C24.3085 14.5679 27.4159 13.2603 30.7046 12.8531C31.4674 12.7562 32.2357 12.7071 33.0049 12.7061C36.8577 12.6936 40.611 13.9127 43.7028 16.1807C46.7946 18.4488 49.0585 21.6439 50.1556 25.2877C50.5839 26.6763 51.3555 27.9379 52.401 28.9591C53.4465 29.9803 54.733 30.7289 56.1447 31.1375C59.5995 32.1456 62.6598 34.1708 64.9158 36.9421C67.1718 39.7134 68.5158 43.0985 68.7679 46.6442C69.0201 50.1899 68.1683 53.7269 66.3266 56.7818C64.4849 59.8367 61.7412 62.2637 58.4628 63.7378C57.9776 63.9828 57.572 64.3575 57.2923 64.8191C57.0125 65.2806 56.8699 65.8105 56.8806 66.3481C56.8748 66.8346 56.993 67.3146 57.2242 67.7443C57.4554 68.1739 57.7924 68.5394 58.204 68.8071C58.6157 69.0749 59.089 69.2363 59.5803 69.2766C60.0717 69.317 60.5655 69.2349 61.0164 69.0379C73.3342 63.1969 80.0264 47.8315 69.5977 32.9865C66.7262 29.2912 62.659 26.677 58.0725 25.5787Z' fill='#7E8299' />
						<path d='M57.7852 55.6547C58.2429 55.1637 58.5 54.4979 58.5 53.8036C58.5 53.1093 58.2429 52.4434 57.7852 51.9524L53.9133 47.7998C52.5398 46.3272 50.6773 45.5 48.7353 45.5C46.7932 45.5 44.9307 46.3272 43.5572 47.7998L39.6853 51.9524C39.2406 52.4463 38.9945 53.1076 39.0001 53.7942C39.0057 54.4807 39.2624 55.1374 39.715 55.6228C40.1677 56.1083 40.78 56.3836 41.4201 56.3896C42.0602 56.3956 42.6769 56.1317 43.1373 55.6547L46.2939 52.2693V72.1317C46.2939 72.8261 46.5512 73.4921 47.009 73.9831C47.4668 74.4741 48.0878 74.75 48.7353 74.75C49.3827 74.75 50.0037 74.4741 50.4615 73.9831C50.9194 73.4921 51.1766 72.8261 51.1766 72.1317V52.2693L54.3332 55.6547C54.791 56.1456 55.4119 56.4213 56.0592 56.4213C56.7066 56.4213 57.3274 56.1456 57.7852 55.6547Z' fill='#7E8299' />
					</svg>
				</div>
				<div class='col-7 px-lg-4 mt-2 content p-0'>
					<span>
            ลากหรือเลือกไฟล์ที่ต้องการอัปโหลด ${extensions !== `` ? `<br />รองรับไฟล์นามสกุล ${extensions}` : ``}
            ${imageSize ? `<br />ขนาดไม่เกิน ${imageSize} x ${imageSize} px` : ``}
            </span>
                    </div>
                    <div class='col-3 p-0' style='text-align: -webkit-center;'>
            <div class='button rounded-4 d-flex justify-content-center'>
							<span>เลือกไฟล์</span>
						</div>
				</div>
			</div>`"
			:allow-file-poster="true"
			:file-poster-height="100"
			:allow-audio-preview="true"
			:allow-video-preview="true"
			:max-files="maxFiles"
			:accepted-file-types="acceptedFileTypes"
			:file-validate-type-label-expected-types-map="getMimeTypes()"
			:allow-file-type-validation="allowFileTypeValidation"
			:file-validate-type-detect-type="fileValidateTypeDetectType"
			:style-item-panel-aspect-ratio="aspectRatio"
			@addfile="handleFile"
			@removefile="handleRemoveFile"
			@warning="handleWarning"
			@update:delete-value="($event: any) => (initFiles = $event)"
		/>
		<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
	</div>
</template>

<style lang="scss">
.filepond--item-panel {
	border-radius: 1rem;
}
.filepond--file-poster-wrapper {
	border-radius: 1rem;
}

// Mode: single
.filepond.single {
	width: 100%;
	overflow: hidden;
	min-height: 100px;
	padding: 1rem 0.2rem 0.8rem;
	border-radius: 1.25rem;
	background-color: #f4f4f4;
	border: 2px var(--kt-gray-300) dashed;

	.filepond--credits {
		display: none;
	}

	svg.icon {
		@media (max-width: 400px) {
			height: 3.5rem;
			margin-top: 8px;
		}
	}

	.content {
		font-size: 1rem;
		align-self: center;
		padding: 0;
	}

	.button {
		border-radius: 1.25rem;
		align-self: center;
		padding: 0.75rem 0.5rem;
		cursor: pointer;
		color: var(--kt-gray-800);
		border: 1px solid var(--kt-primary);
		background-color: white;
		font-size: 1rem !important;
		font-weight: 500 !important;
		margin-top: 13px;
		max-width: 100px;

		&:hover {
			background-color: #e8edf3;
		}

		@media (max-width: 767px) {
			max-width: 100px;
			padding: 0.65rem 0.5rem;
			margin-top: 10px;
		}
	}

	.filepond--panel-root {
		background-color: #f4f4f4;
		border-radius: 1.25rem;
	}

	.filepond--root {
		margin-bottom: -10px;
	}

	.filepond--drop-label > label {
		color: #a1a5b7;
		width: 100% !important;
		padding: 1em !important;
	}

	.filepond--file {
		font-size: 1.25rem !important;
	}

	.filepond--action-remove-item {
		background-color: #383737;
	}

	.filepond--root .filepond--list-scroller {
		margin-top: 0em;
		margin-bottom: 1em;
	}

	.filepond--item {
		min-height: 95px;
		margin-top: -23px;
	}

	.filepond--root .filepond--drop-label {
		padding: 0 10px;
	}

	.filepond--image-bitmap {
		text-align: center;
		display: block;
		width: 100%;
		height: 100%;

		canvas {
			width: 100%;
			height: 100%;
			padding: 1.25rem;
		}
	}

	.filepond--item > .filepond--panel .filepond--panel-bottom {
		box-shadow: none;
	}

	// Poster
	.filepond--file-poster {
		padding: 1.5em;
	}
}

// Mode: multiple
.filepond.multiple {
	.filepond--item {
		border-left: 3px solid #f4f4f4;
		border-right: 3px solid #f4f4f4;

		height: 65px !important;
		width: calc(33.33%);

		@media (max-width: 1199px) {
			width: calc(50%);
		}
	}

	.filepond--file-poster img {
		height: 200%;
	}

	.filepond--file-info-main {
		font-size: 11px !important;
		margin-top: 2px;
	}

	.filepond--file .filepond--file-status {
		margin-right: 0px !important;
		right: 15px !important;
		bottom: 3px;
		transform: unset !important;
		font-size: 14px;
		position: absolute;
		z-index: 999;
	}

	.filepond--file-info {
		margin-top: -3px;
	}

	// .filepond--file-info + .filepond--file-status[style="transform: translate3d(35px, 0px, 0px); opacity: 1;"] {
	// 	display: none;
	// }

	.filepond--panel-center {
		height: 40px !important;
		border-radius: 0px 0px 12px 12px !important;
		transform: none !important;
		top: 25px !important;
	}

	.filepond--file-info-main-container {
		width: 80%;
	}

	.filepond--panel-bottom {
		display: none;
	}

	.filepond--pdf-preview-wrapper:before {
		height: 125%;
	}

	width: 100%;
	min-height: 100px;
	padding: 0.2rem 0.1rem 0.8rem;
	border-radius: 1.25rem;
	background-color: #f4f4f4;
	border: 2px var(--kt-gray-300) dashed;
	overflow: hidden;

	.filepond--drop-label {
		padding-top: 0px;
	}

	.filepond--credits {
		display: none;
	}

	svg.icon {
		@media (max-width: 400px) {
			height: 3.5rem;
			margin-top: 8px;
		}
	}

	.filepond--item > .filepond--panel .filepond--panel-bottom {
		// background-color: #f4f4f4;
		box-shadow: none;
	}

	.content {
		font-size: 1rem;
		align-self: center;
		padding: 0;

		@media (max-width: 767px) {
			font-size: 1rem;
		}

		@media (max-width: 400px) {
			font-size: 1rem;
		}
	}

	.button {
		border-radius: 1.25rem;
		align-self: center;
		padding: 0.75rem 0.5rem;
		cursor: pointer;
		color: var(--kt-primary);
		border: 1px solid var(--kt-primary);
		background-color: transparent;
		font-size: 1rem !important;
		font-weight: 500 !important;
		margin-top: 13px;
		max-width: 100px;

		&:hover {
			color: #fff;
			background-color: var(--kt-primary);
		}

		@media (max-width: 767px) {
			max-width: 100px;
			padding: 0.65rem 0.5rem;
			margin-top: 10px;
		}
	}

	.filepond--panel-root {
		background-color: #f4f4f4;
		border-radius: 1.25rem;
	}

	.filepond--root {
		margin-bottom: -10px;
	}

	.filepond--drop-label > label {
		color: #a1a5b7;
		width: 100% !important;
		padding: 0 2em !important;
	}

	.filepond--file {
		font-size: 1.25rem !important;
	}

	.filepond--action-remove-item {
		background-color: #383737;
	}

	.filepond--list.filepond--list {
		top: 5px;
	}

	.filepond--item {
		margin: 0px 0px 5px 0px;
	}

	.filepond--root .filepond--list-scroller {
		top: 0px;
	}

	// Poster
	.filepond--file-poster {
		padding: 1.5em;
	}

	.filepond--drop-label {
		margin-top: 10px;
		&:has(+ .filepond--list-scroller ul li) {
			margin-top: -2px;
		}
	}
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

.filepond--image-canvas-wrapper {
	transform: unset !important;

	.filepond--image-bitmap {
		canvas {
			width: 100% !important;
			height: 100% !important;
			padding: 8px;
		}
	}
}

.filepond--file-status {
	position: static;
	display: grid;
	align-items: unset !important;
	flex-grow: unset !important;
	flex-shrink: unset !important;
	margin: 0;
	min-width: 2.25em;
	text-align: right;

	* {
		white-space: unset !important;
	}

	.filepond--file-status-main {
		text-align: end;
	}
}

.filepond--panel-bottom,
.filepond--panel-top {
	height: 2em;
}

.filepond--magnify-icon {
	display: none;
}

// Icon
.filepond--download-icon {
	height: 18px;
	width: 17px;
	margin-right: 0.3em;
	-webkit-mask-image: url(/images/icons/svg/cloud-download-alt.svg);
	mask-image: url(/images/icons/svg/cloud-download-alt.svg);
}

// Disabled
.filepond.disabled {
	background-color: var(--kt-input-disabled-bg);
	cursor: no-drop;
}

// Validate
.filepond.is-invalid {
	border: 1px var(--kt-danger) solid;
}

// Max-Limit
.filepond--drop-label[style*="visibility: hidden;"] + .filepond--list-scroller {
	margin-top: 0.35em;
}
</style>
