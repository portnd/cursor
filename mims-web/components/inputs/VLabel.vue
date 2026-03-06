<script setup lang="ts">
defineProps({
	label: {
		type: String,
		required: true,
	},
	required: {
		type: Boolean,
		default: false,
	},
})
</script>

<template>
	<label class="form-text text-gray-800">
		{{ label }}
		<span v-show="required" class="required"></span>
	</label>
</template>

<style scoped></style>
