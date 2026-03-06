<script setup lang="ts">
import { useForm } from "vee-validate"
import { useUsersCreateStore } from "../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useUsersCreateStore()
useStoreLifecycle(store)

onMounted(() => {
	store.getRoles()
})

const validate = computed(() => {
	const validations: IValidate = {}

	validations.firstname = "required"
	validations.lastname = "required"
	validations.tel = "required"
	validations.status = "required"
	validations.username = "username|required"
	validations.password = "required|password"
	validations.roles = "required"

	return validations
})

const { handleSubmit, handleReset, errors, isSubmitting } = useForm({ validationSchema: validate })

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.createUser()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				navigateTo("/users/users")
				handleReset()
			},
		})
	}
})

const onCancel = () => {
	navigateTo("/users/users")
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<form>
			<div class="col-12">
				<VUploadAvatar
					v-model="store.file"
					:file="store.filePath"
					:required="true"
					name="avatar"
					:image-size="250"
					total-file-size="1 MB"
				/>
			</div>
			<div class="col-12 mb-3">
				<div class="row">
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.params.firstname" label="ชื่อ" name="firstname" :required="true" />
					</div>
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.params.lastname" label="นามสกุล" name="lastname" :required="true" />
					</div>
					<div class="col-md-4 mb-3">
						<VPhoneInput v-model="store.params.tel" label="เบอร์โทรศัพท์" name="tel" :required="true" />
					</div>
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.params.email" label="อีเมล" name="email" />
					</div>
					<div class="col-md-4 mb-3">
						<VSelect
							v-model="store.params.ref_user_owner_id"
							:options="store.getUserOwnersOption"
							label="หน่วยงานที่รับผิดชอบ"
							name="ref_user_owner_id"
							placeholder="เลือกหน่วยงานที่รับผิดชอบ"
							:can-deselect="true"
						/>
					</div>
					<div v-if="store.params.ref_user_owner_id === 3" class="col-md-4 mb-3">
						<VSelect
							v-model="store.params.ref_depot_id"
							:options="store.getDepotOption"
							label="หมวด"
							name="depot_id"
							placeholder="เลือกหมวด"
							:can-deselect="true"
						/>
					</div>
					<div class="col-md-4 mb-3">
						<VRadio
							v-model="store.params.status"
							:options="[
								{ label: 'เปิดใช้งาน', value: 1 },
								{ label: 'ปิดใช้งาน', value: 2 },
							]"
							name="status"
							label="สถานะการใช้งาน"
							:required="true"
						/>
					</div>
					<div class="col-12">
						<h4 class="fw-semibold mt-10">ข้อมูลบัญชีผู้ใช้งาน</h4>
					</div>
					<div class="col-md-4">
						<VLabel label="ชื่อผู้ใช้งาน" :required="true" />
						<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
							<i class="fi fi-sr-interrogation fs-5"></i>
							<template #content>
								<span>สามารถประกอบด้วยตัวอักษรภาษาอังกฤษพิมพ์ใหญ่ หรือพิมพ์เล็กและตัวเลข 0 ถึง 9 เท่านั้น</span>
							</template>
						</VPopover>
						<VTextInput v-model="store.params.username" name="username" :required="true" />
					</div>
					<div class="col-md-4">
						<VLabel label="รหัสผ่าน" :required="true" />
						<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
							<i class="fi fi-sr-interrogation fs-5"></i>
							<template #content>
								<ul class="mt-4">
									<li>ต้องมีจำนวน 8 ตัวอักษรขึ้นไป</li>
									<li>ต้องประกอบด้วยตัวอักษร A-Z ตัวพิมพ์ใหญ่</li>
									<li>ต้องประกอบด้วยตัวอักษร a-z ตัวพิมพ์เล็ก</li>
									<li>ต้องประกอบด้วยตัวเลข 0-9 อย่างน้อย 1 ตัว</li>
								</ul>
							</template>
						</VPopover>
						<VPasswordInput v-model="store.params.password" name="password" :required="true" label="" />
					</div>

					<div class="col-md-6 mt-3">
						<VLabel label="สิทธิ์การใช้งาน" :required="true" />
						<div class="ps-3">
							<VCheckbox
								v-model="store.params.role"
								:options="toOptions(store.roles)"
								name="roles"
								:required="true"
								:inline="false"
							/>
						</div>
					</div>
				</div>
			</div>
			<div class="d-flex justify-content-end mt-5">
				<BtnCancel @click="onCancel" />
				<BtnSubmit :disabled="store.loading" :loading="store.loading" label="บันทึก" @click="onSubmit" />
			</div>
		</form>
	</div>
</template>

<style scoped></style>
