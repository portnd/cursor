<script setup lang="ts">
type TType = "submit" | "button" | "reset"

defineProps({
	loading: {
		type: Boolean,
		default: false,
	},
	label: {
		type: String,
		required: true,
	},
	type: {
		type: String as PropType<TType>,
		default: "submit",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
})
</script>

<template>
	<button
		:type="type"
		class="btn btn-primary px-8 py-3"
		:data-kt-indicator="loading"
		:disabled="loading === true || disabled === true"
	>
		<span class="indicator-label">{{ label }}</span>
		<div class="indicator-progress">
			โปรดรอสักครู่...
			<span class="spinner-border spinner-border-sm align-middle ms-2"></span>
		</div>
	</button>
</template>

<style scoped></style>
