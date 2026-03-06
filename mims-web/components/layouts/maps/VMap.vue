<script setup lang="ts">
const emit = defineEmits(["map", "update:modelValue"])

const props = defineProps({
	loading: {
		type: Boolean,
		default: false,
	},
	lastView: {
		type: Boolean,
		default: false,
	},
	btnToggle: {
		type: Boolean,
		default: true,
	},
	modelValue: {
		type: Object,
	},
	height: {
		type: String,
		default: "97vh",
	},
	heightLg: {
		type: String,
		default: "97vh",
	},
	heightMd: {
		type: String,
		default: "97vh",
	},
	heightSm: {
		type: String,
		default: "97vh",
	},
	borderRadius: {
		type: String,
		default: "",
	},
	isSticky: {
		type: Boolean,
		default: false,
	},
	opacityLoad: {
		type: Boolean,
		default: false,
	},
})

const refMap = ref<any>(null)
const isMapLoaded = ref(false)
const isMapShow = ref(true)
const isShowBtnMap = ref(false)
const isModeMap = ref(false)

const setMap = (map: any) => {
	refMap.value = map

	// Disable the context menu
	refMap.value?.Ui.ContextMenu.enableNearPoi(false)
	refMap.value?.Ui.ContextMenu.enableAddress(false)

	emit("map", map)
}

// ซ่อน/แสดง แผนที่
const mapShow = reactive({
	collapsed: false,
})

const toggle = () => {
	mapShow.collapsed = !mapShow.collapsed

	if (mapShow.collapsed === true) {
		document.body.classList.add("longomap-collapsed")
	} else {
		document.body.classList.remove("longomap-collapsed")
	}

	emit("update:modelValue", { collapsed: mapShow.collapsed })
}

onActivated(() => {
	setTimeout(() => {
		isMapLoaded.value = true
	}, 1000)
})

onDeactivated(() => {
	isMapLoaded.value = false
})

// กรณี isSticky
const heightXl = ref(props.height)

const handleScroll = () => {
	// if (props.isSticky) {
	// 	const windowHeight = window.innerHeight
	// 	const fullHeight = document.documentElement.scrollHeight
	// 	const scrollTop =
	// 		window.scrollY ||
	// 		window.pageYOffset ||
	// 		document.body.scrollTop + ((document.documentElement && document.documentElement.scrollTop) || 0)
	// 	if (scrollTop >= fullHeight - windowHeight - 56) {
	// 		heightXl.value = `calc(${props.height} - 56px)`
	// 	} else {
	heightXl.value = `calc(${props.height})`
	// 	}
	// }
	refMap.value?.resize()
}

const borderRadiusLg = computed(() => {
	return props.borderRadius === "" ? "0px 15px 15px" : props.borderRadius
})

const borderRadiusSm = computed(() => {
	return props.borderRadius === "" ? "15px" : props.borderRadius
})

const screenWidth = ref(0)
// สำหรับแสดง/ซ่อนแผนที่
const setShowMap = () => {
	if (screenWidth.value < 992) {
		isMapShow.value = false
		isShowBtnMap.value = true
	} else {
		isMapShow.value = true
		isShowBtnMap.value = false
	}
}

// สำหรับเปิด/ปิดโหมดแผนที่
const handleModeMap = () => {
	// update ค่า
	isModeMap.value = !isModeMap.value

	if (isModeMap.value === true) {
		window.scrollTo({
			top: 0,
		})
		const elementContent = document.getElementById("kt_content") as HTMLElement
		elementContent.style.zIndex = "1"
	} else {
		const elementContent = document.getElementById("kt_content") as HTMLElement
		elementContent.style.zIndex = "unset"
	}

	setModeMap()

	// resize ของ longdo map เพื่อไม่ให้ bug
	setTimeout(() => {
		refMap.value?.resize()
	}, 200)
}

const setModeMap = () => {
	if (isModeMap.value === true) {
		useBody().addBodyClassname("scrollbar-none")
		useBody().addBodyClassname("mode-map")
	} else {
		useBody().removeBodyClassName("scrollbar-none")
		useBody().removeBodyClassName("mode-map")
	}

	if (screenWidth.value > 991) {
		useBody().removeBodyClassName("scrollbar-none")
		useBody().removeBodyClassName("mode-map")
	}
}

const handleWidth = () => {
	screenWidth.value = window.innerWidth

	setShowMap()
	setModeMap()
}

const setDefaultShowMap = (isShow: boolean) => {
	isModeMap.value = isShow
	isShowBtnMap.value = !isShow
}

defineExpose({
	setDefaultShowMap,
})

onMounted(() => {
	document.body.classList.remove("longomap-collapsed")

	handleWidth()

	window.addEventListener("scroll", handleScroll)
	window.addEventListener("resize", handleWidth)
})

onUnmounted(() => {
	window.removeEventListener("scroll", handleScroll)
	window.removeEventListener("resize", handleWidth)
})
</script>

<template>
	<button v-if="mapShow.collapsed" class="btn btn-primary shadow btn-floating-toggle" @click="toggle">
		<i class="fi fi-rr-caret-down"></i>
		<span class="fw-semibold">แสดงแผนที่</span>
	</button>

	<VSkeletonLoader v-if="loading || !isMapLoaded">
		<div class="map-loading w-100 rounded-3"></div>
	</VSkeletonLoader>
	<longdo-map
		v-show="isMapShow || isModeMap"
		v-else-if="!loading && isMapLoaded"
		:last-view="lastView"
		class="longdo-map position-relative"
		:class="opacityLoad ? 'opacity-30' : ''"
		@load="setMap"
	>
		<button v-if="btnToggle" class="bg-primary btn-toggle d-none d-lg-block position-absolute" @click="toggle">
			<i class="fi fi-rr-caret-right" />
		</button>
		<div class="bottom-start">
			<slot name="bottom-start"></slot>
		</div>

		<div class="bottom-end">
			<slot name="bottom-end"></slot>
		</div>

		<div class="bottom-center">
			<slot name="bottom-center"></slot>
		</div>
		<div class="top-end">
			<slot name="top-end"></slot>
		</div>
		<div class="top-end-modal">
			<slot name="top-end-modal"></slot>
		</div>
	</longdo-map>

	<div v-show="isShowBtnMap" class="btn-show-map cursor-pointer" @click="handleModeMap">
		<div class="rounded-circle text-white bg-primary text-center bg-white pt-4">
			<i v-if="!isModeMap" class="fi fi-bs-map lh-md fs-2"></i>
			<i v-else class="fi fi-br-arrow-left lh-md fs-2"></i>
		</div>
	</div>

	<div
		v-show="isModeMap"
		class="footer-map py-3 d-flex flex-column flex-md-row align-items-center justify-content-center container-fluid"
	>
		<!--begin::Copyright-->
		<div class="text-dark text-center">
			<span class="text-muted fw-semobold me-1">{{ useCopyright() }}</span>
		</div>
		<!--end::Copyright-->
	</div>
</template>

<style lang="scss">
.longdo-map {
	height: v-bind(heightXl);
	width: 100% !important;
	@media (max-width: 991px) {
		height: v-bind(heightLg);
	}
	@media (max-width: 767px) {
		height: v-bind(heightMd);
	}
	@media (max-width: 575px) {
		height: v-bind(heightSm);
	}
	div {
		.ldmap_placeholder_fullscreen {
			position: fixed !important;
			top: 0 !important;
			left: 125px !important;
			width: calc(100% - 125px) !important;
			height: 100% !important;
			z-index: 255;
		}
	}
}
.map-loading {
	height: v-bind(heightXl);
	width: 100vh;
	background-color: #d9d9d9;
	@media (max-width: 991px) {
		height: v-bind(heightLg);
	}
	@media (max-width: 767px) {
		height: v-bind(heightMd);
	}
	@media (max-width: 575px) {
		height: v-bind(heightSm);
	}
}

.ldmap_placeholder {
	border-radius: v-bind(borderRadiusLg);

	@media (max-width: 991px) {
		border-radius: v-bind(borderRadiusSm);
	}
}

// Remove Longdo Map Icon from the map body
.ldmap_notice > a {
	display: none !important;
}

@media (max-width: 991px) {
	.longdo-map {
		top: 0px;
	}
}

@media (min-width: 768px) {
	.longdo-map > .btn-toggle {
		height: 50px;
		width: 27px;
		border: 0;
		border-radius: 14px 0 0 14px;
		color: #fff;
		font-size: 1.5rem;
		padding: 5px 0px 0px 3px;
		top: 0;
		left: -27px;
	}
}

.map-sticky {
	margin-top: -53px;
	z-index: 0;
	.widget {
		top: 10px;
		position: sticky !important;
		z-index: 50 !important;
		@media (max-width: 991px) {
			position: initial !important;
			width: 100% !important;
		}
	}
	@media (max-width: 991px) {
		margin-top: 15px;
		padding-right: 10px;
	}
	&:has(.ldmap_placeholder_fullscreen) {
		width: 100%;
	}
}

.map-collapsed {
	height: 45px;
	position: fixed;
	right: -105px;
	text-align: center !important;
	top: 38%;
	transform: rotate(90deg);
	transform-origin: bottom left;
	width: 150px;
	z-index: 99;

	.longdo-map,
	.image {
		visibility: hidden;
		height: 0;
		opacity: 0;
	}

	.btn-floating-toggle {
		border-radius: 0 0 14px 14px !important;
		padding: 0.75rem 1.5rem 0.25rem !important;
		font-size: 14px !important;
	}

	@media (max-width: 991px) {
		transform: rotate(0deg) !important;
		position: initial !important;
		width: 100% !important;
		height: unset;

		.longdo-map,
		.image {
			visibility: visible !important;
		}

		.btn-floating-toggle {
			display: none;
		}
	}
}

.footer-map {
	display: none !important;
}

.longdo-map {
	.square {
		width: 20px;
		height: 20px;
		border-radius: 5px;
		display: inline-block;
	}
	.button {
		display: inline-block;
		margin: 5px;
		padding: 9px 6px 5px 5px;
	}
	.icon-plus {
		position: absolute;
		padding-top: 10px;
		margin-left: -7px;
	}
}

// scrollbar ล่างสุด
.longdo-map.bottomed .ldmap_placeholder .ldmap_crosshair {
	top: 56px;
	@media (max-width: 991px) {
		top: 0px !important;
	}
}

body.modal-open .ldmap_placeholder .ldmap_crosshair {
	right: 5px;
	@media (max-width: 991px) {
		right: 15px !important;
	}
}

.ldmap_placeholder .ldmap_frame .ldmap_popup_close {
	width: 10px;
	height: 10px;
}

// popup
.ldmap_placeholder .ldmap_frame .ldmap_popupholder {
	min-width: 210px;
}
.ldmap_placeholder .popup {
	background-color: #fff;
	padding: 1rem;
	line-height: 1.75;
	display: block;
	border-radius: 1rem;
	margin-top: -120px;
}
.ldmap_placeholder .popup .title {
	font-weight: 500 !important;
	font-size: 1rem;
}
.ldmap_placeholder .popup b {
	font-weight: 400 !important;
}

.ldmap_placeholder .ldmap_contextmenu {
	background: #fff !important;
	border-radius: 0.75rem;
	border: 1px solid #e0e1e9 !important;
}
.ldmap_placeholder .ldmap_contextmenu .ldmap_contextmenu_info {
	background: unset !important;
	padding: 6px 10px !important;
}
.ldmap_placeholder .ldmap_contextmenu .ldmap_contextmenu_extra {
	padding: 4px 12px 6px !important;
}

.ldmap_placeholder .ldmap_frame .ldmap_popup_callout {
	left: 95px !important;
	width: 55px !important;
	height: 53px !important;
	opacity: 0.9;
}

.ldmap_placeholder .ldmap_frame .ldmap_popup_close {
	width: 10px !important;
	height: 10px !important;
	position: absolute;
	top: 10px;
	right: 16px;
}

.ldmap_placeholder .ldmap_frame .ldmap_popup {
	padding: 8px 0px 8px 16px !important;
	line-height: 1.6 !important;
	border: 1px solid #dedede !important;
	box-shadow: 2px 2px 6px #c3c2c2 !important;
	min-width: 292px !important ;
	.colon {
		max-width: 1px;
		padding: 0px;
	}
}
.ldmap_placeholder .ldmap_frame .ldmap_popup_detail {
	max-height: 300px !important;

	.maintenance-popup-view-details:hover {
		opacity: 0.85;
		text-decoration: underline !important;
	}
}
.ldmap_placeholder .ldmap_popup_mini {
	top: 6%;
	left: 2.5%;
	width: 95% !important;
	border-radius: 16px;
	height: max-content !important;
	.ldmap_popup_detail {
		width: unset !important;
		transform: scale(1) !important;
		.colon {
			max-width: 0.5em !important;
		}
		.detail-button-main {
			width: fit-content;
		}
	}
}

// slot
.bottom-start {
	display: block;
	position: absolute;
	z-index: 255;
	bottom: 20px;
	left: 15px;
}
.bottom-end {
	display: block;
	position: absolute;
	z-index: 255;
	bottom: 20px;
	right: 15px;
}
.bottom-center {
	position: absolute;
	z-index: 255;
	bottom: 20px;
	width: 300px;
	bottom: 0;
	left: 50%;
	transform: translate(-50%, -20%);
	@media (max-width: 1100px) {
		width: 250px;
	}
	@media (max-width: 990px) {
		width: 300px;
	}
}
.top-end {
	display: block;
	position: absolute;
	z-index: 255;
	top: 40px;
	right: 15px;
}
.top-end-modal {
	display: block;
	position: absolute;
	z-index: 255;
	top: 85px;
	right: 15px;
}

// btn-show-map
.btn-show-map {
	position: fixed;
	bottom: 57px;
	right: 10px;
	z-index: 999;
	.rounded-circle {
		width: 47px;
		height: 47px;
	}
}

.mode-map {
	.map-sticky {
		position: fixed !important;
		top: 50px;
		z-index: 999;
		padding-right: 0px;
		padding-left: 0px;
		margin-left: -5px;
		background-color: #f3f3f3;
		padding-top: 0px;
	}
	.longdo-map {
		height: 90dvh;
		div {
			.ldmap_placeholder_fullscreen {
				position: fixed !important;
				top: 0 !important;
				left: 0 !important;
				width: 100% !important;
				height: 100% !important;
				z-index: 255;
			}
		}
	}

	.footer-map {
		position: fixed !important;
		background-color: #ffff;
		bottom: 0px;
		z-index: 999;
		display: block !important;
	}

	.ldmap_placeholder {
		border-radius: 0px;
	}
}

.opacity-30 {
	opacity: 0.3;
}
</style>
