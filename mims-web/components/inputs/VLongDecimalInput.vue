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
		default: "",
	},
	align: {
		type: String as PropType<TMode>,
		default: "end",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	readonly: {
		type: Boolean,
		default: false,
	},
	textEnd: {
		type: String,
		default: "",
	},
	precision: {
		type: Number,
		default: 0,
	},
	allowMinus: {
		type: Boolean,
		default: false,
	},
	min: {
		type: Number,
		default: null,
	},
	max: {
		type: Number,
		default: 1000000000000000,
	},
})

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])

// ต้องใช้พร้อมกับ useForm
const validateInput = (value: any) => {
	if (/^[-]?[0-9]+(\.[0-9]*)?$/.test(value) || value === "-") {
		// split จำนวนเต็มกับทศนิยม ออกจากกัน
		const [integerPart, fractionalPart] = value.toString().split(".")

		// เคส min
		if (props.min !== null) {
			if (parseInt(integerPart) <= props.min) {
				return props.min
			}
		}

		// เคส max
		if (parseInt(integerPart) >= props.max) {
			return props.max
		}

		if (fractionalPart) {
			const numDigitsAfterDecimal: number = fractionalPart.length
			// ถ้ามีจำนวนทศนิยมมากกว่าที่กำหนดใน props.precision
			if (numDigitsAfterDecimal > props.precision) {
				const substrFractionalPart = fractionalPart.substring(0, props.precision)

				return String(parseInt(integerPart).toFixed(0)) + (substrFractionalPart ? `.${substrFractionalPart}` : "")
			}
		}

		// Format: 0123 => 123, 00123 => 123
		return value.replace(/^0*(\d+)$/, "$1")
	} else {
		return ""
	}
}

const validateKeypress = (event: any) => {
	const keyCode = event.keyCode || event.which

	if (keyCode === 8 || keyCode === 46) {
		return true
	}

	const char = String.fromCharCode(keyCode)

	const regex = props.allowMinus ? /[0-9.-]/ : /[0-9.]/

	if (regex.test(char)) {
		return true
	} else {
		event.preventDefault()
		return false
	}
}

const validateKeyup = (event: any) => {
	emit("update:modelValue", validateInput(event.target.value))
}

const updateValue = (event: any) => {
	emit("update:modelValue", validateInput(event.target.value))
}

const computedClass = computed(() => {
	let className = ""
	if (!meta.valid && meta.validated) {
		className = "is-invalid"
	}

	className += ` text-${props.align}`
	return className
})

const computedModelValue = computed(() => {
	return props.modelValue === null ? "" : validateInput(String(props.modelValue))
})

onUnmounted(() => {
	resetField({ value: null })
})
</script>

<template>
	<div class="position-relative">
		<VLabel v-show="label !== ''" :label="label" :required="required" />
		<VField
			:id="name"
			as="input"
			type="text"
			class="form-control"
			:class="computedClass"
			:placeholder="placeholder"
			:value="computedModelValue"
			:name="name"
			:disabled="disabled"
			:readonly="readonly"
			@input="updateValue"
			@keypress="validateKeypress"
			@keyup="validateKeyup"
		/>
		<span v-show="textEnd !== ''" class="btn btn-sm btn-icon position-absolute icon-end">
			{{ textEnd }}
		</span>
		<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
	</div>
</template>

<style scoped>
.icon-end {
	right: 0;
	margin: -34px 10px 0px 0px;
	color: #181c32;
	font-size: 13px;
	font-weight: 400;
	width: auto !important;
	text-align: end;
	display: block;
}
</style>
