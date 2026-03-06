<script setup lang="ts">
defineProps({
	label: {
		type: String,
		default: "ยกเลิก",
	},
})
</script>

<template>
	<button type="button" class="btn btn-light rounded-4 px-8 py-3 me-5 fw-semibold text-black">
		{{ label }}
	</button>
</template>

<style scoped>
/* .btn.btn-white:hover {
	border-color: var(--kt-light) !important;
	background-color: var(--kt-light) !important;
} */
</style>
