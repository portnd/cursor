<script setup lang="ts">
defineProps({
	title: {
		type: String,
		default: "",
	},
	message: {
		type: String,
		default: "",
	},
	type: {
		type: String,
		default: "",
	},
})
</script>

<template>
	<div>
		<h5 class="title">{{ title }}</h5>
		<p class="message">{{ message }}</p>
	</div>
</template>

<style scoped>
.title {
	font-size: 1.2rem;
	color: #fff;
	margin-bottom: 0.15rem;
}
.message {
	margin-bottom: 0.15rem;
}
</style>
