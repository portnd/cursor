<script setup lang="ts">
import { useForm } from "vee-validate"
import { useAccountStore } from "../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useAccountStore()
useStoreLifecycle(store)

const validate = computed(() => {
	const validation = {} as IValidate

	validation.firstname = "required"
	validation.lastname = "required"
	validation.tel = "required"

	return validation
})

onMounted(() => {
	store.getAccount()
})

const { handleSubmit } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.updateAccount()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.getAccount()
			},
		})
	}
})

const onCancel = () => {
	navigateTo("/")
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<VSkeletonLoader :loading="store.loading">
			<div class="col-12 mb-3">
				<VUploadAvatar
					v-model="store.params.file"
					:file="store.params.file_path"
					:image-size="250"
					total-file-size="1 MB"
					name="profile-image"
				/>
			</div>
			<div class="col-12">
				<div class="row">
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.account.firstname" label="ชื่อ" name="firstname" :required="true" />
					</div>
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.account.lastname" label="นามสกุล" name="lastname" :required="true" />
					</div>
					<div class="col-md-4 mb-3">
						<VPhoneInput
							v-model="store.account.tel"
							class="text-align-center"
							label="เบอร์โทรศัพท์"
							name="tel"
							:required="true"
						/>
					</div>
				</div>
			</div>
			<div class="col-12">
				<div class="row">
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.account.email" label="อีเมล" name="email" />
					</div>
					<div class="col-md-4 mb-3">
						<VTextInput
							v-model="store.account.ref_user_owner.email"
							label="หน่วยงานที่รับผิดชอบ"
							:disabled="true"
							name="ref_user_owner"
						/>
					</div>
					<div class="col-md-4 mb-3">
						<VTextInput v-model="store.account.ref_depot.name" label="หมวด" :disabled="true" name="ref_depot" />
					</div>
					<div class="col-md-4 mt-3">
						<VLabel label="สิทธิ์การใช้งาน" />
						<div class="ps-3">
							<VCheckbox
								v-model="store.params.roles"
								:disabled="true"
								:options="toOptions(store.account.roles)"
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
		</VSkeletonLoader>
	</div>
</template>

<style scoped></style>
