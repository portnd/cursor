<script setup lang="ts">
const props = defineProps({
	label: {
		type: String,
		required: true,
	},
	name: {
		type: String,
		required: true,
	},
	image: {
		type: String,
		default: "",
	},
	role: {
		type: String,
		required: true,
	},
})

// ตรวจสอบว่า url ของรูปภาพถูกต้อง
const src = ref("")
const isDefault = ref(true)
const checkFile = async () => {
	const defaultImage = "/images/logos/logo.png"

	if (props.image !== "" && props.image !== "undefined") {
		isDefault.value = false
		src.value = (await checkFileExist(props.image)) ? props.image : defaultImage
	} else {
		src.value = defaultImage
	}
}

onMounted(() => {
	checkFile()
})

watch(
	() => props.image,
	() => {
		checkFile()
	}
)
</script>

<template>
	<VPopover class-name="p-1 my-0 text-primary fw-normal form-text cursor-pointer vuser" placement="bottom">
		{{ label }}
		<template #content>
			<div class="row p-0 align-items-center">
				<div class="col-3 text-center">
					<img :class="isDefault ? '' : 'rounded-pill'" :src="src" alt="" width="35" />
				</div>
				<div class="col-9">
					<div class="row">
						<div class="col-12">
							<span class="user-name">{{ name }}</span>
						</div>
						<div class="col-12">
							<span class="user-role">แผนก: {{ role }}</span>
						</div>
					</div>
				</div>
			</div>
		</template>
	</VPopover>
</template>

<style scoped lang="scss">
.vuser {
	text-transform: capitalize;
	margin-top: -1px !important;
}
.user-name {
	text-transform: capitalize;
}
.user-name,
.user-role {
	font-size: 1rem;
}
</style>
