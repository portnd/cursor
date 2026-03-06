<script setup lang="ts">
defineProps({
	label: {
		type: String,
		default: "ค้นหา",
	},
	disabled: {
		type: Boolean,
		default: false,
	},
})
</script>

<template>
	<button type="button" class="btn btn-primary rounded-4 px-8 py-3 mt-3 mt-md-0 fw-semibold" :disabled="disabled">
		{{ label }}
	</button>
</template>

<style scoped></style>
