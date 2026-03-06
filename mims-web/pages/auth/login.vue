<script lang="ts" setup>
import LoginForm from "~~/core/modules/auth/login/ui/LoginForm.vue"

definePageMeta({
	layout: "auth",
	middleware: "guest",
})

useHead({
	title: "เข้าสู่ระบบ",
})

onMounted(() => {
	// กรณีหน้าจอ Mobile ตอนกด Logout จะมี .drawer-overlay
	useBody().removeOverlay()
})
</script>

<template>
	<div class="d-flex h-100">
		<div class="d-flex flex-column flex-root m-auto">
			<div class="d-flex flex-column-fluid justify-content-center">
				<div class="card rounded-5 w-sm-475px m-5 m-md-0">
					<div class="card-body d-flex flex-column pt-13 pb-3">
						<div class="d-flex flex-center flex-column-fluid">
							<LoginForm />
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@media screen and (min-width: 576px) {
	.w-sm-560px {
		width: 560px !important;
	}
}
</style>
