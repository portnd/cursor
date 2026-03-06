<script lang="ts" setup>
import { useErrorStore } from "~~/core/shared/store/ErrorStore"

useHead({
	bodyAttrs: {
		class: "error-bg",
		id: "kt_body",
	},
})

const props = defineProps({
	error: {
		type: Object,
		default: () => ({}),
	},
})

const { error } = toRefs(props)
const errorStore = useErrorStore()
useStoreLifecycle(errorStore)
errorStore.showError(error.value.statusCode, error.value.message)

const title = errorStore.title
useHead({
	meta: [{ content: title }],
	title,
})

const handleError = () => {
	clearError({ redirect: "/" })
}

</script>

<template>
	<div class="page">
		<div class="d-flex h-100">
			<div class="d-flex flex-column flex-root m-auto">
				<div class="d-flex flex-column flex-center flex-column-fluid">
					<div class="d-flex flex-column flex-center text-center p-10">
						<div class="card card-flush w-md-500px px-sm-16">
							<div class="card-body py-14 py-md-16">
								<h2 class="fw-semibold text-gray-800 mb-4">
									{{ errorStore.messageTh }}
								</h2>
								<h3 class="fw-semibold text-gray-400 mb-7">
									{{ errorStore.messageEn }}
								</h3>
								<div class="mb-8">
									<img :src="errorStore.imagePath" class="mw-100 mh-250px theme-light-show" alt="" />
								</div>
								<div class="mb-0">
									<button class="btn btn-primary" @click="handleError">กลับสู่หน้าหลัก</button>
								</div>
								<div class="text-gray-500 text-center fw-normal fs-6 mt-12 mb-1">{{ useCopyright() }}</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
