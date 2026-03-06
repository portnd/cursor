<script setup lang="ts">
// @ts-ignore
import { useCurrencyInput, CurrencyDisplay } from "vue-currency-input"
import { useField } from "vee-validate"

type TMode = "center" | "start" | "end"

interface IPrecision {
	min: number
	max: number
}

const props = defineProps({
	modelValue: {
		type: Number as PropType<number | null>,
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
		type: [Object, Number] as PropType<IPrecision | number>,
		default: 0,
	},
	min: {
		type: Number as PropType<number | null>,
		default: null,
	},
	max: {
		type: Number as PropType<number | null>,
		default: null,
	},
})

const { meta, resetField } = useField(props.name)

const computedClass = computed(() => {
	let className = ""
	if (!meta.valid && meta.validated) {
		className = "is-invalid"
	}

	className += ` text-${props.align}`
	return className
})

const valueRange = () => {
	let range = {}
	if (props.min !== null && props.max !== null) {
		range = {
			min: props.min,
			max: props.max,
		}
	} else if (props.min !== null) {
		range = {
			min: props.min,
		}
	} else if (props.max !== null) {
		range = {
			max: props.max,
		}
	}
	return range
}

const computedPrecision = () => {
	if (typeof props.precision === "number") {
		return {
			min: 0,
			max: props.precision,
		}
	}
	return props.precision
}

const { inputRef, formattedValue, setValue } = useCurrencyInput({
	currency: "THB",
	currencyDisplay: CurrencyDisplay.hidden,
	precision: computedPrecision(),
	hideCurrencySymbolOnFocus: true,
	hideGroupingSeparatorOnFocus: false,
	hideNegligibleDecimalDigitsOnFocus: true,
	autoDecimalDigits: false,
	useGrouping: true,
	accountingSign: false,
	valueRange: valueRange(),
})

watch(
	() => props.modelValue,
	(value) => {
		setValue(value)
	}
)

const emit = defineEmits(["update:modelValue", "change"])
const updateValue = (event: any) => {
	emit("update:modelValue", Number(event.target.value))
}

const computedModelValue = computed(() => {
	let value = props.modelValue === 0 ? "0" : props.modelValue
	value = value === null ? "" : value
	return value
})

onUnmounted(() => {
	resetField({ value: null })
})
</script>

<template>
	<div class="position-relative">
		<VLabel v-show="label !== ''" :label="label" :required="required" />

		<input
			:id="name"
			ref="inputRef"
			:value="formattedValue"
			type="text"
			class="form-control"
			:class="computedClass"
			:placeholder="placeholder"
			:disabled="disabled"
			:readonly="readonly"
			autocomplete="off"
		/>
		<VField
			class="d-none"
			type="text"
			:name="name"
			:value="computedModelValue"
			@input="updateValue"
			@change="updateValue"
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
