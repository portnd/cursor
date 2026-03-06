package helpers

import (
	"context"
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/leekchan/accounting"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
)

func ReportName(HTMLTemplate string) (string, error) {
	switch HTMLTemplate {
	case "TEMPLATE_GENARAL_TYPE1",
		"TEMPLATE_GENARAL_TYPE1_EXCEL":
		return "รายงานสินทรัพย์", nil

	case "TEMPLATE_GENARAL_TYPE2",
		"TEMPLATE_GENARAL_TYPE2_EXCEL":
		return "รายงานแผนที่สินทรัพย์", nil

	case "TEMPLATE_GENARAL_TYPE3",
		"TEMPLATE_GENARAL_TYPE3_EXCEL":
		return "รายงานสรุปสินทรัพย์", nil

	case "TEMPLATE_GENARAL_TYPE4",
		"TEMPLATE_GENARAL_TYPE4_EXCEL":
		return "รายงานการปรับแก้สินทรัพย์ประจำเดือน", nil

	case "TEMPLATE_GENARAL_TYPE6_AEK",
		"TEMPLATE_GENARAL_TYPE6_AEK_EXCEL":
		return "รายงานข้อมูลสภาพทาง", nil

	case "TEMPLATE_GENARAL_TYPE7_AEK_IRI",
		"TEMPLATE_GENARAL_TYPE7_AEK",
		"TEMPLATE_GENARAL_TYPE7_AEK_MPD",
		"TEMPLATE_GENARAL_TYPE7_AEK_RUT",
		"TEMPLATE_GENARAL_TYPE7_AEK_GN",
		"TEMPLATE_GENARAL_TYPE7_AEK_EXCEL_IRI",
		"TEMPLATE_GENARAL_TYPE7_AEK_EXCEL_MPD",
		"TEMPLATE_GENARAL_TYPE7_AEK_EXCEL_RUT",
		"TEMPLATE_GENARAL_TYPE7_AEK_EXCEL_GN":
		return "รายงานสรุปข้อมูลสภาพทาง", nil

	case "TEMPLATE_GENARAL_TYPE11_AEK",
		"TEMPLATE_GENARAL_TYPE11_AEK_EXCEL":
		return "รายงานสรุปข้อมูลความเสียหาย", nil
	case "TEMPLATE_GENARAL_TYPE10_AEK",
		"TEMPLATE_GENARAL_TYPE10_AEK_EXCEL":
		return "รายงานข้อมูลความเสียหาย", nil
	case "TEMPLATE_GENARAL_TYPE5_AEK",
		"TEMPLATE_GENARAL_TYPE5_AEK_EXCEL":
		return "รายงานสรุปรายละเอียดชนิดผิวทาง", nil

	case "TEMPLATE_GENARAL_TYPE8_HTML",
		"TEMPLATE_GENARAL_TYPE8_EXCEL":
		return "รายงานค่าสะท้อนแสงของเส้นจราจร", nil

	case "TEMPLATE_GENARAL_TYPE9_HTML",
		"TEMPLATE_GENARAL_TYPE9_EXCEL":
		return "รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร", nil

	case "TEMPLATE_GENARAL_TYPE13",
		"TEMPLATE_GENARAL_TYPE13_EXCEL":
		return "รายงานประวัติการซ่อมบำรุง", nil

	case "TEMPLATE_GENARAL_TRAFFIC_VOLUME",
		"TEMPLATE_GENARAL_TRAFFIC_VOLUME_EXCEL":
		return "รายงานปริมาณจราจร", nil

	case "TEMPLATE_GENARAL_TYPE12",
		"TEMPLATE_GENARAL_TYPE12_EXCEL":
		return "รายงานอุบัติเหตุ", nil
	case "TEMPLATE_GENARAL_ROAD",
		"TEMPLATE_GENARAL_ROAD_EXCEL":
		return "รายงานสรุปข้อมูลสายทาง", nil
	case "TEMPLATE_GENARAL_TYPE7_AEK_MAP":
		return "รายงานสรุปข้อมูลสายทาง", nil
	case "TEMPLATE_GENARAL_TYPE12_IRI_HTML",
		"TEMPLATE_GENARAL_TYPE12_IFI_HTML",
		"TEMPLATE_GENARAL_TYPE12_RUT_HTML",
		"TEMPLATE_GENARAL_TYPE12_G7_HTML",
		"TEMPLATE_GENARAL_TYPE12_IRI_EXCEL",
		"TEMPLATE_GENARAL_TYPE12_IFI_EXCEL",
		"TEMPLATE_GENARAL_TYPE12_RUT_EXCEL",
		"TEMPLATE_GENARAL_TYPE12_G7_EXCEL":
		return "รายงานการซ่อมบํารุงตามเกณฑ์ KPI ค่าดัชนีความขรุขระสากล", nil
	}

	return "", fmt.Errorf("not found the %s", HTMLTemplate)
}

func FormatFloat(input float64) string {
	ac := accounting.Accounting{Symbol: "", Precision: 3, Thousand: ","}
	if input < 0 {
		return "(" + ac.FormatMoney(-input) + ")"
	}
	if input == 0 {
		return "-"
	}
	return ac.FormatMoney(input)
}

func RequestExport(data interface{}, template string, typ string) (interface{}, error) {
	var pathResult interface{}
	var err error

	switch typ {
	case "html":
		pathResult, err = ExportHTML(data, template)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	case "pdf":
		pathResult, err = ExportPDF(data, template)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}

	return pathResult, nil
}

func ExportHTML(data interface{}, HTMLTemplate string) (interface{}, error) {
	filePath := os.Getenv("GENARAL_HTML")
	templateName := os.Getenv(HTMLTemplate)

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"minus": func(a, b int) int {
			return a - b
		},
		"formatFloat":         FormatFloat,
		"calculateZoom":       CalculateZoom,
		"formatFloatSetPoint": FormatFloatSetPoint,
	}
	templateGen, err := getCachedHtmlTemplate(templateName, func() (*template.Template, error) {
		return template.New(filepath.Base(templateName)).Funcs(funcMap).ParseFiles(templateName)
	})
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName(HTMLTemplate)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	fileName := filePath + code + "_" + reportName + ".html"
	fileWritter, err := os.Create(fileName)

	if err != nil {
		logs.Error(err)
		return nil, err
	}

	if err := templateGen.Execute(fileWritter, data); err != nil {
		logs.Error(err)
		return nil, err
	}

	if err := fileWritter.Close(); err != nil {
		logs.Error(err)
		return nil, err
	}

	return os.Getenv("STORAGE_IP") + "/" + fileName, nil
}

func ExportPDF(data interface{}, HTMLTemplate string) (interface{}, error) {
	filePath := os.Getenv("GENARAL_PDF")
	templateName := os.Getenv(HTMLTemplate)
	// fmt.Println("templateName", templateName)
	result, err := InitDataToHtml(templateName, data, filePath)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return nil, err
	}

	ctx, cancel := NewChromedpContext(context.Background(), log.Printf)
	defer cancel()

	var buf []byte

	isDelay := false

	switch HTMLTemplate {
	case "TEMPLATE_GENARAL_TYPE7_AEK", "TEMPLATE_GENARAL_TYPE6_AEK":
		isDelay = true
	}
	if HTMLTemplate == "TEMPLATE_GENARAL_TYPE10" || HTMLTemplate == "TEMPLATE_GENARAL_TYPE2" {
		isDelay = true
	}
	if HTMLTemplate == "TEMPLATE_GENARAL_TYPE9_HTML" {
		isDelay = true
	}

	if HTMLTemplate == "TEMPLATE_GENARAL_TYPE12_IRI_HTML" ||
		HTMLTemplate == "TEMPLATE_GENARAL_TYPE12_IFI_HTML" ||
		HTMLTemplate == "TEMPLATE_GENARAL_TYPE12_RUT_HTML" ||
		HTMLTemplate == "TEMPLATE_GENARAL_TYPE12_G7_HTML" {
		isDelay = true
	}

	if err := chromedp.Run(ctx, PrintToPDF(string(html), &buf, isDelay)); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName(HTMLTemplate)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	fullFilePath := filePath + code + "_" + reportName + ".pdf"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return nil, err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func PrintToPNG(html string, res *[]byte, isDelay bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}
			delay := 5
			if isDelay {
				delay = 20
			}

			defer chromedp.Run(
				ctx,
				RunWithTimeOut(&ctx, time.Duration(delay), chromedp.Tasks{
					chromedp.WaitVisible("div#success"),
				}),
			)

			wg.Wait()
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			width, height := 455, 340
			if err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false).Do(ctx); err != nil {
				return err
			}

			buf, err := page.CaptureScreenshot().Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func ExportPNG(data interface{}, HTMLTemplate string) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	templateName := os.Getenv(HTMLTemplate)

	result, err := InitDataToHtml(templateName, data, filePath)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return nil, err
	}

	ctx, cancel := NewChromedpContext(context.Background(), log.Printf)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, PrintToPNG(string(html), &buf, false)); err != nil {
		logs.Error(err)
		return nil, err
	}

	unix := strconv.Itoa(int(time.Now().Unix()))

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	fullFilePath := filePath + unix + uuid.NewString() + ".png"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return nil, err
	}

	return fullFilePath, nil
}

func AddImage(sheetName string, j interface{}, alp rune, ext, row, folderPath string, f *excelize.File, merge int) {

	newRow, _ := strconv.Atoi(row)

	// trim img path
	trimmedURL := strings.TrimPrefix(fmt.Sprint(j), os.Getenv("STORAGE_IP")+"/")
	//trimmedURL = strings.ReplaceAll(trimmedURL, ":", "_") //get rid in production

	// img size checker
	file, err := os.Open(trimmedURL)
	if err != nil {
		fmt.Println("can not open file", err)
	}
	defer file.Close()

	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("can not check image size", err)
	}

	picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(newRow+1)

	// colI, _ := f.GetColWidth(sheetName, "I")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("col I:", colI)
	// colJ, _ := f.GetColWidth(sheetName, "J")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("col J:", colJ)

	// fmt.Println("row: ", newRow)

	cellHight, err := f.GetRowHeight(sheetName, newRow)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("height: ", cellHight)
	if cellHight < 20 {
		cellHight = 20
	}
	//fmt.Println(cellHight)
	setY := -1 * int((cellHight*1.3-19)/2+19.5) // 1.3 offset = 1 cell hight

	if ext == ".jpeg" || ext == ".png" {

		if merge%2 == 0 {

			if imageConfig.Height == 19 && imageConfig.Width == 19 {
				f.AddPicture(sheetName, picPosition, trimmedURL, &excelize.GraphicOptions{OffsetX: -9, OffsetY: setY})

			} else {
				f.AddPicture(sheetName, picPosition, trimmedURL, &excelize.GraphicOptions{OffsetX: -9, OffsetY: setY, ScaleX: 19.00 / float64(imageConfig.Height), ScaleY: 19.00 / float64(imageConfig.Width)})
			}

		} else {
			cellWidth, _ := f.GetColWidth("Sheet1", fmt.Sprintf("%c", alp))

			if imageConfig.Height == 19 && imageConfig.Width == 19 {
				f.AddPicture(sheetName, picPosition, trimmedURL, &excelize.GraphicOptions{OffsetX: int((cellWidth * 3) + 0.5), OffsetY: setY})

			} else {
				f.AddPicture(sheetName, picPosition, trimmedURL, &excelize.GraphicOptions{OffsetX: int((cellWidth * 3) + 0.5), OffsetY: setY, ScaleX: 19.00 / float64(imageConfig.Height), ScaleY: 19.00 / float64(imageConfig.Width)})

			}
		}

	} else if ext == ".jpg" {

		// open file
		srcImage, err := imaging.Open(trimmedURL)
		if err != nil {
			fmt.Println("can not open file", err)

		} else {

			fileName := filepath.Base(trimmedURL)

			saveAs := folderPath + fileName

			// convert to .png
			err = imaging.Save(srcImage, saveAs)
			if err != nil {
				fmt.Println("can not save file", err)
			}

			if merge%2 == 0 {
				f.AddPicture(sheetName, picPosition, trimmedURL, &excelize.GraphicOptions{OffsetX: -9, OffsetY: setY, ScaleX: 19.00 / float64(imageConfig.Height), ScaleY: 19.00 / float64(imageConfig.Width)})
			} else {
				cellWidth, _ := f.GetColWidth(sheetName, fmt.Sprintf("%c", alp))
				f.AddPicture(sheetName, picPosition, trimmedURL, &excelize.GraphicOptions{OffsetX: int((cellWidth * 3) + 0.5), OffsetY: setY, ScaleX: 19.00 / float64(imageConfig.Height), ScaleY: 19.00 / float64(imageConfig.Width)})
			}
		}
	}
}

func AddFooter(f *excelize.File, sheetName string) {
	currentTime := time.Now().AddDate(543, 0, 0).Format("2 Jan 2006")
	month := time.Now().Format("Jan")
	monthThai := map[string]string{
		"Jan": "ม.ค.",
		"Feb": "ก.พ.",
		"Mar": "มี.ค.",
		"Apr": "เม.ย.",
		"May": "พ.ค.",
		"Jun": "มิ.ย.",
		"Jul": "ก.ค.",
		"Aug": "ส.ค.",
		"Sep": "ก.ย.",
		"Oct": "ต.ค.",
		"Nov": "พ.ย.",
		"Dec": "ธ.ค.",
	}
	currentTime = strings.Replace(currentTime, month, monthThai[month], 1)

	btm := 0.25
	f.SetPageMargins(sheetName, &excelize.PageLayoutMarginsOptions{
		Footer: &btm,
	})
}

func ExportExcelType1(datas []models.DataReportMap) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE1_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for dataIndex, data := range datas {
		sheetName := "Sheet1 (" + fmt.Sprint(dataIndex+1) + ")"
		// header
		f.SetCellValue(sheetName, "C2", "รายงานสินทรัพย์ : "+datas[0].AssetName)
		f.SetCellValue(sheetName, "C3", "หมายเลขทางหลวง : "+datas[0].RoadGroupNumber+" ตอนควบคุม : "+datas[0].RoadSectionNumber)
		f.SetCellValue(sheetName, "C4", "ชื่อสายทาง : "+datas[0].RoadSectionNameOriginTh+" - "+datas[0].RoadSectionNameDestinationTh)
		f.SetCellValue(sheetName, "C5", "กม.เริ่มต้น "+datas[0].KmStart+" กม.สิ้นสุด "+datas[0].KmEnd+" ระยะทาง "+datas[0].StrRoadLength+" กม.")

		rowPin := 7
		rowData := 8

		// road name
		f.SetCellValue(sheetName, "B"+fmt.Sprint(rowPin), "ชื่อสายทาง: "+data.RoadName)

		headStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: "TH SarabunPSK",
				Bold:   true,
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#E7E6E6"},
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

		// table
		if !data.IsNull {

			folderPath := "storages/asset/temp" + uuid.New().String() + "/"
			// create temp folder
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				logs.Error(err)
				return nil, err
			}

			// remove temp folder
			defer os.RemoveAll(folderPath)

			// header
			alpHead := 'B'
			position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
			f.SetCellValue(sheetName, position, "ลำดับ")
			f.SetCellStyle(sheetName, position, position, headStyle)

			alpHead += 1

			for iIndex, i := range data.Column {

				position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)

				//เงื่อนไข merge หัวตาราง
				switch len(data.Column) {
				case 1:
					{
						alpHead += 8
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)
					}
				case 2:
					{
						if iIndex <= 0 {
							alpHead += 3
						} else {
							alpHead += 4
						}
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)
					}
				case 3:
					{
						alpHead += 2
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)

					}
				case 4:
					{
						if iIndex <= 2 {
							alpHead += 1

						} else {
							alpHead += 2
						}

						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)

					}

				case 5:
					{
						if iIndex <= 3 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)
						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}

					}

				case 6:
					{
						if iIndex == 2 || iIndex == 3 || iIndex == 4 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				case 7:
					{
						if iIndex == 6 {
							alpHead += 2
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}
				case 8:
					{
						if iIndex == 7 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				default:
					{

						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position, headStyle)
					}
				}

				alpHead += 1
			}

			// data
			// style
			style, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:   12,
					Family: "TH SarabunPSK",
				},
				Alignment: &excelize.Alignment{
					WrapText:   true,
					Horizontal: "center",
					Vertical:   "center",
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

			for iIndex, i := range data.Row {
				// road name
				alp := 'B'
				row := fmt.Sprint(rowData + 1 + iIndex)
				position := fmt.Sprintf("%c", alp) + row

				f.SetCellValue(sheetName, position, iIndex+1)
				f.SetCellStyle(sheetName, position, position, style)

				alp += 1

				//เงื่อนไข merge แถว
				for jIndex, j := range i {

					position := fmt.Sprintf("%c", alp) + row
					merge := 0

					switch len(i) {
					case 1:
						{
							alp += 8
							merge = 9
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}
					case 2:
						{
							if jIndex <= 0 {
								alp += 3
								merge = 4
							} else {
								alp += 4
								merge = 5
							}
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}

					case 3:
						{
							alp += 2
							merge = 3
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}
					case 4:
						{
							if jIndex <= 2 {
								alp += 1
								merge = 2
							} else {
								alp += 2
								merge = 3
							}
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}

					case 5:
						{
							if jIndex <= 3 {
								alp += 1
								merge = 2
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {

								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					case 6:
						{
							if jIndex == 2 || jIndex == 3 || jIndex == 4 {
								alp += 1
								merge = 2
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {
								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					case 7:
						{
							if jIndex == 6 {
								alp += 2
								merge = 3
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {
								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					case 8:
						{
							if jIndex == 7 {
								alp += 1
								merge = 2
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {
								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					default:
						{
							merge = 1
							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position, style)
						}
					}

					alp += 1
				}
			}

		} else {

			// header
			alpHead := 'B'
			position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)

			switch len(data.Column) {
			case 2, 3:
				{
					alpHead += 3
					position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
					f.MergeCell(sheetName, position, position2)
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position2, headStyle)
				}

			case 5:
				{
					alpHead += 2
					position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
					f.MergeCell(sheetName, position, position2)
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position2, headStyle)
				}

			case 4, 6, 7, 8:
				{
					alpHead += 1
					position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
					f.MergeCell(sheetName, position, position2)
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position2, headStyle)
				}

			default:
				{
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position, headStyle)
				}

			}

			alpHead += 1

			for iIndex, i := range data.Column {

				position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)

				switch len(data.Column) {
				case 2:
					{
						alpHead += 2
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)
					}

				case 3, 4:
					{
						alpHead += 1
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)

					}

				case 5:
					{
						if iIndex <= 1 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				case 6:
					{
						if iIndex == 2 || iIndex == 3 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				case 7:
					{
						if iIndex <= 0 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				default:
					{

						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position, headStyle)
					}
				}

				alpHead += 1
			}

			// data
			alp := 'B'
			lastAlp := 'B' + len(data.Column)
			position1 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowData+1)
			position2 := fmt.Sprintf("%c", lastAlp) + fmt.Sprint(rowData+1)

			f.MergeCell(sheetName, position1, position2)
			f.SetCellValue(sheetName, position1, "ไม่พบข้อมูล")

			// style
			style, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:   12,
					Family: "TH SarabunPSK",
				},
				Alignment: &excelize.Alignment{
					WrapText:   true,
					Horizontal: "center",
					Vertical:   "center",
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

			f.SetCellStyle(sheetName, position1, position2, style)

		}
		AddFooter(f, sheetName)

	}

	// ตั้งชื่อ sheet
	existingSheetNames := make(map[string]bool)

	for indexPageName, i := range datas {
		var name string
		if len([]rune(i.RoadName)) >= 25 {
			name = string([]rune(i.RoadName)[:25])
		} else {
			name = string(i.RoadName)
		}

		uniqueName := uniqueSheetName(name, existingSheetNames)
		f.SetSheetName("Sheet1 ("+fmt.Sprint(indexPageName+1)+")", uniqueName)
	}

	// ลบ sheet
	for i := len(datas) + 1; i <= 100; i++ {
		f.DeleteSheet("Sheet1 (" + fmt.Sprint(i) + ")")
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE1_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

}

func ExportExcelType2(datas []models.DataReportMap) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	icon := []models.Icon{}
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE2_EXCEL"))
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	imgPath := make([]string, len(datas))

	for imgIndex, imgItem := range datas {
		wg.Add(1)
		go func(imgIndex int, imgItem models.DataReportMap) {
			defer wg.Done()
			// image
			path, _ := ExportPNG(imgItem, "TEMPLATE_GENARAL_TYPE2_MAP")
			imgPath[imgIndex] = fmt.Sprint(path)
		}(imgIndex, imgItem)
	}
	wg.Wait()

	for dataIndex, data := range datas {
		sheetName := "Sheet1 (" + fmt.Sprint(dataIndex+1) + ")"
		// header
		f.SetCellValue(sheetName, "C2", "รายงานแผนที่สินทรัพย์ : "+datas[0].AssetName)
		f.SetCellValue(sheetName, "C3", "หมายเลขทางหลวง : "+datas[0].RoadGroupNumber+" ตอนควบคุม : "+datas[0].RoadSectionNumber)
		f.SetCellValue(sheetName, "C4", "ชื่อสายทาง : "+datas[0].RoadSectionNameOriginTh+" - "+datas[0].RoadSectionNameDestinationTh)
		f.SetCellValue(sheetName, "C5", "กม.เริ่มต้น "+datas[0].KmStart+" กม.สิ้นสุด "+datas[0].KmEnd+" ระยะทาง "+datas[0].StrRoadLength+" กม.")

		rowImg := 7
		rowPin := 20
		rowData := 21

		// img size checker
		file, err := os.Open(fmt.Sprint(imgPath[dataIndex]))
		if err != nil {
			fmt.Println("can not open file", err)
		}
		defer file.Close()

		imageConfig, _, err := image.DecodeConfig(file)
		if err != nil {
			fmt.Println("can not check image size", err)
		}

		f.MergeCell(sheetName, "D"+fmt.Sprint(rowImg), "I"+fmt.Sprint(rowImg+11))
		f.AddPicture(sheetName, "D"+fmt.Sprint(rowImg), (fmt.Sprint(imgPath[dataIndex])), &excelize.GraphicOptions{ScaleX: 420 / float64(imageConfig.Width), ScaleY: 325 / float64(imageConfig.Height)})

		defer os.Remove(fmt.Sprint(imgPath[dataIndex]))

		//pin sign
		if data.PointGeom != nil {

			f.AddPicture(sheetName, "K"+fmt.Sprint(rowPin), "templates/genaral/pin.png", &excelize.GraphicOptions{OffsetY: 5, ScaleX: 15.00 / 512.00, ScaleY: 13.0 / 512.00})
			f.SetCellValue(sheetName, "K"+fmt.Sprint(rowPin), "    สินทรัพย์")

		}

		// road name
		f.MergeCell(sheetName, "B"+fmt.Sprint(rowPin), "J"+fmt.Sprint(rowPin))
		f.SetCellValue(sheetName, "B"+fmt.Sprint(rowPin), "ชื่อสายทาง: "+data.RoadName)

		headStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: "TH SarabunPSK",
				Bold:   true,
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#E7E6E6"},
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

		// table
		if !data.IsNull {

			folderPath := "storages/asset/temp" + uuid.New().String() + "/"
			// create temp folder
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				logs.Error(err)
				return nil, err
			}

			// remove temp folder
			defer os.RemoveAll(folderPath)

			// header
			alpHead := 'B'
			position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
			f.SetCellValue(sheetName, position, "ลำดับ")
			f.SetCellStyle(sheetName, position, position, headStyle)

			alpHead += 1

			for iIndex, i := range data.Column {

				position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)

				//เงื่อนไข merge หัวตาราง
				switch len(data.Column) {
				case 1:
					{
						alpHead += 8
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)
					}
				case 2:
					{
						if iIndex <= 0 {
							alpHead += 3
						} else {
							alpHead += 4
						}
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)
					}
				case 3:
					{
						alpHead += 2
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)

					}
				case 4:
					{
						if iIndex <= 2 {
							alpHead += 1

						} else {
							alpHead += 2
						}

						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)

					}

				case 5:
					{
						if iIndex <= 3 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)
						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}

					}

				case 6:
					{
						if iIndex == 2 || iIndex == 3 || iIndex == 4 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				case 7:
					{
						if iIndex == 6 {
							alpHead += 2
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}
				case 8:
					{
						if iIndex == 7 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				default:
					{

						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position, headStyle)
					}
				}

				alpHead += 1
			}

			// data
			// style
			style, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:   12,
					Family: "TH SarabunPSK",
				},
				Alignment: &excelize.Alignment{
					WrapText:   true,
					Horizontal: "center",
					Vertical:   "center",
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

			for iIndex, i := range data.Row {
				// road name
				alp := 'B'
				row := fmt.Sprint(rowData + 1 + iIndex)
				position := fmt.Sprintf("%c", alp) + row

				f.SetCellValue(sheetName, position, iIndex+1)
				f.SetCellStyle(sheetName, position, position, style)

				alp += 1

				//เงื่อนไข merge แถว
				for jIndex, j := range i {

					position := fmt.Sprintf("%c", alp) + row
					merge := 0

					switch len(i) {
					case 1:
						{
							alp += 8
							merge = 9
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext != ".jpeg" && ext != ".png" && ext != ".jpg" {
								f.SetCellValue(sheetName, position, j)
							} else {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								icon = append(icon, models.Icon{
									SheetName:  sheetName,
									J:          j,
									Alp:        alp,
									Ext:        ext,
									Row:        row,
									FolderPath: folderPath,
									F:          f,
									Merge:      merge,
								})
								//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							}

							f.SetCellStyle(sheetName, position, position2, style)
						}
					case 2:
						{
							if jIndex <= 0 {
								alp += 3
								merge = 4
							} else {
								alp += 4
								merge = 5
							}
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								icon = append(icon, models.Icon{
									SheetName:  sheetName,
									J:          j,
									Alp:        alp,
									Ext:        ext,
									Row:        row,
									FolderPath: folderPath,
									F:          f,
									Merge:      merge,
								})
								//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}

					case 3:
						{
							alp += 2
							merge = 3
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								icon = append(icon, models.Icon{
									SheetName:  sheetName,
									J:          j,
									Alp:        alp,
									Ext:        ext,
									Row:        row,
									FolderPath: folderPath,
									F:          f,
									Merge:      merge,
								})
								//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}
					case 4:
						{
							if jIndex <= 2 {
								alp += 1
								merge = 2
							} else {
								alp += 2
								merge = 3
							}
							position2 := fmt.Sprintf("%c", alp) + row
							f.MergeCell(sheetName, position, position2)

							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								icon = append(icon, models.Icon{
									SheetName:  sheetName,
									J:          j,
									Alp:        alp,
									Ext:        ext,
									Row:        row,
									FolderPath: folderPath,
									F:          f,
									Merge:      merge,
								})
								//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position2, style)
						}

					case 5:
						{
							if jIndex <= 3 {
								alp += 1
								merge = 2
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {

								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					case 6:
						{
							if jIndex == 2 || jIndex == 3 || jIndex == 4 {
								alp += 1
								merge = 2
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {
								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					case 7:
						{
							if jIndex == 6 {
								alp += 2
								merge = 3
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {
								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					case 8:
						{
							if jIndex == 7 {
								alp += 1
								merge = 2
								position2 := fmt.Sprintf("%c", alp) + row
								f.MergeCell(sheetName, position, position2)

								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position2, style)

							} else {
								merge = 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
									// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
									icon = append(icon, models.Icon{
										SheetName:  sheetName,
										J:          j,
										Alp:        alp,
										Ext:        ext,
										Row:        row,
										FolderPath: folderPath,
										F:          f,
										Merge:      merge,
									})
									//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

								} else {
									f.SetCellValue(sheetName, position, j)
								}

								f.SetCellStyle(sheetName, position, position, style)
							}
						}

					default:
						{
							merge = 1
							ext := strings.ToLower(filepath.Ext(fmt.Sprint(j)))
							if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
								// picPosition := fmt.Sprintf("%c", int(alp)-int((merge-1)/2)) + fmt.Sprint(row)
								// AddImage(sheetName,j, alp, ext, picPosition, folderPath, f, merge)
								icon = append(icon, models.Icon{
									SheetName:  sheetName,
									J:          j,
									Alp:        alp,
									Ext:        ext,
									Row:        row,
									FolderPath: folderPath,
									F:          f,
									Merge:      merge,
								})
								//AddImage(sheetName, j, alp, ext, row, folderPath, f, merge)

							} else {
								f.SetCellValue(sheetName, position, j)
							}

							f.SetCellStyle(sheetName, position, position, style)
						}
					}

					alp += 1
				}
			}

		} else {

			// header
			alpHead := 'B'
			position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)

			switch len(data.Column) {
			case 2, 3:
				{
					alpHead += 3
					position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
					f.MergeCell(sheetName, position, position2)
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position2, headStyle)
				}

			case 5:
				{
					alpHead += 2
					position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
					f.MergeCell(sheetName, position, position2)
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position2, headStyle)
				}

			case 4, 6, 7, 8:
				{
					alpHead += 1
					position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
					f.MergeCell(sheetName, position, position2)
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position2, headStyle)
				}

			default:
				{
					f.SetCellValue(sheetName, position, "ลำดับ")
					f.SetCellStyle(sheetName, position, position, headStyle)
				}

			}

			alpHead += 1

			for iIndex, i := range data.Column {

				position := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)

				switch len(data.Column) {
				case 2:
					{
						alpHead += 2
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)
					}

				case 3, 4:
					{
						alpHead += 1
						position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
						f.MergeCell(sheetName, position, position2)
						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position2, headStyle)

					}

				case 5:
					{
						if iIndex <= 1 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				case 6:
					{
						if iIndex == 2 || iIndex == 3 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				case 7:
					{
						if iIndex <= 0 {
							alpHead += 1
							position2 := fmt.Sprintf("%c", alpHead) + fmt.Sprint(rowData)
							f.MergeCell(sheetName, position, position2)
							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position2, headStyle)

						} else {

							f.SetCellValue(sheetName, position, i)
							f.SetCellStyle(sheetName, position, position, headStyle)
						}
					}

				default:
					{

						f.SetCellValue(sheetName, position, i)
						f.SetCellStyle(sheetName, position, position, headStyle)
					}
				}

				alpHead += 1
			}

			// data
			alp := 'B'
			lastAlp := 'B' + len(data.Column)
			position1 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowData+1)
			position2 := fmt.Sprintf("%c", lastAlp) + fmt.Sprint(rowData+1)

			f.MergeCell(sheetName, position1, position2)
			f.SetCellValue(sheetName, position1, "ไม่พบข้อมูล")

			// style
			style, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:   12,
					Family: "TH SarabunPSK",
				},
				Alignment: &excelize.Alignment{
					WrapText:   true,
					Horizontal: "center",
					Vertical:   "center",
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

			f.SetCellStyle(sheetName, position1, position2, style)

		}
		AddFooter(f, sheetName)

	}

	// ลบ sheet
	for i := len(datas) + 1; i <= 100; i++ {
		f.DeleteSheet("Sheet1 (" + fmt.Sprint(i) + ")")
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE2_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}

	if len(icon) > 0 {
		newf, err := excelize.OpenFile(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := newf.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		for _, i := range icon {
			AddImage(i.SheetName, i.J, i.Alp, i.Ext, i.Row, i.FolderPath, newf, i.Merge)
		}

		// ตั้งชื่อ sheet
		existingSheetNames := make(map[string]bool)

		for indexPageName, i := range datas {
			var name string
			if len([]rune(i.RoadName)) >= 25 {
				name = string([]rune(i.RoadName)[:25])
			} else {
				name = string(i.RoadName)
			}

			uniqueName := uniqueSheetName(name, existingSheetNames)
			newf.SetSheetName("Sheet1 ("+fmt.Sprint(indexPageName+1)+")", uniqueName)
		}

		newf.Save()
	} else {
		newf, err := excelize.OpenFile(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := newf.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		// ตั้งชื่อ sheet
		existingSheetNames := make(map[string]bool)

		for indexPageName, i := range datas {
			var name string
			if len([]rune(i.RoadName)) >= 25 {
				name = string([]rune(i.RoadName)[:25])
			} else {
				name = string(i.RoadName)
			}

			uniqueName := uniqueSheetName(name, existingSheetNames)
			newf.SetSheetName("Sheet1 ("+fmt.Sprint(indexPageName+1)+")", uniqueName)
		}
		newf.Save()
	}
	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

}

func ExportExcelType3(data *models.DataReportSummaryAsset) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE3_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//start
	f.SetCellValue("Sheet1", "C3", fmt.Sprintf(`หมายเลขทางหลวง : %s ตอนควบคุม : %s`, data.RoadGroupNumber, data.RoadSectionNumber))
	f.SetCellValue("Sheet1", "C4", fmt.Sprintf(`ชื่อสายทาง : %s`, data.RoadMainName))
	f.SetCellValue("Sheet1", "C5", fmt.Sprintf(`กม.เริ่มต้น %s กม.สิ้นสุด %s ระยะทาง %s กม.`, data.KmStart, data.KmEnd, data.StrRoadLength))

	if !data.IsNull {
		// style
		rowTitelStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: "TH SarabunPSK",
				Bold:   true,
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{{
				Type:  "left",
				Color: "000000",
				Style: 0,
			}, {
				Type:  "top",
				Color: "000000",
				Style: 0,
			}, {
				Type:  "right",
				Color: "000000",
				Style: 0,
			}, {
				Type:  "bottom",
				Color: "000000",
				Style: 0,
			}},
		})

		rowStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: "TH SarabunPSK",
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
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

		LeftRowStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: "TH SarabunPSK",
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "left",
				Vertical:   "center",
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

		headStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: "TH SarabunPSK",
				Bold:   true,
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"#E7E6E6"},
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

		rowIndex := 7
		for index, v := range data.Table {

			positionTitle1 := fmt.Sprintf(`B%d`, rowIndex)
			positionTitle2 := fmt.Sprintf(`C%d`, rowIndex)
			f.MergeCell("Sheet1", positionTitle1, positionTitle2)
			f.SetCellValue("Sheet1", positionTitle1, fmt.Sprintf(`%d. %s`, (index+1), v.Title))
			f.SetCellStyle("Sheet1", positionTitle1, positionTitle2, rowTitelStyle)

			rowIndex++

			for _, i := range v.Table {

				if i.Topic == "1. ข้อมูลเสาไฟฟ้าส่องสว่าง" {
					headIndex := 'C'
					position := fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
					f.SetCellValue("Sheet1", position, i.Topic)

					for _, j := range i.Header {
						position = fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
						f.SetCellValue("Sheet1", position, j)
						f.SetCellStyle("Sheet1", position, position, headStyle)
						headIndex++
					}

					headIndex = 'C'

					for _, k := range i.Row {
						for _, l := range k {
							position = fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
							f.SetCellValue("Sheet1", position, l)

							switch headIndex {
							case 'C', 'D':
								{
									f.SetCellStyle("Sheet1", position, position, LeftRowStyle)
								}
							default:
								{
									f.SetCellStyle("Sheet1", position, position, rowStyle)
								}
							}

							headIndex++

						}
						headIndex = 'C'
						rowIndex++
					}

					rowIndex++

				} else {

					headIndex := 'C'
					position := fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
					position2 := fmt.Sprintf("%c", headIndex+1) + fmt.Sprint(rowIndex)
					f.MergeCell("Sheet1", position, position2)
					f.SetCellValue("Sheet1", position, i.Topic)

					//headIndex++
					rowIndex++

					for jIndex, j := range i.Header {
						switch jIndex {
						case 2:
							{
								position = fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
								position2 = fmt.Sprintf("%c", headIndex+1) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)
								f.SetCellValue("Sheet1", position, j)
								f.SetCellStyle("Sheet1", position, position2, headStyle)

								headIndex += 2
							}
						default:
							{
								position = fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
								f.SetCellValue("Sheet1", position, j)
								f.SetCellStyle("Sheet1", position, position, headStyle)

								headIndex++
							}
						}
					}

					rowIndex++
					headIndex = 'C'

					for _, k := range i.Row {
						for lIndex, l := range k {
							switch lIndex {
							case 2:
								{
									position = fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
									position2 := fmt.Sprintf("%c", headIndex+1) + fmt.Sprint(rowIndex)
									f.MergeCell("Sheet1", position, position2)
									f.SetCellValue("Sheet1", position, l)

									switch headIndex {
									case 'C':
										{
											f.SetCellStyle("Sheet1", position, position2, LeftRowStyle)
										}

									default:
										{
											f.SetCellStyle("Sheet1", position, position2, rowStyle)
										}
									}

									headIndex += 2
								}
							default:
								{
									position = fmt.Sprintf("%c", headIndex) + fmt.Sprint(rowIndex)
									f.SetCellValue("Sheet1", position, l)

									switch headIndex {
									case 'C':
										{
											f.SetCellStyle("Sheet1", position, position, LeftRowStyle)
										}

									default:
										{
											f.SetCellStyle("Sheet1", position, position, rowStyle)
										}
									}

									headIndex++
								}
							}
						}
						headIndex = 'C'
						rowIndex++
					}

					rowIndex++

				}
			}
		}
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	AddFooter(f, "Sheet1")

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE3_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelType4(data *models.DataReportAssetAdjustment) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE4_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//start
	f.SetCellValue("Sheet1", "C2", "รายงานการปรับแก้สินทรัพย์ประจำเดือน "+data.Month+" พ.ศ. "+data.Year)
	f.SetCellValue("Sheet1", "C3", "ชื่อสายทาง : "+data.RoadGroupName)
	f.SetCellValue("Sheet1", "C4", "ช่วง : "+data.RoadName+" รหัส : "+data.RoadCode)
	f.SetCellValue("Sheet1", "C5", "กม.เริ่มต้น "+data.KmStart+" กม.สิ้นสุด "+data.KmEnd+" ระยะทาง "+data.StrRoadLength+" กม.")

	rowIndex := 8

	// style
	topicStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "TH SarabunPSK",
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "left",
			Vertical:   "center",
		},
	})

	rowStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "TH SarabunPSK",
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
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

	headStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "TH SarabunPSK",
			Bold:   true,
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#E7E6E6"},
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

	if !data.IsNull {

		aStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   10,
				Family: "TH SarabunPSK",
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"#c5e0b3"},
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

		mStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   10,
				Family: "TH SarabunPSK",
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"#ffe599"},
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

		dStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   10,
				Family: "TH SarabunPSK",
			},
			Alignment: &excelize.Alignment{
				WrapText:   true,
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"#f39b84"},
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

		folderPath := "storages/adjustment/temp" + uuid.New().String() + "/"
		// create temp folder
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			logs.Error(err)
			return nil, err
		}

		// remove temp folder
		defer os.RemoveAll(folderPath)

		for _, i := range data.Table {

			f.MergeCell("Sheet1", "B"+fmt.Sprint(rowIndex), "P"+fmt.Sprint(rowIndex))
			f.SetCellValue("Sheet1", "B"+fmt.Sprint(rowIndex), i.Topic)
			f.SetCellStyle("Sheet1", "B"+fmt.Sprint(rowIndex), "P"+fmt.Sprint(rowIndex), topicStyle)
			rowIndex++

			rowNo := len(i.Header) - 3
			merge := 12/(len(i.Header)-3) - 1

			switch rowNo {
			case 5:
				{
					alp := 'B'
					for jIndex, j := range i.Header {

						position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

						if len(i.Header)-jIndex <= 3 {
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position, headStyle)

						} else if len(i.Header)-jIndex <= 5 {
							alp += 3
							position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
							f.MergeCell("Sheet1", position, position2)
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position2, headStyle)

						} else if len(i.Header)-jIndex <= 6 {
							alp += 1
							position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
							f.MergeCell("Sheet1", position, position2)
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position2, headStyle)

						} else {

							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position, headStyle)
						}

						alp++
					}

					rowIndex++

					for _, k := range i.Row {

						alp := 'B'
						for lIndex, l := range k {

							position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

							if len(k)-lIndex <= 3 {

								if len(k)-lIndex == 3 {
									switch l {
									case "A":
										{
											f.SetCellValue("Sheet1", position, "เพิ่ม")
											f.SetCellStyle("Sheet1", position, position, aStyle)
										}
									case "M":
										{
											f.SetCellValue("Sheet1", position, "แก้ไข")
											f.SetCellStyle("Sheet1", position, position, mStyle)
										}
									case "D":
										{
											f.SetCellValue("Sheet1", position, "ลบ")
											f.SetCellStyle("Sheet1", position, position, dStyle)
										}
									}

								} else {
									f.SetCellValue("Sheet1", position, l)
									f.SetCellStyle("Sheet1", position, position, rowStyle)
								}

							} else if len(k)-lIndex <= 5 {
								alp += 3
								position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)

								dis := 4
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									//picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									//AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position2, rowStyle)

							} else if len(k)-lIndex <= 6 {
								alp += 1
								position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)

								dis := 2
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									// AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position2, rowStyle)

							} else {
								dis := 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									// AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position, rowStyle)

							}

							alp++
						}

						rowIndex++
					}

					rowIndex++
				}
			case 7:
				{
					alp := 'B'
					for jIndex, j := range i.Header {

						position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

						if len(i.Header)-jIndex <= 7 {
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position, headStyle)

						} else if len(i.Header)-jIndex <= 9 {
							alp += 2
							position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
							f.MergeCell("Sheet1", position, position2)
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position2, headStyle)

						} else {
							alp += 1
							position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
							f.MergeCell("Sheet1", position, position2)
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position2, headStyle)
						}

						alp++
					}

					rowIndex++

					for _, k := range i.Row {

						alp := 'B'
						for lIndex, l := range k {

							position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

							if len(k)-lIndex <= 7 {

								if len(k)-lIndex == 3 {
									switch l {
									case "A":
										{
											f.SetCellValue("Sheet1", position, "เพิ่ม")
											f.SetCellStyle("Sheet1", position, position, aStyle)
										}
									case "M":
										{
											f.SetCellValue("Sheet1", position, "แก้ไข")
											f.SetCellStyle("Sheet1", position, position, mStyle)
										}
									case "D":
										{
											f.SetCellValue("Sheet1", position, "ลบ")
											f.SetCellStyle("Sheet1", position, position, dStyle)
										}
									}

								} else {
									f.SetCellValue("Sheet1", position, l)
									f.SetCellStyle("Sheet1", position, position, rowStyle)
								}

							} else if len(k)-lIndex <= 9 {
								alp += 2
								position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)

								dis := 3
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									// AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position2, rowStyle)

							} else {
								alp += 1
								position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)

								dis := 2
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									// AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position2, rowStyle)
							}

							alp++
						}

						rowIndex++
					}

					rowIndex++
				}
			case 8:
				{
					alp := 'B'
					for jIndex, j := range i.Header {

						position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

						if len(i.Header)-jIndex == 4 || len(i.Header)-jIndex == 7 || len(i.Header)-jIndex == 8 || len(i.Header)-jIndex == 10 {
							alp += 1
							position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
							f.MergeCell("Sheet1", position, position2)
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position2, headStyle)

						} else {
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position, headStyle)

						}

						alp++
					}

					rowIndex++

					for _, k := range i.Row {

						alp := 'B'
						for lIndex, l := range k {

							position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

							if len(k)-lIndex == 3 {
								switch l {
								case "A":
									{
										f.SetCellValue("Sheet1", position, "เพิ่ม")
										f.SetCellStyle("Sheet1", position, position, aStyle)
									}
								case "M":
									{
										f.SetCellValue("Sheet1", position, "แก้ไข")
										f.SetCellStyle("Sheet1", position, position, mStyle)
									}
								case "D":
									{
										f.SetCellValue("Sheet1", position, "ลบ")
										f.SetCellStyle("Sheet1", position, position, dStyle)
									}
								}

							} else if len(k)-lIndex == 4 || len(k)-lIndex == 7 || len(k)-lIndex == 8 || len(k)-lIndex == 10 {
								alp += 1
								position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)

								dis := 2
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									// AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position2, rowStyle)

							} else {
								f.SetCellValue("Sheet1", position, l)
								f.SetCellStyle("Sheet1", position, position, rowStyle)
							}

							alp++
						}

						rowIndex++
					}

					rowIndex++
				}
			default:
				{
					alp := 'B'
					for jIndex, j := range i.Header {

						position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

						if len(i.Header)-jIndex <= 3 {
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position, headStyle)

						} else {
							alp += rune(merge)
							position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
							f.MergeCell("Sheet1", position, position2)
							f.SetCellValue("Sheet1", position, j)
							f.SetCellStyle("Sheet1", position, position2, headStyle)

						}

						alp++
					}

					rowIndex++

					for _, k := range i.Row {

						alp = 'B'
						for lIndex, l := range k {

							position := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)

							if len(k)-lIndex <= 3 {

								if len(k)-lIndex == 3 {
									switch l {
									case "A":
										{
											f.SetCellValue("Sheet1", position, "เพิ่ม")
											f.SetCellStyle("Sheet1", position, position, aStyle)
										}
									case "M":
										{
											f.SetCellValue("Sheet1", position, "แก้ไข")
											f.SetCellStyle("Sheet1", position, position, mStyle)
										}
									case "D":
										{
											f.SetCellValue("Sheet1", position, "ลบ")
											f.SetCellStyle("Sheet1", position, position, dStyle)
										}
									}

								} else {
									f.SetCellValue("Sheet1", position, l)
									f.SetCellStyle("Sheet1", position, position, rowStyle)
								}

							} else {
								alp += rune(merge)
								position2 := fmt.Sprintf("%c", alp) + fmt.Sprint(rowIndex)
								f.MergeCell("Sheet1", position, position2)

								dis := merge + 1
								ext := strings.ToLower(filepath.Ext(fmt.Sprint(l)))
								if ext == ".jpeg" || ext == ".png" || ext == ".jpg" {
									// picPosition := fmt.Sprintf("%c", int(alp)-int((dis-1)/2)) + fmt.Sprint(rowIndex)
									// AddImage("Sheet1",l, alp, ext, picPosition, folderPath, f, dis)
									AddImage("Sheet1", l, alp, ext, fmt.Sprint(rowIndex), folderPath, f, dis)

								} else {
									f.SetCellValue("Sheet1", position, l)
								}

								f.SetCellStyle("Sheet1", position, position2, rowStyle)

							}

							alp++
						}

						rowIndex++
					}

					rowIndex++
				}
			}
		}

	} /* else {
		f.MergeCell("Sheet1", "B9", "C9")
		f.SetCellValue("Sheet1", "B9", "กม. ที่ตั้ง")
		f.SetCellStyle("Sheet1", "B9", "C9", headStyle)

		f.MergeCell("Sheet1", "B10", "C10")
		f.SetCellValue("Sheet1", "B10", "")
		f.SetCellStyle("Sheet1", "B10", "C10", rowStyle)

		f.SetCellValue("Sheet1", "D9", "หมายเลขเสา")
		f.SetCellStyle("Sheet1", "D9", "D9", headStyle)

		f.SetCellValue("Sheet1", "D10", "")
		f.SetCellStyle("Sheet1", "D10", "D10", rowStyle)

		f.MergeCell("Sheet1", "E9", "G9")
		f.SetCellValue("Sheet1", "E9", "ประเภทไฟฟ้าส่งสว่าง")
		f.SetCellStyle("Sheet1", "E9", "G9", headStyle)

		f.MergeCell("Sheet1", "E10", "G10")
		f.SetCellValue("Sheet1", "E10", "")
		f.SetCellStyle("Sheet1", "E10", "G10", rowStyle)

		f.MergeCell("Sheet1", "H9", "J9")
		f.SetCellValue("Sheet1", "H9", "สัญลักษณ์")
		f.SetCellStyle("Sheet1", "H9", "J9", headStyle)

		f.MergeCell("Sheet1", "H10", "J10")
		f.SetCellValue("Sheet1", "H10", "")
		f.SetCellStyle("Sheet1", "H10", "J10", rowStyle)

		f.MergeCell("Sheet1", "K9", "M9")
		f.SetCellValue("Sheet1", "K9", "ชนิดหลอดไฟ")
		f.SetCellStyle("Sheet1", "K9", "M9", headStyle)

		f.MergeCell("Sheet1", "K10", "M10")
		f.SetCellValue("Sheet1", "K10", "")
		f.SetCellStyle("Sheet1", "K10", "M10", rowStyle)

		f.SetCellValue("Sheet1", "N9", "สถานะ")
		f.SetCellStyle("Sheet1", "N9", "N9", headStyle)

		f.SetCellValue("Sheet1", "N10", "")
		f.SetCellStyle("Sheet1", "N10", "N10", rowStyle)

		f.SetCellValue("Sheet1", "O9", "วันที่")
		f.SetCellStyle("Sheet1", "O9", "O9", headStyle)

		f.SetCellValue("Sheet1", "O10", "")
		f.SetCellStyle("Sheet1", "O10", "O10", rowStyle)

		f.SetCellValue("Sheet1", "P9", "ผู้ใช้งาน")
		f.SetCellStyle("Sheet1", "P9", "P9", headStyle)

		f.SetCellValue("Sheet1", "P10", "")
		f.SetCellStyle("Sheet1", "P10", "P10", rowStyle)

	}*/

	AddFooter(f, "Sheet1")

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE4_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelType6(data models.DataReportSummaryCondition, factor string) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE7_AEK_EXCEL_" + factor
	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if data.IsNull {
		f.DeleteSheet("Sheet2")
		f.DeleteSheet("Sheet3")
		f.SetCellValue("Sheet1", "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
		f.SetCellValue("Sheet1", "C4", "ชื่อสายทาง : "+data.RoadSectionName)
		f.SetCellValue("Sheet1", "C5", fmt.Sprintf("กม.เริ่มต้น %s กม.สิ้นสุด %s ระยะทาง %v", data.KmStart, data.KmEnd, data.RoadLength))

		// year
		f.SetCellValue("Sheet1", "C14", "ข้อมูลสภาพทาง "+factor+" ปี "+fmt.Sprint(data.Year))
		f.SetCellValue("Sheet1", "C30", "ข้อมูลสภาพทาง "+factor+" ปี "+fmt.Sprint(data.Year))

		AddFooter(f, "Sheet1")

	} else {
		// excel logic
		switch len(data.Lane) {
		case 1:
			f.DeleteSheet("Sheet2")
			f.DeleteSheet("Sheet3")
		case 2:
			f.DeleteSheet("Sheet1")
			f.DeleteSheet("Sheet3")
		case 3:
			f.DeleteSheet("Sheet1")
			f.DeleteSheet("Sheet2")
		}

		sheetNo := fmt.Sprint(len(data.Lane))

		f.SetCellValue("Sheet1", "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
		f.SetCellValue("Sheet1", "C4", "ชื่อสายทาง : "+data.RoadSectionName)
		f.SetCellValue("Sheet1", "C5", fmt.Sprintf("กม.เริ่มต้น %s กม.สิ้นสุด %s ระยะทาง %v", data.KmStart, data.KmEnd, data.RoadLength))

		// year
		f.SetCellValue("Sheet"+sheetNo, "C14", "ข้อมูลสภาพทาง "+factor+" ปี "+fmt.Sprint(data.Year))
		f.SetCellValue("Sheet"+sheetNo, "C30", "ข้อมูลสภาพทาง "+factor+" ปี "+fmt.Sprint(data.Year))

		// measure
		f.SetCellValue("Sheet"+sheetNo, "J8", data.Grade.A.LeftValue)
		f.SetCellValue("Sheet"+sheetNo, "N8", data.Grade.A.RightValue)
		f.SetCellValue("Sheet"+sheetNo, "J9", data.Grade.B.LeftValue)
		f.SetCellValue("Sheet"+sheetNo, "N9", data.Grade.B.RightValue)
		f.SetCellValue("Sheet"+sheetNo, "J10", data.Grade.C.LeftValue)
		f.SetCellValue("Sheet"+sheetNo, "N10", data.Grade.C.RightValue)

		if factor == "IRI" || factor == "RUT" {
			f.SetCellValue("Sheet"+sheetNo, "J11", data.Grade.D.LeftValue)
			f.SetCellValue("Sheet"+sheetNo, "N11", data.Grade.D.RightValue)
		}

		// variable for last row
		l, _ := strconv.ParseFloat(data.SumRoadLengthStr, 64)
		avg, _ := strconv.ParseFloat(data.SumAvg, 64)
		a, _ := strconv.ParseFloat(data.SumCountA, 64)
		b, _ := strconv.ParseFloat(data.SumCountB, 64)
		c, _ := strconv.ParseFloat(data.SumCountC, 64)

		//table
		for iIndex, i := range data.Lane {

			if factor == "IRI" || factor == "RUT" {
				f.SetCellValue("Sheet"+sheetNo, "B"+fmt.Sprint(47+iIndex), i.No)
				f.SetCellValue("Sheet"+sheetNo, "C"+fmt.Sprint(47+iIndex), i.Length)
				f.SetCellValue("Sheet"+sheetNo, "E"+fmt.Sprint(47+iIndex), i.Avg)
				f.SetCellValue("Sheet"+sheetNo, "G"+fmt.Sprint(47+iIndex), i.StrCountA)
				f.SetCellValue("Sheet"+sheetNo, "I"+fmt.Sprint(47+iIndex), i.StrCountB)
				f.SetCellValue("Sheet"+sheetNo, "K"+fmt.Sprint(47+iIndex), i.StrCountC)
				f.SetCellValue("Sheet"+sheetNo, "M"+fmt.Sprint(47+iIndex), i.StrCountD)

				if iIndex == len(data.Lane)-1 {
					d, _ := strconv.ParseFloat(data.SumCountD, 64)

					f.SetCellValue("Sheet"+sheetNo, "C"+fmt.Sprint(47+iIndex+1), l)
					f.SetCellValue("Sheet"+sheetNo, "E"+fmt.Sprint(47+iIndex+1), avg)
					f.SetCellValue("Sheet"+sheetNo, "G"+fmt.Sprint(47+iIndex+1), a)
					f.SetCellValue("Sheet"+sheetNo, "I"+fmt.Sprint(47+iIndex+1), b)
					f.SetCellValue("Sheet"+sheetNo, "K"+fmt.Sprint(47+iIndex+1), c)
					f.SetCellValue("Sheet"+sheetNo, "M"+fmt.Sprint(47+iIndex+1), d)
				}

			} else if factor == "MPD" { //factor == MPD, GN
				f.SetCellValue("Sheet"+sheetNo, "B"+fmt.Sprint(47+iIndex), i.No)
				f.SetCellValue("Sheet"+sheetNo, "D"+fmt.Sprint(47+iIndex), i.Length)
				f.SetCellValue("Sheet"+sheetNo, "F"+fmt.Sprint(47+iIndex), i.Avg)
				f.SetCellValue("Sheet"+sheetNo, "I"+fmt.Sprint(47+iIndex), i.StrCountA)
				f.SetCellValue("Sheet"+sheetNo, "K"+fmt.Sprint(47+iIndex), i.StrCountB)
				f.SetCellValue("Sheet"+sheetNo, "M"+fmt.Sprint(47+iIndex), i.StrCountC)

				if iIndex == len(data.Lane)-1 {
					f.SetCellValue("Sheet"+sheetNo, "D"+fmt.Sprint(47+iIndex+1), l)
					f.SetCellValue("Sheet"+sheetNo, "F"+fmt.Sprint(47+iIndex+1), avg)
					f.SetCellValue("Sheet"+sheetNo, "I"+fmt.Sprint(47+iIndex+1), a)
					f.SetCellValue("Sheet"+sheetNo, "K"+fmt.Sprint(47+iIndex+1), b)
					f.SetCellValue("Sheet"+sheetNo, "M"+fmt.Sprint(47+iIndex+1), c)
				}
			} else {
				f.SetCellValue("Sheet"+sheetNo, "B"+fmt.Sprint(47+iIndex), i.No)
				f.SetCellValue("Sheet"+sheetNo, "D"+fmt.Sprint(47+iIndex), i.Length)
				f.SetCellValue("Sheet"+sheetNo, "F"+fmt.Sprint(47+iIndex), i.Avg)
				f.SetCellValue("Sheet"+sheetNo, "I"+fmt.Sprint(47+iIndex), i.StrCountA)
				f.SetCellValue("Sheet"+sheetNo, "K"+fmt.Sprint(47+iIndex), i.StrCountB)

				if iIndex == len(data.Lane)-1 {
					f.SetCellValue("Sheet"+sheetNo, "D"+fmt.Sprint(47+iIndex+1), l)
					f.SetCellValue("Sheet"+sheetNo, "F"+fmt.Sprint(47+iIndex+1), avg)
					f.SetCellValue("Sheet"+sheetNo, "I"+fmt.Sprint(47+iIndex+1), a)
					f.SetCellValue("Sheet"+sheetNo, "K"+fmt.Sprint(47+iIndex+1), b)
				}
			}
		}
		AddFooter(f, "Sheet"+sheetNo)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName(template)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

// func ExportExcelType7(data *models.DataReportDamage) (interface{}, error) {
// 	filePath := os.Getenv("GENARAL_EXCEL")
// 	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE7_EXCEL"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()

// 	folderPath := "storages/road/temp" + uuid.New().String() + "/"
// 	// create temp folder
// 	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
// 		logs.Error(err)
// 		return nil, err
// 	}

// 	// remove temp folder
// 	defer os.RemoveAll(folderPath)

// 	var wg sync.WaitGroup

// 	for _, i := range data.Detail {
// 		for _, j := range i.Position {
// 			wg.Add(1) // add number of go routine to wait
// 			go func(imgURL string) {
// 				defer wg.Done() // decrease the number of go routine to wait after go routine finished

// 				// trim img path
// 				trimmedURL := strings.TrimPrefix(imgURL, os.Getenv("STORAGE_IP")+"/")

// 				// open file
// 				srcImage, err := imaging.Open(trimmedURL)
// 				if err != nil {
// 					fmt.Println("can not open file", err)
// 					return
// 				}

// 				fileName := filepath.Base(trimmedURL)
// 				//saveAs := filepath.Join(folderPath, fileName)
// 				saveAs := folderPath + fileName
// 				// convert to .png
// 				err = imaging.Save(srcImage, saveAs)
// 				if err != nil {
// 					fmt.Println("can not save file", err)
// 					return
// 				}
// 			}(j.Image)
// 		}
// 	}

// 	wg.Wait() // waiting for go routine finish

// 	f.SetCellValue("Sheet1", "C3", "ชื่อสายทาง : "+data.RoadGroupName)
// 	f.SetCellValue("Sheet1", "C4", "ช่วง : "+data.RoadName+" รหัส : "+data.RoadCode)
// 	f.SetCellValue("Sheet1", "C5", "กม.เริ่มต้น "+fmt.Sprint(data.StrKmStart)+" กม.สิ้นสุด "+fmt.Sprint(data.StrKmEnd)+" ระยะทาง "+data.RoadLengthStr+" กม.")

// 	if !data.IsNull {
// 		iRow := 7
// 		iStrRow := fmt.Sprint(iRow)
// 		for iIndex, i := range data.Detail {
// 			f.SetCellValue("Sheet1", "B"+iStrRow, "ช่องจราจรที่ "+fmt.Sprint(i.LaneNo))

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+3), i.ACIcrack)          //รอยแตกต่อเนื่อง
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+3), i.CCTransverseCrack) //แตกตามขวาง

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+4), i.ACUcrack)             //รอบแตกไม่ต่อเนื่อง
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+4), i.CCNonTransverseCrack) //แนวทแยง

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+5), i.ACRavelling) //ผิวหลุด
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+5), i.CCFaulting)  //รอยเลื่อน

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+6), i.ACPatching) //รอยปะ
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+6), i.CCSpalling) //รอยบิ่น

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+7), i.ACPothole)      //หลุมบ่อ
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+7), i.CCCornerbreaks) //รอยแตกตามมุม

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+8), i.ACSurfaceDeform)   //เสียรูป
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+8), i.CCJointSealDamage) //รอยปะ

// 			f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+9), i.ACBleeding) //เยิ้ม
// 			f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+9), i.CCPatching) //ผิวหลุด

// 			jRow := iRow + 13
// 			jStrRow := fmt.Sprint(jRow)

// 			for jIndex, j := range i.Position {
// 				f.SetCellValue("Sheet1", "B"+jStrRow, j.No)
// 				f.SetCellValue("Sheet1", "C"+jStrRow, j.StrKm)
// 				f.SetCellValue("Sheet1", "E"+jStrRow, j.Surface)
// 				f.SetCellValue("Sheet1", "G"+jStrRow, j.DamageTypeENG+"\n"+j.DamageType)
// 				f.SetCellValue("Sheet1", "J"+jStrRow, j.Value)
// 				f.SetCellValue("Sheet1", "K"+jStrRow, j.Unit)
// 				// trim img path
// 				fileName := filepath.Base(j.Image)

// 				picRow := fmt.Sprint(jRow + 1)
// 				f.AddPicture("Sheet1", "L"+picRow, folderPath+fileName, &excelize.GraphicOptions{OffsetX: 8, OffsetY: -119, ScaleX: 0.078, ScaleY: 0.075})
// 				f.SetCellHyperLink("Sheet1", "L"+jStrRow, j.Image, "External")

// 				if jIndex < len(i.Position)-1 {
// 					f.DuplicateRow("Sheet1", jRow)
// 					jRow++
// 					jStrRow = fmt.Sprint(jRow)
// 				}
// 			}

// 			if iIndex < len(data.Detail)-1 {

// 				for count := 1; count <= 14; count++ {
// 					f.DuplicateRowTo("Sheet1", 6+count, jRow+1+count)
// 				}
// 				iRow = jRow + 2
// 				iStrRow = fmt.Sprint(iRow)
// 			}
// 		}
// 	}

// 	AddFooter(f, "Sheet1")

// 	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
// 		logs.Error(err)
// 		return nil, err
// 	}

// reportName, err := ReportName("TEMPLATE_GENARAL_TYPE7_EXCEL")
// 	if err != nil {
// 		logs.Error(err)
// 		return nil, err
// 	}
// 	code := fmt.Sprintf("%04d", rand.Intn(10000))

// 	name := code + "_" + reportName

// 	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

// 	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
// }

func ExportExcelType8(data *models.DataReportSurface) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE8_EXCEL"
	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//start

	f.SetCellValue("Sheet1", "C3", "ชื่อสายทาง : "+data.RoadGroupName)
	f.SetCellValue("Sheet1", "C4", "ช่วง : "+data.RoadName+" รหัส : "+data.RoadCode)
	f.SetCellValue("Sheet1", "C5", "กม.เริ่มต้น "+data.KmStart+" กม.สิ้นสุด "+data.KmEnd+" ระยะทาง "+data.StrRoadLength+" กม.")

	f.SetCellValue("Sheet1", "C8", "ข้อมูลสรุปรายละเอียดชนิดผิวทาง ปี "+data.Year)

	if !data.IsNull {

		// chart value
		if data.Surface.PMA.Sum != 0.00 {
			f.SetCellValue("Sheet1", "G13", data.Surface.PMA.Sum)
		}

		if data.Surface.AC.Sum != 0.00 {
			f.SetCellValue("Sheet1", "G14", data.Surface.AC.Sum)
		}

		if data.Surface.Slurry.Sum != 0.00 {
			f.SetCellValue("Sheet1", "G15", data.Surface.Slurry.Sum)
		}

		if data.Surface.Porous.Sum != 0.00 {
			f.SetCellValue("Sheet1", "G16", data.Surface.Porous.Sum)
		}

		if data.Surface.Concrete.Sum != 0.00 {
			f.SetCellValue("Sheet1", "G17", data.Surface.Concrete.Sum)
		}

		// PMA
		f.SetCellValue("Sheet1", "D29", data.Surface.PMA.StrLane5)
		f.SetCellValue("Sheet1", "E29", data.Surface.PMA.StrLane4)
		f.SetCellValue("Sheet1", "F29", data.Surface.PMA.StrLane3)
		f.SetCellValue("Sheet1", "G29", data.Surface.PMA.StrLane2)
		f.SetCellValue("Sheet1", "H29", data.Surface.PMA.StrLane1)
		f.SetCellValue("Sheet1", "I29", data.Surface.PMA.StrSum)

		// AC
		f.SetCellValue("Sheet1", "D30", data.Surface.AC.StrLane5)
		f.SetCellValue("Sheet1", "E30", data.Surface.AC.StrLane4)
		f.SetCellValue("Sheet1", "F30", data.Surface.AC.StrLane3)
		f.SetCellValue("Sheet1", "G30", data.Surface.AC.StrLane2)
		f.SetCellValue("Sheet1", "H30", data.Surface.AC.StrLane1)
		f.SetCellValue("Sheet1", "I30", data.Surface.AC.StrSum)

		// slurry
		f.SetCellValue("Sheet1", "D31", data.Surface.Slurry.StrLane5)
		f.SetCellValue("Sheet1", "E31", data.Surface.Slurry.StrLane4)
		f.SetCellValue("Sheet1", "F31", data.Surface.Slurry.StrLane3)
		f.SetCellValue("Sheet1", "G31", data.Surface.Slurry.StrLane2)
		f.SetCellValue("Sheet1", "H31", data.Surface.Slurry.StrLane1)
		f.SetCellValue("Sheet1", "I31", data.Surface.Slurry.StrSum)

		// porous
		f.SetCellValue("Sheet1", "D32", data.Surface.Porous.StrLane5)
		f.SetCellValue("Sheet1", "E32", data.Surface.Porous.StrLane4)
		f.SetCellValue("Sheet1", "F32", data.Surface.Porous.StrLane3)
		f.SetCellValue("Sheet1", "G32", data.Surface.Porous.StrLane2)
		f.SetCellValue("Sheet1", "H32", data.Surface.Porous.StrLane1)
		f.SetCellValue("Sheet1", "I32", data.Surface.Porous.StrSum)

		// concrete
		f.SetCellValue("Sheet1", "D33", data.Surface.Concrete.StrLane5)
		f.SetCellValue("Sheet1", "E33", data.Surface.Concrete.StrLane4)
		f.SetCellValue("Sheet1", "F33", data.Surface.Concrete.StrLane3)
		f.SetCellValue("Sheet1", "G33", data.Surface.Concrete.StrLane2)
		f.SetCellValue("Sheet1", "H33", data.Surface.Concrete.StrLane1)
		f.SetCellValue("Sheet1", "I33", data.Surface.Concrete.StrSum)

		// summary
		f.SetCellValue("Sheet1", "D34", data.Surface.Sum.StrLane5)
		f.SetCellValue("Sheet1", "E34", data.Surface.Sum.StrLane4)
		f.SetCellValue("Sheet1", "F34", data.Surface.Sum.StrLane3)
		f.SetCellValue("Sheet1", "G34", data.Surface.Sum.StrLane2)
		f.SetCellValue("Sheet1", "H34", data.Surface.Sum.StrLane1)
		f.SetCellValue("Sheet1", "I34", data.Surface.Sum.StrSum)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	AddFooter(f, "Sheet1")

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE8_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelType9(data models.DataResponseReportMaintenanceTracking) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE9_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for i := 1; i <= 100; i++ {
		if i > len(data.Project) {
			f.DeleteSheet("Sheet1 " + "(" + fmt.Sprint(i) + ")")
		}
	}

	if data.IsNull {
		f.SetSheetName("Sheet1 (100)", "ไม่มีโครงการ")
		f.SetCellValue("ไม่มีโครงการ", "F4", "ปีงบประมาณ พ.ศ. "+fmt.Sprint(data.Year))

		AddFooter(f, "ไม่มีโครงการ")

	} else if !data.IsNull {

		for iIndex, i := range data.Project {
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "F4", "ปีงบประมาณ พ.ศ. "+fmt.Sprint(data.Year))
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K8", i.ProjectName)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K9", i.RoadGroupName)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K10", i.ContractNumber)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AC10", i.BudgetYear)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K11", i.BudgetType)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AC11", i.MaintenanceType)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K12", i.ContractorName)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AC12", i.AdviserName)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K13", i.ProjectSecretaryName)
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AC13", i.StrBudgetMaintenance+" บาท")
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K14", i.StrMiddlePrice+" บาท")
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AC14", i.StrContractWorkValue+" บาท")
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "K15", i.StrBudgetProcurement+" บาท")

			Row := 19
			StrRow := fmt.Sprint(Row)

			// create the table that has no MaintenanceStandardName
			if i.MaintenanceDetail[0].MaintenanceStandardName == "" {

				f.MergeCell("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D18", "P18")
				f.MergeCell("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D19", "P19")
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D18", "ลำดับ")
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "Q18", "ช่องจราจร")
			}

			for jIndex, j := range i.MaintenanceDetail {

				// no MaintenanceStandardName
				if j.MaintenanceStandardName == "" && jIndex != len(i.MaintenanceDetail)-1 {

					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+StrRow, j.No)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D"+StrRow, j.RoadName)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "Q"+StrRow, j.Lane)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "X"+StrRow, j.StrKmStart+" - "+j.StrKmEnd)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AE"+StrRow, j.StrDistance)

					f.DuplicateRow("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", Row)
					Row++
					StrRow = fmt.Sprint(Row)

				} else {
					// with MaintenanceStandardName
					if jIndex == len(i.MaintenanceDetail)-1 {
						f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+StrRow, "")
						f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D"+StrRow, "")
						f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "J"+StrRow, "")
						f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "Q"+StrRow, "")
						f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "X"+StrRow, "รวม")
						f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AE"+StrRow, j.StrSumDistance)
						break
					}

					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+StrRow, j.No)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D"+StrRow, j.RoadName)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "J"+StrRow, j.Lane)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "Q"+StrRow, j.MaintenanceStandardName)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "X"+StrRow, j.StrKmStart+" - "+j.StrKmEnd)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AE"+StrRow, j.StrDistance)

					f.DuplicateRow("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", Row)
					Row++
					StrRow = fmt.Sprint(Row)

				}
			}

			Row += 5
			StrRow = fmt.Sprint(Row)

			for _, k := range i.ProgressDetail {

				if k.LastRow {
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+StrRow, "รวม")
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "F"+StrRow, k.DisProgressPlan)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "J"+StrRow, k.DisSumProgressPlan)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "N"+StrRow, k.DisProgress)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "R"+StrRow, k.DisSumProgress)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "V"+StrRow, k.StrDisbursementPlan)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "Z"+StrRow, k.StrSumDisbursementPlan)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AD"+StrRow, k.StrDisbursement)
					f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AH"+StrRow, k.StrSumDisbursement)
					break
				}

				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+StrRow, k.StrSchedule)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "F"+StrRow, k.DisProgressPlan)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "J"+StrRow, k.DisSumProgressPlan)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "N"+StrRow, k.DisProgress)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "R"+StrRow, k.DisSumProgress)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "V"+StrRow, k.StrDisbursementPlan)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "Z"+StrRow, k.StrSumDisbursementPlan)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AD"+StrRow, k.StrDisbursement)
				f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "AH"+StrRow, k.StrSumDisbursement)

				f.DuplicateRow("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", Row)
				Row++
				StrRow = fmt.Sprint(Row)
			}

			f.SetSheetName("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", i.ProjectName)

			AddFooter(f, i.ProjectName)
		}
	}
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE9_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelType10(data []models.DataReportMaintenance) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE10_EXCEL_" + fmt.Sprint((data[0].YearEnd-data[0].YearStart)+1)
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

	for iIndex, i := range data {
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "D5", "ปีงบประมาณ พ.ศ. "+fmt.Sprint(data[0].YearStart)+"-"+fmt.Sprint(data[0].YearEnd))

		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B7", "ชื่อสายทาง : "+i.RoadName)

		yearColumn := []string{"D", "J", "P", "V", "AB", "AH", "AN", "AT", "AZ", "BF"}
		laneColumn := [][]string{
			{"D", "F", "H"},
			{"J", "L", "N"},
			{"P", "R", "T"},
			{"V", "X", "Z"},
			{"AB", "AD", "AF"},
			{"AH", "AJ", "AL"},
			{"AN", "AP", "AR"},
			{"AT", "AV", "AX"},
			{"AZ", "BB", "BD"},
			{"BF", "BH", "BJ"}}

		for yIndex, y := range i.Years {
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", yearColumn[yIndex]+"8", y)
		}

		for jIndex, j := range i.RowExcel {

			if jIndex < len(i.RowExcel)-1 {
				f.DuplicateRow("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", 12+jIndex)
			}

			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", "B"+fmt.Sprint(12+jIndex), j.StrKmStart+" - "+j.StrKmEnd)

			for yearIndex, k := range j.Data {
				for laneIndex, l := range k.Data {

					if l.Data.KmStart != "" {
						f.SetCellRichText("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), []excelize.RichTextRun{
							{
								Text: l.Data.Data.StrKmStart + " - " + l.Data.Data.StrKmEnd + "\n",
							},
							{
								Text: l.Data.Data.Method,
							},
						})
					}

					preCell, _ := f.GetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1))
					currentCell, _ := f.GetCellValue("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex))

					if preCell != "" && currentCell != "" && preCell == currentCell {

						f.MergeCell("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex))

						// style and wrap text
						var color string
						switch l.Data.Data.Method {
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
								Color:   []string{color},
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

						f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), wrapTextStyle)

						if data[iIndex].RowExcel[jIndex-1].Data[yearIndex].Data[laneIndex].Data.Span == 2 {
							f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", 12+jIndex-1, 30)
							f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", 12+jIndex, 30)
						}
					} else if currentCell != "" {

						// style and wrap text
						var color string
						switch l.Data.Data.Method {
						case "OL-Overlay":
							color = "#63b598"
						case "M&OL-Mill&Overlay":
							color = "#ce7d78"
						case "RCL-Recycling":
							color = "#ea9e70"
						case "Rc-Reconstruction":
							color = "#a48a9e"
						case "SS-SlurrySeal":
							color = "#c6e1e8"
						case "FDR":
							color = "#648177"
						case "BCO":
							color = "#0d5ac1"
						case "M&OL":
							color = "#f205e6"
						default:
							color = "#1c0365"
						}

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
								Color:   []string{color},
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

						f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), wrapTextStyle)

						if l.Data.Span == 1 {
							f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", 12+jIndex, 60)
						}
					}
				}
			}
		}

		f.SetSheetName("Sheet1 "+"("+fmt.Sprint(iIndex+1)+")", i.RoadName)
		AddFooter(f, i.RoadName)
		//AddFooter(f, "Sheet1 "+"("+fmt.Sprint(iIndex+1)+")")
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE10_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelTrafficVolume(datas []responses.ReportTrafficVolume) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := os.Getenv("TEMPLATE_GENARAL_TRAFFIC_VOLUME_EXCEL")
	f, err := excelize.OpenFile(template)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	templateSheet := "Sheet1"
	templateSheetIndex, _ := f.GetSheetIndex(templateSheet)
	if templateSheetIndex == -1 {
		return nil, fmt.Errorf("template sheet '%s' does not exist", templateSheet)
	}

	for i, data := range datas {

		sheetName := fmt.Sprintf("Sheet%d", i+1)
		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}

		var newSheetIndex int
		if sheetName != "Sheet1" {
			newSheetIndex, _ = f.NewSheet(sheetName)
			if newSheetIndex == 0 {
				return nil, fmt.Errorf("failed to create new sheet '%s'", sheetName)
			}

			if err := f.CopySheet(templateSheetIndex, newSheetIndex); err != nil {
				return nil, fmt.Errorf("failed to copy template sheet to '%s': %w", sheetName, err)
			}

			// Get pictures from the template sheet
			pictures, err := f.GetPictures("Sheet1", "B2")
			if err != nil {
				return nil, fmt.Errorf("failed to get pictures from template sheet '%s': %w", templateSheet, err)
			}

			for _, pic := range pictures {

				err = f.AddPictureFromBytes(sheetName, "B2", &excelize.Picture{
					File:      pic.File,
					Extension: pic.Extension,
					Format: &excelize.GraphicOptions{
						OffsetX: -20,
						OffsetY: 0,
						ScaleX:  0.5,
						ScaleY:  0.5,
					},
				})
				if err != nil {
					return nil, fmt.Errorf("failed to add picture to sheet '%s': %w", sheetName, err)
				}

			}
		}

		f.SetCellValue(sheetName, "C2", "รายงานข้อมูลปริมาณจราจร ปี พ.ศ. "+data.Year)
		f.SetCellValue(sheetName, "C3", "หมายเลขทางหลวง : "+data.RoadGroupName+" ตอนควบคุม : "+data.RoadSectionName)
		f.SetCellValue(sheetName, "C4", "ชื่อสายทาง : "+data.RoadName)
		f.SetCellValue(sheetName, "C5", "กม.เริ่มต้น : "+data.KmStart+" กม.สิ้นสุด  : "+data.KmEnd+" ระยะทาง : "+data.TotalKm+" กม.")

		f.SetCellValue(sheetName, "E9", data.Veh1)
		f.SetCellValue(sheetName, "H9", data.SurveyedDate)
		f.SetCellValue(sheetName, "E10", data.Veh2)
		f.SetCellValue(sheetName, "H10", data.SurveyedDate)
		f.SetCellValue(sheetName, "E11", data.Veh3)
		f.SetCellValue(sheetName, "H11", data.SurveyedDate)
		f.SetCellValue(sheetName, "E12", data.Total)

		AddFooter(f, sheetName)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := ReportName("TEMPLATE_GENARAL_TRAFFIC_VOLUME_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	name := code + "_" + reportName

	f.SaveAs(filePath + name + ".xlsx")
	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

// func ExportExcelTrafficVolume(datas []responses.ReportTrafficVolume) (interface{}, error) {
// 	filePath := os.Getenv("GENARAL_EXCEL")
// 	template := os.Getenv("TEMPLATE_GENARAL_TRAFFIC_VOLUME_EXCEL")
// 	f, err := excelize.OpenFile(template)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()

// 	// Template sheet name
// 	templateSheet := "Sheet1"

// 	// Ensure the template sheet exists and get its index
// 	templateSheetIndex, _ := f.GetSheetIndex(templateSheet)
// 	if templateSheetIndex == 0 {
// 		return nil, fmt.Errorf("template sheet '%s' does not exist", templateSheet)
// 	}

// 	for i, data := range datas {
// 		sheetName := fmt.Sprintf("Sheet%d", i+1)

// 		// Ensure sheet name is within the 31-character limit
// 		if len(sheetName) > 31 {
// 			sheetName = sheetName[:31]
// 		}

// 		// Create a new sheet and get its index
// 		newSheetIndex, _ := f.NewSheet(sheetName)
// 		if newSheetIndex == 0 {
// 			return nil, fmt.Errorf("failed to create new sheet '%s'", sheetName)
// 		}

// 		fmt.Printf("New sheet '%s' created with index %d\n", sheetName, newSheetIndex)

// 		// Copy the template sheet to the new sheet
// 		if err := f.CopySheet(templateSheetIndex, newSheetIndex); err != nil {
// 			return nil, err
// 		}

// 		f.SetCellValue(sheetName, "B2", "รายงานข้อมูลปริมาณจราจร ปี พ.ศ. "+data.Year)
// 		f.SetCellValue(sheetName, "B3", "หมายเลขทางหลวง : "+data.RoadGroupName+" ตอนควบคุม : "+data.RoadSectionName+" ชื่อสายทาง : "+data.RoadName)
// 		f.SetCellValue(sheetName, "B4", "กม.เริ่มต้น : "+data.KmStart+" กม.สิ้นสุด  : "+data.KmEnd+" ระยะทาง : "+data.TotalKm+" กม.")

// 		f.SetCellValue(sheetName, "E8", data.Veh1)
// 		f.SetCellValue(sheetName, "H8", data.SurveyedDate)
// 		f.SetCellValue(sheetName, "E9", data.Veh2)
// 		f.SetCellValue(sheetName, "H9", data.SurveyedDate)
// 		f.SetCellValue(sheetName, "E10", data.Veh3)
// 		f.SetCellValue(sheetName, "H11", data.SurveyedDate)
// 		f.SetCellValue(sheetName, "E12", data.Total)

// 		AddFooter(f, sheetName)
// 	}

// 	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
// 		logs.Error(err)
// 		return nil, err
// 	}

// 	reportName, err := ReportName("TEMPLATE_GENARAL_TRAFFIC_VOLUME_EXCEL")
// 	if err != nil {
// 		logs.Error(err)
// 		return nil, err
// 	}
// 	code := fmt.Sprintf("%04d", rand.Intn(10000))

// 	name := code + "_" + reportName

// 	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

// 	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
// }

func parseCellRange(cellRange string) (startCol, endCol string, startRow, endRow int, err error) {
	re := regexp.MustCompile(`([A-Z]+)([0-9]+):([A-Z]+)([0-9]+)`)
	matches := re.FindStringSubmatch(cellRange)
	if len(matches) != 5 {
		return "", "", 0, 0, fmt.Errorf("invalid cell range: %s", cellRange)
	}

	startCol, endCol = matches[1], matches[3]
	startRow, err = strconv.Atoi(matches[2])
	if err != nil {
		return "", "", 0, 0, err
	}
	endRow, err = strconv.Atoi(matches[4])
	if err != nil {
		return "", "", 0, 0, err
	}

	return startCol, endCol, startRow, endRow, nil
}

func copyStyles(f *excelize.File, srcSheet, dstSheet, cellRange string) error {
	startCol, endCol, startRow, endRow, err := parseCellRange(cellRange)
	if err != nil {
		return err
	}

	// Convert column letters to numbers
	startColNum, _ := excelize.ColumnNameToNumber(startCol)
	endColNum, _ := excelize.ColumnNameToNumber(endCol)

	for row := startRow; row <= endRow; row++ {
		for col := startColNum; col <= endColNum; col++ {
			cell, err := excelize.CoordinatesToCellName(col, row)
			if err != nil {
				return err
			}

			styleID, err := f.GetCellStyle(srcSheet, cell)
			if err != nil {
				return err
			}

			if err := f.SetCellStyle(dstSheet, cell, cell, styleID); err != nil {
				return err
			}
		}
	}
	return nil
}

func ExportExcelType12(data *models.DataReportAccidentVolume) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE12_EXCEL"
	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	f.SetCellValue("Sheet1", "B3", "รายงานข้อมูลอุบัติเหตุ ปี พ.ศ. "+fmt.Sprint(data.Year))
	f.SetCellValue("Sheet1", "D4", "รหัส : "+data.Code+" ชื่อสายทาง : "+data.Name)

	if !data.IsNull {
		f.SetCellValue("Sheet1", "E8", data.StrAcc1)
		f.SetCellValue("Sheet1", "H8", data.StrSurveyedDate)
		f.SetCellValue("Sheet1", "E9", data.StrAcc2)
		f.SetCellValue("Sheet1", "H9", data.StrSurveyedDate)
		f.SetCellValue("Sheet1", "E10", data.StrAcc3)
		f.SetCellValue("Sheet1", "H10", data.StrSurveyedDate)
		f.SetCellValue("Sheet1", "E11", data.StrAcc4)
		f.SetCellValue("Sheet1", "H11", data.StrSurveyedDate)
		//f.SetCellValue("Sheet1", "E12", data.StrSum)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	AddFooter(f, "Sheet1")

	reportName, err := ReportName("TEMPLATE_GENARAL_TYPE12_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelRoad(data []models.RoadListReport) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_ROAD_EXCEL"
	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	styleHead, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   16,
			Family: "TH SarabunPSK",
			Bold:   true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#E6E6E6"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
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

	styleSubHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   16,
			Family: "TH SarabunPSK",
			Bold:   true,
		},
	})

	styleBody, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   16,
			Family: "TH SarabunPSK",
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
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

	// f.SetCellValue("Sheet1", "B3", "รายงานข้อมูลอุบัติเหตุ ปี พ.ศ. "+fmt.Sprint(data.Year))
	// f.SetCellValue("Sheet1", "D4", "รหัส : "+data.Code+" ชื่อสายทาง : "+data.Name)
	rowStart := 5
	for _, roadGroup := range data {
		row := fmt.Sprintf("%d", rowStart)
		f.SetCellValue("Sheet1", "B"+row, "ทางหลวงหมายเลข "+roadGroup.Number+" ถนน "+roadGroup.ShortName+" กม.เริ่มต้น "+roadGroup.KmStart+" กม.สิ้นสุด "+roadGroup.KmEnd+" ระยะทาง "+roadGroup.StrDistance+" กม.")
		f.SetCellStyle("Sheet1", "B"+row, "B"+row, styleSubHeader)
		rowStart++
		row = fmt.Sprintf("%d", rowStart)
		f.SetCellValue("Sheet1", "B"+row, "ลำดับ")
		f.SetCellStyle("Sheet1", "B"+row, "B"+row, styleHead)

		f.SetCellValue("Sheet1", "C"+row, "หมายเลขทางหลวง")
		f.SetCellStyle("Sheet1", "c"+row, "C"+row, styleHead)

		f.SetCellValue("Sheet1", "D"+row, "ตอนควบคุม")
		f.SetCellStyle("Sheet1", "D"+row, "D"+row, styleHead)

		f.SetCellValue("Sheet1", "E"+row, "ชื่อตอนควบคุม")
		f.SetCellStyle("Sheet1", "E"+row, "E"+row, styleHead)

		f.SetCellValue("Sheet1", "F"+row, "กม.เริ่มต้น")
		f.SetCellStyle("Sheet1", "F"+row, "F"+row, styleHead)

		f.SetCellValue("Sheet1", "G"+row, "กม.สิ้นสุด")
		f.SetCellStyle("Sheet1", "G"+row, "G"+row, styleHead)

		f.SetCellValue("Sheet1", "H"+row, "ระยะทาง(กม.)")
		f.SetCellStyle("Sheet1", "H"+row, "H"+row, styleHead)

		f.SetCellValue("Sheet1", "I"+row, "หมวดทางหลวง พิเศษระหว่างเมือง")
		f.SetCellStyle("Sheet1", "I"+row, "I"+row, styleHead)
		for _, section := range roadGroup.Sections {
			rowStart++
			row := fmt.Sprintf("%d", rowStart)
			f.SetCellValue("Sheet1", "B"+row, section.No)
			f.SetCellStyle("Sheet1", "B"+row, "B"+row, styleBody)

			f.SetCellValue("Sheet1", "C"+row, roadGroup.Number)
			f.SetCellStyle("Sheet1", "C"+row, "C"+row, styleBody)

			f.SetCellValue("Sheet1", "D"+row, section.Number)
			f.SetCellStyle("Sheet1", "D"+row, "D"+row, styleBody)

			if section.NameOriginTH == section.NameDestinationTH {
				f.SetCellValue("Sheet1", "E"+row, section.NameOriginTH)
			} else {
				f.SetCellValue("Sheet1", "E"+row, section.NameOriginTH+" - "+section.NameDestinationTH)
			}
			f.SetCellStyle("Sheet1", "E"+row, "E"+row, styleBody)

			f.SetCellValue("Sheet1", "F"+row, section.KmStart)
			f.SetCellStyle("Sheet1", "F"+row, "F"+row, styleBody)

			f.SetCellValue("Sheet1", "G"+row, section.KmEnd)
			f.SetCellStyle("Sheet1", "G"+row, "G"+row, styleBody)

			f.SetCellValue("Sheet1", "H"+row, section.StrDistance)
			f.SetCellStyle("Sheet1", "H"+row, "H"+row, styleBody)

			f.SetCellValue("Sheet1", "I"+row, section.RefDepot.Name)
			f.SetCellStyle("Sheet1", "I"+row, "I"+row, styleBody)

		}
		rowStart++
		rowStart++
		//f.SetCellValue("Sheet1", "E12", data.StrSum)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	AddFooter(f, "Sheet1")

	reportName, err := ReportName("TEMPLATE_GENARAL_ROAD_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
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
