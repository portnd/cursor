<script setup lang="ts">
const scrollElRef = ref<null | HTMLElement>(null)

onMounted(() => {
	if (scrollElRef.value) {
		scrollElRef.value.scrollTop = 0
	}
})

const route = useRoute()
const currentActive = (current: string) => {
	return route.path.startsWith(current)
}

const hasActiveChildren = (match: string) => {
	return route.path.includes(match)
}

const removeOverlay = () => {
	useBody().removeOverlay()
}
</script>

<template>
	<div
		id="kt_aside_menu_wrapper"
		class="hover-scroll-overlay-y my-2 my-lg-0 pe-lg-n1"
		data-kt-scroll="true"
		data-kt-scroll-height="auto"
		data-kt-scroll-dependencies="#kt_aside_logo, #kt_aside_footer"
		data-kt-scroll-wrappers="#kt_aside, #kt_aside_menu"
		data-kt-scroll-offset="50px"
	>
		<div
			id="kt_aside_menu"
			class="menu menu-column menu-title-gray-700 menu-state-title-primary menu-state-icon-primary menu-state-bullet-primary menu-arrow-gray-500 fw-normal"
			data-kt-menu="true"
		>
			<template v-for="item in useInitMenu()" :key="item.id">
				<!-- กรณีไม่มี Sub เมนู -->
				<template v-if="item.is_children">
					<div :class="{ 'show here': currentActive(item.route) }" class="menu-item py-2">
						<NuxtLink class="menu-link menu-center" :to="item.route" @click="removeOverlay">
							<span class="menu-icon mb-2 me-0">
								<i :class="item.icon" class="fs-1 text-gray-600" />
							</span>
							<span class="menu-title text-gray-600">{{ item.name }}</span>
						</NuxtLink>
					</div>
				</template>

				<!-- กรณีมี Sub เมนู -->
				<template v-if="!item.is_children">
					<div
						:key="item.id"
						:class="{ 'show here': hasActiveChildren(item.route) }"
						data-kt-menu-trigger="{default: 'click', lg: 'hover'}"
						data-kt-menu-placement="right-start"
						class="menu-item py-2"
					>
						<span class="menu-link menu-center">
							<span class="menu-icon mb-2 me-0">
								<i :class="item.icon" class="fs-1 text-gray-600" />
							</span>
							<span class="menu-title text-gray-600">{{ item.name }}</span>
						</span>
						<div class="menu-sub menu-sub-dropdown px-2 py-3 w-225px mh-75 overflow-auto">
							<div class="menu-item">
								<div class="menu-content">
									<span class="menu-section fs-5 fw-semibold ps-1 py-1">{{ item.name }}</span>
								</div>
							</div>
							<div v-for="subMenu in item.children" :key="subMenu.id" class="menu-item">
								<NuxtLink
									v-show="!subMenu.name.includes('สินทรัพย์นอกเขตทาง')"
									class="menu-link"
									:class="{ active: currentActive(subMenu.route) }"
									:to="subMenu.route"
									@click="removeOverlay"
								>
									<span class="menu-bullet">
										<span class="bullet bullet-dot"></span>
									</span>
									<span class="menu-title">{{ subMenu.name }}</span>
								</NuxtLink>
							</div>
						</div>
					</div>
				</template>
			</template>
		</div>
	</div>
</template>

<style scoped>
.menu > .menu-item > .menu-link .menu-title {
	margin-top: -8px;
	font-size: 1rem;
	text-align: center;
}
.menu-sub-dropdown > .menu-item > .menu-link .menu-title {
	font-size: 1rem;
}

.menu > .menu-item.here > .menu-link,
.menu > .menu-item:not(.here) > .menu-link:hover:not(.disabled):not(.active):not(.here) {
	background: var(--kt-primary);
	border-radius: 24px;
}

.menu > .menu-item.here > .menu-link .menu-icon i {
	color: var(--kt-gray-800) !important;
}

.menu-item .menu-link {
	padding: 0.55rem 0rem !important;
}

.menu > .menu-item.here > .menu-link .menu-icon i,
.menu > .menu-item.here > .menu-link .menu-title,
.menu > .menu-item:not(.here) > .menu-link:hover:not(.disabled):not(.active):not(.here) .menu-title,
.menu-state-icon-primary .menu-item:not(.here) .menu-link:hover:not(.disabled):not(.active):not(.here) .menu-icon i {
	color: var(--kt-gray-800) !important;
}

.show.menu-dropdown > .menu-sub-dropdown,
.menu-sub-dropdown.menu.show,
.menu-sub-dropdown.show[data-popper-placement] {
	margin-left: -30px !important;
	border-radius: 1rem;
}

.menu-item .menu-link .menu-bullet {
	margin-left: 1.75rem;
}

.menu-item .menu-content {
	color: #4c4e6f !important;
}
</style>
