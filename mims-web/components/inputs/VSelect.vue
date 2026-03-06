<script setup lang="ts">
import { useField } from "vee-validate"
import { IOption } from "~~/core/shared/types/Option"

type TMode = "single" | "multiple" | "tags"

const props = defineProps({
	modelValue: {
		type: [Array<String>, Array<Number>, String, Number],
		default: "",
	},
	mode: {
		type: String as PropType<TMode>,
		default: "single",
	},
	options: {
		type: (Array as PropType<IOption[]>) || [],
		required: true,
	},
	searchable: {
		type: Boolean,
		default: false,
	},
	closeOnSelect: {
		type: Boolean,
		default: false,
	},
	placeholder: {
		type: String,
		default: "ทั้งหมด",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
	label: {
		type: String,
		default: "",
	},
	name: {
		type: String,
		required: true,
	},
	required: {
		type: Boolean,
		default: false,
	},
	noOptionsText: {
		type: String,
		default: "ไม่มีรายการ",
	},
	noResultsText: {
		type: String,
		default: "ไม่พบข้อมูล",
	},
	hideSelected: {
		type: Boolean,
		default: false,
	},
	canClear: {
		type: Boolean,
		default: true,
	},
	canDeselect: {
		type: Boolean,
		default: true,
	},
	autoHeight: {
		type: Boolean,
		default: false,
	},
})

const { meta, resetField } = useField(props.name)

const emit = defineEmits(["update:modelValue"])
const updateValue = (event: any) => {
	emit("update:modelValue", event)
}

onUnmounted(() => {
	resetField({ value: "" })
})
</script>

<template>
	<VLabel v-show="label !== ''" :label="label" :required="required" />

	<VField style="display: none" :name="name" :model-value="modelValue" />
	<div class="min-height">
		<VueMultiselect
			:id="name"
			ref="multiselectRef"
			class="multiselect"
			:mode="mode"
			:value="modelValue"
			:placeholder="placeholder"
			:close-on-select="mode !== 'single' ? closeOnSelect : true"
			:searchable="searchable"
			:options="options"
			:disabled="disabled"
			:no-options-text="noOptionsText"
			:no-results-text="noResultsText"
			:class="{ 'is-invalid': !meta.valid && meta.validated, 'auto-height': autoHeight }"
			:hide-selected="hideSelected"
			:can-clear="canClear"
			:can-deselect="canDeselect"
			:append-to-body="true"
			label="label"
			@input="updateValue"
		>
			<template v-if="mode !== 'single'" #multiplelabel="{ values }">
				<div v-if="values.length > 1" class="multiselect-multiple-label">
					<span class="multiselect-tag mb-0">
						{{ `${values[0].label}` }}
						<b v-if="values[0].description" class="description">{{ values[0].description }}</b>
					</span>
					<span class="more">{{ `และอีก ${values.length - 1} รายการ` }}</span>
				</div>
				<div v-else class="multiselect-multiple-label">
					<span class="multiselect-tag mb-0">
						{{ `${values[0].label}` }}
						<b v-if="values[0].description" class="description">{{ values[0].description }}</b>
					</span>
				</div>
			</template>
			<template v-else #singlelabel="{ value }">
				<div class="multiselect-single-label">
					<span>
						<img
							v-if="value.image"
							:src="value.image"
							:width="value.image?.width ? value.image?.width : 25"
							height="25"
							class="character-label-icon me-2"
						/>
						{{ value.label }} <b v-if="value.description" class="description">{{ value.description }}</b>
					</span>
				</div>
			</template>

			<template #option="{ option, isSelected }">
				<div v-if="mode !== 'single'" class="d-flex">
					<i v-if="isSelected(option)" class="fi fi-sr-checkbox align-middle fs-4 lh-0 me-3 text-primary"></i>
					<i v-else class="fi fi-rr-square align-middle fs-4 lh-0 me-3 text-gray-400"></i>
					<label :for="option.id">{{ option.label }}</label>
				</div>
				<div v-else>
					<img
						v-if="option.image"
						:src="option.image"
						:width="option.image?.width ? option.image?.width : 25"
						height="25"
						class="me-1"
					/>
					{{ option.label }} <span v-if="option.description" class="description">{{ option.description }}</span>
				</div>
			</template>
		</VueMultiselect>
	</div>
	<VErrorMessage v-show="!meta.valid && meta.validated" class="invalid-feedback" :name="name" />
</template>

<style lang="scss">
.multiselect {
	padding: 0;
	height: var(--kt-input-height);
	// border: 1px solid var(--kt-gray-300);
	border: 1px solid #eaecf3;

	--ms-font-size: var(--kt-input-font-size);
	--ms-bg: var(--kt-input-bg);
	--ms-option-color-selected: var(--kt-black);
	--ms-option-bg-selected: var(--kt-primary-light);
	--ms-option-bg-selected-pointed: var(--kt-primary-light);
	--ms-option-color-selected-pointed: var(--kt-black);
	--ms-dropdown-radius: var(--kt-form-border-radius);
	--ms-radius: var(--kt-form-border-radius);
	--ms-max-height: 22rem;
	--ms-option-py: 0.85rem;
	--ms-option-px: 1rem;
	--ms-px: 0.85rem;
	--ms-option-bg-pointed: e;
	--ms-border-color-active: var(--kt-gray-300);
	--ms-dropdown-border-color: #f1f0f0;

	--ms-tag-font-size: 0.975rem;
	--ms-tag-bg: var(--kt-primary-light);
	--ms-tag-color: var(--kt-primary);
	--ms-tag-font-weight: 500;
	--ms-tag-radius: 8px;
	--ms-tag-py: 0.225rem;
	--ms-tag-py: 0.225rem;
	--ms-tag-px: 0.75rem;
	--ms-tag-my: 0.35rem;
}

.multiselect-option,
.multiselect-option span {
	font-weight: 400;
	color: var(--kt-input-color);
	&.description {
		color: var(--kt-primary);
		font-size: 0.95rem;
	}
}

.multiselect-caret {
	height: 1.5rem;
	width: 1.5rem;
	z-index: 2 !important;
}
.multiselect-caret {
	background-color: var(--ms-caret-color, var(--kt-gray-700));
}
.multiselect:has(.multiselect-wrapper > .multiselect-tags) {
	height: auto;
	min-height: var(--kt-input-height);
	padding: 5px 3px;
}
.multiselect-clear {
	padding: 0 0.25rem 0 0;
}
.multiselect-clear-icon {
	height: 1rem;
	width: 0.75rem;
}
.multiselect.is-active {
	box-shadow: none;
	border-bottom: 1px solid #e7e7e7;
}
.multiselect-tags-search {
	background-color: var(--kt-input-bg);
}
.multiselect.is-disabled {
	color: var(--kt-text-gray-400);
	background-color: var(--kt-input-disabled-bg);
	cursor: no-drop;
	.multiselect-wrapper {
		cursor: no-drop;
	}
}

.multiselect-single-label span {
	display: block;
	max-width: 100%;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	color: var(--kt-gray-800);
	.description {
		color: var(--kt-primary);
		font-size: 0.95rem;
		font-weight: 400 !important;
	}
}

.multiselect-multiple-label .more {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 100%;

	display: block;
	margin: 0px 0;
	padding: 2px 0;
	color: #bdbdbd;
	font-size: 12px;
	font-weight: 500;
}

.multiselect-option.is-selected {
	background-color: var(--kt-primary-light) !important;
}

.multiselect-option.is-selected.is-pointed {
	background-color: var(--kt-primary-light) !important;
}

// auto-height
.multiselect {
	&.auto-height {
		height: auto;
		min-height: var(--kt-input-height);
		.multiselect-multiple-label,
		.multiselect-placeholder,
		.multiselect-single-label {
			width: 100%;
			position: unset;
			padding: 8px 0px 8px 10px;
		}

		.multiselect-single-label span {
			overflow: unset;
			text-overflow: unset;
			white-space: unset;
		}

		.multiselect-wrapper {
			justify-content: flex-end;
		}
	}
}

.multiselect-multiple-label .multiselect-tag {
	padding: var(--ms-tag-py) var(--ms-tag-px);
}

/* begin::Validate */
.multiselect.is-invalid {
	border: 1px solid var(--kt-danger);
}
/* end::Validate */

.customize-rows-per-page {
	width: 70px;
	.multiselect {
		height: 36px;
		background: #ffffff;
		border: 1px solid #ffffff;
		.multiselect-single-label span {
			overflow: unset;
		}
		.multiselect-caret {
			margin: 0 0.5rem 0 0;
		}
		.multiselect-option {
			padding: 0.75rem;
		}
	}
}

.multiselect-dropdown {
	--ms-dropdown-radius: var(--kt-form-border-radius);
	--ms-dropdown-border-color: var(--kt-gray-300);
	--ms-option-color-selected: var(--kt-black);
	--ms-option-bg-selected: var(--kt-primary-light);
	--ms-option-bg-selected-pointed: var(--kt-primary-light);
	--ms-option-color-selected-pointed: var(--kt-black);
	--ms-radius: var(--kt-form-border-radius);
	--ms-max-height: 18rem;
	--ms-option-py: 0.85rem;
	--ms-option-px: 1rem;
	--ms-option-bg-pointed: #f5f5f5;

	&.is-top {
		border-bottom: 0px;
	}

	&[data-popper-placement="bottom"] {
		border-top: 0px;
	}

	z-index: 9999;
	max-height: 18rem;

	&::-webkit-scrollbar {
		width: 4px;
	}
	&::-webkit-scrollbar-thumb {
		background: #b5b5b5 !important;
	}
	&::-webkit-scrollbar-thumb:hover {
		background: #555 !important;
	}
	&::-webkit-scrollbar-track {
		background: rgba(218, 218, 218, 0.6) !important;
	}
}
</style>
