<script setup lang="ts">
import { useDashboardStore } from "../store"
defineProps({
	mapCollapsed: {
		type: Boolean,
		default: false,
	},
})

const store = useDashboardStore()
const refAssetsChart = ref()

const summaryChart = reactive({
	chartOptions: {
		chart: {
			type: "bar",
			toolbar: {
				show: false,
			},
		},
		plotOptions: {
			bar: {
				horizontal: true,
				barHeight: "75%",
				dataLabels: {
					position: "top",
				},
			},
		},
		dataLabels: {
			enabled: true,
			offsetY: -1,
			offsetX: -28,
			formatter: function (value: any, opt: any) {
				const comineData = opt.w.globals.series.flatMap((value: number[]) => value)
				const max = Math.max(...comineData)
				const threshold = max / 4

				return value >= threshold ? opt.w.globals.seriesNames[opt.seriesIndex] : ""
			},
			style: {
				fontSize: "10px",
				fontWeight: 400,
				colors: ["#000"],
			},
		},
		xaxis: {
			categories: [],
			show: false,
			axisBorder: {
				show: true,
				color: "#EAEAEA",
			},
			axisTicks: {
				show: false,
			},
			labels: {
				show: true,
				formatter: function (value: any) {
					return toNumber(value)
					// if (value !== null) {
					// 	return toNumber(value)
					// }
				},
			},
		},
		yaxis: {
			axisBorder: {
				show: false,
				color: "#EAEAEA",
			},
			axisTicks: {
				show: false,
			},
		},
		legend: {
			show: false,
			// position: "bottom",
			// offsetX: -10,
			// offsetY: -5,
			// horizontalAlign: "center",
			// itemMargin: {
			// 	horizontal: 5,
			// },
			// markers: {
			// 	radius: 2,
			// },
		},
		grid: {
			borderColor: "#EAEAEA",
			strokeDashArray: 5,
			xaxis: {
				lines: {
					show: false,
				},
			},
			yaxis: {
				lines: {
					show: true,
				},
			},
		},
		fill: {
			opacity: 1,
		},
		tooltip: {
			shared: true,
			intersect: false,
			custom: function ({ series, dataPointIndex, w }: any) {
				let html = `<div class="apexcharts-result">`
				html += `<label class="fs-7 border-bottom  mb-2">${w.config.xaxis?.categories[dataPointIndex]}</label>`
				w.globals.seriesNames?.forEach((label: string, index: number) => {
					html += `<div class="d-flex align-items-center">`
					html += ``
					html += `<div class="fs-7 d-flex align-items-center "><div class="me-2" style="background-color: ${
						w.globals.colors[index]
					}; width: 14px; height: 14px; border-radius: 50%;"></div>${label}:</div>
          <span class="fs-7 ms-1">${toNumber(series[index][dataPointIndex]) || 0} </span>`
					html += `</div>`
				})
				html += `</div>`
				return html
			},
		},
		colors: ["#FDB833", "#FFDA6A"],
		responsive: [
			{
				breakpoint: 576,
				options: {
					dataLabels: {
						enabled: false,
					},
				},
			},
		],
	},
})

const showAsset = (id: number) => {
	const index = store.params.ref_asset_id.indexOf(id)
	if (index === -1) {
		store.params.ref_asset_id.push(id)
	} else {
		store.params.ref_asset_id.splice(index, 1)
	}

	store.params.ref_asset_id = [...new Set(store.params.ref_asset_id)]
}

watch(
	() => store.data.assets,
	() => {
		refAssetsChart.value?.updateOptions({
			xaxis: {
				categories: store.getAssetCategories,
				offsetX: 50,
			},
		})
	}
)
</script>
<template>
	<div class="row">
		<div class="col-auto d-flex align-items-center"><span>การแสดงผล</span></div>
		<div class="col-xl-3 col-md-4 col-sm-4 col-6">
			<VSelect
				v-model="store.display_type"
				:options="[
					{ label: 'ภาพรวม', value: 1 },
					{ label: 'รายละเอียด', value: 2 },
				]"
				label=""
				name="display_type"
				:close-on-select="true"
				:can-deselect="false"
				:can-clear="false"
			/>
		</div>
	</div>
	<div v-show="store.display_type === 1" class="row">
		<div class="col-sm-1 col-2 align-items-center d-flex" style="margin-top: 2.5em; margin-bottom: 3.5em">
			<div class="row d-flex align-self-stretch">
				<div class="col-12 checkbox-group">
					<VCheckbox
						v-model="store.params.ref_asset_id"
						:options="store.getAssetItemsOptions"
						:name="`ref_asset_id`"
						@update:model-value="store.onChecked"
					/>
				</div>
			</div>
		</div>
		<div class="col-sm-11 col-10 ps-0">
			<ClientOnly>
				<apexchart
					id="summary"
					ref="refAssetsChart"
					height="100%"
					:options="summaryChart.chartOptions"
					:series="store.getAssetSeries"
				/>
			</ClientOnly>
		</div>
	</div>
	<div v-show="store.display_type === 2" class="row pt-5">
		<template v-for="(item, key) in store.data?.assets_details" :key="key">
			<div class="col-12 mb-8">
				<div class="card card-group px-0">
					<div
						class="card-header fw-semibold cursor-pointer"
						:style="{
							backgroundColor: item.is_active === true ? item.colors : '#D9DBE4',
							color: item.is_active === true ? '#3F4254' : '#7E8299',
						}"
						@click="store.setParentActive(item)"
					>
						{{ item.name }}
					</div>
					<div
						class="card-body pt-0 pb-6 px-4"
						:style="
							item.is_active === true ? `background: ${convertHexToRGBA(item.colors, 0.2)}` : `background: #F9F9F9`
						"
					>
						<div class="row">
							<template v-for="(child, childKey) in item.items" :key="childKey">
								<div
									class="col-12 px-7"
									:class="[!mapCollapsed ? 'col-xl-4 col-lg-6 col-md-4 col-sm-6' : 'col-xl-3 col-md-4 col-sm-6']"
									@click="
										() => {
											store.setChildActive(child)
										}
									"
								>
									<div
										class="card card-item m-auto p-3 mt-5"
										:style="
											child.is_active ? `border-left: 5px solid ${item.colors}` : `border-left: 5px solid #D9DBE4`
										"
										@click="showAsset(child.id)"
									>
										<span class="fs-6" :style="{ color: child.is_active ? '#3F4254' : '#D9DBE4' }">
											{{ child.name }}</span
										>
										<div class="d-flex justify-content-between my-auto">
											<img :src="store.generateImage(child)" alt="" />
											<div class="row text-end">
												<div class="col-12">
													<h1
														class="fw-bold mb-0"
														:style="child.is_active ? `color: ${item.colors}` : `color: #D9DBE4`"
													>
														{{ toNumber(child.value) }}
													</h1>
													<span
														class="fs-6"
														:style="{
															color: child.is_active ? '#3F4254' : '#D9DBE4',
														}"
													>
														จุด
													</span>
												</div>
											</div>
										</div>
									</div>
								</div>
							</template>
						</div>
					</div>
				</div>
			</div>
		</template>
	</div>
</template>

<style lang="scss" scoped>
.checkbox-group {
	display: flex;
	align-items: center;
}
.card-group {
	.card-header {
		width: 100%;
		min-height: 40px;
		border-radius: 10px 10px 0px 0px;
		align-items: center;
		font-size: 13px;
		color: var(--kt-white);
		font-weight: 400;
		justify-content: center;
		border-bottom: 0;
	}
	.card-body {
		.card-item {
			border-radius: 0px 10px 10px 0px !important;
			cursor: pointer;
			background: #fff;
			height: 110px;
			filter: drop-shadow(0px 4px 4px rgba(0, 0, 0, 0.25));
			@media screen and (max-width: 992px) {
				width: auto;
			}
			img {
				width: 50px;
			}
		}
	}
}
.col-lg-12.col-xl-12 {
	.card-item {
		width: 220px !important;
	}
}
@media screen and (min-width: 992px) and (max-width: 1096px) {
	.card-body .col-lg-6 {
		width: 100%;
	}
	.col-xl-4.col-lg-6 {
		.card-item {
			width: 220px;
		}
	}
}
</style>
