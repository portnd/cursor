package repositories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/hsms/handlers"
	"gitlab.com/mims-api-service/src/hsms/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRepositoryHandler(conn *gorm.DB) *handlers.Handler {
	useCase := usecases.NewUseCase(&Repository{conn})
	handler := handlers.NewHandler(useCase)
	return handler
}

func (r *Repository) GetHsmsBridge() ([]models.Hsms01Bridge, error) {
	var hsms01Bridges []models.Hsms01Bridge
	hsms01Table := "hsms_01_bridge"
	refAssetTableID := 359

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_bridge
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Bridge{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Bridge{}, err
		}

		var Hsms01Bridge []models.Hsms01Bridge
		if err := query.Where("road_id = ?", roadID).Find(&Hsms01Bridge).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Bridge{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayFootbridge{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Bridge{}, err
		}

		for _, item := range Hsms01Bridge {
			var hsmsMotorwayFootbridge models.HsmsMotorwayFootbridge
			hsmsMotorwayFootbridge.RoadId = item.RoadID
			hsmsMotorwayFootbridge.RoadAssetID = roadAssetId
			hsmsMotorwayFootbridge.TheGeom = "SRID=4326;" + "POINT(" + r.FloatToString(item.Longitude) + " " + r.FloatToString(item.Latitude) + ")"
			hsmsMotorwayFootbridge.IsDeleted = false
			hsmsMotorwayFootbridge.KM = r.StringToFloat(item.KM) * 1000
			hsmsMotorwayFootbridge.RoadCode = item.RoadCode
			hsmsMotorwayFootbridge.SectionCode = item.SectionCode
			hsmsMotorwayFootbridge.TypeBridge = item.BridgeType
			hsmsMotorwayFootbridge.BridgeLength = item.BridgeLength
			hsmsMotorwayFootbridge.SpanNum = item.SpanNum
			hsmsMotorwayFootbridge.BridgeWidth = item.BridgeWidth
			hsmsMotorwayFootbridge.BridgeHeight = item.BridgeHeight
			hsmsMotorwayFootbridge.SpanWidth = item.SpanWidth
			hsmsMotorwayFootbridge.PlateHeight = item.PlateHeight
			hsmsMotorwayFootbridge.FinishDate = time.Now()
			hsmsMotorwayFootbridge.BudgetOwner = item.BudgetOwner
			hsmsMotorwayFootbridge.Contractor = item.Contractor
			hsmsMotorwayFootbridge.Budget = float64(item.Budget)
			hsmsMotorwayFootbridge.DepotName = item.DepotName
			hsmsMotorwayFootbridge.ApproveStatus = item.ApproveStatus

			if err := query.Create(&hsmsMotorwayFootbridge).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Bridge{}, err
			}

			hsms01Bridges = append(hsms01Bridges, item)
		}
	}
	query.Commit()
	return hsms01Bridges, nil
}

func (r *Repository) GetHsmsGuard() ([]models.Hsms01Guard, error) {
	var hsms01Guard []models.Hsms01Guard
	hsms01Table := "hsms_01_guard"
	refAssetTableID := 357

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_guard
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Guard{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Guard{}, err
		}

		var Hsms01Guard []models.Hsms01Guard
		if err := query.Where("road_id = ?", roadID).Find(&Hsms01Guard).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Guard{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayGuardrail{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Guard{}, err
		}

		for _, item := range Hsms01Guard {
			var hsmsMotorwayGuardrail models.HsmsMotorwayGuardrail

			hsmsMotorwayGuardrail.RoadId = item.RoadId
			hsmsMotorwayGuardrail.RoadAssetId = roadAssetId
			hsmsMotorwayGuardrail.TheGeom = "SRID=4326;" + item.Geom
			hsmsMotorwayGuardrail.IsDeleted = false
			hsmsMotorwayGuardrail.KmStart = r.StringToFloat(item.KmStart) * 1000
			hsmsMotorwayGuardrail.KmEnd = r.StringToFloat(item.KmEnd) * 1000
			hsmsMotorwayGuardrail.RoadCode = item.RoadCode
			hsmsMotorwayGuardrail.SectionCode = item.SectionCode
			hsmsMotorwayGuardrail.LocationType = item.LocationTypeText
			hsmsMotorwayGuardrail.GuardType = item.GuardTypeText
			hsmsMotorwayGuardrail.GuardLeft = r.IntegerToString(item.GuardLeft)
			hsmsMotorwayGuardrail.GuardLeftLength = float64(item.GuardLeftLength)
			hsmsMotorwayGuardrail.GuardRight = r.IntegerToString(item.GuardRight)
			hsmsMotorwayGuardrail.GuardRightLength = float64(item.GuardRightLength)
			hsmsMotorwayGuardrail.GuardCenter = r.IntegerToString(item.GuardCenter)
			hsmsMotorwayGuardrail.GuardCenterLength = float64(item.GuardCenterLength)
			hsmsMotorwayGuardrail.SetupDate = item.SetupDate
			hsmsMotorwayGuardrail.PlanYear = r.IntegerToString(item.PlanYear)
			hsmsMotorwayGuardrail.Contractor = item.Contractor
			hsmsMotorwayGuardrail.Budget = float64(item.Budget)
			hsmsMotorwayGuardrail.DepotName = item.DepotName
			hsmsMotorwayGuardrail.ApproveStatus = item.ApproveStatus
			hsmsMotorwayGuardrail.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayGuardrail).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Guard{}, err
			}

			hsms01Guard = append(hsms01Guard, item)
		}
	}

	query.Commit()
	return hsms01Guard, nil
}

func (r *Repository) GetHsmsInterchange() ([]models.Hsms01Interchange, error) {
	var hsms01Interchanges []models.Hsms01Interchange
	hsms01Table := "hsms_01_interchange"
	refAssetTableID := 362

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_interchange
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Interchange{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Interchange{}, err
		}

		var hsms01Interchange []models.Hsms01Interchange
		if err := query.Where("road_id = ?", roadID).Find(&hsms01Interchange).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Interchange{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayInterchange{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Interchange{}, err
		}

		for _, item := range hsms01Interchange {
			var hsmsMotorwayInterchange models.HsmsMotorwayInterchange

			hsmsMotorwayInterchange.RoadId = item.RoadId
			hsmsMotorwayInterchange.RoadAssetId = roadAssetId
			hsmsMotorwayInterchange.TheGeom = "SRID=4326;" + item.Geom
			hsmsMotorwayInterchange.IsDeleted = false
			hsmsMotorwayInterchange.Km = r.StringToFloat(item.Km) * 1000
			hsmsMotorwayInterchange.RoadCode = item.RoadCode
			hsmsMotorwayInterchange.SectionCode = item.SectionCode
			hsmsMotorwayInterchange.Route2 = r.StringToFloat(item.Route2)
			hsmsMotorwayInterchange.Control2 = r.StringToFloat(item.Control2)
			hsmsMotorwayInterchange.Km2 = r.StringToFloat(item.Km2) * 1000
			hsmsMotorwayInterchange.Route3 = r.StringToFloat(item.Route3)
			hsmsMotorwayInterchange.Control3 = r.StringToFloat(item.Control3)
			hsmsMotorwayInterchange.Km3 = r.StringToFloat(item.Km3) * 1000
			hsmsMotorwayInterchange.Overpass = r.IntegerToString(item.Overpass)
			hsmsMotorwayInterchange.Underpass = r.IntegerToString(item.Underpass)
			hsmsMotorwayInterchange.InterchangeType = item.InterchangeTypeText
			hsmsMotorwayInterchange.DepotName = item.DepotName
			hsmsMotorwayInterchange.ApproveStatus = item.ApproveStatus
			hsmsMotorwayInterchange.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayInterchange).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Interchange{}, err
			}

			hsms01Interchanges = append(hsms01Interchanges, item)
		}
	}

	query.Commit()
	return hsms01Interchanges, nil
}

func (r *Repository) GetHsmsIntersection() ([]models.Hsms01Intersection, error) {
	var hsms01Intersections []models.Hsms01Intersection
	hsms01Table := "hsms_01_intersection"
	refAssetTableID := 361

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_intersection
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Intersection{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Intersection{}, err
		}

		var hsms01Intersection []models.Hsms01Intersection
		if err := query.Where("road_id = ?", roadID).Find(&hsms01Intersection).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Intersection{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayIntersection{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Intersection{}, err
		}

		for _, item := range hsms01Intersection {
			var hsmsMotorwayIntersection models.HsmsMotorwayIntersection

			hsmsMotorwayIntersection.RoadId = item.RoadId
			hsmsMotorwayIntersection.RoadAssetId = roadAssetId
			hsmsMotorwayIntersection.TheGeom = "SRID=4326;" + item.Geom
			hsmsMotorwayIntersection.IsDeleted = false
			hsmsMotorwayIntersection.Km = r.StringToFloat(item.Km) * 1000
			hsmsMotorwayIntersection.RoadCode = item.RoadCode
			hsmsMotorwayIntersection.SectionCode = item.SectionCode
			hsmsMotorwayIntersection.IntersectType = item.IntersectTypeText
			hsmsMotorwayIntersection.Junction = item.Junction
			hsmsMotorwayIntersection.SurfaceType = item.SurfaceTypeText
			hsmsMotorwayIntersection.Width = item.Width
			hsmsMotorwayIntersection.ShoulderWidth = item.ShoulderWidth
			hsmsMotorwayIntersection.ResponsibleTrafficlightType = item.ResponsibleTrafficlightTypeText
			hsmsMotorwayIntersection.ResponsibleFlashlightType = item.ResponsibleFlashlightTypeText
			hsmsMotorwayIntersection.ResponsibleTrafficSignType = item.ResponsibleTrafficSignTypeText
			hsmsMotorwayIntersection.ResponsibleLightType = item.ResponsibleLightTypeText
			hsmsMotorwayIntersection.ResponsibleLightFlType = item.ResponsibleLightFlTypeText
			hsmsMotorwayIntersection.ResponsibleLightHmType = item.ResponsibleLightHmTypeText
			hsmsMotorwayIntersection.ResponsibleLightOtherType = item.ResponsibleLightOtherTypeText
			hsmsMotorwayIntersection.FeatureType = item.FeatureTypeText
			hsmsMotorwayIntersection.DepotName = item.DepotName
			hsmsMotorwayIntersection.ApproveStatus = item.ApproveStatus
			hsmsMotorwayIntersection.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayIntersection).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Intersection{}, err
			}

			hsms01Intersections = append(hsms01Intersections, item)
		}
	}

	query.Commit()
	return hsms01Intersections, nil
}

func (r *Repository) GetHsmsStreetlight() ([]models.Hsms01Light, error) {
	var hsms01Light []models.Hsms01Light
	hsms01Table := "hsms_01_light"
	refAssetTableID := 358

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_light
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Light{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Light{}, err
		}

		var hsms01Light []models.Hsms01Light
		if err := query.Where("road_id = ?", roadID).Find(&hsms01Light).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Light{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayStreetlight{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Light{}, err
		}

		for _, item := range hsms01Light {
			var hsmsMotorwayStreetlight models.HsmsMotorwayStreetlight

			hsmsMotorwayStreetlight.RoadId = item.RoadId
			hsmsMotorwayStreetlight.RoadAssetId = roadAssetId
			hsmsMotorwayStreetlight.TheGeom = "SRID=4326;" + item.Geom
			hsmsMotorwayStreetlight.IsDeleted = false
			hsmsMotorwayStreetlight.KmStart = r.StringToFloat(item.KmStart) * 1000
			hsmsMotorwayStreetlight.KmEnd = r.StringToFloat(item.KmEnd) * 1000
			hsmsMotorwayStreetlight.RoadCode = item.RoadCode
			hsmsMotorwayStreetlight.SectionCode = item.SectionCode
			hsmsMotorwayStreetlight.LocationType = item.LocationTypeText
			hsmsMotorwayStreetlight.LampType = item.LampTypeText
			hsmsMotorwayStreetlight.Watt = r.StringToFloat(item.Watt)
			hsmsMotorwayStreetlight.PoleType = item.PoleTypeText
			hsmsMotorwayStreetlight.SetupDate = item.SetupDate
			hsmsMotorwayStreetlight.PlanYear = r.IntegerToString(item.PlanYear)
			hsmsMotorwayStreetlight.Contractor = item.Contractor
			hsmsMotorwayStreetlight.Budget = r.StringToFloat(item.Budget)
			hsmsMotorwayStreetlight.DepotName = item.DepotName
			hsmsMotorwayStreetlight.ApproveStatus = item.ApproveStatus
			hsmsMotorwayStreetlight.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayStreetlight).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Light{}, err
			}

			hsms01Light = append(hsms01Light, item)
		}
	}

	query.Commit()
	return hsms01Light, nil
}

func (r *Repository) GetHsmsRailwaycrossing() ([]models.Hsms01Railwaycrossing, error) {
	var hsms01Railwaycrossings []models.Hsms01Railwaycrossing
	hsms01Table := "hsms_01_railwaycrossing"
	refAssetTableID := 363

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_railwaycrossing
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Railwaycrossing{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Railwaycrossing{}, err
		}

		var hsms01Railwaycrossing []models.Hsms01Railwaycrossing
		if err := query.Where("road_id = ?", roadID).Find(&hsms01Railwaycrossing).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Railwaycrossing{}, err
		}

		if err := query.Table("hsms_motorway_railwaycrossing").Where("road_id = ?", roadID).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Railwaycrossing{}, err
		}

		for _, item := range hsms01Railwaycrossing {
			var hsmsMotorwayRailwaycrossing models.HsmsMotorwayRailwaycrossing

			hsmsMotorwayRailwaycrossing.RoadId = item.RoadId
			hsmsMotorwayRailwaycrossing.RoadAssetId = roadAssetId
			hsmsMotorwayRailwaycrossing.TheGeom = "SRID=4326;" + item.Geom
			hsmsMotorwayRailwaycrossing.IsDeleted = false
			hsmsMotorwayRailwaycrossing.Km = r.StringToFloat(item.Km) * 1000
			hsmsMotorwayRailwaycrossing.RoadCode = item.RoadCode
			hsmsMotorwayRailwaycrossing.SectionCode = item.SectionCode
			hsmsMotorwayRailwaycrossing.CrossType = item.CrossTypeText
			hsmsMotorwayRailwaycrossing.HighwayType = item.HighwayTypeText
			hsmsMotorwayRailwaycrossing.SurfaceType = item.SurfaceTypeText
			hsmsMotorwayRailwaycrossing.Width = item.Width
			hsmsMotorwayRailwaycrossing.ShoulderSurfaceType = item.ShoulderSurfaceTypeText
			hsmsMotorwayRailwaycrossing.ShoulderWidth = item.ShoulderWidth
			hsmsMotorwayRailwaycrossing.IslandWidth = item.IslandWidth
			hsmsMotorwayRailwaycrossing.RailwayDivision = item.RailwayDivision
			hsmsMotorwayRailwaycrossing.RailwayDistrict = item.RailwayDistrict
			hsmsMotorwayRailwaycrossing.RailwayKm = r.StringToFloat(item.RailwayKm) * 1000
			hsmsMotorwayRailwaycrossing.RailwayWidth = item.RailwayWidth
			hsmsMotorwayRailwaycrossing.RailwayLandWidth = item.RailwayLandWidth
			hsmsMotorwayRailwaycrossing.RailwayAadt = float64(item.RailwayAadt)
			dohInventorys := []string{}
			if item.HasTrafficsign {
				dohInventorys = append(dohInventorys, "ป้่ายจราจร")
			}

			if item.HasTrafficpavement {
				dohInventorys = append(dohInventorys, "เครื่องหมายจราจร")
			}

			if item.HasLight {
				dohInventorys = append(dohInventorys, "ไฟฟ้าแสงสว่าง")
			}

			if item.HasTrafficlight {
				dohInventorys = append(dohInventorys, "สัญญาณไฟจราจร")
			}

			if item.HasOtherInventory {
				dohInventorys = append(dohInventorys, item.OtherInventory)
			}

			railwayInventorys := []string{}
			if item.HasLiftBeam {
				railwayInventorys = append(railwayInventorys, "คานยก")
			}

			if item.HasStraightBeam {
				railwayInventorys = append(railwayInventorys, "คานยกตรง")
			}

			if item.HasHaul {
				railwayInventorys = append(railwayInventorys, "แผงเข้นล้อเลื่อน")
			}

			if item.HasFhb {
				railwayInventorys = append(railwayInventorys, "เอฟ.เอช.บี")
			}

			if item.HasFlb {
				railwayInventorys = append(railwayInventorys, "เอฟ.แอล.ปี")
			}

			if item.HasWarningsign {
				railwayInventorys = append(railwayInventorys, "ป้ายเตือน")
			}

			if item.HasStopsign {
				railwayInventorys = append(railwayInventorys, "ป้ายหยุด")
			}

			if item.HasWinch {
				railwayInventorys = append(railwayInventorys, "มือหมุน")
			}

			if item.HasPanelup {
				railwayInventorys = append(railwayInventorys, "แผงขึ้นลง")
			}

			if item.HasOtherSign {
				railwayInventorys = append(railwayInventorys, item.OtherSign)
			}

			dohInventoryString := strings.Join(dohInventorys, ", ")
			railwayInventoryString := strings.Join(railwayInventorys, ", ")

			hsmsMotorwayRailwaycrossing.DohInventory = dohInventoryString
			hsmsMotorwayRailwaycrossing.RailwayInventory = railwayInventoryString
			hsmsMotorwayRailwaycrossing.DepotName = item.DepotName
			hsmsMotorwayRailwaycrossing.ApproveStatus = item.ApproveStatus
			hsmsMotorwayRailwaycrossing.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayRailwaycrossing).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Railwaycrossing{}, err
			}

			hsms01Railwaycrossings = append(hsms01Railwaycrossings, item)
		}
	}

	query.Commit()
	return hsms01Railwaycrossings, nil
}

func (r *Repository) GetHsmsTrafficlight() ([]models.Hsms01Signal, error) {
	var hsms01Signals []models.Hsms01Signal
	hsms01Table := "hsms_01_signal"
	refAssetTableID := 356

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_signal
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Signal{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Signal{}, err
		}

		var hsms01Signal []models.Hsms01Signal
		if err := query.Where("road_id = ?", roadID).Find(&hsms01Signal).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Signal{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayTrafficlight{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Signal{}, err
		}

		for _, item := range hsms01Signal {
			var hsmsMotorwayTrafficlight models.HsmsMotorwayTrafficlight

			hsmsMotorwayTrafficlight.RoadId = item.RoadId
			hsmsMotorwayTrafficlight.RoadAssetId = roadAssetId
			hsmsMotorwayTrafficlight.TheGeom = "SRID=4326;" + "POINT(" + r.FloatToString(item.Longitude) + " " + r.FloatToString(item.Latitude) + ")"
			hsmsMotorwayTrafficlight.IsDeleted = false
			hsmsMotorwayTrafficlight.Km = r.StringToFloat(item.Km) * 1000
			hsmsMotorwayTrafficlight.RoadCode = item.RoadCode
			hsmsMotorwayTrafficlight.SectionCode = item.SectionCode
			hsmsMotorwayTrafficlight.Location = item.Location
			hsmsMotorwayTrafficlight.LocationType = item.LocationTypeText
			hsmsMotorwayTrafficlight.LampType = item.LampTypeText
			hsmsMotorwayTrafficlight.SystemType = item.SystemTypeText
			hsmsMotorwayTrafficlight.PhaseType = item.PhaseTypeText
			hsmsMotorwayTrafficlight.NumLight = float64(item.NumLight)
			hsmsMotorwayTrafficlight.NumPole = float64(item.NumPole)
			hsmsMotorwayTrafficlight.ControlType = item.ControlTypeText
			hsmsMotorwayTrafficlight.ExpireDate = item.ExpireDate
			hsmsMotorwayTrafficlight.Contractor = item.Contractor
			hsmsMotorwayTrafficlight.Budget = r.StringToFloat(item.Budget)
			hsmsMotorwayTrafficlight.DepotName = item.DepotName
			hsmsMotorwayTrafficlight.ApproveStatus = item.ApproveStatus
			hsmsMotorwayTrafficlight.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayTrafficlight).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Signal{}, err
			}

			hsms01Signals = append(hsms01Signals, item)

		}
	}

	query.Commit()
	return hsms01Signals, nil
}

func (r *Repository) GetHsmsUturnbridge() ([]models.Hsms01Uturnbridge, error) {
	var hsms01Uturnbridges []models.Hsms01Uturnbridge
	hsms01Table := "hsms_01_uturnbridge"
	refAssetTableID := 360

	query := r.StartTransSection()

	// หา group road ID จาก tabel hsms_01_uturnbridge
	roadIDs, err := r.GetRoadIdByTableName(hsms01Table, query)
	if err != nil {
		query.Rollback()
		return []models.Hsms01Uturnbridge{}, err
	}

	// loop road ids
	for _, roadID := range roadIDs {
		roadAssetId, err := r.CheckData(refAssetTableID, roadID, query)
		if err != nil {
			query.Rollback()
			return []models.Hsms01Uturnbridge{}, err
		}

		var hsms01Uturnbridge []models.Hsms01Uturnbridge
		if err := query.Where("road_id = ?", roadID).Find(&hsms01Uturnbridge).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Uturnbridge{}, err
		}

		if err := query.Where("road_id = ?", roadID).Delete(&models.HsmsMotorwayUTurn{}).Error; err != nil {
			query.Rollback()
			return []models.Hsms01Uturnbridge{}, err
		}

		for _, item := range hsms01Uturnbridge {
			var hsmsMotorwayUTurn models.HsmsMotorwayUTurn

			hsmsMotorwayUTurn.RoadId = item.RoadId
			hsmsMotorwayUTurn.RoadAssetId = roadAssetId
			hsmsMotorwayUTurn.TheGeom = "SRID=4326;" + item.Geom
			hsmsMotorwayUTurn.IsDeleted = false
			hsmsMotorwayUTurn.Km = r.StringToFloat(item.Km) * 1000
			hsmsMotorwayUTurn.RoadCode = item.RoadCode
			hsmsMotorwayUTurn.SectionCode = item.SectionCode
			hsmsMotorwayUTurn.DirectionType = item.DirectionTypeText
			hsmsMotorwayUTurn.ConnectedBuilding = item.ConnectedBuildingText
			hsmsMotorwayUTurn.IslandType = item.IslandTypeText
			hsmsMotorwayUTurn.IslandWidth = item.IslandWidth
			hsmsMotorwayUTurn.IslandWidthArea = item.IslandWidthAreaText
			hsmsMotorwayUTurn.LaneWidth = item.LaneWidth
			hsmsMotorwayUTurn.AscentLength = item.AscentLength
			hsmsMotorwayUTurn.AscentSlope = item.AscentSlope
			hsmsMotorwayUTurn.CurveLength = item.CurveLength
			hsmsMotorwayUTurn.CurveRadias = item.CurveRadias
			hsmsMotorwayUTurn.DescentLength = item.DescentLength
			hsmsMotorwayUTurn.DescentSlope = item.DescentSlope
			hsmsMotorwayUTurn.SignHeight = r.FloatToString(item.SignHeight)
			hsmsMotorwayUTurn.SignDistance = item.SignDistance
			hsmsMotorwayUTurn.FlashlightHeight = r.FloatToString(item.FlashlightHeight)
			hsmsMotorwayUTurn.FlashlightDistance = item.FlashlightDistance
			hsmsMotorwayUTurn.SpeedLimit = r.FloatToString(item.SpeedLimit)
			hsmsMotorwayUTurn.HeightLimit = r.FloatToString(item.HeightLimit)
			HasLight := "ไม่มี"
			if item.HasLight {
				HasLight = "มี"
			}
			hsmsMotorwayUTurn.HasLight = HasLight
			hsmsMotorwayUTurn.RailBridgeHeight = item.RailBridgeHeight
			hsmsMotorwayUTurn.DepotName = item.DepotName
			hsmsMotorwayUTurn.ApproveStatus = item.ApproveStatus
			hsmsMotorwayUTurn.UpdateBy = item.UpdateBy

			if err := query.Create(&hsmsMotorwayUTurn).Error; err != nil {
				query.Rollback()
				return []models.Hsms01Uturnbridge{}, err
			}

			hsms01Uturnbridges = append(hsms01Uturnbridges, item)

		}
	}

	query.Commit()
	return hsms01Uturnbridges, nil
}

func (r *Repository) StringToFloat(s string) float64 {
	f := 0.0
	s2 := strings.ReplaceAll(s, "+", ".")
	f, err := strconv.ParseFloat(s2, 32)
	if err != nil {
		f = 0.0
	}
	return f
}

func (r *Repository) IntegerToString(i int) string {
	intger := strconv.Itoa(i)
	return intger
}

func (r *Repository) FloatToString(f float64) string {
	return fmt.Sprintf("%f", f)
}

func (r *Repository) GetRoadAssetID(refAssetTableID, roadID int) (models.RoadAsset, error) {
	query := r.conn
	var roadAssets models.RoadAsset
	subQuery := query.Model(&roadAssets).
		Select("max(revision)").
		Where("ref_asset_table_id = ?", refAssetTableID).
		Where("status = ?", "A").
		Where("road_id = ?", roadID)

	if err := query.Model(&roadAssets).
		Where("ref_asset_table_id = ?", refAssetTableID).
		Where("status = ?", "A").
		Where("revision = (?)", subQuery).
		Where("road_id = ?", roadID).
		First(&roadAssets).Error; err != nil {
		return models.RoadAsset{}, err
	}
	return roadAssets, nil
}

func (r *Repository) CheckData(refAssetTableID, roadId int, query *gorm.DB) (int, error) {
	data := models.RoadAsset{}
	haveData := false
	// เช็คข้อมูล road_asset
	roadAsset, err := r.GetRoadAssetID(refAssetTableID, roadId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			haveData = false
		} else {
			query.Rollback()
			return 0, err
		}
	} else {
		haveData = true
	}

	if haveData {
		// ถ้าเคยมีข้อมูลใร table road_asset จะเปลี่ยนสถานะจาก A -> D
		if err := query.Where("id = ?", roadAsset.Id).Updates(models.RoadAsset{Status: "D"}).Error; err != nil {
			query.Rollback()
			return 0, err
		}
		// insert data road asset status = A
		data = models.RoadAsset{
			RoadId:          roadId,
			RefAssetTableId: refAssetTableID,
			UpdatedDate:     time.Now(),
			CreatedDate:     time.Now(),
			CreatedBy:       99999,
			UpdatedBy:       99999,
			IdParent:        roadAsset.IdParent,
			Status:          "A",
			Revision:        roadAsset.Revision + 1,
			IsExclusiveLock: false,
		}
	} else {
		// กรณีไม่เคยมีข้อมูลใน table road_asset ให้ใส่ status = A, revision = 0
		data = models.RoadAsset{
			RoadId:          roadId,
			RefAssetTableId: refAssetTableID,
			UpdatedDate:     time.Now(),
			CreatedDate:     time.Now(),
			CreatedBy:       99999,
			UpdatedBy:       99999,
			IdParent:        roadAsset.IdParent,
			Status:          "A",
			Revision:        0,
			IsExclusiveLock: false,
		}
	}

	if err := query.Create(&data).Error; err != nil {
		query.Rollback()
		return 0, err
	}

	return data.Id, nil
}

func (r *Repository) GetRoadIdByTableName(tabel string, query *gorm.DB) ([]int, error) {
	var roadIDs []int
	if err := query.Table(tabel).
		Select("road_id").
		Where("road_id IS NOT NULL").
		Group("road_id").
		Find(&roadIDs).Error; err != nil {
		query.Rollback()
		return []int{}, err
	}
	return roadIDs, nil
}

func (t *Repository) StartTransSection() *gorm.DB {
	tx := t.conn.Begin()
	return tx
}

func (t *Repository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (t *Repository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
