<script setup lang="ts">
import { useField } from "vee-validate"
import { ITree } from "~~/core/shared/types/Tree"

type TMode = "ALL" | "BRANCH_PRIORITY" | "LEAF_PRIORITY" | "ALL_WITH_INDETERMINATE"
type TDirection = "auto" | "below" | "bottom" | "above" | "top"

const props = defineProps({
	modelValue: {
		type: [Array<String>, Array<Number>, Number],
		default: null,
	},
	mode: {
		type: String as PropType<TMode>,
		default: "BRANCH_PRIORITY",
	},
	options: {
		type: (Array as PropType<ITree[]>) || [],
		required: true,
	},
	searchable: {
		type: Boolean,
		default: true,
	},
	placeholder: {
		type: String,
		default: "เลือก",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	label: {
		type: String,
		default: "",
	},
	name: {
		type: String,
		required: true,
	},
	required: {
		type: Boolean,
		default: false,
	},
	noOptionsText: {
		type: String,
		default: "ไม่มีรายการ",
	},
	noResultsText: {
		type: String,
		default: "ไม่พบข้อมูล",
	},
	limit: {
		type: Number,
		default: Infinity,
	},
	multiple: {
		type: Boolean,
		default: false,
	},
	showCount: {
		type: Boolean,
		default: true,
	},
	disableBranchNodes: {
		type: Boolean,
		default: false,
	},
	defaultExpandLevel: {
		type: Number,
		default: 1,
	},
	alwaysOpen: {
		type: Boolean,
		default: false,
	},
	reposition: {
		type: Boolean,
		default: false,
	},
	openDirection: {
		type: String as PropType<TDirection>,
		default: "auto",
	},
})

const treeselect = ref()
const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const updateValue = (value: any) => {
	emit("update:modelValue", value)
	console.log("props.modelValue 0= ", value)
	// setPlaceholder()
}

const limitText = (count: string) => {
	if (props.limit > 0) {
		return `และอีก ${count} รายการ`
	} else {
		return `เลือก ${count} รายการ`
	}
}

const handleOpen = () => {
	if (props.reposition) {
		repositionDropDown()
	}

	const treeselectPlaceholder = treeselect.value.querySelector(".vue-treeselect__placeholder")
	if (treeselectPlaceholder) {
		treeselectPlaceholder.classList.add("vue-treeselect-helper-hide")
	}
}

const handleClose = () => {
	if (props.reposition) {
		resetPositionDropDown()
	}

	setPlaceholder()
}

const setPlaceholder = () => {
	setTimeout(() => {
		const treeselectPlaceholder = treeselect.value.querySelector(".vue-treeselect__placeholder")
		const treeselectInput = treeselect.value.querySelector(".vue-treeselect__input")

		if (Array.isArray(props.modelValue)) {
			if (props.modelValue.length === 0) {
				if (treeselectPlaceholder) {
					treeselectPlaceholder.classList.remove("vue-treeselect-helper-hide")
				}
				if (treeselectInput) {
					treeselectInput.value = ""
				}
			}
		} else if (!props.modelValue && typeof props.modelValue !== "number") {
			if (treeselectPlaceholder) {
				treeselectPlaceholder.classList.remove("vue-treeselect-helper-hide")
			}
			if (treeselectInput) {
				treeselectInput.value = ""
			}
		}
	}, 150)
}

const direction = props.reposition ? ref("bottom") : ref(props.openDirection)
const handleDirection = () => {
	const screenHeight = window.innerHeight
	const ref = vueTreeselectRef.value.$el
	const { bottom } = ref.getBoundingClientRect()

	const spaceBelow = bottom + 400 // ความสูงของ option

	if (spaceBelow > screenHeight) {
		direction.value = "top"
	} else {
		direction.value = "bottom"
	}
}

onMounted(() => {
	handleDirection()
	window.addEventListener("scroll", handleDirection)
})

// กรณี VTree อยู่ในตาราง
const vueTreeselectRef = ref()

const resetPositionDropDown = () => {
	const ref = vueTreeselectRef.value.$el

	if (vueTreeselectRef) {
		ref.style.width = `100%`
		ref.style.position = "relative"
		ref.style.top = `auto`
		ref.style.left = `auto`
		ref.style.zIndex = `auto`

		document.body.classList.remove("scrollbar-stopped")
	}
}

const repositionDropDown = () => {
	const ref = vueTreeselectRef.value.$el
	const { top, left, width } = ref.getBoundingClientRect()

	if (vueTreeselectRef) {
		ref.style.width = `${width}px`
		ref.style.position = "fixed"
		ref.style.top = `${top}px`
		ref.style.left = `${left}px`
		ref.style.zIndex = `999`

		document.body.classList.add("scrollbar-stopped")
	}
}

onUnmounted(() => {
	window.removeEventListener("scroll", handleDirection)

	resetField({ value: null })
})
</script>

<template>
	<VLabel v-if="label !== ''" :label="label" :required="required" />

	<VField style="display: none" :name="name" :model-value="modelValue" />

	<div ref="treeselect">
		<VueTreeselect
			ref="vueTreeselectRef"
			:model-value="modelValue"
			:name="name"
			:placeholder="placeholder"
			:searchable="searchable"
			:options="options"
			:disabled="disabled"
			:no-options-text="noOptionsText"
			:no-results-text="noResultsText"
			:limit="limit"
			:multiple="multiple"
			:limit-text="limitText"
			:show-count="showCount"
			:max-height="300"
			:class="{
				'is-invalid': !meta.valid && meta.validated,
				'vue-treeselect_multiple': multiple,
				'vue-treeselect_limit': limit > 0,
			}"
			:disable-branch-nodes="disableBranchNodes"
			:default-expand-level="defaultExpandLevel"
			:value-consists-of="mode"
			:always-open="alwaysOpen"
			:open-direction="direction"
			@update:model-value="updateValue"
			@close="handleClose"
			@open="handleOpen"
		/>

		<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
	</div>
</template>

<style lang="scss">
.vue-treeselect__menu {
	border: 1px solid var(--kt-input-border-color);
	padding-top: 0px;
	padding-bottom: 0px;
	border-top-left-radius: 1rem !important;
	border-top-right-radius: 1rem !important;
	max-height: 400px !important;
}
.vue-treeselect--open .vue-treeselect__control .vue-treeselect__value-container {
	border: 1px solid var(--kt-input-border-color);
}
.vue-treeselect--open-above .vue-treeselect__menu-container .vue-treeselect__menu {
	border-top-left-radius: 1.25rem !important;
	border-top-right-radius: 1.25rem !important;
}
.vue-treeselect--open-below .vue-treeselect__menu-container .vue-treeselect__menu {
	border-top-left-radius: 0rem !important;
	border-top-right-radius: 0rem !important;
	border-bottom-left-radius: 1.25rem !important;
	border-bottom-right-radius: 1.25rem !important;
}
.vue-treeselect__control {
	color: var(--kt-input-color);
	background-color: var(--kt-input-bg);
	border: 1px solid #eaecf3;
	box-shadow: none !important;
	border-radius: var(--kt-form-border-radius);
	height: var(--kt-input-height);
}
.vue-treeselect__multi-value,
.vue-treeselect__value-container {
	padding: 0px 3px;
	margin-bottom: 3px !important;
}
.vue-treeselect__placeholder,
.vue-treeselect__single-value {
	margin-left: 10px;
	font-size: 1.05rem !important;
}
.vue-treeselect--searchable:not(.vue-treeselect--disabled) .vue-treeselect__value-container {
	height: var(--kt-input-height);
	border: 0px;
}
.vue-treeselect__option {
	border-collapse: initial;
	padding: 0.5rem 0.75rem;
}
.vue-treeselect_multiple .vue-treeselect__option {
	padding: 0.25rem 0.75rem;
}
.vue-treeselect__checkbox--checked > .vue-treeselect__check-mark {
	background-image: var(--kt-form-check-input-checked-bg-image) !important;
	background-size: 50% 50%;
	width: 20px;
	height: 20px;
	left: 2px;
	top: 2px;
}
.vue-treeselect--single .vue-treeselect__option--selected {
	background: #fff0d9;
}
.vue-treeselect--single .vue-treeselect__option--selected:hover {
	background: #fff0d9;
}
.vue-treeselect__minus-mark {
	width: 20px;
	height: 20px;
	background-size: 8px 6px;
	left: 3px;
	top: 4px;
}
.vue-treeselect__checkbox-container {
	width: 30px;
	height: 30px;
}
.vue-treeselect__checkbox {
	width: 1.3rem !important;
	height: 1.3rem !important;
	border-radius: 0.25rem !important;
}
.vue-treeselect__option-arrow {
	margin-left: -2px;
}
.vue-treeselect__option-arrow-container,
.vue-treeselect__option-arrow-placeholder {
	width: 15px;
}
.vue-treeselect__checkbox--checked,
.vue-treeselect__checkbox--indeterminate,
.vue-treeselect__label-container:hover .vue-treeselect__checkbox--checked,
.vue-treeselect__label-container:hover .vue-treeselect__checkbox--indeterminate {
	border-color: var(--kt-form-check-input-checked-bg-color-solid);
	background: var(--kt-form-check-input-checked-bg-color-solid);
}
.vue-treeselect__label-container:hover .vue-treeselect__checkbox--unchecked {
	border-color: var(--kt-form-check-input-checked-bg-color-solid);
}
.vue-treeselect__label {
	font-size: 0.975rem;
	font-weight: 400;
	text-overflow: initial !important;
	white-space: pre-line !important;
}
.vue-treeselect__x-container {
	width: 8px;
	color: #999;
}
.vue-treeselect__control-arrow-container {
	width: 30px;
}
.vue-treeselect__control-arrow {
	width: 11px;
	height: 11px;
	color: #5e6278;
}
.vue-treeselect__limit-tip-text {
	cursor: default;
	display: block;
	margin: 0px 0px 0px 5px;
	padding: 5px 0;
	color: #181c32;
	font-size: 1.05rem !important;
	font-weight: 400;
}
.vue-treeselect_limit .vue-treeselect__limit-tip-text {
	color: #5e6278;
	font-size: 0.95rem !important;
}
.vue-treeselect__multi-value-label {
	font-size: 0.975rem;
	padding: 3px 8px;
	// white-space: nowrap;
	// overflow: hidden;
	// text-overflow: ellipsis;
	// max-width: 100px;
}
.vue-treeselect:not(.vue-treeselect--disabled)
	.vue-treeselect__multi-value-item:not(.vue-treeselect__multi-value-item-disabled):hover
	.vue-treeselect__multi-value-item:not(.vue-treeselect__multi-value-item-new)
	.vue-treeselect__multi-value-item:not(.vue-treeselect__multi-value-item-new):hover,
.vue-treeselect__multi-value-item {
	cursor: pointer;
	background: var(--kt-primary-light);
	color: var(--kt-primary);
	border-radius: 8px;
	font-weight: 500;
}
.vue-treeselect__value-remove {
	color: var(--kt-primary);
	padding-left: 0px;
	border-left: 0px;
	line-height: 0;
}
.vue-treeselect__multi-value-item-container {
	padding-top: 4px;
}
.vue-treeselect--focused:not(.vue-treeselect--open) .vue-treeselect__control {
	border: 1px solid var(--kt-input-border-color);
}
.vue-treeselect__control {
	padding-left: 0px !important;
}
.vue-treeselect__count {
	font-weight: 400;
	opacity: 0.5;
}

.vue-treeselect__control-arrow-container {
	svg {
		&.vue-treeselect__control-arrow--rotated {
			transform: rotate(180deg);
		}
		height: 1.5rem;
		width: 1.5rem;

		-webkit-mask-position: center;
		mask-position: center;
		-webkit-mask-repeat: no-repeat;
		mask-repeat: no-repeat;
		-webkit-mask-size: contain;
		mask-size: contain;

		flex-grow: 0;
		flex-shrink: 0;
		pointer-events: none;
		position: relative;
		transform: rotate(0deg);
		transition: transform 0.3s;
		z-index: 10;
		background-color: #5e6278;
		-webkit-mask-image: url("data:image/svg+xml;charset=utf-8,%3Csvg viewBox='0 0 320 512' fill='currentColor' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M31.3 192h257.3c17.8 0 26.7 21.5 14.1 34.1L174.1 354.8c-7.8 7.8-20.5 7.8-28.3 0L17.2 226.1C4.6 213.5 13.5 192 31.3 192z'/%3E%3C/svg%3E");
		mask-image: url("data:image/svg+xml;charset=utf-8,%3Csvg viewBox='0 0 320 512' fill='currentColor' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M31.3 192h257.3c17.8 0 26.7 21.5 14.1 34.1L174.1 354.8c-7.8 7.8-20.5 7.8-28.3 0L17.2 226.1C4.6 213.5 13.5 192 31.3 192z'/%3E%3C/svg%3E");
	}
}

.vue-treeselect__x-container {
	svg {
		height: 1rem;
		width: 0.75rem;

		-webkit-mask-position: center;
		mask-position: center;
		-webkit-mask-repeat: no-repeat;
		mask-repeat: no-repeat;
		-webkit-mask-size: contain;
		mask-size: contain;

		flex-grow: 0;
		flex-shrink: 0;
		pointer-events: none;
		position: relative;
		transform: rotate(0deg);
		transition: transform 0.3s;
		z-index: 10;
		background-color: #5e6278;
		-webkit-mask-image: url("data:image/svg+xml;charset=utf-8,%3Csvg viewBox='0 0 320 512' fill='currentColor' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='m207.6 256 107.72-107.72c6.23-6.23 6.23-16.34 0-22.58l-25.03-25.03c-6.23-6.23-16.34-6.23-22.58 0L160 208.4 52.28 100.68c-6.23-6.23-16.34-6.23-22.58 0L4.68 125.7c-6.23 6.23-6.23 16.34 0 22.58L112.4 256 4.68 363.72c-6.23 6.23-6.23 16.34 0 22.58l25.03 25.03c6.23 6.23 16.34 6.23 22.58 0L160 303.6l107.72 107.72c6.23 6.23 16.34 6.23 22.58 0l25.03-25.03c6.23-6.23 6.23-16.34 0-22.58L207.6 256z'/%3E%3C/svg%3E");
		mask-image: url("data:image/svg+xml;charset=utf-8,%3Csvg viewBox='0 0 320 512' fill='currentColor' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='m207.6 256 107.72-107.72c6.23-6.23 6.23-16.34 0-22.58l-25.03-25.03c-6.23-6.23-16.34-6.23-22.58 0L160 208.4 52.28 100.68c-6.23-6.23-16.34-6.23-22.58 0L4.68 125.7c-6.23 6.23-6.23 16.34 0 22.58L112.4 256 4.68 363.72c-6.23 6.23-6.23 16.34 0 22.58l25.03 25.03c6.23 6.23 16.34 6.23 22.58 0L160 303.6l107.72 107.72c6.23 6.23 16.34 6.23 22.58 0l25.03-25.03c6.23-6.23 6.23-16.34 0-22.58L207.6 256z'/%3E%3C/svg%3E");
	}
}

.vue-treeselect__x-container:hover {
	color: #000;
}

.vue-treeselect__indent-level-1 .vue-treeselect__option {
	padding-left: 30px;
}
.vue-treeselect__indent-level-2 .vue-treeselect__option {
	padding-left: 40px;
}
.vue-treeselect__indent-level-3 .vue-treeselect__option {
	padding-left: 50px;
}
.vue-treeselect__indent-level-4 .vue-treeselect__option {
	padding-left: 60px;
}
.vue-treeselect__indent-level-5 .vue-treeselect__option {
	padding-left: 70px;
}
.vue-treeselect__indent-level-6 .vue-treeselect__option {
	padding-left: 80px;
}
.vue-treeselect__indent-level-7 .vue-treeselect__option {
	padding-left: 90px;
}

.vue-treeselect__no-results-tip .vue-treeselect__tip-text {
	padding: 10px 0;
}

.vue-treeselect--multi .vue-treeselect__input {
	padding-top: 5px;
	font-size: 1.05rem;
}
.vue-treeselect__placeholder,
.vue-treeselect__single-value {
	line-height: 3.2;
}
.vue-treeselect__placeholder {
	color: #9ca3af;
	margin-left: 7px !important;
}

// Validate
.is-invalid .vue-treeselect__control {
	border: 1px var(--kt-danger) solid;
}

// Hover
.vue-treeselect__option--highlight {
	background: #fff;
	&:hover {
		background: #f5f5f5;
	}
}

// กรณี VTree อยู่ในตาราง
.scrollbar-stopped {
	overflow: hidden !important;
	padding-right: 15px !important;
}
</style>
