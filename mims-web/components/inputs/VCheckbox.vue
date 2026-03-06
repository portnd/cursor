<script setup lang="ts">
import { useField } from "vee-validate"
import { IOption } from "~~/core/shared/types/Option"

type TMode = "single" | "multiple"

const props = defineProps({
	modelValue: {
		type: [String, Number, Boolean, Array],
		default: "",
	},
	mode: {
		type: String as PropType<TMode>,
		default: "multiple",
	},
	name: {
		type: String,
		required: true,
	},
	option: {
		type: Object as PropType<IOption>,
		default: null,
	},
	options: {
		type: Array<IOption>,
		default: [],
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
const result = ref<Boolean>()

const emit = defineEmits(["update:modelValue", "click"])
const updateValue = (value: any) => {
	// กรณี mode ที่เป็น single
	if (props.mode === "single") {
		result.value = !(value === "" || value === undefined || value === false)
		emit("update:modelValue", result.value)
	} else if (result.value !== value) {
		// กันการ update ซ้ำซ้อน
		result.value = value
		emit("update:modelValue", result.value)
	}
}

const onClick = () => {
	emit("click", true)
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />
	<div>
		<template v-if="mode === 'single'">
			<div class="form-check form-check-custom me-3 mb-2 mt-1" :class="{ 'form-check-inline': inline }">
				<VField
					:id="name"
					class="form-check-input"
					type="checkbox"
					:name="name"
					:value="!result"
					:model-value="modelValue"
					:checked="modelValue === true ? true : false"
					:disabled="disabled ? disabled : option.disabled"
					:class="{ 'is-invalid': !meta.valid && meta.validated }"
					:style="{ backgroundColor: option.color, borderColor: option.color }"
					@update:model-value="updateValue"
					@click="onClick"
				/>
				<img v-show="option.image" :src="option.image?.src" :width="option.image?.width" class="ms-4 me-1" />
				<label :for="name" class="form-check-label">{{ option.label }}</label>
			</div>
		</template>
		<template v-else>
			<div
				v-for="(option, key) in options"
				:key="`${name}_${option.value}_${key}`"
				class="form-check form-check-custom me-3 mb-3"
				:class="{ 'form-check-inline': inline }"
			>
				<div class="d-flex">
					<VField
						:id="`${name}${option.value}`"
						as="input"
						class="form-check-input"
						type="checkbox"
						:name="`${name}[]`"
						:value="option.value"
						:model-value="modelValue"
						:disabled="disabled ? disabled : option.disabled"
						:class="{ 'is-invalid': !meta.valid && meta.validated }"
						:style="{
							backgroundColor: !option.isSquare ? option.color : '',
							borderColor: !option.isSquare ? option.color : '',
						}"
						@update:model-value="updateValue"
					/>
					<label class="d-block cursor-pointer" :for="`${name}${option.value}`">
						<div v-if="option.isSquare" class="square ms-4 me-1" :style="`background: ${option.color}`"></div>
						<img
							v-else-if="option.image && option.image.src !== ''"
							:src="option.image?.src"
							:width="option.image?.width"
							class="ms-4 me-1"
						/>
					</label>
					<label :for="`${name}${option.value}`" class="form-check-label">{{ option.label }}</label>
				</div>
			</div>
		</template>
	</div>
	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style scoped lang="scss">
.form-check-input {
	border-radius: 0.25rem !important;
	margin-top: 1px;
	width: 1.5rem !important;
	height: 1.5rem !important;
}
.form-check-input:not(:checked) {
	border: 1px solid var(--kt-gray-300) !important;
	background-color: var(--kt-form-check-input-bg) !important;
}
img {
	margin-top: -5px;
}
.square {
	width: 1.8rem;
	height: 1.8rem;
	// border-radius: 0.25rem;
	margin-top: -1px;
}
</style>
