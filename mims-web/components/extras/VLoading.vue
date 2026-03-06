<script setup lang="ts">
const props = defineProps({
	loading: {
		type: Boolean,
		default: false,
	},
	marginTop: {
		type: String,
		default: "2",
	},
	float: {
		type: String,
		default: "start",
	},
})

const isLoading = ref(false)
watch(
	() => props.loading,
	(newValue) => {
		if (newValue === true) {
			isLoading.value = newValue
		} else {
			setTimeout(() => {
				isLoading.value = newValue
			}, 500)
		}
	}
)
</script>

<template>
	<div
		class="fs-6 float-none"
		:class="`mt-${marginTop} float-md-${float}`"
		:style="{ visibility: !loading ? 'hidden' : 'unset' }"
	>
		โปรดรอสักครู่...
		<span class="spinner-border spinner-border-sm align-middle ms-2"></span>
	</div>
</template>

<style scoped></style>
