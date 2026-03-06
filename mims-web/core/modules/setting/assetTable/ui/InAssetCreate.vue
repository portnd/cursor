<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAssetTableCreateStore } from "../store"
import { GeomTypeForm, ColumnForm } from "./index"
import { useInitDataStore } from "~/core/modules/initData/store"

const store = useAssetTableCreateStore()
useStoreLifecycle(store)

onMounted(() => {
	store.reset("km")
	store.asset.data.geom_type = "km"
})

const onCancel = () => {
	return navigateTo("/settings/in-assets")
}

const { handleSubmit, handleReset } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.assetCreate("in")

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: function () {
				const initDataStore = useInitDataStore()
				initDataStore.initData()
				navigateTo("/settings/in-assets")
				handleReset()
			},
		})
	}
})

watch(
	() => store.asset.data.geom_type,
	(_) => {
		store.storeAsset()
	}
)

watch(
	() => store.asset.data.table_name,
	() => {
		if (store.asset.data?.table_name) {
			store.asset.data.table_name = store.asset.data.table_name.toLowerCase()
		}
	}
)

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
									<VLabel label="กำหนดชื่อตารางข้อมูล" :required="true" />

									<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
										<i class="fi fi-sr-interrogation fs-5"></i>
										<template #content>
											<div>กรอกได้เฉพาะ a-z, 0-9 และเครื่องหมาย _</div>
										</template>
									</VPopover>

									<VTextInput
										v-model="store.asset.data.table_name"
										name="table_name"
										:required="true"
										:validate-english="true"
										:paste="true"
									/>
								</div>

								<div class="col-12 mb-3">
									<GeomTypeForm
										v-model="store.asset.data.geom_type"
										asset-type="in"
										:image="store.imageAsset"
										:colors="store.asset.data.line_color_code"
										name="geom_type"
										@update:line-color-code="($colorCode) => (store.asset.data.line_color_code = $colorCode)"
										@update:image="($image) => (store.imageAsset = $image)"
									/>
								</div>
							</div>
						</div>
						<div class="col-md-6">
							<div class="row">
								<div class="col-12 mb-3">
									<VTextInput
										v-model="store.asset.data.table_label"
										label="ชื่อสินทรัพย์"
										name="table_label"
										:required="true"
									/>
								</div>
								<div class="col-12 mb-3">
									<VSelect
										v-model="store.asset.data.asset_group"
										:options="toOptions(useInitData().refAsset())"
										label="กลุ่มสินทรัพย์"
										name="asset_group"
										placeholder="เลือก"
										:required="true"
									/>
								</div>
							</div>
							<!-- <div class="col-12 mb-3">
								<VSelect
									v-model="store.asset.data.approver_id"
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
									v-model="store.asset.data.viewer_id"
									:options="toOptions(useInitData().refDepartment())"
									mode="tags"
									label="ผู้ดูข้อมูล"
									name="viewer"
									placeholder="เลือก"
									:required="true"
								/>
							</div> -->
						</div>

						<div class="col-lg-12">
							<ColumnForm
								v-model="store.asset.data.columns"
								@update:column="($column) => store.addColumnAsset($column)"
								@update:delete-column="($column) => store.deleteColumnAsset($column)"
								@update:edit-column="($column) => store.editColumnAsset($column)"
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
