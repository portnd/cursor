<script setup lang="ts">
import { useDashboardConditionStore } from "../store"

const props = defineProps({
	collapsed: {
		type: Boolean,
		default: false,
	},
	searchCall: {
		type: Boolean,
		default: false,
	},
})

const histChart = ref()
const store = useDashboardConditionStore()

watch(
	() => [store.conditionType, props.searchCall],
	() => {
		// store.getCheckBox()
		// store.updateConditionChart(pieChart, barChart)

		store.getCondition()

		if (store.toggle.hist) {
			setTimeout(() => {
				store.updatehisChart(histChart)
			}, 100)
		}
	}
)

watch(
	() => store.toggle.hist,
	() => {
		if (store.toggle.hist) {
			setTimeout(() => {
				store.updatehisChart(histChart)
			}, 100)
		}
	}
)

// watch(
// 	() => store.menu,
// 	() => {
// 		if (store.menu === "condition") {
// 			// store.getCheckBox()
// 			// store.updateConditionChart(pieChart, barChart)
// 			if (dashboardStore.params.road_id.length === 1) {
// 				store.getRoads()
// 				console.log("test")
// 			} else {
// 				store.getCondition()
// 			}
// 		}
// 	}
// )

const handleCondition = (id: number) => {
	const index = store.conditionArray.indexOf(id)
	if (index === -1) {
		store.conditionArray.push(id)
	} else {
		store.conditionArray.splice(index, 1)
	}
	store.conditionArray.sort((a, b) => a - b)
	store.colors = store.conditionColors
		.flatMap((item, key) => {
			return store.conditionArray.map((id) => {
				if (key === id) {
					return item
				}
				return undefined
			})
		})
		.filter((item: any) => item !== undefined)
	store.labels = store.conditionLabel
		.flatMap((item, key) => {
			return store.conditionArray.map((id) => {
				if (key === id) {
					return item
				}
				return undefined
			})
		})
		.filter((item: any) => item !== undefined)
}

const handleHeader = computed(() => {
	const ChartArray = ref<Array<any>>([])

	if (store.data.chart) {
		store.data.chart.lable.forEach((item: string) => {
			ChartArray.value.push({ text: item, color: "" })
		})
		store.data.chart.color.forEach((item: string, index: number) => {
			if (ChartArray.value[index]) {
				ChartArray.value[index].color = item
			}
		})
	}
	return ChartArray
})
</script>
<template>
	<div class="row">
		<div class="col-xl-6 col-lg-12 col-md-6 col-12 mb-5">
			<div class="card-chart h-100 text-center p-2 pt-4">
				<ClientOnly>
					<apexchart type="pie" height="275" :options="store.pieOptions()" :series="store.pieSeries()" />
				</ClientOnly>
				<div v-if="store.data.table" class="row justify-content-center mt-8">
					<template v-for="(item, key) in handleHeader.value" :key="key">
						<div
							class="col-auto d-flex align-items-center cursor-pointer"
							:class="store.conditionArray.includes(key) ? '' : 'selected'"
							@click="handleCondition(key)"
						>
							<div class="square my-2 me-2" :style="`background: ${item.color}`"></div>
							<span>{{ item.text }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
		<div class="col-xl-6 col-lg-12 col-md-6 col-12 mb-5">
			<div class="card-chart h-100 text-center p-2 pt-4">
				<ClientOnly>
					<apexchart type="bar" height="315" :options="store.barOption()" :series="store.barSeries()" />
				</ClientOnly>
				<div v-if="store.data.table" class="row justify-content-center">
					<template v-for="(item, key) in handleHeader.value" :key="key">
						<div
							class="col-auto d-flex align-items-center cursor-pointer"
							:class="store.conditionArray.includes(key) ? '' : 'selected'"
							@click="handleCondition(key)"
						>
							<div class="square my-2 me-2" :style="`background: ${item.color}`"></div>
							<span>{{ item.text }}</span>
						</div>
					</template>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped lang="scss">
.nav-line-tabs .nav-item .nav-link {
	padding: 10px 16px;
}

.card-chart {
	border: solid 1px var(--kt-gray-300);
	border-radius: 1rem;
}

.square {
	width: 20px;
	height: 20px;
	border-radius: 5px;
	display: inline-block;
}

.selected {
	opacity: 0.2;
}
</style>
