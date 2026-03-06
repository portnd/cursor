<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAssetTableEditStore } from "../store"
import { GeomTypeForm, ColumnForm } from "./index"
import { useInitDataStore } from "~/core/modules/initData/store"

const store = useAssetTableEditStore()
useStoreLifecycle(store)
onMounted(async () => {
	const route = useRoute()
	await store.getAsset(route.params.id as string)
})

const onCancel = () => {
	return navigateTo("/settings/out-assets")
}

const { handleSubmit, handleReset } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.editAsset("out")
	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: function () {
				const initDataStore = useInitDataStore()
				initDataStore.initData()
				navigateTo("/settings/out-assets")
				handleReset()
			},
		})
	}
})

</script>

<template>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<!-- begin::Form -->
				<form @submit="onSubmit">
					<div class="row">
						<div class="col-md-6">
							<div class="row">
								<div class="col-12 mb-3">
									<VLabel label="ชื่อตารางข้อมูล" :required="true" />

									<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
										<i class="fi fi-sr-interrogation fs-5"></i>
										<template #content>
											<div>กรอกได้เฉพาะ a-z, A-Z, 0-9 และเครื่องหมาย _</div>
										</template>
									</VPopover>

									<VTextInput v-model="store.asset.table_name" name="table_name" :required="true" :disabled="true" />
								</div>

								<div class="col-12 mb-3 disabled">
									<GeomTypeForm
										v-model="store.asset.geom_type"
										:image="store.imageAsset"
										:files="store.asset.icon_filepath"
										name="geom_type"
										:disabled="true"
										:colors="store.asset.line_color_code"
										@update:lineColorCode="($colorCode) => (store.asset.line_color_code = $colorCode)"
										@update:image="($image) => (store.imageAsset = $image)"
										@update:deleteImageDefault="($image) => (store.asset.icon_filepath = $image)"
									/>
								</div>
							</div>
						</div>
						<div class="col-md-6">
							<div class="row">
								<div class="col-12 mb-3">
									<VTextInput
										v-model="store.asset.table_label"
										label="ชื่อสินทรัพย์"
										name="table_label"
										:required="true"
									/>
								</div>
								<div class="col-12 mb-3">
									<VSelect
										v-model="store.asset.asset_group.id"
										:options="toOptions(useInitData().refAsset())"
										label="กลุ่มสินทรัพย์"
										name="asset_group"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
							</div>
							<div class="col-12 mb-3">
								<VSelect
									v-model="store.asset.approver"
									:options="toOptions(useInitData().refDepartment())"
									mode="tags"
									label="ผู้อนุมัติ/แก้ไขข้อมูล"
									name="approver"
									placeholder="เลือก"
									:required="true"
								/>
							</div>
							<div class="col-12 mb-3">
								<VSelect
									v-model="store.asset.viewer"
									:options="toOptions(useInitData().refDepartment())"
									mode="tags"
									label="ผู้ดูข้อมูล"
									name="viewer"
									placeholder="เลือก"
									:required="true"
								/>
							</div>
						</div>

						<div class="col-lg-12">
							<ColumnForm
								v-model="store.asset.columns"
								:type-edit="true"
								@update:column="($column) => store.addColumnAsset($column)"
								@update:deleteColumn="($column) => store.deleteColumnAsset($column)"
								@update:editColumn="($column) => store.editColumnAsset($column)"
							/>
						</div>
					</div>
					<div class="d-flex justify-content-end mt-5">
						<BtnCancel @click="onCancel" />
						<BtnSubmit :loading="store.loading" label="บันทึก" />
					</div>
				</form>
				<!-- end::Form -->
			</div>
		</div>
	</div>
</template>

<style scoped></style>
