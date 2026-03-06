<script setup lang="ts">
defineProps({
	id: {
		type: Number,
		required: true,
	},
	icon: {
		type: String,
		default: "",
	},
	width: {
		type: String,
		default: "28px",
	},
	heigth: {
		type: String,
		default: "23px",
	},
	customClass: {
		type: String,
		default: "",
	},
	color: {
		type: String,
		required: true,
	},
})
</script>

<template>
	<div v-if="icon && icon !== ''" class="road-type-icon">
		<span :class="customClass" v-html="icon"></span>
	</div>
	<div v-else class="road-type-icon">
		<span :class="customClass" v-html="getRoadTypeIcon(id)"></span>
	</div>
</template>

<style>
.road-type-icon > span {
	background-color: v-bind(color);
	border-radius: 50%;
	display: inline-block;
	padding: 13px 10px;
	color: #fff;
	line-height: 0;
	margin: 0.25rem;
}

.road-type-icon > span svg {
	width: v-bind(width);
	height: v-bind(height);
}
</style>
