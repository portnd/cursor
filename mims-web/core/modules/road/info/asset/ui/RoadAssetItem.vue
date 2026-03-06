<script setup lang="ts">
import { IRoadsAssetItem } from "../infrastructure"

const route = useRoute()
const roadId = route.params.roadId

const props = defineProps({
	dataList: {
		type: Array as PropType<IRoadsAssetItem[]>,
		request: true,
		default: [] as PropType<IRoadsAssetItem[]>,
	},
	assetType: {
		type: String,
	},
	loading: {
		type: Boolean,
		required: true,
	},
})

const path = computed(() => {
	return props.assetType === "assetin" ? "in-asset" : "out-asset"
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
			}, 1000)
		}
	}
)
</script>

<template>
	<div class="accordion mt-5">
		<template v-if="isLoading">
			<div v-for="n in 6" :key="n" class="mb-5">
				<VSkeletonLoader>
					<div class="card d-flex flex-row py-1 align-items-center justify-content-between w-100 px-4">
						<div class="w-50 h-50">
							<p class="mb-0">ชื่อสินทรัพย์</p>
						</div>
						<div class="accordion-button collapsed fs-4 py-4 w-25 h-50"></div>
					</div>
				</VSkeletonLoader>
			</div>
		</template>
		<template v-else-if="dataList.length > 0">
			<div v-for="(item, index) in dataList" :key="index" class="accordion-item overflow-hidden shadow rounded-2 mb-5">
				<h2 :id="`assetIn${index}`" class="accordion-header">
					<button
						class="accordion-button fs-5 p-5 fw-semibold"
						:class="index !== 0 && 'collapsed'"
						type="button"
						data-bs-toggle="collapse"
						:data-bs-target="`#assetInCollapse${index}`"
						:aria-expanded="index === 0 ? 'true' : 'false'"
						:aria-controls="`assetInCollapse${index}`"
					>
						{{ item.name }}
					</button>
				</h2>
				<div
					:id="`assetInCollapse${index}`"
					class="accordion-collapse collapse"
					:class="index === 0 && 'show'"
					:aria-labelledby="`assetIn${index}`"
				>
					<div class="accordion-body pt-1 pb-4">
						<div v-for="(items, index) in item.items" :key="index" class="row">
							<div class="col-md-12 mb-3">
								<NuxtLink :to="`/roads/${roadId}/${path}/${items.id}`" class="text-black">
									<div class="rounded-3 border border-1 p-4 bg-hover-gray-200">{{ items.name }}</div>
								</NuxtLink>
							</div>
						</div>
					</div>
				</div>
			</div>
		</template>
		<template v-else>
			<VNotFound height="65vh" />
		</template>
	</div>
</template>

<style scoped>
.accordion-button:not(.collapsed) {
	background-color: #fff;
	color: var(--kt-accordion-color);
	box-shadow: none;
}
.accordion-collapse .border {
	border: 1px solid #e0e0e0 !important;
}
</style>
