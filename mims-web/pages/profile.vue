<script lang="ts" setup>
import { AccountForm } from "~~/core/modules/account/ui/index"

useHead({
	title: "ข้อมูลส่วนตัว",
})
</script>

<template>
	<TheBreadcrumb :breadcrumbs="[{ name: useInitUser()?.username, to: '/' }]" title="ข้อมูลส่วนตัว" />

	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<AccountForm />
			</div>
		</div>
	</div>
</template>

<style scoped></style>
