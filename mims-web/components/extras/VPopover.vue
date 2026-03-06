<script setup lang="ts">
import * as bootstrap from "bootstrap"
import { onClickOutside } from "@vueuse/core"

type TPlacement = "auto" | "top" | "bottom" | "left" | "right"

const props = defineProps({
	className: {
		type: String,
		default: "",
	},
	hover: {
		type: Boolean,
		default: false,
	},
	isHover: {
		type: Boolean,
		default: false,
	},
	placement: {
		type: String as PropType<TPlacement>,
		default: "auto",
	},
})

const button = ref()
const content = ref()
const popover = ref<bootstrap.Popover | null>(null)

onMounted(() => {
	popover.value = new bootstrap.Popover(button.value!, {
		content: () => {
			return content.value as string
		},
		html: true,
		placement: props.placement,
	})
})

const togglePopover = () => {
	popover.value?.toggle()
}

const hoverPopover = () => {
	if (props.isHover) {
		popover.value?.toggle()
	}
	if (props.hover) {
		popover.value?.show()
	}
}

const outPopover = () => {
	if (props.isHover) {
		popover.value?.hide()
	} else if (!props.isHover) {
		if (props.hover) {
			popover.value?.hide()
		}
	}
}

onClickOutside(
	content,
	() => {
		popover.value?.hide()
	},
	{ ignore: [button] }
)
</script>

<template>
	<button
		v-if="!props.isHover"
		ref="button"
		type="button"
		class="btn btn-popover"
		:class="className"
		@click="togglePopover"
		@mouseover="hoverPopover"
		@mouseout="outPopover"
	>
		<slot />

		<div ref="content" class="popover-content">
			<slot name="content" />
		</div>
	</button>
	<button
		v-else
		ref="button"
		type="button"
		class="btn btn-popover"
		:class="className"
		@mouseover="hoverPopover"
		@mouseout="outPopover"
	>
		<slot />

		<div ref="content" class="popover-content">
			<slot name="content" />
		</div>
	</button>
</template>

<style scoped lang="scss">
button {
	height: 30px;
	.popover-content {
		display: none;
	}
}
</style>
