<script setup lang="ts">
import { api as viewerApi } from "v-viewer"
import HistoryProjectTab from "../HistoryProjectTab.vue"
import { useMaintenanceHistoryAttachmentStore } from "../../../store/MaintenanceHistoryAttachmentStore"

const route = useRoute()
const id = Number(route.params.id)
const store = useMaintenanceHistoryAttachmentStore()
useStoreLifecycle(store)

const collapsed = ref(false)
onMounted(() => {
	useIntervalFn(() => {
		const body = document.querySelector("body")
		if (body && body.classList.contains("longomap-collapsed")) {
			collapsed.value = true
		} else {
			collapsed.value = false
		}
	}, 500)
})

watchEffect(() => {
	if (store.params) {
		store.getMaintenanceHistoryFile(id)
	}
})

const showImage = (path: string) => {
	viewerApi({
		images: [path],
	})
}

</script>

<template>
	<div class="col-12">
		<div class="card p-5">
			<div class="row">
				<HistoryProjectTab />
				<div :class="collapsed ? 'col-md-3 col-12' : 'col-md-4 col-12'">
					<VSelect
						v-model="store.params.file_type"
						:options="[
							{ label: 'เอกสาร', value: 'doc' },
							{ label: 'รูปภาพ', value: 'pic' },
						]"
						:can-deselect="false"
						label="ประเภท"
						name="file_type"
					/>
				</div>
				<div :class="collapsed ? 'col-md-3 col-12' : 'col-md-4 col-12'">
					<VSelect
						v-model="store.params.order"
						:options="[
							{ label: 'วันที่', value: 'date' },
							{ label: 'ชื่อ', value: 'name' },
						]"
						:can-deselect="false"
						:can-clear="false"
						label="จัดเรียงตาม"
						name="sort"
					/>
				</div>
				<div :class="collapsed ? 'col-12 col-md-6 text-end' : 'col-12 col-md-4 text-end'">
					<VLoading :loading="store.loading" float="end" class="mt-0 mt-md-10" />
				</div>
			</div>

			<VSkeletonLoader :loading="store.loading">
				<div class="row mt-0 mt-md-3">
					<div
						v-for="item of store.attachmentData"
						:key="item.id"
						:class="
							collapsed ? 'col-xxl-3 col-xl-4 col-md-4 col-sm-6 col-12' : 'col-xxl-4 col-xl-6 col-md-6 col-sm-6 col-12'
						"
					>
						<div class="row">
							<NuxtLink
								:to="['png', 'jpg', 'jpeg'].includes(item.file_type) ? 'javascript:void(0)' : item.path"
								:target="['png', 'jpg', 'jpeg'].includes(item.file_type) ? '' : '_blank'"
								@click=";['png', 'jpg', 'jpeg'].includes(item.file_type) ? showImage(item.path) : null"
							>
								<div class="col-12 mt-5">
									<div class="card card-file px-3 py-4">
										<div class="row">
											<div class="col-3 pe-1 text-center align-self-center">
												<img
													v-if="['png', 'jpg', 'jpeg'].includes(item.file_type)"
													:src="item.path"
													class="cursor-pointer"
													width="50"
												/>
												<img v-else-if="item.file_type === 'pdf'" :src="`/images/files/pdf.png`" width="50" />
												<img v-else-if="item.file_type === 'docx'" :src="`/images/files/docx.png`" width="50" />
												<img v-else-if="item.file_type === 'xlsx'" :src="`/images/files/xlsx.png`" width="50" />
												<img v-else-if="item.file_type === 'dwg'" :src="`/images/files/dwg.png`" width="50" />
												<img v-else :to="item.path" target="_blank" :src="`/images/files/unknown.png`" width="50" />
											</div>
											<div class="col align-self-center cursor-pointer">
												<h6 class="fw-semibold fs-6 filename cursor-pointer">{{ item.file_name }}</h6>
												<label class="fw-normal text-gray-500 fs-7 mb-1 cursor-pointer">
													วันที่ {{ buddhistFormatDate(item.created_at, "dd mmm yyyy HH:ii น.") }}
												</label>
												<label class="fw-normal text-gray-500 fs-7 cursor-pointer">
													อัปโหลดโดย {{ item.created_by.firstname }} {{ item.created_by.lastname }}
												</label>
											</div>
										</div>
									</div>
								</div>
							</NuxtLink>
						</div>
					</div>
				</div>
			</VSkeletonLoader>
		</div>
		<VNotFound v-show="store.attachmentData.length === 0" class="mt-5" height="65dvh" />
	</div>
</template>

<style scoped>
.card-file {
	border: 1px solid var(--kt-gray-300) !important;
	box-shadow: none !important;
}

.filename {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 150px;
	@media (max-width: 576px) {
		max-width: 100%;
	}
}
</style>
