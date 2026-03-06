<script setup lang="ts">
import { ScrollTopComponent } from "~~/assets/themes/ts/components/_ScrollTopComponent"

onMounted(() => {
	ScrollTopComponent.reinitialization()
})
</script>

<template>
	<div id="kt_scrolltop" ref="kt_scrolltop" class="scrolltop" data-kt-scrolltop="true">
		<span class="text-white fs-1">
			<i class="fi fi-br-arrow-small-up"></i>
		</span>
	</div>
</template>

<style scoped>
.scrolltop > span > i.fi {
	vertical-align: sub;
}
</style>
