<script setup lang="ts">
defineProps({
	label: {
		type: String,
		default: "เพิ่มรายการ",
	},
})
</script>

<template>
	<a type="button" class="btn btn-primary rounded-4 px-8 py-3 fw-semibold lh-xl">
		{{ label }}
	</a>
</template>

<style scoped></style>
