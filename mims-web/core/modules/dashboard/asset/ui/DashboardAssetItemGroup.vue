<script setup lang="ts">
import { useDashboardAssetStore } from "../store"

const props = defineProps({
	assetType: {
		type: String,
		default: "IN",
	},
})

const store = useDashboardAssetStore()
const assetID = ref<Array<Number>>([])

watch(assetID, () => {
	store.assetId = assetID.value.join(",")
	store.clearMarker()
	if (store.assetId !== "") {
		setTimeout(async () => {
			await store.getLocation()
			if (store.asset) {
				store.setLocation()
			}
		}, 400)
	}
})

watch(props, async () => {
	store.clearMarker()
	await store.get(props.assetType)
	store.assetId = ""
	assetID.value = []
})

onMounted(async () => {
	await store.get(props.assetType)
})
</script>

<template>
	<div class="tab-content">
		<div
			:id="`asset${props.assetType}`"
			class="tab-pane fade active show"
			role="tabpanel"
			:aria-labelledby="`asset${props.assetType}-tab`"
		>
			<div class="accordion">
				<template v-if="store.loading">
					<div v-for="n in 6" :key="n" class="mb-5">
						<VSkeletonLoader>
							<div class="card d-flex flex-row py-1 align-items-center justify-content-between w-100 px-4">
								<div class="w-50 h-50">
									<p class="mb-0">ชื่อสายทาง</p>
								</div>
								<div class="accordion-button collapsed fs-4 py-4 w-25 h-50"></div>
							</div>
						</VSkeletonLoader>
					</div>
				</template>
				<template v-else-if="store.data.length > 0">
					<div v-for="(item, i) in store.data" :key="i" class="accordion-item overflow-hidden rounded-2 shadow mb-5">
						<h2 :id="`asset_header${props.assetType}_${i}`" class="accordion-header">
							<button
								class="accordion-button fs-5 p-5 fw-semibold"
								:class="{ collapsed: i > 0 }"
								type="button"
								data-bs-toggle="collapse"
								:data-bs-target="`#asset${props.assetType}Collapse-${i}`"
								:aria-controls="`asset${props.assetType}Collapse-${i}`"
							>
								{{ item.asset_group.name }}
							</button>
						</h2>
						<div
							:id="`asset${props.assetType}Collapse-${i}`"
							class="accordion-collapse collapse"
							:class="{ show: i <= 0 }"
							:aria-labelledby="`#asset_header${props.assetType}_${i}`"
						>
							<div class="accordion-body pt-1 pb-4">
								<div v-if="item.asset_list.length > 0" class="row">
									<div v-for="(list, j) in item.asset_list" :key="j" class="col-md-4 col-6 mb-5">
										<div v-if="list.asset.default_color !== ''" class="d-flex">
											<VCheckbox
												v-model="assetID"
												:options="[
													{
														color: list.asset.default_color,
														label: list.asset.name + ' (' + list.value + ')',
														value: list.asset.id,
														isSquare: true,
													},
												]"
												:name="`asset-${list.asset.id}`"
											/>
										</div>
										<VCheckbox
											v-else-if="list.asset.default_color === ''"
											v-model="assetID"
											:options="[
												{
													image: {
														src: `${list.asset.default_icon_url}`,
														width: 35,
													},
													label: list.asset.name + ' (' + list.value + ')',
													value: list.asset.id,
												},
											]"
											:name="`asset-${list.asset.id}`"
										/>
										<VCheckbox
											v-else
											v-model="assetID"
											:options="[
												{
													image: {
														src: `/images/icons/png/location-pin.png`,
														width: 35,
													},
												},
											]"
											:name="`asset-${list.asset.id}`"
										/>
									</div>
								</div>
								<div v-else class="text-center">
									<span>ไม่พบข้อมูล</span>
								</div>
							</div>
						</div>
					</div>
				</template>
				<template v-else>
					<VNotFound height="65vh" />
				</template>
			</div>
		</div>
	</div>
</template>

<style scoped>
.accordion-button:not(.collapsed) {
	background-color: #fff;
}
.accordion-button:not(.collapsed) {
	color: var(--kt-accordion-color);
}
</style>
