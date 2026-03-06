<script setup lang="ts">
import { useDashboardSurfaceStore } from "../store"

const store = useDashboardSurfaceStore()

const roadId = ref<Array<number>>([])

onMounted(async () => {
	await store.getRoadsTree()
	store.createRoadsOptions()
})

const searchData = async () => {
	store.roadId = roadId.value.join(",")
	await store.get()
	store.defaultLocation()
	store.createLine()
}
</script>

<template>
	<div class="row mb-5">
		<div class="col-12">
			<div class="card p-5 pt-2">
				<div class="row">
					<div class="col-md">
						<VTree
							v-model="roadId"
							placeholder="ทั้งหมด"
							label="สายทาง"
							:multiple="true"
							:limit="3"
							:options="store.roadsOptions"
							name="select"
							mode="LEAF_PRIORITY"
						/>
					</div>
					<div class="col-md-2 align-self-end">
						<BtnSearch :disabled="store.loading" @click="searchData()" />
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
