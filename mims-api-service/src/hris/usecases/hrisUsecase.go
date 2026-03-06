package usecases

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/hris/domains"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type Usecase struct {
	Repository domains.Repository
}

func NewUsecase(repo domains.Repository) domains.UseCase {
	return &Usecase{Repository: repo}
}

///////////////// SectionGeom /////////////////

func (u *Usecase) GetSectionGeom() error {

	run := 1
	limit := 10
	var err error

	for run <= limit {
		data, errLogic := u.LogicSectionGeom()
		if errLogic == nil {
			message := "บันทึกข้อมูลลง MongoDb จำนวน " + strconv.Itoa(len(data)) + " ข้อมูล"
			errMongoDbLog := helpers.MongoDbLog("Section Geom", message, "MONGODB_HRIS_ERROR_LOG", true)
			if errMongoDbLog != nil {
				return errMongoDbLog
			}

			return nil
		} else {
			message := errLogic.Error()
			errMongoDbLog := helpers.MongoDbLog("Section Geom", message, "MONGODB_HRIS_ERROR_LOG", false)
			if errMongoDbLog != nil {
				return errMongoDbLog
			}

			err = errLogic
		}

		run++
		time.Sleep(3 * time.Second)
	}

	return err
}

func (u *Usecase) LogicSectionGeom() ([]interface{}, error) {
	var GetSectionGeom models.SectionGeom

	url := os.Getenv("HRIS2_SECTION_GEOM_URL")
	method := "GET"

	response, err := u.Request(url, method)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(response, &GetSectionGeom)
	if err != nil {
		return nil, err
	}

	var insertData []interface{}
	for _, item := range GetSectionGeom.Item {
		item.IsLatested = true
		insertData = append(insertData, item)
	}

	err = u.Repository.InsertSectionGeom(insertData)
	if err != nil {
		return nil, err
	}

	return insertData, nil
}

///////////////// RoadLatest /////////////////

func (u *Usecase) GetRoadLatest() error {

	run := 1
	limit := 10
	var err error

	for run <= limit {
		data, errLogic := u.LogicRoadLatest()
		if errLogic == nil {
			message := "บันทึกข้อมูลลง MongoDb จำนวน " + strconv.Itoa(len(data)) + " ข้อมูล"
			errMongoDbLog := helpers.MongoDbLog("Road Latest", message, "MONGODB_HRIS_ERROR_LOG", true)
			if errMongoDbLog != nil {
				return errMongoDbLog
			}

			return nil
		} else {

			message := errLogic.Error()
			errMongoDbLog := helpers.MongoDbLog("Road Latest", message, "MONGODB_HRIS_ERROR_LOG", false)
			if errMongoDbLog != nil {
				return errMongoDbLog
			}

			err = errLogic
		}

		run++
		time.Sleep(3 * time.Second)
	}

	return err
}

func (u *Usecase) LogicRoadLatest() ([]interface{}, error) {
	var roadLatest []models.RoadLatest

	url := os.Getenv("HRIS2_ROAD_LATEST_URL") + "?date=2014-08-31"
	method := "GET"

	response, err := u.Request(url, method)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &roadLatest)
	if err != nil {
		return nil, err
	}

	var insertData []interface{}
	for _, item := range roadLatest {
		item.IsLatested = true
		if item.Status == "A" {
			insertData = append(insertData, item)
		}
	}

	err = u.Repository.InsertRoadLatest(insertData)
	if err != nil {
		return nil, err
	}

	return insertData, nil
}

///////////////// MatchData /////////////////

func (u *Usecase) MatchData() (interface{}, error) {

	tx := u.Repository.StartTransSection()

	refHirs, err := u.Repository.GetRefHirs(tx)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	sectionGeoms, err := u.SectionGeom(refHirs)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	roadLatest, err := u.RoadLatest()
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	roadGroups, err := u.Repository.GetRoadGroup(tx)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	roadSections, err := u.Repository.GetRoadSection(tx)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	dataModlue1, err := u.MatchDataModule1(sectionGeoms, roadGroups, roadLatest, tx)
	if err != nil {

		err = helpers.MongoDbLog("Road Group (Match Data)", err.Error(), "MONGODB_HRIS_ERROR_LOG", false)
		if err != nil {
			u.Repository.RollBack(tx)
			return nil, err
		}

		u.Repository.RollBack(tx)
		return nil, err
	}

	roadGroups, err = u.Repository.GetRoadGroup(tx)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	dataModlue2, err := u.MatchDataModule2(sectionGeoms, roadGroups, roadSections, tx)
	if err != nil {

		err = helpers.MongoDbLog("Road Section (Match Data)", err.Error(), "MONGODB_HRIS_ERROR_LOG", false)
		if err != nil {
			u.Repository.RollBack(tx)
			return nil, err
		}

		u.Repository.RollBack(tx)
		return nil, err
	}

	u.Repository.Commit(tx)

	err = helpers.MongoDbLog("Road Group (Match Data)", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(dataModlue1))+" ข้อมูล", "MONGODB_HRIS_ERROR_LOG", true)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	err = helpers.MongoDbLog("Road Section (Match Data)", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(dataModlue2))+" ข้อมูล", "MONGODB_HRIS_ERROR_LOG", true)
	if err != nil {
		u.Repository.RollBack(tx)
		return nil, err
	}

	return nil, nil
}

///////////////// OTHER /////////////////

func (u *Usecase) Request(url, method string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

func (u *Usecase) SectionGeom(refHirs []models.RefHris) ([]models.Item, error) {
	var level3 bson.A
	for _, item := range refHirs {
		var level1 bson.A

		if item.OfficeOfHighwaysCode != "" {
			level1 = append(level1, bson.M{"office_of_highways_code": item.OfficeOfHighwaysCode})
		}

		if item.RoadNumber != "" {
			level1 = append(level1, bson.M{"road_number": item.RoadNumber})
		}

		if item.SectionRoadNumber != "" {
			level1 = append(level1, bson.M{"section_road_number": item.SectionRoadNumber})
		}

		if len(level1) > 0 {
			level2 := bson.D{
				{"$and", level1}}

			level3 = append(level3, level2)

		}
	}

	filter := bson.D{
		{"$or", level3},
	}

	sectionGeom, err := u.Repository.GetSectionGeomWithFilter(filter)
	if err != nil {
		return nil, err
	}

	return sectionGeom, nil
}

func (u *Usecase) MatchDataModule1(sectionGeoms []models.Item, roadGroups []models.RoadGroup, roadLatest []models.RoadLatest, xt *gorm.DB) ([]models.RoadLatest, error) {

	isDuplicate := map[string]bool{}
	sectionGeomRoadCodes := []string{}
	for _, item := range sectionGeoms {
		if !isDuplicate[item.RoadNumber] {
			sectionGeomRoadCodes = append(sectionGeomRoadCodes, item.RoadNumber)
			isDuplicate[item.RoadNumber] = true
		}
	}

	roadGroupMap := map[string]models.RoadGroup{}
	for _, item := range roadGroups {
		roadGroupMap[item.Number] = item
	}

	roadLatestMap := map[string]models.RoadLatest{}
	for _, item := range roadLatest {
		roadLatestMap[item.RoadCode] = item
	}

	var roadLatests []models.RoadLatest

	for _, item := range sectionGeomRoadCodes {

		roadName := ""
		roadNames := strings.Split(roadLatestMap[item].RoadName, " ")
		if len(roadNames) >= 3 {
			roadName = roadNames[0] + " " + roadNames[1] + " " + roadNames[2]
		} else {
			roadName = roadNames[0]
		}

		kmStart, err := u.KmToM(roadLatestMap[item].KmStart)
		if err != nil {
			kmStart = 0
		}

		kmEnd, err := u.KmToM(roadLatestMap[item].KmEnd)
		if err != nil {
			kmEnd = 0
		}

		length, err := strconv.ParseFloat(roadLatestMap[item].Length, 32)
		if err != nil {
			length = 0
		}

		if (roadGroupMap[item] == models.RoadGroup{}) {
			var insert models.InsertRoadGroup
			insert.Number = roadLatestMap[item].RoadCode
			insert.Name = roadLatestMap[item].RoadName
			insert.ShortName = roadName
			insert.KmStart = float32(kmStart)
			insert.KmEnd = float32(kmEnd)
			insert.Distance = float32(length)
			err = u.Repository.InsertRoadGroup(insert, xt)
			if err != nil {
				return []models.RoadLatest{}, err
			}

		} else {
			if roadGroupMap[item].Name != roadLatestMap[item].RoadName {
				err := u.Repository.UpdateRoadGroupNameByNumber(roadLatestMap[item].RoadCode, roadLatestMap[item].RoadName, roadName, xt)
				if err != nil {
					return []models.RoadLatest{}, err
				}
			}
		}

		roadLatests = append(roadLatests, roadLatestMap[item])

	}

	return roadLatests, nil
}

func (u *Usecase) MatchDataModule2(sectionGeoms []models.Item, roadGroups []models.RoadGroup, roadSections []models.RoadSection, xt *gorm.DB) ([]models.Item, error) {

	roadSectionsMap := map[string]models.RoadSection{}
	for _, item := range roadSections {
		roadSectionsMap[strconv.Itoa(item.RoadGroupId)+"-"+item.Number] = item
	}

	roadGroupNumberMap := map[string]int{}
	for _, item := range roadGroups {
		roadGroupNumberMap[item.Number] = item.Id
	}

	var items []models.Item

	for _, item := range sectionGeoms {
		nameOriginTH := ""
		nameDestinationTH := ""
		nameOriginEn := ""
		nameDestinationEn := ""

		nameTh := strings.Split(item.SectionRoadThName, " - ")
		if len(nameTh) == 1 {
			nameOriginTH = item.SectionRoadThName
			nameDestinationTH = item.SectionRoadThName
		} else {
			nameOriginTH = nameTh[0]
			nameDestinationTH = nameTh[1]
		}

		nameEn := strings.Split(item.SectionRoadEngName, " - ")
		if len(nameEn) == 1 {
			nameOriginEn = item.SectionRoadEngName
			nameDestinationEn = item.SectionRoadEngName
		} else {
			nameOriginEn = nameEn[0]
			nameDestinationEn = nameEn[1]
		}

		provinceCord, err := strconv.Atoi(item.ProvinceCord)
		if err != nil {
			provinceCord = 0
		}

		kmEnd, err := u.KmToM(item.KmEnd)
		if err != nil {
			kmEnd = 0
		}

		kmStart, err := u.KmToM(item.KmBegin)
		if err != nil {
			kmStart = 0
		}

		length, err := strconv.ParseFloat(item.Distance, 32)
		if err != nil {
			length = 0
		}

		if roadSectionsMap[strconv.Itoa(roadGroupNumberMap[item.RoadNumber])+"-"+item.SectionRoadNumber].Id == 0 {
			var insert models.InsertRoadSection

			sectionPartID, err := strconv.Atoi(item.SectionPartID)
			if err != nil {
				return []models.Item{}, err
			}

			insert.Id = sectionPartID
			insert.NameOriginTH = nameOriginTH
			insert.NameDestinationTH = nameDestinationTH
			insert.NameOriginEn = nameOriginEn
			insert.NameDestinationEn = nameDestinationEn
			insert.KmStart = float32(kmStart)
			insert.KmEnd = float32(kmEnd)
			insert.RoadGroupId = roadGroupNumberMap[item.RoadNumber]
			insert.Number = item.SectionRoadNumber
			insert.Distance = float32(length)
			insert.RefDivisionCode = item.OfficeOfHighwaysCode
			insert.RefDistrictCode = item.HighwayDistrictCode
			insert.RefDepotCode = item.DepotCode
			insert.ProvinceCode = append(insert.ProvinceCode, int64(provinceCord))

			err = u.Repository.InsertRoadSection(insert, xt)
			if err != nil {
				return []models.Item{}, err
			}

		} else {
			var update models.InsertRoadSection
			sectionPartID, err := strconv.Atoi(item.SectionPartID)
			if err != nil {
				return []models.Item{}, err
			}

			update.Id = sectionPartID
			update.NameOriginTH = nameOriginTH
			update.NameDestinationTH = nameDestinationTH
			update.NameOriginEn = nameOriginEn
			update.NameDestinationEn = nameDestinationEn
			update.RoadGroupId = roadGroupNumberMap[item.RoadNumber]
			update.Number = item.SectionRoadNumber
			update.KmStart = float32(kmStart)
			update.KmEnd = float32(kmEnd)
			update.Distance = float32(length)
			update.RefDivisionCode = item.OfficeOfHighwaysCode
			update.RefDistrictCode = item.HighwayDistrictCode
			update.RefDepotCode = item.DepotCode
			update.ProvinceCode = append(update.ProvinceCode, int64(provinceCord))

			err = u.Repository.InsertRoadSection(update, xt)
			if err != nil {
				return []models.Item{}, err
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func (u *Usecase) RoadLatest() ([]models.RoadLatest, error) {

	filter := bson.D{{}}
	roadLatest, err := u.Repository.GetRoadLatest(filter)
	if err != nil {
		return nil, err
	}

	return roadLatest, nil
}

func (u *Usecase) KmToM(value string) (int, error) {

	stringKm := strings.ReplaceAll(strings.ReplaceAll(value, ".", ""), "+", "")

	stringInt, err := strconv.Atoi(stringKm)
	if err != nil {
		return 0, err
	}

	return stringInt, nil
}
