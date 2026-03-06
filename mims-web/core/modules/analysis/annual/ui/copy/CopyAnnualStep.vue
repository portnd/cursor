<script setup lang="ts">
import { useAnnualAnalyseCopyStore } from "../../store"

const route = useRoute()
const tab = route.path.split("/")[3]
const store = useAnnualAnalyseCopyStore()
</script>

<template>
	<div id="kt_stepper_example_basic" class="stepper stepper-pills">
		<div class="stepper-nav .flex-md-start flex-center flex-wrap mb-5">
			<div class="card py-2 px-1 me-md-3 mb-md-0 mb-3">
				<div class="stepper-item mx-8 my-4" :class="{ current: store.step === 1 }" data-kt-stepper-element="nav">
					<div class="stepper-wrapper d-flex align-items-center">
						<div class="stepper-icon w-40px h-40px">
							<template v-if="tab === 'set-condition'">
								<i class="fi fi-br-check text-primary fs-5 mt-1"></i>
							</template>
							<template v-else>
								<span class="stepper-number">1</span>
							</template>
						</div>
						<div class="stepper-label">
							<h3 class="stepper-title">ขั้นตอนที่ 1</h3>
							<div class="stepper-desc">เลือกสายทาง</div>
						</div>
					</div>
				</div>
			</div>
			<hr class="step-line d-none d-md-block" />
			<div class="card py-2 px-0 ms-md-3">
				<div class="stepper-item mx-8 my-4" :class="{ current: store.step === 2 }" data-kt-stepper-element="nav">
					<div class="stepper-wrapper d-flex align-items-center">
						<div class="stepper-icon w-40px h-40px">
							<i class="stepper-check fas fa-check"></i>
							<span class="stepper-number">2</span>
						</div>
						<div class="stepper-label">
							<h3 class="stepper-title">ขั้นตอนที่ 2</h3>
							<div class="stepper-desc">กำหนดเงื่อนไข</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.card {
	box-shadow: none !important;
}
.step-line {
	border: 0.2px solid #000000;
	width: 3%;
}

.stepper.stepper-pills .stepper-item.current .stepper-icon {
	background-color: #fdb833 !important;
}

.stepper.stepper-pills .stepper-item.current .stepper-icon .stepper-number {
	color: #ffffff !important;
}

.stepper.stepper-pills .stepper-item .stepper-icon .stepper-number {
	color: #a1a5b7;
}

.stepper.stepper-pills .stepper-item .stepper-icon {
	background-color: #f4f4f4;
}
</style>
