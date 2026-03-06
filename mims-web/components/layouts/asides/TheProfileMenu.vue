<script setup lang="ts">
const removeOverlay = () => {
	useBody().removeOverlay()
}
</script>

<template>
	<div class="d-flex justify-content-center profile-menu">
		<button
			type="button"
			class="btn btm-sm btn-icon"
			data-kt-menu-trigger="{default: 'click', lg: 'hover'}"
			data-kt-menu-overflow="true"
			data-kt-menu-placement="top-start"
			data-bs-toggle="tooltip"
			data-bs-placement="right"
			data-bs-dismiss="click"
			:title="useInitUser()?.firstname"
		>
			<i class="fi fi-rr-user fs-2 text-gray-600"></i>
			<span class="profile-title text-gray-400"> {{ useInitUser()?.firstname }} {{ useInitUser()?.lastname }}</span>
		</button>

		<!--begin::Menu-->
		<div
			class="menu menu-sub menu-sub-dropdown menu-column menu-rounded menu-title-gray-700 menu-state-title-primary menu-state-icon-primary menu-state-bullet-primary menu-arrow-gray-500 fw-semobold w-225px px-2 py-3 pb-5 ms-10 profile-submenu"
			data-kt-menu="true"
		>
			<div class="menu-item px-3">
				<div class="menu-content fs-6 fw-bold px-1 py-2">
					{{ useInitUser()?.firstname }} {{ useInitUser()?.lastname }}<br />
					<span class="fs-6 fw-normal"
						>{{ useInitUser()?.ref_user_owner.email }} {{ useInitUser()?.ref_depot.name }}</span
					>
				</div>
			</div>
			<div class="menu-item px-3 py-0">
				<NuxtLink to="/profile" class="menu-link ps-4" @click="removeOverlay">
					<span class="menu-bullet">
						<span class="bullet bullet-dot"></span>
					</span>
					ข้อมูลส่วนตัว
				</NuxtLink>
			</div>
			<div class="menu-item px-3 py-0">
				<NuxtLink to="/change-password" class="menu-link ps-4" @click="removeOverlay">
					<span class="menu-bullet">
						<span class="bullet bullet-dot"></span>
					</span>
					เปลี่ยนรหัสผ่าน
				</NuxtLink>
			</div>

			<div class="menu-item px-3 py-0">
				<NuxtLink to="/auth/logout" class="menu-link ps-4" @click="removeOverlay">
					<span class="menu-bullet">
						<span class="bullet bullet-dot"></span>
					</span>
					ออกจากระบบ
				</NuxtLink>
			</div>
		</div>
		<!--end::Menu-->
	</div>
</template>

<style scoped>
.profile-title,
.menu-content {
	text-transform: capitalize;
}

.profile-menu .btn .profile-title {
	margin-top: 8px;
	font-size: 1rem;
	font-weight: 400;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 105px;
	display: inline-block;
	text-align: left;
}

.profile-menu .btn {
	width: 105px;
	height: 74px;
	margin: 0 auto;
	flex-direction: column;
	display: flex;
	border-radius: 0.475rem;
}

.profile-submenu.menu-sub-dropdown {
	border-radius: 1rem;
	z-index: 115 !important;
}
</style>
