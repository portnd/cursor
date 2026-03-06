<script setup lang="ts">
import { useField } from "vee-validate"

const props = defineProps({
	modelValue: {
		type: String,
		default: "#dddddd",
	},
	name: {
		type: String,
		required: true,
	},
	label: {
		type: String,
		required: true,
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

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const onChange = (value: string | string[]) => {
	emit("update:modelValue", value)
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel :label="label" :required="required" />
	<div class="color-picker-container" :class="{ 'is-invalid': !meta.valid && meta.validated, disabled: disabled }">
		<VField style="display: none" :name="name" :model-value="modelValue" />
		<VueColorPicker v-show="!disabled" :value="modelValue" format="hex" :size="30" @change="onChange" />

		<div v-show="disabled" class="color-picker-color" :style="{ backgroundColor: modelValue }"></div>
		<div class="color-picker-value">
			{{ modelValue }}
		</div>
	</div>
	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style>
.picker {
	z-index: 999 !important;
}
.color-picker {
	cursor: pointer;
}
.color-picker .color-item {
	border: 0px !important;
}
.color-picker-container {
	display: flex;
	height: var(--kt-input-height);
	border: 1px solid var(--kt-gray-300);
	border-radius: var(--kt-form-border-radius);
	padding: 0.2rem 0.5rem 0.25rem;
	max-width: 150px;
}
.color-picker-container.is-invalid {
	border: 1px solid var(--kt-danger);
}
.color-picker-value {
	padding: 0.775rem 0.65rem;
	font-size: var(--kt-input-font-size);
}
.color-picker-container.disabled {
	color: var(--kt-text-gray-400);
	background-color: var(--kt-input-disabled-bg);
	cursor: no-drop;
}
.color-picker-color {
	width: 30px;
	height: 30px;
	border-radius: 5px;
	margin: 5px;
}
</style>
