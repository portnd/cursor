<script setup lang="ts">
const props = defineProps({
	imageNo: {
		type: Number,
		default: 1,
	},
	imageData: {
		type: String,
		default: "",
	},
	mapShow: {
		type: Boolean,
	},
})

watch(
	() => props.imageData,
	(newValue, oldValue) => {
		if (newValue !== oldValue) {
			noImage.value = false
		}
	}
)

const noImage = ref(false)
const handleImageError = (e: any) => {
	if (e.type === "error") {
		noImage.value = true
	}
}
</script>

<template>
	<div class="card card-rounded p-4 mt-5" :class="mapShow ? 'd-none' : ''">
		<div class="row align-items-center">
			<template v-if="imageData === '' || noImage">
				<VNotFound height="33dvh" :is-not-shadow="true" />
			</template>
			<template v-else>
				<div v-viewer class="col-12 cursor-pointer justify-content-center text-center">
					<img :src="`${imageData}`" class="img" @error="(e) => handleImageError(e)" />
				</div>
			</template>
		</div>
	</div>
</template>

<style scoped lang="scss">
img {
	width: 100%;
	// max-height: calc(38vh - 56px);
	// @media (max-width: 991px) {
	// 	max-height: 100%;
	// }
}
</style>
