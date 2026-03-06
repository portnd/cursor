<script setup lang="ts">
import { useMaintenanceHistoryDetailsStore } from "../../../store/MaintenanceHistoryDetailsStore"

const store = useMaintenanceHistoryDetailsStore()
</script>

<template>
	<div class="row">
		<div class="col-12 mt-4 mt-md-4">
			<div class="border border-1 border-gray-300 rounded-1 p-4">
				<h5 class="fw-semibold text-start mb-3">สถานะโครงการ</h5>
				<VSelect
					v-model="store.params.planName"
					:options="store.getPlanOptions"
					name="plan"
					:can-clear="false"
					:can-deselect="false"
					เลือก
					@update:model-value="(e: any) => store.setSchedulePlan(e)"
				/>
				<div class="order-responive mt-6 ps-5">
					<div v-if="store.getScheduleList.length > 0" class="order-track">
						<div v-for="(item, index) of store.getScheduleList" :key="index" class="order-track-step mb-2">
							<div class="order-track-status">
								<span class="order-track-status-dot" :class="item.is_checked ? 'active' : ''"></span>
								<span class="order-track-status-line"></span>
							</div>
							<div class="order-track-text">
								<span class="order-track-text-stat"> {{ item.status }} </span>
								<br />
								<span class="order-track-text-sub text-gray-600">{{
									buddhistFormatDate(item.disbursement_date, "dd mmm yyyy")
								}}</span>
							</div>
						</div>
					</div>
					<VNotFound v-else :is-not-shadow="true" height="50dvh" />
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
$primary: #fdb833;

@mixin dot {
	display: block;
	width: 1.5rem;
	height: 1.5rem;
	border-radius: 50%;
}

.order-track {
	padding: 0 0.25rem;
	display: flex;
	flex-direction: column;

	&-step {
		display: flex;
		height: 4rem;

		&:last-child {
			overflow: hidden;
			height: 3rem;
			& .order-track-status span:last-of-type {
				display: none;
			}
		}
	}

	&-status {
		margin-right: 1.5rem;
		position: relative;
		&-dot {
			@include dot;
			border: 2px solid rgba(0, 0, 0, 0.06);
			background: #ffffff;
		}
		&-line {
			display: block;
			margin: 0 auto;
			width: 2px;
			height: 3.5rem;
			background: rgba(0, 0, 0, 0.06);
		}
	}

	.active {
		@include dot;
		border: 2px solid $primary;
		background: $primary;
	}

	&-text {
		&-stat {
			font-size: 1rem;
			font-weight: 400;
			margin-bottom: 3px;
		}

		&-sub {
			font-size: 0.9rem;
			font-weight: 400;
		}
	}
}

.order-track {
	transition: all 0.3s height 0.3s;
	transform-origin: top center;
	// transform: scale(0);
	// height: 0;
}

.order-track-text {
	// display: flex;
	// justify-content: space-between;
	// width: 100%;
	display: block;
	width: 100%;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 100%;
	display: inline-block;
	text-align: left;
}

.order-responive {
	overflow-y: auto;
	height: 430px;
}
</style>
