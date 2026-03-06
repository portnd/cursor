<script setup lang="ts">
import { useField } from "vee-validate"

const props = defineProps({
	modelValue: {
		type: String,
		default: "",
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
	placeholder: {
		type: String,
		default: "",
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

const { meta, resetField, errorMessage } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const updateValue = (event: any) => {
	emit("update:modelValue", event.target.value)
}

// ปุ่มซ่อน/แสดงรหัสผ่าน
const passwordFieldType = ref("password")
const icon = ref("fi-rr-eye-crossed")
const switchVisibility = () => {
	if (passwordFieldType.value === "password") {
		passwordFieldType.value = "text"
		icon.value = "fi-rr-eye"
	} else {
		passwordFieldType.value = "password"
		icon.value = "fi-rr-eye-crossed"
	}
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<div class="position-relative">
		<VLabel v-show="label" :label="label" :required="required" />
		<VField
			:type="passwordFieldType"
			:name="name"
			class="form-control"
			:placeholder="placeholder"
			:value="modelValue"
			:class="{ 'is-invalid': !meta.valid && meta.validated }"
			:disabled="disabled"
			:readonly="readonly"
			@input="updateValue"
		/>
		<span class="btn btn-sm btn-icon position-absolute translate-middle icon-end" @click="switchVisibility">
			<i class="fs-2 text-gray-500" :class="icon"></i>
		</span>
		<span v-if="!meta.valid && meta.validated" class="invalid-feedback" :name="name" v-html="errorMessage"></span>

		<!-- <VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" /> -->
	</div>
</template>

<style scoped>
input[type="password"] {
	font-family: Verdana !important;
	font-size: 1.5rem;
}
.icon-end {
	right: 0;
	margin: -16px -5px 0 0;
	width: 40px !important;
	text-align: end;
	display: block;
}
</style>
