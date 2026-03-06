<script lang="ts" setup>
import { ChangePasswordForm } from "~/core/modules/account/ui"

useHead({
	title: "เปลี่ยนรหัสผ่าน",
})
</script>

<template>
	<TheBreadcrumb title="เปลี่ยนรหัสผ่าน" :breadcrumbs="[{ name: useInitUser()?.username, to: '/' }]" />
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<ChangePasswordForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
