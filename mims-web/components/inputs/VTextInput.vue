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
	textEnd: {
		type: String,
		default: "",
	},
	validateEnglish: {
		type: Boolean,
		default: false,
	},
	validateNumber: {
		type: Boolean,
		default: false,
	},
	paste: {
		type: Boolean,
		default: false,
	},
	validatePosition: {
		type: Boolean,
		default: false,
	},
})

const { meta, resetField, errorMessage } = useField(props.name)
const selectStart = ref()
const selectEnd = ref()

const emit = defineEmits(["update:modelValue"])
const updateValue = (event: any) => {
	emit("update:modelValue", event.target.value)
}

const onKeyPressed = (event: any) => {
	if (props.validateNumber) {
		const keyCode = event.keyCode
		const char = String.fromCharCode(keyCode)
		const numberRegex = /^[0-9_]+$/
		console.log("keyCode =", keyCode)

		if (keyCode >= 48 && keyCode <= 57) {
			// Get the lowercase version of the pressed key
			const lowercaseKey = char.toLowerCase()

			// Prevent the default action (inserting the uppercase letter)
			event.preventDefault()

			// Update the modelValue with the lowercase key
			emit("update:modelValue", event.target.value + lowercaseKey)
		}

		// Check if the character matches the English regex pattern
		if (!char.match(numberRegex)) {
			event.preventDefault()
		}
	} else if (props.validateEnglish) {
		const keyCode = event.keyCode
		const char = String.fromCharCode(keyCode)
		const englishRegex = /^[a-z0-9_]+$/

		if (keyCode >= 65 && keyCode <= 90) {
			// Get the lowercase version of the pressed key
			const lowercaseKey = char.toLowerCase()

			// Prevent the default action (inserting the uppercase letter)
			event.preventDefault()

			// Update the modelValue with the lowercase key
			emit("update:modelValue", event.target.value + lowercaseKey)
		}

		// Check if the character matches the English regex pattern
		if (!char.match(englishRegex)) {
			event.preventDefault()
		}
	}
}

// Prevent Control + V (paste)
const handlePaste = (event: any) => {
	if (props.paste) {
		const Regex1 = /^[0-9_]+$/
		const Regex2 = /^[a-zA-Z]+$/
		const clipboard = event.clipboardData.getData("text")
		const resultListChar = []

		for (let i = 0; i < clipboard?.length; i++) {
			const char = clipboard[i]

			if (!char.match(Regex1) && !char.match(Regex2)) {
				event.preventDefault()
			} else if (char.match(Regex1)) {
				resultListChar.push(char)
			} else if (char.match(Regex2)) {
				resultListChar.push(char.toLowerCase())
			}
		}

		event.preventDefault()
		const totalSelect = selectStart.value + selectEnd.value
		if (selectEnd.value === props.modelValue?.length && selectStart.value === 0) {
			emit("update:modelValue", resultListChar.join(""))
		} else if (totalSelect === 0 || props.modelValue?.length === 0 || props.modelValue === null) {
			emit("update:modelValue", event.target.value + resultListChar.join(""))
		} else {
			const selectedText = props.modelValue?.substring(selectStart.value, selectEnd.value)
			const result = [props.modelValue[selectStart.value - 1], selectedText, props.modelValue[selectEnd.value]]
			const index = ref<number>()
			for (let i = 0; i < result?.length; i++) {
				if (result[i]?.length > 1) {
					index.value = i
				}
			}
			if (index) {
				result[Number(index.value)] = resultListChar.join("")
			}
			emit("update:modelValue", result.join(""))
		}
		selectStart.value = 0
		selectEnd.value = 0
	}
}

const handleSelect = () => {
	const input = document.getElementById(props.name) as HTMLInputElement
	if (input) {
		selectStart.value = input.selectionStart
		selectEnd.value = input.selectionEnd
	}
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<div class="position-relative">
		<VLabel v-show="label !== ''" :label="label" :required="required" />
		<VField
			:id="name"
			type="text"
			:name="name"
			class="form-control"
			:placeholder="placeholder"
			:value="modelValue"
			:model-value="modelValue"
			:class="{ 'is-invalid': !meta.valid && meta.validated }"
			:disabled="disabled"
			:readonly="readonly"
			autocomplete="off"
			@input="updateValue"
			@keypress="onKeyPressed"
			@paste="handlePaste"
			@select="handleSelect"
		/>
		<span v-show="textEnd !== ''" class="btn btn-sm btn-icon position-absolute icon-end">
			{{ textEnd }}
		</span>
		<span
			v-if="!meta.valid && meta.validated"
			class="invalid-feedback"
			:class="{ 'position-absolute': validatePosition }"
			:name="name"
			v-html="errorMessage"
		></span>
		<!-- <VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" /> -->
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
