<script setup lang="ts">
import { useField } from "vee-validate"

const props = defineProps({
	modelValue: {
		type: Boolean,
		default: false,
	},
	label: {
		type: String,
		default: "",
	},
	required: {
		type: Boolean,
		default: false,
	},
	name: {
		type: String,
		required: true,
	},
	value: {
		type: String,
		required: true,
	},
	disabled: {
		type: Boolean,
		default: false,
	},
})

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue", "update:value"])
const updateValue = (event: any) => {
	emit("update:modelValue", event.target.checked)
}

onUnmounted(() => {
	resetField({ value: null })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />
	<div class="form-check form-switch form-check-custom form-check-solid">
		<input
			class="form-check-input cursor-pointer"
			type="checkbox"
			:checked="modelValue"
			:model-value="modelValue"
			:name="name"
			:class="{ 'is-invalid': !meta.valid && meta.validated }"
			:disabled="disabled"
			@input="updateValue"
		/>
		<VField style="display: none" :name="name" :checked="modelValue" :model-value="modelValue" />
		<label class="form-check-label cursor-pointer">{{ value }}</label>
	</div>
	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style scoped></style>
