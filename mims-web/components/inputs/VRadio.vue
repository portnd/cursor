<script setup lang="ts">
import { useField } from "vee-validate"
import { IOption } from "~~/core/shared/types/Option"

const props = defineProps({
	modelValue: {
		type: [String, Number],
		default: "",
	},
	name: {
		type: String,
		default: "",
	},
	options: {
		type: Array<IOption>,
		required: true,
	},
	label: {
		type: String,
		default: "",
	},
	required: {
		type: Boolean,
		default: false,
	},
	inline: {
		type: Boolean,
		default: true,
	},
	disabled: {
		type: Boolean,
		default: false,
	},
})

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const updateValue = (value: any) => {
	emit("update:modelValue", value)
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />

	<div class="mt-3">
		<div
			v-for="(option, key) in options"
			:key="key"
			class="form-check form-check-custom mb-2 me-6"
			:class="{ 'form-check-inline': inline }"
		>
			<VField
				:id="`${name}${option.value}`"
				as="input"
				class="form-check-input cursor-pointer"
				type="radio"
				:name="`${name}${option.value}`"
				:value="option.value"
				:model-value="modelValue"
				:disabled="disabled ? disabled : option.disabled"
				:class="{ 'is-invalid': !meta.valid && meta.validated }"
				@update:model-value="updateValue"
			/>
			<label
				:for="`${name}${option.value}`"
				:style="`color: ${option.color}`"
				class="form-check-label ms-3 cursor-pointer"
				>{{ option.label }}</label
			>
		</div>
	</div>
	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style scoped>
.form-check-input {
	line-height: 1.9;
	width: 1.5rem;
	height: 1.5rem;
}
</style>
