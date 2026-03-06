<script setup lang="ts">
import { useDashboardAssetStore } from "../store"

const store = useDashboardAssetStore()

const roadId = ref([])

const props = defineProps({
	assetType: {
		type: String,
		default: "IN",
	},
})

onMounted(async () => {
	await store.getRoadsTree()
	await store.createRoadsOptions()
})

const searchData = async () => {
	store.roadId = roadId.value.join(",")
	await store.get(props.assetType)
	store.clearMarker()
	if (store.assetId !== "") {
		setTimeout(async () => {
			await store.getLocation()
			if (store.asset) {
				store.setLocation()
			}
		}, 400)
	}
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
