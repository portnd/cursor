<script setup lang="ts">
import { useField } from "vee-validate"

type TMode = "center" | "start" | "end"

const props = defineProps({
	modelValue: {
		type: Number as PropType<string | number | null>,
		default: null,
	},
	name: {
		type: String,
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
	placeholder: {
		type: String,
		default: "0669999999",
	},
	align: {
		type: String as PropType<TMode>,
		default: "start",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	readonly: {
		type: Boolean,
		default: false,
	},
})

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])

const updateValue = (event: any) => {
	emit("update:modelValue", event.target.value)
}

const computedClass = computed(() => {
	let className = ""
	if (!meta.valid && meta.validated) {
		className = "is-invalid"
	}

	className += ` text-${props.align}`
	return className
})

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<div class="position-relative">
		<VLabel v-show="label !== ''" :label="label" :required="required" />
		<VField
			:id="name"
			v-maska
			as="input"
			type="text"
			class="form-control"
			:class="computedClass"
			:placeholder="placeholder"
			:value="modelValue"
			:name="name"
			data-maska="0#########"
			:disabled="disabled"
			:readonly="readonly"
			autocomplete="off"
			@input="updateValue"
		/>
		<img src="/images/icons/png/thailand.png" class="icon" alt="" />

		<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
	</div>
</template>

<style scoped>
.form-control {
	padding-left: 3.25rem;
}
.icon {
	position: absolute;
	width: 23px;
	top: 40px;
	left: 10px;
}
</style>
