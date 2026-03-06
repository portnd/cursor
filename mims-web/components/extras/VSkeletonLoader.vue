<script setup lang="ts">
const props = defineProps({
	loading: {
		type: Boolean,
		default: true,
	},
})

const isLoading = ref(props.loading)
let hideTimer: ReturnType<typeof setTimeout> | null = null

watch(
	() => props.loading,
	(newValue) => {
		if (newValue === true) {
			if (hideTimer) {
				clearTimeout(hideTimer)
				hideTimer = null
			}
			isLoading.value = true
		} else {
			hideTimer = setTimeout(() => {
				isLoading.value = false
				hideTimer = null
			}, 500)
		}
	}
)
</script>

<template>
	<div :class="{ 'skeleton-loading': isLoading }">
		<slot></slot>
	</div>
</template>

<style scoped></style>
