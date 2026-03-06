<script setup lang="ts">
import { useField } from "vee-validate"

const props = defineProps({
	modelValue: {
		type: String,
		default: "",
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
	disabled: {
		type: Boolean,
		default: false,
	},
	readonly: {
		type: Boolean,
		default: false,
	},
	minHeight: {
		type: String,
		default: "46px",
	},
	className: {
		type: String,
		default: "",
	},
})

const { textarea, triggerResize } = useTextareaAutosize()

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const updateValue = (event: any) => {
	triggerResize()
	emit("update:modelValue", event.target.value)
}

const computedClass = computed(() => {
	let className = ""
	if (!meta.valid && meta.validated) {
		className = "is-invalid"
	}

	className += ` ${props.className}`
	return className
})

const computedModelValue = computed(() => {
	return decodeHTML(props.modelValue)
})

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />

	<VField style="display: none" :name="name" :model-value="modelValue" />
	<textarea
		ref="textarea"
		class="form-control resize-none"
		:value="computedModelValue"
		:name="name"
		:class="computedClass"
		:disabled="disabled"
		:readonly="readonly"
		@input="updateValue"
	/>

	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style scoped>
textarea.form-control {
	min-height: v-bind(minHeight) !important;
	overflow-y: hidden;
}
.resize-none {
	resize: none;
}
</style>
