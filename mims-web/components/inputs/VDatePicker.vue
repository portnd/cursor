<script setup lang="ts">
import { useField } from "vee-validate"

const props = defineProps({
	modelValue: {
		type: [Date, Array],
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
	enableTimePicker: {
		type: Boolean,
		default: false,
	},
	range: {
		type: Boolean,
		default: false,
	},
	minRange: {
		type: [String, Number],
		default: "",
	},
	maxRange: {
		type: [String, Number],
		default: "",
	},
	monthPicker: {
		type: Boolean,
		default: false,
	},
	timePicker: {
		type: Boolean,
		default: false,
	},
	yearPicker: {
		type: Boolean,
		default: false,
	},
	weekPicker: {
		type: Boolean,
		default: false,
	},
	textInput: {
		type: Boolean,
		default: false,
	},
	minDate: {
		type: Date,
	},
	maxDate: {
		type: Date,
		default: new Date(),
	},
	teleportCenter: {
		type: Boolean,
		default: false,
	},
})

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const updateValue = (event: any) => {
	emit("update:modelValue", event)
}

const buddhistFormat = (date: Date): String => {
	if (props.yearPicker) {
		return buddhistFormatDate(date, "yyyy")
	} else if (props.monthPicker) {
		return buddhistFormatDate(date, "mmm yyyy")
	} else if (props.timePicker) {
		return buddhistFormatDate(date, "HH:mm")
	} else {
		return buddhistFormatDate(date, "dd mmm yyyy")
	}
}

const format = (data: any) => {
	// กรณีใช้  Mode: Range
	if (data instanceof Array) {
		const dates = data

		const ranges: Array<String> = []
		dates.forEach((date: Date) => {
			ranges.push(buddhistFormat(date))
		})
		return ranges.join(" - ")
	} else {
		// กรณีใช้  Mode: single
		return buddhistFormat(data)
	}
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />
	<div>
		<VField style="display: none" :name="name" :model-value="modelValue" />
		<VueDatepicker
			:name="name"
			class="form-control"
			auto-apply
			locale="th"
			:format="format"
			:model-value="modelValue"
			:placeholder="placeholder"
			:enable-time-picker="enableTimePicker"
			:disabled="disabled"
			:readonly="readonly"
			:range="range"
			:min-range="minRange"
			:max-range="maxRange"
			:month-picker="monthPicker"
			:time-picker="timePicker"
			:year-picker="yearPicker"
			:week-picker="weekPicker"
			:text-input="textInput"
			:min-date="minDate"
			:max-date="maxDate"
			:teleport-center="teleportCenter"
			:class="{ 'is-invalid': !meta.valid && meta.validated }"
			@update:model-value="updateValue"
		>
			<template #year="{ year }">
				{{ year + 543 }}
			</template>
			<template #year-overlay-value="{ value }">
				{{ value + 543 }}
			</template>
			<template #input-icon>
				<i class="fi fi-rr-calendar fs-2 text-gray-500"></i>
			</template>
		</VueDatepicker>
	</div>
	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style>
.dp__theme_light {
	--dp-background-color: #fff;
	--dp-border-color: var(--kt-input-bg);
	--dp-border-color-hover: var(--kt-input-bg);
	--dp-menu-border-color: var(--kt-input-bg);
	--dp-primary-color: var(--kt-primary);
	--dp-icon-color: var(--kt-text-gray-700);
	--dp-disabled-color: var(--kt-input-bg);
}
.dp__input {
	font-size: var(--kt-input-font-size);
	background-color: var(--kt-input-bg);
}
.dp__input_icons {
	stroke-width: 2px;
}
.dp__input_icon {
	top: 65%;
}
.dp__clear_icon {
	top: 55%;
}
.dp__clear_icon {
	right: -10px;
}
.dp__input_icon_pad {
	padding: 3px 20px 0 30px;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 100%;
}
.dp__menu {
	padding: 5px;
}
.dp__month_year_row {
	font-size: 1.1rem;
}
.dp__time_display {
	letter-spacing: 2px;
}
.dp__overlay_cell {
	font-size: 1rem;
}
.dp__calendar_header_item {
	font-weight: 600;
}
.dp__input:disabled {
	color: var(--kt-text-gray-400);
	background-color: var(--kt-input-disabled-bg);
	cursor: no-drop;
}
.dp__overlay_cell_active {
	color: var(--kt-text-gray-800);
}

.dp__overlay_cell_disabled {
	color: var(--ms-placeholder-color, #9ca3af) !important;
}

.dp__overlay_cell_disabled:hover {
	color: var(--ms-placeholder-color, #9ca3af) !important;
}
</style>
