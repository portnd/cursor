<script setup lang="ts">
import { useConditionStore } from "../store"

defineProps({
	mapShow: {
		type: Boolean,
	},
})

const store = useConditionStore()

watch(
	() => store.image.imageID,
	(newImageID, oldImageID) => {
		if (newImageID !== oldImageID) {
			const matchingChild = store.conditionDataTable
				.flatMap((parent) => parent.child || [])
				.find((child) => child.id === store.image.imageID)

			if (matchingChild) {
				store.image.path = matchingChild.image
			}
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
	<div class="card card-rounded p-5 mt-5" :class="mapShow ? 'd-none' : ''">
		<div class="row align-items-center">
			<div class="col col-title">
				<h6 class="fw-normal fs-5 mt-0">กล้องหน้า</h6>
			</div>
			<div class="col col-image text-end">
				<button
					type="button"
					class="btn btn-outline btn-sm me-2 tools"
					:class="{ active: store.playBtn.play }"
					@click="store.play(1, 'play')"
				>
					<i class="fi fi-sr-play pe-0 fs-7"></i>
				</button>
				<button
					type="button"
					class="btn btn-outline btn-sm fw-semibold me-2 tools"
					:class="{ active: store.playBtn.speed2x }"
					@click="store.play(2, 'speed2x')"
				>
					2X
				</button>
				<button
					type="button"
					class="btn btn-outline btn-sm fw-semibold me-2 tools"
					:class="{ active: store.playBtn.speed4x }"
					@click="store.play(4, 'speed4x')"
				>
					4X
				</button>
				<button
					type="button"
					class="btn btn-outline btn-sm me-2"
					:class="{ active: store.playBtn.stop }"
					@click="store.stop('stop')"
				>
					<i class="fi fi-sr-stop pe-0 fs-7"></i>
				</button>
				<button
					type="button"
					class="btn btn-outline btn-sm tools"
					:class="{ active: store.playBtn.pause }"
					@click="store.pause('pause')"
				>
					<i class="fi fi-sr-pause pe-0 fs-7"></i>
				</button>
			</div>
		</div>
		<div class="row align-items-center mt-3 not-found mt-6">
			<div class="col-12 h-100 d-flex justify-content-center">
				<template v-if="store.image.path === '' || noImage">
					<VNotFound height="100%" :is-not-shadow="true" message="ไม่พบข้อมูล" />
				</template>
				<template v-else>
					<img v-viewer :src="`${store.image.path}`" class="img cursor-pointer" @error="(e) => handleImageError(e)" />
				</template>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.tools {
	height: 38px;
}
img {
	width: 90%;
	// max-height: calc(38vh - 56px);
}
.btn {
	width: 45px;
	height: 38px;
	padding: 0;
	border: solid 1px var(--kt-primary) !important;
	color: var(--kt-text-gray-800) !important;
	&:hover {
		background-color: var(--kt-primary) !important;
	}
}
.btn.active {
	background-color: var(--kt-primary) !important;
	color: var(--kt-text-gray-800) !important;
}

.col-title {
	flex: 0 0 auto;
	width: 25%;
	@media only screen and (max-width: 1200px) {
		width: 100%;
	}
	@media only screen and (max-width: 991px) {
		width: 25%;
	}
	@media only screen and (max-width: 445px) {
		width: 100%;
	}
}

.col-image {
	flex: 0 0 auto;
	width: 75%;
	@media only screen and (max-width: 1200px) {
		width: 100%;
	}
	@media only screen and (max-width: 991px) {
		width: 75%;
	}
	@media only screen and (max-width: 445px) {
		width: 100%;
	}
}

.not-found {
	min-height: 30vh;
}
</style>
