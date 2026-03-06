<script setup lang="ts">
import ForgotPasswordForm from "~~/core/modules/auth/forgotPassword/ui/index"

definePageMeta({
	layout: "auth",
	middleware: "guest",
})

useHead({
	title: "ลืมรหัสผ่าน ?",
})
</script>

<template>
	<div class="d-flex h-100">
		<div class="d-flex flex-column flex-root m-auto">
			<div class="d-flex flex-column-fluid justify-content-center">
				<div class="card rounded-5 w-sm-500px m-5 m-md-0">
					<div class="card-body d-flex flex-column pt-13">
						<div class="d-flex flex-center flex-column-fluid">
							<ForgotPasswordForm />
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
