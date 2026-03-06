<script setup lang="ts">
import { useHRISPreviewStore } from "../store"

const store = useHRISPreviewStore()
useStoreLifecycle(store)

const onClickImport = async () => {
	const res = await store.import()
	if (res?.status === false) {
		useHandlerError(res.code, res.error, { showAlert: true })
	} else {
		useHandlerSuccess(res!.code, {
			showAlert: true,
			fn: () => {
				navigateTo("/settings/hris")
			},
		})
	}
}

onMounted(() => {
	store.get()
})

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<div class="row">
			<div class="col-12">
				<div class="card p-5">
					<div class="row">
						<div class="col-12 mt-3">
							<div class="row align-items-center">
								<div class="col-6">
									<h5 class="fw-semibold">ทางหลวง</h5>
								</div>
							</div>
							<div class="table-responsive mt-3">
								<table class="table customize-basic-table mb-0 text-truncate table-hover">
									<thead>
										<tr>
											<th class="text-center">ลำดับ</th>
											<th class="text-center">รหัสสายทาง</th>
											<th class="text-center">ชื่อสายทาง</th>
										</tr>
									</thead>
									<tbody>
										<tr v-for="(item, index) of store.data.road_group" :key="index" class="cursor-pointer">
											<td class="text-center">{{ index + 1 }}</td>
											<td class="text-center">{{ item.number }}</td>
											<td class="text-center">{{ item.name }}</td>
										</tr>
										<tr v-show="store.data.road_group.length === 0" class="text-center">
											<td colspan="8">ไม่พบข้อมูล</td>
										</tr>
									</tbody>
								</table>
							</div>
						</div>
					</div>
					<div class="row mt-5">
						<div class="col-12 mt-5">
							<div class="row align-items-center">
								<div class="col-6">
									<h5 class="fw-semibold">ตอนควบคุม</h5>
								</div>
							</div>
							<div class="table-responsive mt-3">
								<table class="table customize-basic-table mb-0 text-truncate table-hover">
									<thead>
										<tr>
											<th class="text-center">ลำดับ</th>
											<th class="text-center">รหัสสายทาง</th>
											<th class="text-center">รหัสตอนควบคุม</th>
											<th class="text-center">ชื่อตอนควบคุม (ไทย)</th>

											<th class="text-center">ชื่อตอนควบคุม (อังกฤษ)</th>
											<th class="text-center">กม. เริ่มต้น</th>
											<th class="text-center">กม. สิ้นสุด</th>
										</tr>
									</thead>
									<tbody>
										<tr
											v-for="(item, index) of store.data.road_section"
											:key="index"
											class="align-middle cursor-pointer"
										>
											<td class="text-center">{{ index + 1 }}</td>
											<td class="text-center">{{ item.road_group_number }}</td>
											<td class="text-center">{{ item.section_road_number }}</td>
											<td class="text-center">{{ item.section_road_th_name }}</td>
											<td class="text-center">{{ item.section_road_eng_name }}</td>

											<td class="text-center">{{ item.km_start }}</td>
											<td class="text-center">{{ item.km_end }}</td>
										</tr>
										<tr v-show="store.data.road_section.length === 0" class="text-center">
											<td colspan="8">ไม่พบข้อมูล</td>
										</tr>
									</tbody>
								</table>
							</div>
							<div class="row align-items-center mt-5">
								<div class="col-12 text-end">
									<button
										type="button"
										class="btn btn-primary rounded-4 ms-5 mt-md-0 mt-sm-2 mt-2 fw-semibold"
										:disabled="store.importing"
										@click="onClickImport()"
									>
										<span v-if="!store.importing">นำเข้าข้อมูล</span>
										<span v-else class="spinner-border spinner-border-sm align-middle mx-8"></span>
									</button>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</VSkeletonLoader>
</template>

<style scoped></style>
