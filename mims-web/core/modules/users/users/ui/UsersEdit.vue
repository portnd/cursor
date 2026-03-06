<script setup lang="ts">
import { useForm } from "vee-validate"
import { useUsersEditStore } from "../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useUsersEditStore()
useStoreLifecycle(store)
const route = useRoute()
const id = Number(route.params.id)

const isChangePassword = ref<boolean>(false)

onMounted(() => {
	store.getDefaultUsers(id)
})

watch(
	() => store.params.status,
	() => {
		store.updateStatus()
	}
)

const onChangePassword = () => {
	isChangePassword.value = !isChangePassword.value
}

const validate = computed(() => {
	const validations: IValidate = {}

	validations.firstname = "required"
	validations.lastname = "required"
	validations.tel = "required"
	validations.status = "required"
	validations.username = "required|username"
	validations.password = isChangePassword.value ? "required|password" : ""
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

const onSubmit = handleSubmit(async (_, acitons) => {
	useAction(acitons)

	const res = await store.updateUser(id)

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

watch(isChangePassword, () => {
	if (isChangePassword.value === false) {
		store.params.password = ""
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
		<div class="col-12">
			<VUploadAvatar
				v-model="store.file"
				:file="store.data.profile_img_path"
				:required="true"
				name="avatar"
				:image-size="250"
				total-file-size="1MB"
			/>
		</div>
		<div class="col-12 mb-3">
			<div class="row">
				<div class="col-md-4 mb-3">
					<VTextInput v-model="store.data.firstname" label="ชื่อ" name="firstname" :required="true" />
				</div>
				<div class="col-md-4 mb-3">
					<VTextInput v-model="store.data.lastname" label="นามสกุล" name="lastname" :required="true" />
				</div>
				<div class="col-md-4 mb-3">
					<VPhoneInput v-model="store.data.tel" label="เบอร์โทรศัพท์" name="tel" :required="true" />
				</div>
				<div class="col-md-4 mb-3">
					<VTextInput v-model="store.data.email" label="อีเมล" name="emailt" />
				</div>
				<div class="col-md-4 mb-3">
					<VSelect
						v-model="store.data.ref_user_owner_id"
						:options="store.getUserOwnersOption"
						label="หน่วยงานที่รับผิดชอบ"
						name="ref_user_owner_id"
						placeholder="เลือกหน่วยงานที่รับผิดชอบ"
						:can-deselect="true"
					/>
				</div>
				<div v-if="store.data.ref_user_owner_id === 3" class="col-md-4 mb-3">
					<VSelect
						v-model="store.data.ref_depot_id"
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
					<VTextInput v-model="store.data.username" label="ชื่อผู้ใช้งาน" name="username" :disabled="true" />
				</div>
				<div class="col-md-8">
					<div class="row">
						<div v-show="isChangePassword" class="col-md-6 col-sm-6 mb-3">
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
							<VPasswordInput v-model="store.params.password" name="password" :required="true" />
						</div>
						<div class="col-md-6 col-sm-6 pt-5" :class="isChangePassword === false ? 'align-self-center' : 'mt-4'">
							<template v-if="!isChangePassword">
								<button
									type="button"
									class="btn btn-primary rounded-4 px-8 py-3 fw-semibold mt-5"
									@click="onChangePassword"
								>
									เปลี่ยนรหัสผ่าน
								</button>
							</template>
							<template v-if="isChangePassword">
								<BtnCancel label="ยกเลิกการเปลี่ยนรหัสผ่าน" @click="onChangePassword" />
							</template>
						</div>
					</div>
				</div>

				<div class="col-md-6 mt-3">
					<VLabel label="สิทธิ์การใช้งาน" :required="true" />
					<div class="ps-3">
						<VCheckbox
							v-model="store.params.roles"
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
	</div>
</template>

<style scoped></style>
