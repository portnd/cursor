<script setup lang="ts">
import { useDashboardSurfaceStore } from "../store"

const store = useDashboardSurfaceStore()
useStoreLifecycle(store)
const refBar = ref()

const handleSurface = (id: number) => {
	const index = store.surfaceArray.indexOf(id)
	if (index === -1) {
		store.surfaceArray.push(id)
	} else {
		store.surfaceArray.splice(index, 1)
	}
	store.surfaceArray.sort((a, b) => a - b)
	store.colors = store.surfaceColors.flatMap((item, key) => {
		return store.surfaceArray
			.map((id) => {
				if (key === id) {
					return item
				}
				return undefined
			})
			.filter((item: any) => item !== undefined)
	})
	store.createLine()
}

</script>

<template>
	<div class="row mb-5">
		<div class="col-12">
			<div class="card h-100 text-center p-5">
				<ClientOnly>
					<apexchart
						ref="refBar"
						:type="'pie'"
						height="320"
						:options="store.barOptions()"
						:series="store.barSeries()"
					/>
				</ClientOnly>
				<div class="row justify-content-center">
					<div
						v-for="(value, index) of store.data.summary"
						:key="index"
						class="col-auto d-flex align-items-center cursor-pointer"
						:class="store.surfaceArray.includes(index) ? '' : 'selected'"
						@click="handleSurface(index)"
					>
						<div class="square my-2 me-2" :style="{ backgroundColor: store.surfaceColors[index] }"></div>
						<span>{{ value.surface.name }}</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.square {
	width: 20px;
	height: 20px;
	border-radius: 5px;
	display: inline-block;
}

.selected {
	opacity: 0.2;
}
</style>
