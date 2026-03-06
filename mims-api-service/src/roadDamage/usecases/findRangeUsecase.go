package usecases

import (
	"fmt"
	"math"
	"strings"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
)

func FindRange1(kmData []float64, key string, m_start_r, m_end_r, direction int, rcData requests.RcData, mItem requests.MItem, roadGeom models.RoadGeom, mIdx string, item requests.RoadDamageCsv, rangeItems map[string]requests.MItem) requests.RoadRangeItems {
	var roadRange requests.RoadRangeItems
	m_start_r = int(kmData[0])
	m_end_r = int(kmData[len(kmData)-1])
	subRoad_start := math.Abs(float64(m_start_r)-float64(rcData.KmStart)) / math.Abs(rcData.KmStart-rcData.KmEnd)
	subRoad_end := math.Abs(float64(m_end_r)-float64(rcData.KmEnd)) / math.Abs(rcData.KmStart-rcData.KmEnd)
	fmt.Println(m_start_r, m_end_r, key)
	subRoad_min := helpers.Min(subRoad_start, subRoad_end)
	subRoad_max := helpers.Max(subRoad_start, subRoad_end)
	// helpers.PrintlnJson(subRoad_min, subRoad_max)

	mItem.KmStart = float64(m_start_r)
	mItem.KmEnd = float64(m_end_r)
	if subRoad_max > 1 {
		subRoad_max = 1
	}

	mItem.TheGeom = "ST_LineSubstring('" + roadGeom.TheGeom + "', " + fmt.Sprintf("%f", subRoad_min) + ", " + fmt.Sprintf("%f", subRoad_max) + ")"
	rangeItems[mIdx] = mItem
	roadRange.MItem = mItem
	surfaceType := strings.Split(key, "_")
	roadRange.SurveyType = surfaceType[0]
	roadRange.AcUcrack += item.AcUcrack
	switch direction {
	case 1:

		if roadRange.KmEnd == roadGeom.KmStart {
			roadRange.KmStart = roadGeom.KmStart
			roadRange.KmEnd = roadGeom.KmEnd + 100
		}

		if float64(m_end_r) <= rcData.KmEnd {
		} else {
			if roadRange.KmStart == roadGeom.KmEnd {
				roadRange.KmStart = roadRange.KmStart - 100
			}
			roadRange.KmEnd = roadGeom.KmEnd
		}
	case 2:

		if roadRange.KmStart == roadGeom.KmEnd {
			roadRange.KmEnd = roadGeom.KmEnd
			roadRange.KmStart = roadGeom.KmStart - 100
		}
		// helpers.PrintlnJson("m_start_r", m_start_r, m_end_r, rcData.KmStart, rcData.KmEnd)
		if float64(m_start_r) <= rcData.KmStart {
			// roadRanges = append(roadRanges, roadRange)
		} else {
			if roadRange.KmStart == roadGeom.KmEnd {
				roadRange.KmStart = roadRange.KmStart + 100
			}
			roadRange.KmStart = roadGeom.KmStart
			// roadRanges = append(roadRanges, roadRange)
		}

	}
	return roadRange
}

func FindRange2(isFixData bool, key string, m_start_r, m_end_r, direction int, rcData requests.RcData, mItem requests.MItem, roadGeom models.RoadGeom, mIdx string, item requests.RoadDamageCsv, rangeItems map[string]requests.MItem) requests.RoadRangeItems {
	var roadRange requests.RoadRangeItems

	subRoad_start := math.Abs(float64(m_start_r)-float64(rcData.KmStart)) / math.Abs(rcData.KmStart-rcData.KmEnd)
	subRoad_end := math.Abs(float64(m_end_r)-float64(rcData.KmEnd)) / math.Abs(rcData.KmStart-rcData.KmEnd)

	subRoad_min := helpers.Min(subRoad_start, subRoad_end)
	subRoad_max := helpers.Max(subRoad_start, subRoad_end)
	// helpers.PrintlnJson(subRoad_min, subRoad_max)

	mItem.KmStart = float64(m_start_r)
	mItem.KmEnd = float64(m_end_r)
	if subRoad_max > 1 {
		subRoad_max = 1
	}

	mItem.TheGeom = "ST_LineSubstring('" + roadGeom.TheGeom + "', " + fmt.Sprintf("%f", subRoad_min) + ", " + fmt.Sprintf("%f", subRoad_max) + ")"
	rangeItems[mIdx] = mItem
	roadRange.MItem = mItem
	surfaceType := strings.Split(key, "_")
	helpers.PrintlnJson("surfaceTypesurfaceTypesurfaceType", surfaceType)
	roadRange.SurveyType = surfaceType[0]
	roadRange.AcUcrack += item.AcUcrack
	switch direction {
	case 1:

		if roadRange.KmEnd == roadGeom.KmStart {
			roadRange.KmStart = roadGeom.KmStart
			if !isFixData {
				roadRange.KmEnd = roadGeom.KmEnd + 100 // ตรงนี้
			}

		}

		if float64(m_end_r) <= rcData.KmEnd {
			// roadRanges = append(roadRanges, roadRange)
		} else {
			if roadRange.KmStart == roadGeom.KmEnd {
				roadRange.KmStart = roadRange.KmStart - 100
			}
			roadRange.KmEnd = roadGeom.KmEnd
			// roadRanges = append(roadRanges, roadRange)
		}
	case 2:

		if roadRange.KmStart == roadGeom.KmEnd {
			roadRange.KmEnd = roadGeom.KmEnd
			roadRange.KmStart = roadGeom.KmStart - 100
		}
		// helpers.PrintlnJson("m_start_r", m_start_r, m_end_r, rcData.KmStart, rcData.KmEnd)
		if float64(m_start_r) <= rcData.KmStart {
			// roadRanges = append(roadRanges, roadRange)
		} else {
			if roadRange.KmStart == roadGeom.KmEnd {
				roadRange.KmStart = roadRange.KmStart + 100
			}
			roadRange.KmStart = roadGeom.KmStart
			// roadRanges = append(roadRanges, roadRange)
		}

	}
	return roadRange
}

func RoadRangeDatas(roadRanges []requests.RoadRangeItems) []requests.RoadRangeItems {
	var data []requests.RoadRangeItems
	for index, item := range roadRanges {
		kmEnd := int(item.KmEnd)
		originalNumber := kmEnd
		modTarget := 100
		addition := modTarget - (originalNumber % modTarget)
		newNumber := originalNumber + addition
		if index == len(roadRanges)-1 {
			item.KmEnd = float64(kmEnd)
		} else {
			if roadRanges[index+1].SurveyType == item.SurveyType { // && (roadRanges[index+1].KmStart != float64(newNumber)) {
				if int(item.KmEnd)%100 != 0 {
					item.KmEnd = float64(newNumber)
				} else {
					item.KmEnd = float64(kmEnd)
				}
			}
		}
		data = append(data, item)
	}
	return data
}

// func Test3(roadRanges []requests.RoadRangeItems) []requests.RoadRangeItems {
// 	var data []requests.RoadRangeItems
// 	for index, item := range roadRanges {
// 		kmStart := int(item.KmStart)
// 		if index == 0 {
// 			item.KmStart = float64(kmStart)
// 		} else {
// 			originalNumber := kmStart
// 			subtract := originalNumber % 100
// 			newNumber := originalNumber - subtract
// 			if newNumber == int(data[index-1].KmStart) {
// 				if data[index-1].SurveyType != item.SurveyType {
// 					item.KmStart = item.KmStart
// 				} else {
// 					item.KmStart = data[index-1].KmEnd
// 				}

// 			} else {
// 				// helpers.PrintlnJson(data[index-1].KmStart, item.KmStart)
// 				if data[index-1].KmStart == item.KmStart {
// 					item.KmStart = float64(kmStart)
// 				} else {
// 					if data[index-1].SurveyType == item.SurveyType {
// 						item.KmStart = float64(newNumber)
// 					} else {
// 						item.KmStart = float64(kmStart)
// 					}

// 				}

// 			}

// 		}
// 		data = append(data, item)
// 	}

// 	var data2 []requests.RoadRangeItems
// 	for index, item := range data {

// 		kmEnd := int(item.KmEnd)
// 		// originalNumber := kmEnd
// 		// modTarget := 100                                     // We want the result after adding to be a multiple of 100
// 		// addition := modTarget - (originalNumber % modTarget) // Calculate the required addition to get to the next multiple of 100
// 		// newNumber := originalNumber + addition

// 		if index == len(roadRanges)-1 {
// 			item.KmEnd = float64(kmEnd)
// 		} else {

// 			// if item.KmEnd == data[index+1].KmStart {
// 			// helpers.PrintlnJson(item.KmStart, item.KmEnd, data[index+1].SurveyType, item.SurveyType)
// 			// if data[index+1].SurveyType == item.SurveyType {
// 			// 	item.KmEnd = float64(newNumber)
// 			// } else {
// 			// 	// if item.KmEnd == data[index+1].KmStart && item.KmEnd == data[index+1].KmEnd {
// 			// 	item.KmEnd = float64(newNumber)
// 			// } else {
// 			item.KmEnd = float64(kmEnd)
// 			// }

// 			// }
// 			// } else {
// 			// 	// if data[index+1].SurveyType == item.SurveyType {

// 			// 	// 	item.KmEnd = float64(newNumber)
// 			// 	// } else {
// 			// 	item.KmEnd = float64(kmEnd)
// 			// 	// }
// 			// }

// 		}
// 		data2 = append(data2, item)
// 	}

// 	var data3 []requests.RoadRangeItems
// 	for index, item := range data2 {
// 		if item.KmEnd == item.KmStart {
// 			if index == 0 {
// 				item.KmStart = item.KmStart
// 			} else {
// 				if item.KmStart == data[index-1].KmEnd {
// 					item.KmStart = item.KmStart
// 				} else {
// 					item.KmStart = item.KmStart - 100
// 				}

// 			}
// 		}

// 		data3 = append(data3, item)
// 	}

// 	// var data3 []requests.RoadRangeItems
// 	// for index, item := range data2 {
// 	// 	if item.KmEnd == item.KmStart {

// 	// 		if index == 0 {
// 	// 			item.KmStart = item.KmStart
// 	// 		} else {
// 	// 			if item.KmStart == data[index-1].KmEnd {
// 	// 				item.KmStart = item.KmStart
// 	// 			} else {
// 	// 				item.KmStart = item.KmStart - 100
// 	// 			}

// 	// 		}
// 	// 	}

// 	// 	data3 = append(data3, item)
// 	// }

// 	return data3
// }

func Test3(roadRanges []requests.RoadRangeItems) []requests.RoadRangeItems {
	var data []requests.RoadRangeItems
	for index, item := range roadRanges {
		if item.KmEnd == item.KmStart {

			if index == 0 {
				item.KmStart = item.KmStart
			} else {
				if item.KmStart == roadRanges[index-1].KmEnd {
					item.KmStart = item.KmStart
				} else {
					item.KmStart = item.KmStart - 100
				}

			}
		}

		data = append(data, item)
	}

	return data
}
