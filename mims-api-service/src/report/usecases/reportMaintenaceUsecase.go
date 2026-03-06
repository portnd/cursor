package usecases

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/jinzhu/copier"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (u *UseCase) Report13(req requests.Report13) (interface{}, error) {
	start := req.YearStart
	end := req.YearEnd
	typ := req.TypeReport
	groupYear := req.Group
	dataRes, detail, err := u.Repo.GetReportMaintenance(req.YearStart, req.YearEnd, req.RoadSectionId)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return nil, err
	}

	var data []models.DataReportMaintenance
	copier.Copy(&data, &dataRes)

	if len(data) == 0 || len(detail) == 0 {
		newData, err := u.Repo.GetMultiRoadInfo(req.RoadSectionId)
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			logs.Error(err)
			return responses.RoadConditionDetails{}, err
		}

		for _, i := range newData {
			data = append(data, models.DataReportMaintenance{
				RoadID:     i.RoadID,
				RoadName:   i.RoadName,
				SecKmStrat: i.KmStart,
				SecKmEnd:   i.KmStart,
				StrKmStart: helpers.FormatKM(int64(i.KmStart)),
				StrKmEnd:   helpers.FormatKM(int64(i.KmEnd)),
				KmStart:    i.KmStart,
				KmEnd:      i.KmEnd,
				IsNull:     true,
			})
		}

		for index, i := range data {

			yearGroup := start
			group := groupYear

			if typ == "excel" {
				var row []models.DataRowMaintenance
				s := i.KmStart

				if i.KmStart < i.KmEnd {
					for j := 0; s < i.KmEnd; j++ {
						e := s + 200

						// e must be not higher than real km_end
						if e > i.KmEnd {
							e = i.KmEnd
						}

						row = append(row, models.DataRowMaintenance{
							KmStart:    s,
							KmEnd:      e,
							StrKmStart: helpers.FormatKM(int64(s)),
							StrKmEnd:   helpers.FormatKM(int64(e)),
						})

						// add year
						yearStart := yearGroup
						var years []models.DataYearMaintenance

						for k := 0; yearStart <= end; k++ {
							years = append(years, models.DataYearMaintenance{
								Year: yearStart + 543,
							})

							// add lane
							lane := 1
							var lanes []models.DataLaneMaintenance

							for lane <= 6 {

								// add deatil
								var methodDetail models.DataDetailMaintenance

								lanes = append(lanes, models.DataLaneMaintenance{
									LaneNo: lane,
									Data:   methodDetail,
								})
								lane++
							}

							years[k].Data = lanes
							yearStart++
						}

						row[j].Data = years

						s += 200
					}

				} else {
					for j := 0; i.KmEnd < s; j++ {
						e := s - 200

						// e must be not lower than real km_end
						if e < i.KmEnd {
							e = i.KmEnd
						}

						row = append(row, models.DataRowMaintenance{
							KmStart:    s,
							KmEnd:      e,
							StrKmStart: helpers.FormatKM(int64(s)),
							StrKmEnd:   helpers.FormatKM(int64(e)),
						})

						// add year
						yearStart := yearGroup
						var years []models.DataYearMaintenance

						for k := 0; yearStart < yearGroup+group && yearStart <= end; k++ {
							years = append(years, models.DataYearMaintenance{
								Year: yearStart + 543,
							})

							// add lane
							lane := 1
							var lanes []models.DataLaneMaintenance

							for lane <= 6 {

								// add deatil
								var methodDetail models.DataDetailMaintenance
								lanes = append(lanes, models.DataLaneMaintenance{
									LaneNo: lane,
									Data:   methodDetail,
								})
								lane++
							}

							years[k].Data = lanes
							yearStart++
						}

						row[j].Data = years
						s -= 200
					}
				}

				data[index].RowExcel = append(data[index].RowExcel, row...)

			} else {

				// loop for years group
				for yearGroup <= end {
					var row []models.DataRowMaintenance
					s := i.KmStart

					if i.KmStart < i.KmEnd {

						//main logic
						for j := 0; s < i.KmEnd; j++ {
							e := s + 200

							// e must be not higher than real km_end
							if e > i.KmEnd {
								e = i.KmEnd
							}

							row = append(row, models.DataRowMaintenance{
								KmStart:    s,
								KmEnd:      e,
								StrKmStart: helpers.FormatKM(int64(s)),
								StrKmEnd:   helpers.FormatKM(int64(e)),
							})

							// add year
							yearStart := yearGroup
							var years []models.DataYearMaintenance

							for k := 0; yearStart < yearGroup+group && yearStart <= end; k++ {
								years = append(years, models.DataYearMaintenance{
									Year: yearStart + 543,
								})

								// add lane
								lane := 1
								var lanes []models.DataLaneMaintenance

								for lane <= 6 {

									// add deatil
									var methodDetail models.DataDetailMaintenance
									lanes = append(lanes, models.DataLaneMaintenance{
										LaneNo: lane,
										Data:   methodDetail,
									})
									lane++
								}

								years[k].Data = lanes
								yearStart++
							}

							row[j].Data = years

							s += 200
						}

					} else {
						for j := 0; i.KmEnd < s; j++ {
							e := s - 200

							// e must be not lower than real km_end
							if e < i.KmEnd {
								e = i.KmEnd
							}

							row = append(row, models.DataRowMaintenance{
								KmStart:    s,
								KmEnd:      e,
								StrKmStart: helpers.FormatKM(int64(s)),
								StrKmEnd:   helpers.FormatKM(int64(e)),
							})

							// add year
							yearStart := yearGroup
							var years []models.DataYearMaintenance

							for k := 0; yearStart < yearGroup+group && yearStart <= end; k++ {
								years = append(years, models.DataYearMaintenance{
									Year: yearStart + 543,
								})

								// add lane
								lane := 1
								var lanes []models.DataLaneMaintenance

								for lane <= 6 {

									// add deatil
									var methodDetail models.DataDetailMaintenance
									lanes = append(lanes, models.DataLaneMaintenance{
										LaneNo: lane,
										Data:   methodDetail,
									})
									lane++
								}

								years[k].Data = lanes
								yearStart++
							}

							row[j].Data = years
							s -= 200
						}
					}

					yearGroup += group

					newRows := [][]models.DataRowMaintenance{row}
					data[index].Rows = append(data[index].Rows, newRows[0])
				}
			}

			budgetYear := start
			for budgetYear <= end {
				data[index].Years = append(data[index].Years, budgetYear+543)
				budgetYear++
			}

			data[0].YearStart = start + 543
			data[0].YearEnd = end + 543
		}

	} else {
		// sort by range
		for iIndex, i := range detail {
			detail[iIndex].Range = int(math.Abs(float64(i.KmEnd) - float64(i.KmStart)))
		}

		sort.Slice(detail, func(i, j int) bool {
			return detail[i].Range < detail[j].Range
		})

		for index, i := range data {

			yearGroup := start
			group := groupYear

			if typ == "excel" {
				var row []models.DataRowMaintenance
				s := i.KmStart

				if i.KmStart < i.KmEnd {

					//reservation for span
					var reserv []models.RangeMaintenance
					for r := i.KmStart; r < i.KmEnd; r += 200 {

						e := r + 200
						// e must be not higher than real km_end
						if e > i.KmEnd {
							e = i.KmEnd
						}

						for y := start; y <= end; y++ {
							for l := 1; l <= 6; l++ {
								reserv = append(reserv, models.RangeMaintenance{
									Start: r,
									End:   e,
									Value: 0,
									Year:  y,
									Lane:  l,
								})
							}
						}
					}

					// add index reserveration
					for _, d := range detail {
						for rIndex, r := range reserv {
							if r.Year == d.Year && r.Lane == d.Lane {
								if d.KmStart < r.End && (r.End-200) < d.KmEnd { //เช็คว่าอยู่ใน range ไหม
									if r.Value == 0 {
										reserv[rIndex].Value = 1
										reserv[rIndex].Range = d.Range
										reserv[rIndex].Method = d.Method
									}
								}

							}
						}
					}

					//main logic
					for j := 0; s < i.KmEnd; j++ {
						e := s + 200

						// e must be not higher than real km_end
						if e > i.KmEnd {
							e = i.KmEnd
						}

						row = append(row, models.DataRowMaintenance{
							KmStart:    s,
							KmEnd:      e,
							StrKmStart: helpers.FormatKM(int64(s)),
							StrKmEnd:   helpers.FormatKM(int64(e)),
						})

						// add year
						yearStart := yearGroup
						var years []models.DataYearMaintenance

						for k := 0; yearStart <= end; k++ {
							years = append(years, models.DataYearMaintenance{
								Year: yearStart + 543,
							})

							// add lane
							lane := 1
							var lanes []models.DataLaneMaintenance

							for lane <= 6 {

								// add deatil
								var methodDetail models.DataDetailMaintenance

								for _, idetail := range detail {

									if lane == idetail.Lane && i.RoadID == idetail.RoadID && yearStart == idetail.Year {

										if (idetail.KmStart-s) < 200 && idetail.KmEnd > s { //เช็คระยะ

											span := 0
											for rIndex, r := range reserv {
												if r.Start >= s {
													if r.Year == yearStart && r.Lane == lane {
														if idetail.KmStart < r.End && (r.End-200) < idetail.KmEnd {

															if r.Range != idetail.Range || r.Method != idetail.Method {
																break
															}

															if r.Range == idetail.Range && r.Method == idetail.Method {
																if r.Writed == 0 {
																	span++
																	reserv[rIndex].Writed = 1
																}
															}
														}
													}
												}
											}

											if span != 0 {
												methodDetail = models.DataDetailMaintenance{
													Km:      s,
													KmStart: helpers.FormatKM(int64(s)),
													KmEnd:   helpers.FormatKM(int64(e)),
													LaneNo:  lane,
													Span:    0,
													IsWrite: 1,
													Year:    yearStart,
													Data: models.MethodMaintenance{
														KmStart:    idetail.KmStart,
														KmEnd:      idetail.KmEnd,
														StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
														StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
														Method:     idetail.Method,
														Color:      ICColour(idetail.Method),
														Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
													},
												}

											} else if methodDetail.KmStart == "" {
												methodDetail = models.DataDetailMaintenance{
													Km:      s,
													KmStart: helpers.FormatKM(int64(s)),
													KmEnd:   helpers.FormatKM(int64(e)),
													LaneNo:  lane,
													Span:    0,
													Year:    yearStart,
													Data: models.MethodMaintenance{
														KmStart:    idetail.KmStart,
														KmEnd:      idetail.KmEnd,
														StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
														StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
														Method:     idetail.Method,
														Color:      ICColour(idetail.Method),
														Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
													},
												}
											}
										}
									}
								}
								lanes = append(lanes, models.DataLaneMaintenance{
									LaneNo: lane,
									Data:   methodDetail,
								})
								lane++
							}

							years[k].Data = lanes
							yearStart++
						}

						row[j].Data = years

						s += 200
					}

				} else {

					// reservation for span
					var reserv []models.RangeMaintenance
					for r := i.KmStart; r > i.KmEnd; r -= 200 {

						e := r - 200
						// e must be not higher than real km_end
						if e < i.KmEnd {
							e = i.KmEnd
						}

						for y := start; y <= end; y++ {
							for l := 1; l <= 6; l++ {
								reserv = append(reserv, models.RangeMaintenance{
									Start: r,
									End:   e,
									Value: 0,
									Year:  y,
									Lane:  l,
								})
							}
						}
					}

					// add index reserveration
					for _, d := range detail {
						for rIndex, r := range reserv {
							if r.Year == d.Year && r.Lane == d.Lane {
								if d.KmStart > r.End && (r.End+200) > d.KmEnd { //เช็คว่าอยู่ใน range ไหม
									if r.Value == 0 {
										reserv[rIndex].Value = 1
										reserv[rIndex].Range = d.Range
										reserv[rIndex].Method = d.Method
									}
								}
							}
						}
					}

					for j := 0; i.KmEnd < s; j++ {
						e := s - 200

						// e must be not lower than real km_end
						if e < i.KmEnd {
							e = i.KmEnd
						}

						row = append(row, models.DataRowMaintenance{
							KmStart:    s,
							KmEnd:      e,
							StrKmStart: helpers.FormatKM(int64(s)),
							StrKmEnd:   helpers.FormatKM(int64(e)),
						})

						// add year
						yearStart := yearGroup
						var years []models.DataYearMaintenance

						for k := 0; yearStart <= end; k++ {
							years = append(years, models.DataYearMaintenance{
								Year: yearStart + 543,
							})

							// add lane
							lane := 1
							var lanes []models.DataLaneMaintenance

							for lane <= 6 {

								// add deatil
								var methodDetail models.DataDetailMaintenance

								for _, idetail := range detail {
									if lane == idetail.Lane && i.RoadID == idetail.RoadID && yearStart == idetail.Year {

										if (s-idetail.KmStart) < 200 && idetail.KmEnd < s { //เช็คระยะ

											span := 0
											for rIndex, r := range reserv {
												if r.Start <= s {
													if r.Year == yearStart && r.Lane == lane {
														if idetail.KmStart > r.End && (r.End+200) > idetail.KmEnd {

															if r.Range != idetail.Range || r.Method != idetail.Method {
																break
															}

															if r.Range == idetail.Range && r.Method == idetail.Method {
																if r.Writed == 0 {
																	span++
																	reserv[rIndex].Writed = 1
																}
															}
														}
													}
												}
											}

											if span != 0 {
												methodDetail = models.DataDetailMaintenance{
													Km:      s,
													KmStart: helpers.FormatKM(int64(s)),
													KmEnd:   helpers.FormatKM(int64(e)),
													LaneNo:  lane,
													Span:    span,
													IsWrite: 1,
													Year:    yearStart,
													Data: models.MethodMaintenance{
														KmStart:    idetail.KmStart,
														KmEnd:      idetail.KmEnd,
														StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
														StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
														Method:     idetail.Method,
														Color:      ICColour(idetail.Method),
														Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
													},
												}
											} else if methodDetail.KmStart == "" {
												methodDetail = models.DataDetailMaintenance{
													Km:      s,
													KmStart: helpers.FormatKM(int64(s)),
													KmEnd:   helpers.FormatKM(int64(e)),
													LaneNo:  lane,
													Span:    0,
													Year:    yearStart,
													Data: models.MethodMaintenance{
														KmStart:    idetail.KmStart,
														KmEnd:      idetail.KmEnd,
														StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
														StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
														Method:     idetail.Method,
														Color:      ICColour(idetail.Method),
														Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
													},
												}
											}
										}
									}
								}

								lanes = append(lanes, models.DataLaneMaintenance{
									LaneNo: lane,
									Data:   methodDetail,
								})
								lane++
							}

							years[k].Data = lanes
							yearStart++
						}

						row[j].Data = years
						s -= 200
					}
				}

				data[index].RowExcel = append(data[index].RowExcel, row...)

			} else {

				// loop for years group
				for yearGroup <= end {
					var row []models.DataRowMaintenance
					s := i.KmStart

					if i.KmStart < i.KmEnd {

						//reservation for span
						var reserv []models.RangeMaintenance
						for r := i.KmStart; r < i.KmEnd; r += 200 {

							e := r + 200
							// e must be not higher than real km_end
							if e > i.KmEnd {
								e = i.KmEnd
							}

							for y := start; y <= end; y++ {
								for l := 1; l <= 6; l++ {
									reserv = append(reserv, models.RangeMaintenance{
										Start: r,
										End:   e,
										Value: 0,
										Year:  y,
										Lane:  l,
									})
								}
							}
						}

						// add index reserveration
						for _, d := range detail {
							for rIndex, r := range reserv {
								if r.Year == d.Year && r.Lane == d.Lane {
									if d.KmStart < r.End && (r.End-200) < d.KmEnd { //เช็คว่าอยู่ใน range ไหม
										if r.Value == 0 {
											reserv[rIndex].Value = 1
											reserv[rIndex].Range = d.Range
											reserv[rIndex].Method = d.Method
										}
									}

								}
							}
						}

						//main logic
						for j := 0; s < i.KmEnd; j++ {
							e := s + 200

							// e must be not higher than real km_end
							if e > i.KmEnd {
								e = i.KmEnd
							}

							row = append(row, models.DataRowMaintenance{
								KmStart:    s,
								KmEnd:      e,
								StrKmStart: helpers.FormatKM(int64(s)),
								StrKmEnd:   helpers.FormatKM(int64(e)),
							})

							// add year
							yearStart := yearGroup
							var years []models.DataYearMaintenance

							for k := 0; yearStart < yearGroup+group && yearStart <= end; k++ {
								years = append(years, models.DataYearMaintenance{
									Year: yearStart + 543,
								})

								// add lane
								lane := 1
								var lanes []models.DataLaneMaintenance

								for lane <= 6 {

									// add deatil
									var methodDetail models.DataDetailMaintenance

									for _, idetail := range detail {

										if lane == idetail.Lane && i.RoadID == idetail.RoadID && yearStart == idetail.Year {

											if (idetail.KmStart-s) < 200 && idetail.KmEnd > s { //เช็คระยะ

												span := 0
												for rIndex, r := range reserv {
													if r.Start >= s {
														if r.Year == yearStart && r.Lane == lane {
															if idetail.KmStart < r.End && (r.End-200) < idetail.KmEnd {

																if r.Range != idetail.Range || r.Method != idetail.Method {
																	break
																}

																if r.Range == idetail.Range && r.Method == idetail.Method {
																	if r.Writed == 0 {
																		span++
																		reserv[rIndex].Writed = 1
																	}
																}
															}
														}
													}
												}

												if span != 0 {
													methodDetail = models.DataDetailMaintenance{
														Km:      s,
														KmStart: helpers.FormatKM(int64(s)),
														KmEnd:   helpers.FormatKM(int64(e)),
														LaneNo:  lane,
														Span:    span,
														IsWrite: 1,
														Year:    yearStart,
														Data: models.MethodMaintenance{
															KmStart:    idetail.KmStart,
															KmEnd:      idetail.KmEnd,
															StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
															StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
															Method:     idetail.Method,
															Color:      ICColour(idetail.Method),
															Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
														},
													}
												} else if methodDetail.KmStart == "" {
													methodDetail = models.DataDetailMaintenance{
														Km:      s,
														KmStart: helpers.FormatKM(int64(s)),
														KmEnd:   helpers.FormatKM(int64(e)),
														LaneNo:  lane,
														Span:    0,
														Year:    yearStart,
														Data: models.MethodMaintenance{
															KmStart:    idetail.KmStart,
															KmEnd:      idetail.KmEnd,
															StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
															StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
															Method:     idetail.Method,
															Color:      ICColour(idetail.Method),
															Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
														},
													}
												}
											}
										}
									}
									lanes = append(lanes, models.DataLaneMaintenance{
										LaneNo: lane,
										Data:   methodDetail,
									})
									lane++
								}

								years[k].Data = lanes
								yearStart++
							}

							row[j].Data = years

							s += 200
						}

					} else {

						// reservation for span
						var reserv []models.RangeMaintenance
						for r := i.KmStart; r > i.KmEnd; r -= 200 {

							e := r - 200
							// e must be not higher than real km_end
							if e < i.KmEnd {
								e = i.KmEnd
							}

							for y := start; y <= end; y++ {
								for l := 1; l <= 6; l++ {
									reserv = append(reserv, models.RangeMaintenance{
										Start: r,
										End:   e,
										Value: 0,
										Year:  y,
										Lane:  l,
									})
								}
							}
						}

						// add index reserveration
						for _, d := range detail {
							for rIndex, r := range reserv {
								if r.Year == d.Year && r.Lane == d.Lane {
									if d.KmStart > r.End && (r.End+200) > d.KmEnd { //เช็คว่าอยู่ใน range ไหม
										if r.Value == 0 {
											reserv[rIndex].Value = 1
											reserv[rIndex].Range = d.Range
											reserv[rIndex].Method = d.Method
										}
									}
								}
							}
						}

						for j := 0; i.KmEnd < s; j++ {
							e := s - 200

							// e must be not lower than real km_end
							if e < i.KmEnd {
								e = i.KmEnd
							}

							row = append(row, models.DataRowMaintenance{
								KmStart:    s,
								KmEnd:      e,
								StrKmStart: helpers.FormatKM(int64(s)),
								StrKmEnd:   helpers.FormatKM(int64(e)),
							})

							// add year
							yearStart := yearGroup
							var years []models.DataYearMaintenance

							for k := 0; yearStart < yearGroup+group && yearStart <= end; k++ {
								years = append(years, models.DataYearMaintenance{
									Year: yearStart + 543,
								})

								// add lane
								lane := 1
								var lanes []models.DataLaneMaintenance

								for lane <= 6 {

									// add deatil
									var methodDetail models.DataDetailMaintenance

									for _, idetail := range detail {
										if lane == idetail.Lane && i.RoadID == idetail.RoadID && yearStart == idetail.Year {

											if (s-idetail.KmStart) < 200 && idetail.KmEnd < s { //เช็คระยะ

												span := 0
												for rIndex, r := range reserv {
													if r.Start <= s {
														if r.Year == yearStart && r.Lane == lane {
															if idetail.KmStart > r.End && (r.End+200) > idetail.KmEnd {

																if r.Range != idetail.Range || r.Method != idetail.Method {
																	break
																}

																if r.Range == idetail.Range && r.Method == idetail.Method {
																	if r.Writed == 0 {
																		span++
																		reserv[rIndex].Writed = 1
																	}
																}
															}
														}
													}
												}

												if span != 0 {
													methodDetail = models.DataDetailMaintenance{
														Km:      s,
														KmStart: helpers.FormatKM(int64(s)),
														KmEnd:   helpers.FormatKM(int64(e)),
														LaneNo:  lane,
														Span:    span,
														IsWrite: 1,
														Year:    yearStart,
														Data: models.MethodMaintenance{
															KmStart:    idetail.KmStart,
															KmEnd:      idetail.KmEnd,
															StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
															StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
															Method:     idetail.Method,
															Color:      ICColour(idetail.Method),
															Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
														},
													}
												} else if methodDetail.KmStart == "" {
													methodDetail = models.DataDetailMaintenance{
														Km:      s,
														KmStart: helpers.FormatKM(int64(s)),
														KmEnd:   helpers.FormatKM(int64(e)),
														LaneNo:  lane,
														Span:    0,
														Year:    yearStart,
														Data: models.MethodMaintenance{
															KmStart:    idetail.KmStart,
															KmEnd:      idetail.KmEnd,
															StrKmStart: helpers.FormatKM(int64(idetail.KmStart)),
															StrKmEnd:   helpers.FormatKM(int64(idetail.KmEnd)),
															Method:     idetail.Method,
															Color:      ICColour(idetail.Method),
															Unique:     fmt.Sprint(yearStart) + fmt.Sprint(lane) + fmt.Sprint(idetail.KmStart) + fmt.Sprint(idetail.KmEnd),
														},
													}
												}
											}
										}
									}

									lanes = append(lanes, models.DataLaneMaintenance{
										LaneNo: lane,
										Data:   methodDetail,
									})
									lane++
								}

								years[k].Data = lanes
								yearStart++
							}

							row[j].Data = years
							s -= 200
						}
					}

					yearGroup += group

					newRows := [][]models.DataRowMaintenance{row}
					data[index].DistanceStr = fmt.Sprintf("%.3f", i.Distance)
					data[index].StrKmStart = helpers.FormatKM(int64(i.KmStart))
					data[index].StrKmEnd = helpers.FormatKM(int64(i.KmEnd))
					data[index].Rows = append(data[index].Rows, newRows[0])
				}
			}
			budgetYear := start
			for budgetYear <= end {
				data[index].Years = append(data[index].Years, budgetYear+543)
				budgetYear++
			}
			data[0].YearStart = start + 543
			data[0].YearEnd = end + 543

		}
	}

	// generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = ExportExcelType13(data)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(data, "TEMPLATE_GENARAL_TYPE13", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}

	return pathResult, nil
}

func ExportExcelType13(data []models.DataReportMaintenance) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE13_EXCEL_" + fmt.Sprint((data[0].YearEnd-data[0].YearStart)+1)

	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for i := 1; i <= 50; i++ {
		if i > len(data) {
			f.DeleteSheet("Sheet1 " + "(" + fmt.Sprint(i) + ")")
		}
	}
	startIndexData := 14
	for iIndex, i := range data {
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D5", "ปีงบประมาณ พ.ศ. "+fmt.Sprint(data[0].YearStart)+"-"+fmt.Sprint(data[0].YearEnd))
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D6", fmt.Sprintf(`หมายเลขทางหลวง : %s ตอนควบคุม : %s ชื่อสายทาง : %s`, i.RoadGroupNumber, i.RoadSectionNumber, i.RoadMainName))
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D7", fmt.Sprintf(`กม.เริ่มต้น %s กม.สิ้นสุด %s ระยะทาง %s กม.`, helpers.FormatKM(int64(i.SecKmStrat)), helpers.FormatKM(int64(i.SecKmEnd)), fmt.Sprintf("%.3f", i.Distance)))

		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B9", "ชื่อสายทาง : "+i.RoadName)

		yearColumn := []string{"D", "P", "AB", "AN", "AZ", "BL", "BX", "CJ", "CV", "DH"}
		laneColumn := [][]string{
			{"D", "F", "H", "J", "L", "N"},
			{"P", "R", "T", "V", "X", "Z"},
			{"AB", "AD", "AF", "AH", "AJ", "AL"},
			{"AN", "AP", "AR", "AT", "AV", "AX"},
			{"AZ", "BB", "BD", "BF", "BH", "BJ"},
			{"BL", "BN", "BP", "BR", "BT", "BV"},
			{"BX", "BZ", "CB", "CD", "CF", "CH"},
			{"CJ", "CL", "CN", "CP", "CR", "CT"},
			{"CV", "CX", "CZ", "DB", "DD", "DF"},
			{"DH", "DJ", "DL", "DN", "DP", "DR"}}

		for yIndex, y := range i.Years {
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", yearColumn[yIndex]+"10", y)
		}

		for jIndex, j := range i.RowExcel {

			if jIndex < len(i.RowExcel)-1 {
				f.DuplicateRow("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", startIndexData+jIndex)
			}

			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+fmt.Sprint(startIndexData+jIndex), j.StrKmStart+" - "+j.StrKmEnd)

			for yearIndex, k := range j.Data {
				for laneIndex, l := range k.Data {

					if l.Data.KmStart != "" {
						f.SetCellRichText("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex), []excelize.RichTextRun{
							{
								Text: l.Data.Data.StrKmStart + " - " + l.Data.Data.StrKmEnd + "\n",
							},
							{
								Text: l.Data.Data.Method,
							},
						})
					}

					preCell, _ := f.GetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex-1))
					currentCell, _ := f.GetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex))

					if preCell != "" && currentCell != "" && preCell == currentCell {

						f.MergeCell("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex))

						// style and wrap text

						wrapTextStyle, _ := f.NewStyle(&excelize.Style{
							Font: &excelize.Font{
								Size:   10,
								Family: "TH SarabunPSK",
							},
							Alignment: &excelize.Alignment{
								WrapText:   true,
								Horizontal: "center",
								Vertical:   "top",
							},
							Fill: excelize.Fill{
								Type:    "pattern",
								Color:   []string{l.Data.Data.Color},
								Pattern: 1,
							},
							Border: []excelize.Border{{
								Type:  "left",
								Color: "000000",
								Style: 1,
							}, {
								Type:  "top",
								Color: "000000",
								Style: 1,
							}, {
								Type:  "right",
								Color: "000000",
								Style: 1,
							}, {
								Type:  "bottom",
								Color: "000000",
								Style: 1,
							}},
						})

						f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex), wrapTextStyle)

						if jIndex != 0 {
							if data[iIndex].RowExcel[jIndex-1].Data[yearIndex].Data[laneIndex].Data.Span == 2 {
								f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", startIndexData+jIndex-1, 30)
								f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", startIndexData+jIndex, 30)
							}
						}

					} else if currentCell != "" {

						// style and wrap text
						wrapTextStyle, _ := f.NewStyle(&excelize.Style{
							Font: &excelize.Font{
								Size:   10,
								Family: "TH SarabunPSK",
							},
							Alignment: &excelize.Alignment{
								WrapText:   true,
								Horizontal: "center",
								Vertical:   "top",
							},
							Fill: excelize.Fill{
								Type:    "pattern",
								Color:   []string{l.Data.Data.Color},
								Pattern: 1,
							},
							Border: []excelize.Border{{
								Type:  "left",
								Color: "000000",
								Style: 1,
							}, {
								Type:  "top",
								Color: "000000",
								Style: 1,
							}, {
								Type:  "right",
								Color: "000000",
								Style: 1,
							}, {
								Type:  "bottom",
								Color: "000000",
								Style: 1,
							}},
						})

						f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex), laneColumn[yearIndex][laneIndex]+fmt.Sprint(startIndexData+jIndex), wrapTextStyle)

						if l.Data.Span == 1 {
							f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", startIndexData+jIndex, 60)
						}
					}
				}
			}
		}

		// f.SetSheetName("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", i.RoadName)
		helpers.AddFooter(f, i.RoadName)
		//AddFooter(f, "Sheet1 "+"("+fmt.Sprint(iIndex+1)+")")
	}
	existingSheetNames := make(map[string]bool)

	for indexPageName, i := range data {
		var name string
		if len([]rune(i.RoadName)) >= 25 {
			name = string([]rune(i.RoadName)[:25])
		} else {
			name = string(i.RoadName)
		}

		uniqueName := uniqueSheetName(name, existingSheetNames)
		f.SetSheetName("Sheet1 ("+fmt.Sprint(indexPageName+1)+")", uniqueName)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE13_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ICColour(ic string) string {
	// style and wrap text
	var color string
	switch ic {
	case "OL-Overlay":
		color = "#ff7a8d"
	case "M&OL-Mill&Overlay":
		color = "#ffb800"
	case "RCL-Recycling":
		color = "#af85ff"
	case "Rc-Reconstruction":
		color = "#b22727"
	case "SS-SlurrySeal":
		color = "#82E0AA"
	case "FDR":
		color = "#418fff"
	case "BCO":
		color = "#FF69B4"
	case "M&OL":
		color = "#7FFFD4"
	default:
		color = "#ff7f33"
	}

	return color
}

func uniqueSheetName(baseName string, existingNames map[string]bool) string {
	if !existingNames[baseName] {
		existingNames[baseName] = true
		return baseName
	}

	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s (%d)", baseName, i)
		if !existingNames[newName] {
			existingNames[newName] = true
			return newName
		}
	}
}
