<script setup lang="ts">
defineProps({
	message: {
		type: String,
		default: "ไม่พบข้อมูล",
	},
	height: {
		type: String,
		default: "30vh",
	},
	isNotShadow: {
		type: Boolean,
		default: false,
	},
})
</script>

<template>
	<div class="card" :class="{ 'shadow-none': isNotShadow }">
		<div class="text-center not-found align-middle">
			<inline-svg src="/images/icons/svg/file-circle-info.svg" :width="50" :height="50" />
			<br />
			<div class="fs-5 mt-2 text-gray-600">{{ message }}</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.card {
	display: flex;
	align-items: center;
	justify-content: center;
	height: v-bind(height);
}

.not-found {
	margin-top: -30px;
}

@media (max-width: 991px) {
	.card {
		height: 55vh;
	}
}
</style>
