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
		default: "",
	},
	btnLabel: {
		type: String,
		default: "เลือก",
	},
	required: {
		type: Boolean,
		default: false,
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	placeholder: {
		type: String,
		default: "",
	},
})

const { meta, resetField, errorMessage } = useField(props.name)

const emit = defineEmits(["update:modelValue", "click"])
const updateValue = (event: any) => {
	emit("update:modelValue", event.target.value)
}

const onClick = () => {
	emit("click")
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />
	<div class="input-group mb-3">
		<input
			type="text"
			class="form-control position-relative"
			:class="{ 'is-invalid': !meta.valid && meta.validated }"
			:readonly="true"
			:name="name"
			:value="modelValue"
			:placeholder="placeholder"
			aria-label="Recipient's username"
			aria-describedby="button-addon2"
			@input="updateValue"
		/>
		<!-- <span class="badge badge-light-secondary fs-6">{{ modelValue }}</span> -->
		<button id="button-addon2" class="btn btn-primary mt-0" :disabled="disabled" type="button" @click="onClick">
			{{ btnLabel }}
		</button>

		<span v-if="!meta.valid && meta.validated" class="invalid-feedback" :name="name" v-html="errorMessage"></span>
	</div>
</template>

<style scoped>
.form-control:focus {
	background-color: transparent;
}

input {
	font-weight: 500;
}

.btn-primary {
	height: var(--kt-input-height);
}
</style>
