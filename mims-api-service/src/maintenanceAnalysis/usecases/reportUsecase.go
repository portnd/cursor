package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"golang.org/x/sync/errgroup"
)

func (u *UseCase) GetReportType1(ID, userID int, typeFile string) (interface{}, error) {
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	interventionCriteria, _ := u.Repo.GetInterventionCriteriaParamsById(analysisData.InterventionCriteriaParmasID)

	// var jsonObject map[string]interface{}
	var jsonObject responses.InterventionCriteriaParams
	err = json.Unmarshal([]byte(interventionCriteria.Params), &jsonObject)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	// return interventionCriteriaParams, nil
	// return jsonObject, nil
	// var dataRes responses.AnalyzeReport1
	var interventionCriteriaData []responses.InterventionCriteriaReportRes
	// =============================== asphalt ===============================
	criteriaMethodAS, err := u.Repo.GetRefCriteriaMethodBySurface("as")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	// return criteriaMethodAS, nil
	asphalts := jsonObject.Asphalt
	asphaltData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range asphalts {
		asphaltData[item.Name] = item
	}
	for _, method := range criteriaMethodAS {
		val, isVal := asphaltData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}
		// return asphalt.(map[string]interface{})["maintenance_sequence"], nil
		interventionCriterias := asphaltData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReportRes
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes
			for _, item2 := range conditions {

				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}

			name := item.MaintenanceStandardName
			description := item.MaintenanceDescription
			data.Seq = item.MaintenanceSequence
			data.StandardName = name
			data.Description = description
			data.Conditions = conditionDatas
			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
		}
	}

	// =============================== concrete ===============================
	criteriaMethodCC, err := u.Repo.GetRefCriteriaMethodBySurface("cc")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	concrete := jsonObject.Concrete
	concreteData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range concrete {
		concreteData[item.Name] = item
	}

	for _, method := range criteriaMethodCC {
		val, isVal := concreteData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}

		// return yyy, nil
		interventionCriterias := concreteData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReportRes
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes
			for _, item2 := range conditions {
				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}
			name := item.MaintenanceStandardName
			description := item.MaintenanceDescription
			data.Seq = item.MaintenanceSequence
			data.StandardName = name
			data.Description = description
			data.Conditions = conditionDatas
			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
		}
	}

	sort.Slice(interventionCriteriaData, func(i, j int) bool {
		return interventionCriteriaData[i].Seq < interventionCriteriaData[j].Seq
	})
	var interventionCriteriaReportData responses.InterventionCriteriaReportData
	userUpdate, _ := u.Repo.GetUserByID(uint(interventionCriteria.CreatedBy))
	userPrint, _ := u.Repo.GetUserByID(uint(userID))
	interventionCriteriaReportData.UpdatedBy = userUpdate.Firstname + " " + userUpdate.Lastname
	interventionCriteriaReportData.UpdatedAt = helpers.ConvertToThaiFullCalendar(interventionCriteria.CreatedAt)
	interventionCriteriaReportData.User = userPrint.Firstname + " " + userPrint.Lastname
	interventionCriteriaReportData.PrintDate = helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour * 7))

	interventionCriteriaReportData.InterventionCriteria = interventionCriteriaData
	filePath := os.Getenv("MAINTENANCE_ANALYSIS_PDF") + fmt.Sprintf("%d/", ID)

	templateName := os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE1")

	result, err := helpers.InitDataToHtml(templateName, interventionCriteriaReportData, filePath)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	if typeFile == "html" {
		texts := strings.Split(result, "/")
		text := texts[len(texts)-1]

		new := strings.ReplaceAll(result, text, "เงื่อนไขการซ่อมบำรุง.html")

		os.Rename(result, new)

		return os.Getenv("STORAGE_IP") + "/" + new, nil
	}

	// return result, nil
	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return "", err
	}

	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, helpers.PrintToPDF(string(html), &buf, false)); err != nil {
		logs.Error(err)
		return "", err
	}

	// unix := strconv.Itoa(int(time.Now().Unix()))

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}
	fullFilePath := filePath + "เงื่อนไขการซ่อมบำรุง" + ".pdf"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}
	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *UseCase) GetReportType2(id, userID int, typeFile string) (interface{}, error) {
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	interventionCriteria, _ := u.Repo.GetInterventionCriteriaParamsById(analysisData.InterventionCriteriaParmasID)
	var jsonObject responses.InterventionCriteriaParams
	err = json.Unmarshal([]byte(interventionCriteria.Params), &jsonObject)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	var interventionCriteriaData []responses.InterventionCriteriaReport2Res
	// =============================== asphalt ===============================
	criteriaMethodAS, err := u.Repo.GetRefCriteriaMethodBySurface("as")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	asphalts := jsonObject.Asphalt
	asphaltData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range asphalts {
		asphaltData[item.Name] = item
	}
	i := 1
	j := 1
	for _, method := range criteriaMethodAS {

		val, isVal := asphaltData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}

		interventionCriterias := asphaltData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReport2Res
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes

			for _, item2 := range conditions {
				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}

			name := item.MaintenanceStandardName
			costPerUnit := item.MaintenanceCostPerUnit
			data.No = j
			data.Seq = float64(item.MaintenanceSequence)
			data.StandardName = name
			data.MaintenanceCostPerUnit = helpers.FormatNumberFloat(costPerUnit)

			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
			j++
		}
		i++
	}

	// =============================== concrete ===============================
	criteriaMethodCC, err := u.Repo.GetRefCriteriaMethodBySurface("cc")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	concrete := jsonObject.Concrete
	concreteData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range concrete {
		concreteData[item.Name] = item
	}
	for _, method := range criteriaMethodCC {
		val, isVal := concreteData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}

		interventionCriterias := concreteData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReport2Res
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes

			for _, item2 := range conditions {
				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}

			name := item.MaintenanceStandardName
			costPerUnit := item.MaintenanceCostPerUnit
			data.No = j
			data.Seq = float64(item.MaintenanceSequence)
			data.StandardName = name
			data.MaintenanceCostPerUnit = helpers.FormatNumberFloat(costPerUnit)

			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
			j++
		}
	}

	sort.Slice(interventionCriteriaData, func(i, j int) bool {
		return interventionCriteriaData[i].Seq < interventionCriteriaData[j].Seq
	})
	var interventionCriteriaDatas []responses.InterventionCriteriaReport2Res
	for index, item := range interventionCriteriaData {
		var interventionCriteriaData responses.InterventionCriteriaReport2Res
		interventionCriteriaData.No = index + 1
		interventionCriteriaData.MaintenanceCostPerUnit = item.MaintenanceCostPerUnit
		interventionCriteriaData.StandardName = item.StandardName
		interventionCriteriaDatas = append(interventionCriteriaDatas, interventionCriteriaData)
	}
	// return interventionCriteriaDatas, nil
	// dataRes.Data = interventionCriteriaDatas
	filePath := os.Getenv("MAINTENANCE_ANALYSIS_PDF") + fmt.Sprintf("%d/", id)
	templateName := os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE2")
	// var data interface{}
	// return interventionCriteriaData, nil
	var interventionCriteriaRes responses.InterventionCriteriaReport2Data

	userUpdate, _ := u.Repo.GetUserByID(uint(interventionCriteria.CreatedBy))
	userPrint, _ := u.Repo.GetUserByID(uint(userID))
	interventionCriteriaRes.UpdatedBy = userUpdate.Firstname + " " + userUpdate.Lastname
	interventionCriteriaRes.UpdatedAt = helpers.ConvertToThaiFullCalendar(interventionCriteria.CreatedAt)
	interventionCriteriaRes.User = userPrint.Firstname + " " + userPrint.Lastname
	interventionCriteriaRes.PrintDate = helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour * 7))
	interventionCriteriaRes.InterventionCriteria = interventionCriteriaDatas
	result, err := helpers.InitDataToHtml(templateName, interventionCriteriaRes, filePath)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	if typeFile == "html" {
		texts := strings.Split(result, "/")
		text := texts[len(texts)-1]

		new := strings.ReplaceAll(result, text, "เงื่อนไขค่าใช้จ่ายการซ่อมบำรุง.html")

		fmt.Println(new)

		os.Rename(result, new)

		return os.Getenv("STORAGE_IP") + "/" + new, nil
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return "", err
	}

	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, helpers.PrintToPDF(string(html), &buf, false)); err != nil {
		logs.Error(err)
		return "", err
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}
	fullFilePath := filePath + "เงื่อนไขค่าใช้จ่ายการซ่อมบำรุง" + ".pdf"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *UseCase) GetReportType3(id, userID int, typeFile string) (interface{}, error) {
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	condition, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(analysisData.Condition)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	target, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*analysisData.Target)
	if err != nil {
		logs.Error(err)
		// return "", responses.NewAppErr(400, err.Error())
	}
	var reportYear interface{}
	var report3 interface{}
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		reportYear, _ = u.Report3YearData(id)
		report3, _ = u.Report3Data(id)
	} else {
		reportYear, _ = u.Report3YearDataAnnual(id)
		report3, _ = u.Report3DataAnnual(id)
	}

	modelResults, err := u.Repo.GetModelResultDataById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	yearMinData := []int{}
	for _, item := range modelResults {
		yearMinData = append(yearMinData, item.AnalystYear)
	}

	yeatMin := helpers.FindMin(yearMinData)

	// analystYear := prepareData.AnalystYear
	year := analysisData.Year
	user, _ := u.Repo.GetUserByID(uint(userID))

	var data responses.Report3Data
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		data.Title = fmt.Sprintf("พ.ศ. %d - %d", yeatMin+543, yeatMin+*year-1+543)
	} else {
		data.Title = fmt.Sprintf("พ.ศ. %d", yeatMin+543)
	}

	data.Report3 = report3.([]responses.Report3)
	data.Year = reportYear
	data.Date = helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour * 7))
	data.User = user.Firstname + " " + user.Lastname
	data.Condition = condition.Name
	data.Target = target.Name

	// idInteger, err := helpers.ConvertStringToInt(id)

	filePath := os.Getenv("MAINTENANCE_ANALYSIS_PDF") + fmt.Sprintf("%d/", id)
	templateName := os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE3")

	// var data interface{}
	result, err := helpers.InitDataToHtmlFunc(templateName, data, filePath, "type3.html")
	if err != nil {
		logs.Error(err)
		return "", err
	}

	if typeFile == "html" {
		texts := strings.Split(result, "/")
		text := texts[len(texts)-1]

		new := strings.ReplaceAll(result, text, "สรุปค่าซ่อมบำรุงและค่าIRIของแต่ละปี.html")

		os.Rename(result, new)

		return os.Getenv("STORAGE_IP") + "/" + new, nil
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return "", err
	}

	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, helpers.PrintToPDF(string(html), &buf, false)); err != nil {
		logs.Error(err)
		return "", err
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}
	fullFilePath := filePath + "สรุปค่าซ่อมบำรุงและค่าIRIของแต่ละปี.pdf"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *UseCase) Report3Data(id int) (interface{}, error) {
	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisPlan, err := u.Repo.GetMaintenanceAnalysisPlanByAnalysisID(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisPlanMap := make(map[int]map[string]float64)
	for _, item := range maintenanceAnalysisPlan {
		maintenanceAnalysisPlanMapSub := map[string]float64{}
		for indexPlan := 1; indexPlan < *maintenanceAnalysis.NumberPlan+1; indexPlan++ {
			switch indexPlan {
			case 1:
				maintenanceAnalysisPlanMapSub["plan_1"] = *item.Plan1
			case 2:
				maintenanceAnalysisPlanMapSub["plan_2"] = *item.Plan2
			case 3:
				maintenanceAnalysisPlanMapSub["plan_3"] = *item.Plan3
			}
		}
		maintenanceAnalysisPlanMap[int(item.PlanYear)] = maintenanceAnalysisPlanMapSub
	}

	numberPlanDefault := 2
	numberPlanWithDefault := *maintenanceAnalysis.NumberPlan + numberPlanDefault
	numberPlan := *maintenanceAnalysis.NumberPlan

	var dataTables []responses.Report3
	for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
		var dataTable responses.Report3
		plan := ""
		seq := 0.0
		if indexPlan > numberPlan {
			if (indexPlan - numberPlan) == 2 {
				plan = "ซ่อมบำรุงปกติ"
				seq = 0
			} else {
				plan = "ไม่จำกัดงบประมาณ"
				seq = 0.1
			}
		} else {
			plan = "แผนที่ " + strconv.Itoa(indexPlan)
			seq = float64(indexPlan)
		}

		title := []string{"งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)", "B/C", "ประโยชน์ผู้ใช้ทาง (ล้านบาท)"}
		var dataTables2 []responses.Report3
		for _, item := range title {
			values := []float64{}
			var dataTable2 responses.Report3

			dataTable2.Value = values
			dataTable2.Name = item
			dataTables2 = append(dataTables2, dataTable2)
		}
		dataTable.Seq = seq
		dataTable.Value = dataTables2
		dataTable.Name = plan
		dataTables = append(dataTables, dataTable)
		sort.Slice(dataTables, func(i, j int) bool {
			return dataTables[i].Seq < dataTables[j].Seq
		})
	}
	return dataTables, nil
}

func (u *UseCase) Report3YearData(id int) (interface{}, error) {
	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisPlan, err := u.Repo.GetMaintenanceAnalysisPlanByAnalysisID(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	modelResults, err := u.Repo.GetModelResultDataById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	yearMinData := []int{}
	for _, item := range modelResults {
		yearMinData = append(yearMinData, item.AnalystYear)
	}

	yeatMin := helpers.FindMin(yearMinData)

	// latestPrepareData, err := u.Repo.GetLatestPrepareDataByMaintenanceAnalysisId(id)
	// if err != nil {
	// 	logs.Error(err)
	// 	return models.DashboardStrategicMaintenance{}, err
	// }
	// return latestPrepareData, nil
	maintenanceAnalysisPlanMap := make(map[int]map[string]float64)
	for _, item := range maintenanceAnalysisPlan {
		maintenanceAnalysisPlanMapSub := map[string]float64{}
		for indexPlan := 1; indexPlan < *maintenanceAnalysis.NumberPlan+1; indexPlan++ {
			switch indexPlan {
			case 1:
				maintenanceAnalysisPlanMapSub["plan_1"] = *item.Plan1
			case 2:
				maintenanceAnalysisPlanMapSub["plan_2"] = *item.Plan2
			case 3:
				maintenanceAnalysisPlanMapSub["plan_3"] = *item.Plan3
			}
		}
		maintenanceAnalysisPlanMap[int(item.PlanYear)] = maintenanceAnalysisPlanMapSub
	}

	numberPlanDefault := 2
	numberPlanWithDefault := *maintenanceAnalysis.NumberPlan + numberPlanDefault
	numberPlan := *maintenanceAnalysis.NumberPlan
	year := *maintenanceAnalysis.Year
	analystYear := yeatMin //*&latestPrepareData.AnalystYear

	yearData := make(map[int]interface{})
	for indexYear := 1; indexYear < year+1; indexYear++ {
		var dataTables []responses.Report3
		indexPlan := 1
		for indexPlan = 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {

			var dataTable responses.Report3
			plan := ""
			seq := 0.0
			if indexPlan > numberPlan {
				if (indexPlan - numberPlan) == 2 {
					plan = "ซ่อมบำรุงปกติ"
					seq = 0
				} else {
					plan = "ไม่จำกัดงบประมาณ"
					seq = 0.1
				}
			} else {
				plan = "แผนที่ " + strconv.Itoa(indexPlan)
				seq = float64(indexPlan)
			}

			title := []string{"งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)", "B/C", "ประโยชน์ผู้ใช้ทาง (ล้านบาท)"}
			var dataTables2 []responses.Report3
			for _, item := range title {
				// for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
				values := []float64{}
				var dataTable2 responses.Report3

				switch item {
				case "งบประมาณซ่อมบำรุง (บาท)":
					sum := 0.0
					for _, item := range modelResults {
						if (indexYear) == item.Year && indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sum = sum + item.Data.UsedBudget
							}
						}
					}
					values = append(values, sum)
				case "IRI (ก่อนวิเคราะห์)":
					sumLength := 0.0
					sumIri := 0.0
					for _, item := range modelResults {
						if (indexYear) == item.Year && indexPlan == item.PlanSequence {
							sumIri = sumIri + (item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length))
							sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
						}
					}
					values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
				case "IRI (หลังวิเคราะห์)":
					sumLength := 0.0
					sumIri := 0.0
					for _, item := range modelResults {
						if (indexYear) == item.Year && indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sumIri = sumIri + (item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length))
								sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
							} else {
								sumIri = sumIri + (item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length))
								sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
							}
						}
					}
					values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
				case "B/C":
					sumLength := 0.0
					sumBC := 0.0
					for _, item := range modelResults {
						if (indexYear) == item.Year && indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sumBC += item.Data.OptimizationResult.BC * float64(item.Data.PrepareDataBefore.Length)
								sumLength += float64(item.Data.PrepareDataBefore.Length)
							} else {
								sumBC += item.Data.DeteriorationResult.Result.Bc * float64(item.Data.PrepareDataBefore.Length)
								sumLength += float64(item.Data.PrepareDataBefore.Length)
							}
						}
					}
					values = append(values, helpers.RoundFloat(sumBC/sumLength, 2))
				case "ประโยชน์ผู้ใช้ทาง (ล้านบาท)":
					sumCost := 0.0
					for _, item := range modelResults {
						if (indexYear) == item.Year && indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sumCost += item.Data.OptimizationResult.Benefit / 1000000
							} else {
								sumCost += 0
							}
						}
					}
					values = append(values, helpers.RoundFloat(sumCost, 2))
				}
				sumVal := 0.0
				for _, value := range values {
					sumVal += value
				}

				var val interface{}
				if helpers.CountDecimalPlaces(sumVal) > 0 {
					val = fmt.Sprintf("%.2f", sumVal)
				} else {
					val = strings.Split(fmt.Sprintf("%.2f", sumVal), ".")[0]
				}
				dataTable2.Value = helpers.AddCommasToNumber(val.(string))
				dataTable2.Name = item
				dataTables2 = append(dataTables2, dataTable2)
				// }
			}
			dataTable.Seq = seq
			dataTable.Value = dataTables2
			dataTable.Name = plan
			dataTables = append(dataTables, dataTable)
		}
		sort.Slice(dataTables, func(i, j int) bool {
			return dataTables[i].Seq < dataTables[j].Seq
		})
		yearData[analystYear+indexYear-1] = dataTables
		// var yearData responses.Report3Year
		// yearData.Year = dataTables
		// yearDatas = append(yearDatas, yearData)
	}

	return yearData, nil
}

func (u *UseCase) Report3DataAnnual(id int) (interface{}, error) {

	numberPlanDefault := 1
	numberPlanWithDefault := numberPlanDefault

	var dataTables []responses.Report3
	for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
		var dataTable responses.Report3
		plan := ""
		plan = "ซ่อมบำรุงปกติ"

		title := []string{"งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)", "B/C", "ประโยชน์ผู้ใช้ทาง (ล้านบาท)"}
		var dataTables2 []responses.Report3
		for _, item := range title {
			values := []float64{}
			var dataTable2 responses.Report3

			dataTable2.Value = values
			dataTable2.Name = item
			dataTables2 = append(dataTables2, dataTable2)
		}

		dataTable.Value = dataTables2
		dataTable.Name = plan
		dataTables = append(dataTables, dataTable)
	}
	return dataTables, nil
}

func (u *UseCase) Report3YearDataAnnual(id int) (interface{}, error) {
	modelResults, err := u.Repo.GetModelResultDataById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	yearMinData := []int{}
	for _, item := range modelResults {
		yearMinData = append(yearMinData, item.AnalystYear)
	}

	yeatMin := helpers.FindMin(yearMinData)
	// latestPrepareData, err := u.Repo.GetLatestPrepareDataByMaintenanceAnalysisId(id)
	// if err != nil {
	// 	logs.Error(err)
	// 	return models.DashboardStrategicMaintenance{}, err
	// }

	numberPlanDefault := 1
	numberPlanWithDefault := numberPlanDefault
	year := 0
	analystYear := yeatMin

	yearData := make(map[int]interface{})
	for indexYear := 0; indexYear < year+1; indexYear++ {
		var dataTables []responses.Report3
		for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
			var dataTable responses.Report3
			plan := ""

			plan = "ซ่อมบำรุงปกติ"

			title := []string{"งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)", "B/C", "ประโยชน์ผู้ใช้ทาง (ล้านบาท)"}
			var dataTables2 []responses.Report3
			for _, item := range title {
				values := []float64{}
				var dataTable2 responses.Report3

				switch item {
				case "งบประมาณซ่อมบำรุง (บาท)":
					sum := 0.0
					for _, item := range modelResults {
						if indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sum = sum + item.Data.UsedBudget
							}
						}
					}
					values = append(values, sum)
				case "IRI (ก่อนวิเคราะห์)":
					sumLength := 0.0
					sumIri := 0.0
					for _, item := range modelResults {
						if indexPlan == item.PlanSequence {
							sumIri = sumIri + (item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length))
							sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
						}
					}
					values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
				case "IRI (หลังวิเคราะห์)":
					sumLength := 0.0
					sumIri := 0.0
					for _, item := range modelResults {
						if indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sumIri = sumIri + (item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length))
								sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
							} else {
								sumIri = sumIri + (item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length))
								sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
							}
						}
					}
					values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
				case "B/C":
					sumLength := 0.0
					sumBC := 0.0
					for _, item := range modelResults {
						if indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sumBC += item.Data.OptimizationResult.BC * float64(item.Data.PrepareDataBefore.Length)
								sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
							} else {
								sumBC += item.Data.DeteriorationResult.Result.Bc * float64(item.Data.PrepareDataBefore.Length)
								sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
							}
						}
					}
					values = append(values, helpers.RoundFloat(sumBC/sumLength, 2))
				case "ประโยชน์ผู้ใช้ทาง (ล้านบาท)":
					sumCost := 0.0
					for _, item := range modelResults {
						if indexPlan == item.PlanSequence {
							if item.Data.Repair {
								sumCost += item.Data.OptimizationResult.Benefit / 1000000
							} else {
								sumCost += 0
							}
						}
					}
					values = append(values, helpers.RoundFloat(sumCost, 2))
				}
				sumVal := 0.0
				for _, value := range values {
					sumVal += value
				}

				var val interface{}
				if helpers.CountDecimalPlaces(sumVal) > 0 {
					val = fmt.Sprintf("%.2f", sumVal)
				} else {
					val = strings.Split(fmt.Sprintf("%.2f", sumVal), ".")[0]
				}
				dataTable2.Value = helpers.AddCommasToNumber(val.(string))
				dataTable2.Name = item
				dataTables2 = append(dataTables2, dataTable2)
			}

			dataTable.Value = dataTables2
			dataTable.Name = plan
			dataTables = append(dataTables, dataTable)
		}
		yearData[analystYear] = dataTables
	}
	return yearData, nil
}

func (u *UseCase) GetReportType4(id, userID, plan int, typeFile string) (interface{}, error) {
	data, _ := u.GetReport4Data(id, userID, plan)
	filePath := os.Getenv("MAINTENANCE_ANALYSIS_PDF") + fmt.Sprintf("%d/", id)
	templateName := os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE4")
	result, err := helpers.InitDataToHtmlFunc(templateName, data, filePath, "type4.html")
	if err != nil {
		logs.Error(err)
		return "", err
	}

	if typeFile == "html" {
		texts := strings.Split(result, "/")
		text := texts[len(texts)-1]

		new := strings.ReplaceAll(result, text, "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง.html")

		os.Rename(result, new)

		return os.Getenv("STORAGE_IP") + "/" + new, nil
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return "", err
	}

	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, helpers.PrintToPDF(string(html), &buf, false)); err != nil {
		logs.Error(err)
		return "", err
	}
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}
	fullFilePath := filePath + "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง" + ".pdf"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *UseCase) GetReport4Data(id, userID, plan int) (responses.Report4Data, error) {
	var roads []models.Road
	var roadGroup []models.RoadGroup
	var roadSection []models.RoadSection
	var analysis models.MaintenanceAnalysis
	g, gctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		var err error
		roads, err = u.Repo.GetRoadAll()
		return err
	})
	g.Go(func() error {
		var err error
		roadGroup, err = u.Repo.GetRoadGroups()
		return err
	})
	g.Go(func() error {
		var err error
		roadSection, err = u.Repo.GetRoadSections()
		return err
	})
	g.Go(func() error {
		var err error
		analysis, err = u.Repo.GetMaintenanceAnalysisById(id)
		return err
	})
	if err := g.Wait(); err != nil {
		logs.Error(err)
		return responses.Report4Data{}, responses.NewAppErr(400, err.Error())
	}
	_ = gctx

	roadCode := make(map[int]models.Road)
	for _, item := range roads {
		roadCode[item.Id] = item
	}
	roadGroupMap := map[int]models.RoadGroup{}
	for _, item := range roadGroup {
		roadGroupMap[item.Id] = item
	}
	roadSectionMap := map[int]models.RoadSection{}
	for _, item := range roadSection {
		roadSectionMap[item.Id] = item
	}
	analysisTypeID := analysis.MaintenanceAnalysisTypeId

	if analysisTypeID == 2 {
		if analysis.Condition == 2 || analysis.Condition == 3 {
			plan = 1
		} else {
			plan = 0
		}
	}

	data2, _ := u.Repo.GetReport4(id, plan)

	var reports []responses.Report4
	roadIndex := make(map[int]int)
	years := []int{}
	for _, item := range data2 {
		var report responses.Report4
		report.Year = item.AnalystYear
		report.RoadID = item.Data.PrepareDataBefore.Road.RoadID
		report.RoadCode = roadGroupMap[roadCode[item.Data.PrepareDataBefore.Road.RoadID].RoadGroupId].Number
		report.RoadName = roadSectionMap[roadCode[item.Data.PrepareDataBefore.Road.RoadID].RoadSectionId].Number
		report.RoadInfoName = item.Data.PrepareDataBefore.Road.RoadName
		report.KmStart = item.KmStart
		report.KmEnd = item.KmEnd
		report.KmTotal = (math.Abs(float64(item.KmStart) - float64(item.KmEnd))) / 1000
		report.LaneNo = item.Data.PrepareDataBefore.RoadGeom.LaneNo
		if item.Data.Repair {
			report.InterventionCriteria = item.Data.IcResult.Name
			report.IriBefore = item.Data.DeteriorationResult.Result.Iri
			report.IriAfter = item.Data.RweResult.Iri
			report.Acc = math.Round(item.Data.RucResult.RucAfterResult.Summary.ACCFINAL)
			report.Voc = math.Round(item.Data.RucResult.RucAfterResult.Summary.VOCFINAL)
			report.Vot = math.Round(item.Data.RucResult.RucAfterResult.Summary.VOTFINAL)
			report.Ruc = math.Round(item.Data.RucResult.RucAfterResult.Summary.TOTALRUCWITHACC)

			report.Benefit = item.Data.OptimizationResult.Benefit

			report.AccRm = item.Data.RucResult.RucBeforeResult.Summary.ACCFINAL
			report.VocRm = item.Data.RucResult.RucBeforeResult.Summary.VOCFINAL
			report.VotRm = item.Data.RucResult.RucBeforeResult.Summary.VOTFINAL

			report.RucRm = item.Data.RucResult.RucBeforeResult.Summary.TOTALRUCWITHACC

		} else {
			report.InterventionCriteria = "ซ่อมบำรุงปกติ"
			report.IriBefore = item.Data.PrepareDataBefore.RoadCondition.Iri
			report.IriAfter = item.Data.DeteriorationResult.Result.Iri
			report.Acc = 0
			report.Voc = 0
			report.Vot = 0
			report.Ruc = 0
			report.Benefit = 0
			report.AccRm = 0
			report.VocRm = 0
			report.VotRm = 0
			report.RucRm = 0
		}

		report.Area = float64(item.Data.PrepareDataBefore.Area)
		report.Budget = item.Data.UsedBudget
		report.BC = item.Data.OptimizationResult.BC
		report.Aadt = item.Data.PrepareDataBefore.Aadt

		roadIndex[item.RoadID] = item.Data.PrepareDataBefore.Road.RefDirectionID
		years = append(years, item.AnalystYear)
		reports = append(reports, report)
	}
	// index := 0
	var directionLT []responses.Report4Res
	var directionRT []responses.Report4Res
	for road, direction := range roadIndex {
		for _, item := range reports {
			var report responses.Report4Res

			if item.RoadID == road {

				report.Year = item.Year
				report.RoadID = item.RoadID
				report.RoadCode = item.RoadCode
				report.RoadName = item.RoadName
				report.RoadInfoName = item.RoadInfoName
				report.KmStart = item.KmStart
				report.KmEnd = item.KmEnd
				report.KmTotal = item.KmTotal
				report.LaneNo = item.LaneNo
				report.InterventionCriteria = item.InterventionCriteria
				report.Area = helpers.AddCommasToNumber(fmt.Sprintf("%.2f", item.Area))
				report.Budget = helpers.AddCommasToNumber(fmt.Sprintf("%.2f", item.Budget))
				report.BC = helpers.AddCommasToNumber(fmt.Sprintf("%.2f", item.BC))
				report.Aadt = helpers.AddCommasToNumber(fmt.Sprintf("%.2f", item.Aadt))
				report.IriBefore = helpers.AddCommasToNumber(fmt.Sprintf("%.2f", item.IriBefore))
				report.IriAfter = helpers.AddCommasToNumber(fmt.Sprintf("%.2f", item.IriAfter))
				report.Acc = helpers.FormatFloatNumber(math.Round(item.Acc))
				report.Voc = helpers.FormatFloatNumber(math.Round(item.Voc))
				report.Vot = helpers.FormatFloatNumber(math.Round(item.Vot))
				report.Ruc = helpers.FormatFloatNumber(math.Round(item.Ruc))

				report.Benefit = helpers.FormatFloatNumber(math.Round(item.Benefit))
				report.AccRm = helpers.FormatFloatNumber(math.Round(item.AccRm))
				report.VocRm = helpers.FormatFloatNumber(math.Round(item.VocRm))
				report.VotRm = helpers.FormatFloatNumber(math.Round(item.VotRm))
				report.RucRm = helpers.FormatFloatNumber(math.Round(item.RucRm))
				if direction == 1 {
					directionLT = append(directionLT, report)
				} else {
					directionRT = append(directionRT, report)
				}
			}
		}

		// if direction == 1 {
		// 	sort.Slice(directionLT, func(i, j int) bool {
		// 		return directionLT[i].KmStart < directionLT[j].KmEnd
		// 	})
		// } else {
		// 	sort.Slice(directionRT, func(i, j int) bool {
		// 		return directionRT[i].KmStart > directionRT[j].KmEnd
		// 	})

		// }

	}
	unique := make(map[int]bool)
	yearData := []int{}
	for _, num := range years {
		// Check if the element is already present in the map
		if !unique[num] {
			// If not present, add it to the map and result slice
			unique[num] = true
			yearData = append(yearData, num)
		}
	}

	data := append(directionLT, directionRT...)
	var dataRes []responses.Report4Res
	uniqueKMData := make(map[string]bool)
	for _, year := range years {
		for _, item := range data {
			if year == item.Year {
				if !uniqueKMData[fmt.Sprintf("%d%d%d%d%d", item.RoadID, item.LaneNo, item.KmStart, item.KmEnd, item.Year)] {
					uniqueKMData[fmt.Sprintf("%d%d%d%d%d", item.RoadID, item.LaneNo, item.KmStart, item.KmEnd, item.Year)] = true
					dataRes = append(dataRes, item)
				}
			}
		}
	}
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return responses.Report4Data{}, err
	}

	condition, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(analysisData.Condition)
	if err != nil {
		return responses.Report4Data{}, responses.NewAppErr(400, err.Error())
	}

	target, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*analysisData.Target)
	if err != nil {
		// return responses.Report4Data{}, responses.NewAppErr(400, err.Error())
	}
	user, _ := u.Repo.GetUserByID(uint(userID))

	// analystYear := prepareData.AnalystYear
	year := analysisData.Year

	formattedTime := helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour * 7))
	var report4Data responses.Report4Data
	report4Data.Report4 = dataRes
	planName := ""
	yearMinData := []int{}
	for _, item := range data2 {
		yearMinData = append(yearMinData, item.AnalystYear)
	}
	yeatMin := helpers.FindMin(yearMinData)
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		if plan == 0 {
			planName = "(แผนไม่จำกัดงบประมาณ)"
		} else {
			planName = fmt.Sprintf("แผนที่ %d", plan)
		}
		report4Data.Title = fmt.Sprintf("%s ปี พ.ศ. %d - %d", planName, yeatMin+543, yeatMin+*year-1+543)
	} else {
		report4Data.Title = fmt.Sprintf("ปี พ.ศ. %d", yeatMin+543)
	}

	report4Data.Date = formattedTime
	report4Data.User = user.Firstname + " " + user.Lastname
	report4Data.Condition = condition.Name
	report4Data.Target = target.Name
	report4Data.PlanName = planName
	return report4Data, nil
}

func (u *UseCase) Report1Excel(ID, userID int) (interface{}, error) {
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	interventionCriteria, _ := u.Repo.GetInterventionCriteriaParamsById(analysisData.InterventionCriteriaParmasID)

	// var jsonObject map[string]interface{}
	var jsonObject responses.InterventionCriteriaParams
	err = json.Unmarshal([]byte(interventionCriteria.Params), &jsonObject)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	// return interventionCriteriaParams, nil
	// return jsonObject, nil
	// var dataRes responses.AnalyzeReport1
	var interventionCriteriaData []responses.InterventionCriteriaReportRes
	// =============================== asphalt ===============================
	criteriaMethodAS, err := u.Repo.GetRefCriteriaMethodBySurface("as")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	// return criteriaMethodAS, nil
	asphalts := jsonObject.Asphalt
	asphaltData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range asphalts {
		asphaltData[item.Name] = item
	}
	for _, method := range criteriaMethodAS {
		val, isVal := asphaltData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}
		// return asphalt.(map[string]interface{})["maintenance_sequence"], nil
		interventionCriterias := asphaltData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReportRes
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes
			for _, item2 := range conditions {

				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}

			name := item.MaintenanceStandardName
			description := item.MaintenanceDescription
			data.Seq = item.MaintenanceSequence
			data.StandardName = name
			data.Description = description
			data.Conditions = conditionDatas
			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
		}
	}

	// =============================== concrete ===============================
	criteriaMethodCC, err := u.Repo.GetRefCriteriaMethodBySurface("cc")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	concrete := jsonObject.Concrete
	concreteData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range concrete {
		concreteData[item.Name] = item
	}

	for _, method := range criteriaMethodCC {
		val, isVal := concreteData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}

		// return yyy, nil
		interventionCriterias := concreteData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReportRes
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes
			for _, item2 := range conditions {
				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}
			name := item.MaintenanceStandardName
			description := item.MaintenanceDescription
			data.Seq = item.MaintenanceSequence
			data.StandardName = name
			data.Description = description
			data.Conditions = conditionDatas
			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
		}
	}

	sort.Slice(interventionCriteriaData, func(i, j int) bool {
		return interventionCriteriaData[i].Seq < interventionCriteriaData[j].Seq
	})
	var interventionCriteriaReportData responses.InterventionCriteriaReportData
	userUpdate, _ := u.Repo.GetUserByID(uint(interventionCriteria.CreatedBy))
	userPrint, _ := u.Repo.GetUserByID(uint(userID))
	interventionCriteriaReportData.UpdatedBy = userUpdate.Firstname + " " + userUpdate.Lastname
	interventionCriteriaReportData.UpdatedAt = helpers.ConvertToThaiFullCalendar(interventionCriteria.CreatedAt)
	interventionCriteriaReportData.User = userPrint.Firstname + " " + userPrint.Lastname
	interventionCriteriaReportData.PrintDate = helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour * 7))

	interventionCriteriaReportData.InterventionCriteria = interventionCriteriaData
	filePath := os.Getenv("ANALYSIS_EXCEL") + fmt.Sprintf("%d", ID)
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE1_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Font: &excelize.Font{
			Size:   16,
			Family: "TH SarabunPSK",
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
	if err != nil {
		log.Fatal(err)
	}

	textCenterHeader, err := f.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Font:      &excelize.Font{Bold: true, Color: "000000", Size: 16, Family: "TH SarabunPSK"},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#d9d9d9"},
				Pattern: 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	textLeftHeader, err := f.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
			Font:      &excelize.Font{Bold: true, Color: "000000", Size: 16, Family: "TH SarabunPSK"},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#e7e7e7"},
				Pattern: 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// return interventionCriteriaData, nil

	f.SetCellValue("Sheet1", "C3", "ผู้ประมวลผล : "+userPrint.Firstname+" "+userPrint.Lastname)
	f.SetCellValue("Sheet1", "C4", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7))+" น.")
	f.SetCellValue("Sheet1", "A5", "แก้ไขล่าสุดโดย : "+userUpdate.Firstname+" "+userUpdate.Lastname+" เมื่อ "+helpers.ConvertToThaiFullCalendar(interventionCriteria.CreatedAt))
	f.MergeCell("Sheet1", "A5", "H5")
	// f.SetCellValue("Sheet1", "D12", "แก้ไขล่าสุดโดย : ")
	// f.MergeCell("Sheet1", "D12", "O12")

	// if !data.IsNull {
	row := 6
	strRow := fmt.Sprint(row)
	for _, item := range interventionCriteriaData {
		f.SetCellValue("Sheet1", "A"+strRow, "มาตรฐานการซ่อมบํารุง "+item.StandardName)
		f.MergeCell("Sheet1", "A"+strRow, "H"+strRow)
		f.SetCellStyle("Sheet1", "A"+strRow, "H"+strRow, textCenterHeader)
		row++
		strRow = fmt.Sprint(row)
		f.SetCellValue("Sheet1", "A"+strRow, "คำอธิบาย : "+item.Description)
		f.MergeCell("Sheet1", "A"+strRow, "H"+strRow)
		f.SetCellStyle("Sheet1", "A"+strRow, "H"+strRow, textLeftHeader)
		row++
		strRow = fmt.Sprint(row)
		for _, condition := range item.Conditions {
			if condition.ConditionLink == "" {
				f.SetCellValue("Sheet1", "A"+strRow, "เงื่อนไข")
			} else {
				f.SetCellValue("Sheet1", "A"+strRow, condition.ConditionLink)
			}
			f.MergeCell("Sheet1", "A"+strRow, "B"+strRow)
			f.SetCellStyle("Sheet1", "A"+strRow, "B"+strRow, style)

			f.SetCellValue("Sheet1", "C"+strRow, condition.ConditionValue1)
			f.MergeCell("Sheet1", "C"+strRow, "D"+strRow)
			f.SetCellStyle("Sheet1", "C"+strRow, "D"+strRow, style)

			f.SetCellValue("Sheet1", "E"+strRow, condition.ConditionOperation1+" "+condition.ConditionCriterion+" "+condition.ConditionOperation2)
			f.MergeCell("Sheet1", "E"+strRow, "F"+strRow)
			f.SetCellStyle("Sheet1", "E"+strRow, "F"+strRow, style)

			f.SetCellValue("Sheet1", "G"+strRow, condition.ConditionValue2)
			f.MergeCell("Sheet1", "G"+strRow, "H"+strRow)
			f.SetCellStyle("Sheet1", "G"+strRow, "H"+strRow, style)
			row++
			strRow = fmt.Sprint(row)
		}
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	name := "เงื่อนไขการซ่อมบำรุง"
	f.SaveAs(filePath + "/" + name + ".xlsx")
	return os.Getenv("STORAGE_IP") + "/" + filePath + "/" + name + ".xlsx", nil
}

func (u *UseCase) Report2Excel(ID, userID int) (interface{}, error) {
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	interventionCriteria, _ := u.Repo.GetInterventionCriteriaParamsById(analysisData.InterventionCriteriaParmasID)

	var jsonObject responses.InterventionCriteriaParams
	err = json.Unmarshal([]byte(interventionCriteria.Params), &jsonObject)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	var interventionCriteriaData []responses.InterventionCriteriaReport2Res
	// =============================== asphalt ===============================
	criteriaMethodAS, err := u.Repo.GetRefCriteriaMethodBySurface("as")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	asphalts := jsonObject.Asphalt
	asphaltData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range asphalts {
		asphaltData[item.Name] = item
	}
	i := 1
	j := 1
	for _, method := range criteriaMethodAS {

		val, isVal := asphaltData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}

		interventionCriterias := asphaltData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReport2Res
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes

			for _, item2 := range conditions {
				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}

			name := item.MaintenanceStandardName
			costPerUnit := item.MaintenanceCostPerUnit
			data.No = j
			data.Seq = float64(item.MaintenanceSequence)
			data.StandardName = name
			data.MaintenanceCostPerUnit = helpers.FormatNumberFloat(costPerUnit)

			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
			j++
		}
		i++
	}

	// =============================== concrete ===============================
	criteriaMethodCC, err := u.Repo.GetRefCriteriaMethodBySurface("cc")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	concrete := jsonObject.Concrete
	concreteData := make(map[string]responses.InterventionCriteriaSurfaceParams)
	for _, item := range concrete {
		concreteData[item.Name] = item
	}
	for _, method := range criteriaMethodCC {
		val, isVal := concreteData[method.Name]
		if !isVal {
			continue
		}

		if val.Name == "" {
			continue
		}

		interventionCriterias := concreteData[method.Name].InterventionCriteriaMaintenances
		for _, item := range interventionCriterias {
			var data responses.InterventionCriteriaReport2Res
			conditions := item.MaintenanceCondition
			var conditionDatas []responses.InterventionCriteriaConditionsReportRes

			for _, item2 := range conditions {
				var condition responses.InterventionCriteriaConditionsReportRes
				condition.ConditionCriterion = item2.ConditionCriterion
				condition.ConditionLink = item2.ConditionLink
				condition.ConditionOperation1 = item2.ConditionOperation1
				condition.ConditionOperation2 = item2.ConditionOperation2
				condition.ConditionValue1 = helpers.FormatNumberFloat(item2.ConditionValue1)
				condition.ConditionValue2 = helpers.FormatNumberFloat(item2.ConditionValue2)
				// return item2.(map[string]interface{})["condition_criterion"], nil
				// return item2, nil
				conditionDatas = append(conditionDatas, condition)
			}

			name := item.MaintenanceStandardName
			costPerUnit := item.MaintenanceCostPerUnit
			data.No = j
			data.Seq = float64(item.MaintenanceSequence)
			data.StandardName = name
			data.MaintenanceCostPerUnit = helpers.FormatNumberFloat(costPerUnit)

			if len(conditionDatas) == 0 {
				continue
			}
			interventionCriteriaData = append(interventionCriteriaData, data)
			j++
		}
	}

	sort.Slice(interventionCriteriaData, func(i, j int) bool {
		return interventionCriteriaData[i].Seq < interventionCriteriaData[j].Seq
	})
	var interventionCriteriaDatas []responses.InterventionCriteriaReport2Res
	for index, item := range interventionCriteriaData {
		var interventionCriteriaData responses.InterventionCriteriaReport2Res
		interventionCriteriaData.No = index + 1
		interventionCriteriaData.MaintenanceCostPerUnit = item.MaintenanceCostPerUnit
		interventionCriteriaData.StandardName = item.StandardName
		interventionCriteriaDatas = append(interventionCriteriaDatas, interventionCriteriaData)
	}

	sort.Slice(interventionCriteriaData, func(i, j int) bool {
		return interventionCriteriaData[i].Seq < interventionCriteriaData[j].Seq
	})
	var interventionCriteriaRes responses.InterventionCriteriaReport2Data
	interventionCriteriaRes.UpdatedBy = ""
	// interventionCriteriaRes.UpdatedAt

	layout := "2006-01-02 15:04:05.999999 -0700 -0700"
	// Parse the input time with the given layout
	t, err := time.Parse(layout, time.Now().String())
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}

	// Set the time zone to Asia/Bangkok (Thailand)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}
	t = t.In(loc)

	// Format the time as desired
	formattedTime := t.Format("2006-01-02 15:04:05")
	userUpdate, _ := u.Repo.GetUserByID(uint(interventionCriteria.CreatedBy))
	userPrint, _ := u.Repo.GetUserByID(uint(userID))
	interventionCriteriaRes.UpdatedBy = userUpdate.Firstname + " " + userUpdate.Lastname
	interventionCriteriaRes.UpdatedAt = "2023-09-19 10:56:28"
	interventionCriteriaRes.User = userPrint.Firstname + " " + userPrint.Lastname
	interventionCriteriaRes.PrintDate = formattedTime
	interventionCriteriaRes.InterventionCriteria = interventionCriteriaData

	filePath := os.Getenv("ANALYSIS_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE2_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Size:   16,
			Family: "TH SarabunPSK",
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
	if err != nil {
		log.Fatal(err)
	}

	textCenterHeader, err := f.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Font:      &excelize.Font{Bold: true, Color: "000000", Size: 16, Family: "TH SarabunPSK"},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#d9d9d9"},
				Pattern: 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// return interventionCriteriaData, nil

	f.SetCellValue("Sheet1", "C3", "ผู้ประมวลผล : "+userPrint.Firstname+" "+userPrint.Lastname)
	f.SetCellValue("Sheet1", "C4", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7))+" น.")
	f.SetCellValue("Sheet1", "F5", "แก้ไขล่าสุดโดย : "+userUpdate.Firstname+" "+userUpdate.Lastname+" เมื่อ "+helpers.ConvertToThaiFullCalendar(interventionCriteria.CreatedAt))
	// f.SetCellValue("Sheet1", "D12", "แก้ไขล่าสุดโดย : "
	// f.MergeCell("Sheet1", "D12", "O12")

	f.SetCellStyle("Sheet1", "B6", "B6", textCenterHeader)

	f.SetCellStyle("Sheet1", "C6", "C6", textCenterHeader)

	f.SetCellStyle("Sheet1", "G6", "G6", textCenterHeader)

	row := 7
	strRow := fmt.Sprint(row)
	for index, item := range interventionCriteriaData {
		f.SetCellValue("Sheet1", "B"+strRow, index+1)
		f.SetCellStyle("Sheet1", "B"+strRow, "B"+strRow, style)

		f.SetCellValue("Sheet1", "C"+strRow, item.StandardName)
		f.MergeCell("Sheet1", "C"+strRow, "F"+strRow)
		f.SetCellStyle("Sheet1", "C"+strRow, "F"+strRow, style)

		f.SetCellValue("Sheet1", "G"+strRow, item.MaintenanceCostPerUnit)
		f.MergeCell("Sheet1", "G"+strRow, "J"+strRow)
		f.SetCellStyle("Sheet1", "G"+strRow, "J"+strRow, style)
		row++
		strRow = fmt.Sprint(row)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	name := "เงื่อนไขค่าใช้จ่ายการซ่อมบำรุง"
	f.SaveAs(filePath + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func (u *UseCase) GetReportType3Excel(id, userID int) (interface{}, error) {
	analysisData, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	condition, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(analysisData.Condition)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	target, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*analysisData.Target)
	if err != nil {
		logs.Error(err)
		// return "", responses.NewAppErr(400, err.Error())
	}
	var report3Year interface{}
	var report3 interface{}
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		report3Year, _ = u.Report3YearData(id)
		report3, _ = u.Report3Data(id)
	} else {
		report3Year, _ = u.Report3YearDataAnnual(id)
		report3, _ = u.Report3DataAnnual(id)
	}

	// analystYear := prepareData.AnalystYear
	year := analysisData.Year
	currentTime := time.Now()
	formattedTime := currentTime.Format("02/01/2006 15:04")
	user, _ := u.Repo.GetUserByID(uint(userID))
	var data responses.Report3Data

	data.Report3 = report3.([]responses.Report3)
	data.Year = report3Year
	data.Date = formattedTime
	data.User = user.Firstname + " " + user.Lastname
	data.Condition = condition.Name
	data.Target = target.Name

	filePath := os.Getenv("ANALYSIS_EXCEL") + fmt.Sprintf("%d", id) + "/"
	template := ""
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		if *analysisData.NumberPlan == 1 {
			template = os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE3_1_EXCEL")
		} else if *analysisData.NumberPlan == 2 {
			template = os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE3_2_EXCEL")
		} else if *analysisData.NumberPlan == 3 {
			template = os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE3_3_EXCEL")
		} else {
			template = os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE3_EXCEL")
		}
	} else {
		template = os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE3_EXCEL")
	}

	f, err := excelize.OpenFile(template)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logs.Error(err)
		}
	}()

	modelResults, err := u.Repo.GetModelResultDataById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	yearMinData := []int{}
	for _, item := range modelResults {
		yearMinData = append(yearMinData, item.AnalystYear)
	}

	yeatMin := helpers.FindMin(yearMinData)

	if analysisData.MaintenanceAnalysisTypeId == 1 {
		f.SetCellValue("Sheet1", "A1", "สรุปค่าซ่อมบำรุงและค่า IRI ปี พ.ศ. "+fmt.Sprintf("%d - %d", yeatMin+543, yeatMin+*year-1+543))
	} else {
		f.SetCellValue("Sheet1", "A1", "สรุปค่าซ่อมบำรุงและค่า IRI ปี พ.ศ. "+fmt.Sprintf("%d ", yeatMin+543))
	}
	f.SetCellValue("Sheet1", "A2", "เงื่อนไข : "+condition.Name+" เป้าหมาย : "+target.Name)
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		if *analysisData.NumberPlan == 1 {
			f.SetCellValue("Sheet1", "C4", "ผู้ประมวลผล : "+user.Firstname+" "+user.Lastname)
			f.SetCellValue("Sheet1", "C5", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7)))
		} else if *analysisData.NumberPlan == 2 {
			f.SetCellValue("Sheet1", "C4", "ผู้ประมวลผล : "+user.Firstname+" "+user.Lastname)
			f.SetCellValue("Sheet1", "C5", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7)))
		} else if *analysisData.NumberPlan == 3 {
			f.SetCellValue("Sheet1", "D4", "ผู้ประมวลผล : "+user.Firstname+" "+user.Lastname)
			f.SetCellValue("Sheet1", "D5", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7)))
		} else {
			f.SetCellValue("Sheet1", "B4", "ผู้ประมวลผล : "+user.Firstname+" "+user.Lastname)
			f.SetCellValue("Sheet1", "B5", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7)))
		}
	} else {
		f.SetCellValue("Sheet1", "B4", "ผู้ประมวลผล : "+user.Firstname+" "+user.Lastname)
		f.SetCellValue("Sheet1", "B5", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7)))
	}

	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Size:   12,
			Family: "TH SarabunPSK",
			Bold:   true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#d9d9d9"},
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
	if err != nil {
		logs.Error(err)
	}

	styleBody, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Size:   12,
			Family: "TH SarabunPSK",
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
	if err != nil {
		logs.Error(err)
	}

	planColumn := [][]string{
		{"B", "C", "D", "E", "F"},
		{"G", "H", "I", "J", "K"},
		{"L", "M", "N", "O", "P"},
		{"Q", "R", "S", "T", "U"},
		{"V", "W", "X", "Y", "Z"},
	}

	f.SetCellStyle("Sheet1", "A7", "A7", style)

	for index, report := range report3.([]responses.Report3) {
		f.SetCellValue("Sheet1", planColumn[index][0]+"7", report.Name)
		f.MergeCell("Sheet1", planColumn[index][0]+"7", planColumn[index][4]+"7")
		f.SetCellStyle("Sheet1", planColumn[index][0]+"7", planColumn[index][4]+"7", style)
		for indexCol, col := range report.Value.([]responses.Report3) {
			f.SetCellValue("Sheet1", planColumn[index][indexCol]+"8", col.Name)
			f.SetCellStyle("Sheet1", planColumn[index][indexCol]+"8", planColumn[index][indexCol]+"8", style)
		}
	}
	row := 9
	strRow := fmt.Sprint(row)
	index := 0
	for year, item := range report3Year.(map[int]interface{}) {
		f.SetCellValue("Sheet1", "A"+strRow, year+543)
		f.SetCellStyle("Sheet1", "A"+strRow, "A"+strRow, styleBody)
		for col, item2 := range item.([]responses.Report3) {
			for col3, item3 := range item2.Value.([]responses.Report3) {
				f.SetCellValue("Sheet1", planColumn[col][col3]+strRow, item3.Value)
				f.SetCellStyle("Sheet1", planColumn[col][col3]+strRow, planColumn[col][col3]+strRow, styleBody)
			}
		}
		index++
		row++
		strRow = fmt.Sprint(row)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	name := "สรุปค่าซ่อมบำรุงและค่าIRIของแต่ละปี"
	f.SaveAs(filePath + "/" + name + ".xlsx")
	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

}

func (u *UseCase) Report4Excel(id, userID int, plan int) (interface{}, error) {
	data, _ := u.GetReport4Data(id, userID, plan)
	// return data, nil
	filePath := os.Getenv("ANALYSIS_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE4_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   9,
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
	if err != nil {
		log.Fatal(err)
	}

	textCenterHeader, err := f.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Font:      &excelize.Font{Bold: true, Color: "000000", Size: 9, Family: "TH SarabunPSK"},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#d9d9d9"},
				Pattern: 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	analysisData, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	condition, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(analysisData.Condition)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	target, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*analysisData.Target)
	if err != nil {
		// return "", responses.NewAppErr(400, err.Error())
	}
	user, _ := u.Repo.GetUserByID(uint(userID))

	sort.SliceStable(data.Report4, func(i, j int) bool {
		return data.Report4[i].Year < data.Report4[j].Year
	})

	planName := ""
	title := ""
	if analysisData.MaintenanceAnalysisTypeId == 1 {
		if plan == 0 {
			planName = "(แผนไม่จำกัดงบประมาณ)"
		} else {
			planName = fmt.Sprintf("แผนที่ %d", plan)
		}

		title = fmt.Sprintf(" ปี พ.ศ. %d - %d", data.Report4[0].Year+543, data.Report4[len(data.Report4)-1].Year+543)
	} else {
		if len(data.Report4) != 0 {
			title = fmt.Sprintf(" ปี พ.ศ. %d", data.Report4[0].Year+543)
		}
	}

	if analysisData.MaintenanceAnalysisTypeId == 1 {
		f.SetCellValue("Sheet1", "D1", "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง "+planName+" "+title)
	} else {
		f.SetCellValue("Sheet1", "D1", "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง "+title)
	}

	f.SetCellValue("Sheet1", "D2", "เงื่อนไข : "+condition.Name+" เป้าหมาย : "+target.Name)
	f.SetCellValue("Sheet1", "C4", "ผู้ประมวลผล : "+user.Firstname+" "+user.Lastname)
	f.SetCellValue("Sheet1", "C5", "วันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysisData.CreatedAt.Add(time.Hour*7))+"น.")
	// return interventionCriteriaData, nil

	f.SetCellStyle("Sheet1", "A7", "Y7", textCenterHeader)

	// f.SetCellValue("Sheet1", "E6", "ผู้ประมวลผล : "+userPrint.Firstname+" "+userPrint.Lastname)
	// f.SetCellValue("Sheet1", "E7", "วันที่วิเคราะห์ : "+formattedTime)
	// f.SetCellValue("Sheet1", "M10", "แก้ไขล่าสุดโดย : "+userUpdate.Firstname+" "+userUpdate.Lastname)
	// f.SetCellValue("Sheet1", "D12", "แก้ไขล่าสุดโดย : ")
	// f.MergeCell("Sheet1", "D12", "O12")

	// if !data.IsNull {
	row := 8
	strRow := fmt.Sprint(row)
	for index, item := range data.Report4 {
		f.SetCellValue("Sheet1", "A"+strRow, index+1)
		f.SetCellValue("Sheet1", "B"+strRow, item.Year+543)
		f.SetCellValue("Sheet1", "C"+strRow, item.RoadCode)
		f.SetCellValue("Sheet1", "D"+strRow, item.RoadName)
		f.SetCellValue("Sheet1", "E"+strRow, item.RoadInfoName)
		f.SetCellValue("Sheet1", "F"+strRow, helpers.FormatKM(int64(item.KmStart)))
		f.SetCellValue("Sheet1", "G"+strRow, helpers.FormatKM(int64(item.KmEnd)))
		f.SetCellValue("Sheet1", "H"+strRow, item.KmTotal)
		f.SetCellValue("Sheet1", "I"+strRow, item.LaneNo)
		f.SetCellValue("Sheet1", "J"+strRow, item.InterventionCriteria)
		f.SetCellValue("Sheet1", "K"+strRow, item.Area)
		f.SetCellValue("Sheet1", "L"+strRow, item.Budget)
		f.SetCellValue("Sheet1", "M"+strRow, item.BC)
		f.SetCellValue("Sheet1", "N"+strRow, item.Aadt)
		f.SetCellValue("Sheet1", "O"+strRow, item.IriBefore)
		f.SetCellValue("Sheet1", "P"+strRow, item.IriAfter)
		f.SetCellValue("Sheet1", "Q"+strRow, item.Acc)
		f.SetCellValue("Sheet1", "R"+strRow, item.Voc)
		f.SetCellValue("Sheet1", "S"+strRow, item.Vot)
		f.SetCellValue("Sheet1", "T"+strRow, item.Ruc)
		f.SetCellValue("Sheet1", "U"+strRow, item.Benefit)
		f.SetCellValue("Sheet1", "V"+strRow, item.AccRm)
		f.SetCellValue("Sheet1", "W"+strRow, item.VocRm)
		f.SetCellValue("Sheet1", "X"+strRow, item.VotRm)
		f.SetCellValue("Sheet1", "Y"+strRow, item.RucRm)
		f.SetCellStyle("Sheet1", "A"+strRow, "A"+strRow, style)
		f.SetCellStyle("Sheet1", "B"+strRow, "B"+strRow, style)
		f.SetCellStyle("Sheet1", "C"+strRow, "C"+strRow, style)
		f.SetCellStyle("Sheet1", "D"+strRow, "D"+strRow, style)
		f.SetCellStyle("Sheet1", "E"+strRow, "E"+strRow, style)
		f.SetCellStyle("Sheet1", "F"+strRow, "F"+strRow, style)
		f.SetCellStyle("Sheet1", "G"+strRow, "G"+strRow, style)
		f.SetCellStyle("Sheet1", "H"+strRow, "H"+strRow, style)
		f.SetCellStyle("Sheet1", "I"+strRow, "I"+strRow, style)
		f.SetCellStyle("Sheet1", "J"+strRow, "J"+strRow, style)
		f.SetCellStyle("Sheet1", "K"+strRow, "K"+strRow, style)
		f.SetCellStyle("Sheet1", "L"+strRow, "L"+strRow, style)
		f.SetCellStyle("Sheet1", "M"+strRow, "M"+strRow, style)
		f.SetCellStyle("Sheet1", "N"+strRow, "N"+strRow, style)
		f.SetCellStyle("Sheet1", "O"+strRow, "O"+strRow, style)
		f.SetCellStyle("Sheet1", "P"+strRow, "P"+strRow, style)
		f.SetCellStyle("Sheet1", "Q"+strRow, "Q"+strRow, style)
		f.SetCellStyle("Sheet1", "R"+strRow, "R"+strRow, style)
		f.SetCellStyle("Sheet1", "S"+strRow, "S"+strRow, style)
		f.SetCellStyle("Sheet1", "T"+strRow, "T"+strRow, style)
		f.SetCellStyle("Sheet1", "U"+strRow, "U"+strRow, style)
		f.SetCellStyle("Sheet1", "V"+strRow, "V"+strRow, style)
		f.SetCellStyle("Sheet1", "W"+strRow, "W"+strRow, style)
		f.SetCellStyle("Sheet1", "X"+strRow, "X"+strRow, style)
		f.SetCellStyle("Sheet1", "Y"+strRow, "Y"+strRow, style)
		row++
		strRow = fmt.Sprint(row)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	name := "รายละเอียดแผนงานซ่อมบำรุงตามสายทาง"
	f.SaveAs(filePath + name + ".xlsx")
	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

}

type Cache struct {
	mu    sync.RWMutex
	items map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	return item, found
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func contains(list []responses.KMReport, x responses.KMReport) bool {
	for i := range list {
		if eq(x, list[i]) {
			return true
		}
	}
	return false
}

func eq(a, b responses.KMReport) bool {
	return a.KmStart == b.KmStart && a.KmEnd == b.KmEnd && a.LaneNo == b.LaneNo && a.Year == b.Year
	//a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Peek() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}
