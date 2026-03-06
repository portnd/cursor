package repositories

import (
	"log"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/ref/handlers"
	"gitlab.com/mims-api-service/src/ref/usecases"

	servicesDB "gitlab.com/mims-api-service/services/database"
	"gorm.io/gorm"
)

type refRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRefRepositoryHandler(conn *gorm.DB) *handlers.RefHandler {
	servicesDB := servicesDB.NewServicesDatabase(conn)
	useCase := usecases.NewRefUseCase(&refRepository{conn}, servicesDB)
	handler := handlers.NewRefHandler(useCase)
	return handler
}

func (rr *refRepository) GetRef(result interface{}) error {
	query := rr.conn
	err := query.Order("id ASC").Find(result).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (rr *refRepository) GetRefAssetSignImage(result interface{}) error {
	query := rr.conn
	err := query.Where("status = ?", 1).Order("id ASC").Find(result).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (rr *refRepository) GetRefStatus(result interface{}) error {
	query := rr.conn
	err := query.Where("status = ?", 1).Order("id ASC").Find(result).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (rr *refRepository) GetRefTableList() ([]models.RefTableList, error) {
	tableLists := []models.RefTableList{}

	err := rr.conn.Raw(`
	SELECT DISTINCT pgc.relname AS ref_name,
	(SELECT c2.description FROM pg_description c2
	 WHERE c2.objoid = pgc.oid AND c2.objsubid = 0) AS ref_desc
FROM pg_attribute a
JOIN pg_class pgc ON pgc.oid = a.attrelid
LEFT JOIN pg_index i ON (pgc.oid = i.indrelid AND i.indkey[0] = a.attnum)
LEFT JOIN pg_description com ON (pgc.oid = com.objoid AND a.attnum = com.objsubid)
LEFT JOIN pg_attrdef def ON (a.attrelid = def.adrelid AND a.attnum = def.adnum)
LEFT JOIN pg_catalog.pg_namespace ns ON pgc.relnamespace = ns.oid
LEFT JOIN pg_constraint p ON p.conrelid = pgc.oid AND a.attnum = ANY (p.conkey)  
WHERE a.attnum > 0
AND pgc.oid = a.attrelid
AND pg_table_is_visible(pgc.oid)
AND NOT a.attisdropped
AND pgc.relkind IN ('r', 'v')
AND ns.nspname IN ('public')
AND pgc.relname LIKE 'ref_asset_%' 
AND pgc.relname NOT IN ('ref_asset_table', 'ref_asset_table_staff', 'ref_asset_table_columns')
ORDER BY ref_desc;

	`).Scan(&tableLists).Error
	if err != nil {
		return tableLists, err
	}

	return tableLists, nil
}

func (rr *refRepository) GetRefColorList() ([]models.RefColorList, error) {
	var colorList []models.RefColorList

	query := rr.conn
	err := query.Find(&colorList).Error
	if err != nil {
		return nil, err
	}

	return colorList, nil
}

func (rr *refRepository) GetRefRoadTypeIcon() ([]models.RefRoadTypeIcon, error) {
	var icons []models.RefRoadTypeIcon

	query := rr.conn
	err := query.Find(&icons).Error
	if err != nil {
		return nil, err
	}

	return icons, nil
}

func (rr *refRepository) GetRoadConditionGrades() ([]models.ParamsConditionPreload, error) {
	var raramsCondition []models.ParamsConditionPreload

	if err := rr.conn.Preload("RefOwner").Preload("RefGrade").Find(&raramsCondition).Error; err != nil {
		return raramsCondition, err
	}

	return raramsCondition, nil
}

func (rr *refRepository) GetParameterVehicleTypeListByRoadGroupId(road_group_id int) ([]models.RefAadtParameterVehicleType, error) {
	var refAadtParameterVehicleTypeList []models.RefAadtParameterVehicleType
	query := rr.conn
	if road_group_id == 7 {
		query = query.Where("for_road_group_id = ?", 7)
	} else {
		query = query.Where("for_road_group_id != ?", 7)
	}
	query = query.Order("id")
	err := query.Find(&refAadtParameterVehicleTypeList).Error
	if err != nil {
		return nil, err
	}

	return refAadtParameterVehicleTypeList, nil
}

func (rr *refRepository) GetRefCriteriaType() ([]models.RefCriteriaType, error) {
	var criteriaType []models.RefCriteriaType

	query := rr.conn
	err := query.Find(&criteriaType).Error
	if err != nil {
		return nil, err
	}

	return criteriaType, nil

}

func (rr *refRepository) GetRoadUserCostAcc() ([]models.RefRoadUserCostAcc, error) {
	var refRoadUserCostAcc []models.RefRoadUserCostAcc

	query := rr.conn
	err := query.Find(&refRoadUserCostAcc).Error
	if err != nil {
		return nil, err
	}

	return refRoadUserCostAcc, nil
}

func (rr *refRepository) GetRoadUserCostRuc() ([]models.RefRoadUserCostRuc, error) {
	var refRoadUserCostRuc []models.RefRoadUserCostRuc

	query := rr.conn
	err := query.Find(&refRoadUserCostRuc).Error
	if err != nil {
		return nil, err
	}

	return refRoadUserCostRuc, nil

}

func (rr *refRepository) GetMaintenanceAnalysisStrategicBudgetType() ([]models.RefMaintenanceAnalysisCondition, error) {
	var maintenanceAnalysisStrategicBudgetType []models.RefMaintenanceAnalysisCondition
	query := rr.conn
	err := query.Find(&maintenanceAnalysisStrategicBudgetType).Error
	if err != nil {
		return nil, err
	}
	return maintenanceAnalysisStrategicBudgetType, nil

}
func (rr *refRepository) GetMaintenanceAnalysisStrategicTargetType() ([]models.RefMaintenanceAnalysisTarget, error) {
	var maintenanceAnalysisStrategicTagetType []models.RefMaintenanceAnalysisTarget
	query := rr.conn
	err := query.Find(&maintenanceAnalysisStrategicTagetType).Error
	if err != nil {
		return nil, err
	}
	return maintenanceAnalysisStrategicTagetType, nil
}
func (rr *refRepository) GetMaintenanceAnalysisStrategic() ([]models.MaintenanceAnalysisStrategicTypePreload, error) {
	var maintenanceAnalysisStrategicTypePreload []models.MaintenanceAnalysisStrategicTypePreload
	query := rr.conn
	query = query.Preload("Budget")
	query = query.Preload("Budget.Target")
	err := query.Find(&maintenanceAnalysisStrategicTypePreload).Error
	if err != nil {
		return nil, err
	}
	return maintenanceAnalysisStrategicTypePreload, nil

}

func (r *refRepository) GetDataList(model interface{}, where string) error {
	query := r.conn
	if where != "" {
		query = query.Where(where)
	}
	err := query.Order("id ASC").Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *refRepository) GetParamsCondition(ownerId int) ([]models.ParamsConditionPreload, error) {
	records := []models.ParamsConditionPreload{}

	err := sr.conn.Where("ref_owner_id = ?", ownerId).
		Preload("RefOwner").
		Preload("RefOwner.RefConditionRange").
		Preload("RefGrade").
		Order("id ASC").Find(&records).Error
	if err != nil {
		return []models.ParamsConditionPreload{}, err
	}

	return records, nil
}

func (sr *refRepository) GetParamsRoadLine(ownerId int) ([]models.ParamsRoadLinePreload, error) {
	records := []models.ParamsRoadLinePreload{}

	err := sr.conn.Where("ref_owner_road_line_id = ?", ownerId).Preload("RefOwnerRoadLine").Preload("RefGrade").Order("id ASC").Find(&records).Error
	if err != nil {
		return []models.ParamsRoadLinePreload{}, err
	}

	return records, nil
}

func (rr *refRepository) GetRefDistrictsList() (interface{}, error) {
	var districts []models.RefDistrictData

	query := rr.conn
	query = query.Select("id, district_code, name, name_en, division_code, st_astext(the_geom) as the_geom")

	query = query.Preload("Depots", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, depot_code, name, district_code,st_astext(the_geom) as the_geom")
	})
	if err := query.Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

func (rr *refRepository) GetRefDistrictsInitList(isAllData, isOwnerData bool, depotCode string) (interface{}, error) {
	var districts []models.RefDistrictInitData

	query := rr.conn
	query = query.Select("ref_district.id, ref_district.district_code, ref_district.name").
		Joins("RIGHT JOIN road_section on road_section.ref_district_code = ref_district.district_code").
		Group("ref_district.id, ref_district.district_code, ref_district.name").
		Order("ref_district.id ASC ")

	query = query.Preload("Depots", func(db *gorm.DB) *gorm.DB {
		db = db.Select("ref_depot.id, ref_depot.depot_code, ref_depot.name, ref_depot.district_code, COALESCE(ARRAY_AGG(DISTINCT sect.id), '{}') AS road_section_id").
			Joins("LEFT JOIN road_section AS sect ON ref_depot.depot_code = sect.ref_depot_code").
			Group("ref_depot.id, ref_depot.depot_code, ref_depot.name, ref_depot.district_code")
		if isOwnerData && !isAllData {
			db = db.Where("sect.ref_depot_code = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("sect.ref_depot_code = ?", "no_data")
		}
		return db
	})
	if err := query.Order("id ASC").Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

func (rr *refRepository) GetRoadGroupList() ([]models.RoadGroupInitData, error) {
	var roadGroup []models.RoadGroupInitData

	query := rr.conn

	query = query.Select("road_group.id, road_group.number, road_group.short_name, COALESCE(ARRAY_AGG(DISTINCT road_section.ref_division_code) FILTER (WHERE road_section.ref_division_code IS NOT NULL), '{}') as ref_division_codes, COALESCE(ARRAY_AGG(DISTINCT road_section.ref_district_code) FILTER (WHERE road_section.ref_district_code IS NOT NULL), '{}') as ref_district_codes").
		Joins("LEFT JOIN road_section ON road_group.id = road_section.road_group_id").
		Group("road_group.id, road_group.number, road_group.short_name")

	if err := query.Order("road_group.id ASC").Find(&roadGroup).Error; err != nil {
		return nil, err
	}

	return roadGroup, nil
}

func (rr *refRepository) GetRoadSectionList(isAllData, isOwnerData bool, depotCode string) ([]models.RoadSectionInitData, error) {
	var roadSection []models.RoadSectionInitData

	query := rr.conn

	query = query.Select("road_section.id,road_section.road_group_id ,road_section.number, road_section.name_origin_th , road_section.name_destination_th, ARRAY_AGG(DISTINCT road_surface_lane.ref_surface_id) as ref_surface_id").
		Joins("LEFT Join road on road.road_section_id = road_section.id and road.is_active = true").
		Joins("LEFT join road_surface on road_surface.road_id = road.id and road_surface.status = 'A'").
		Joins("LEFT join road_surface_lane on road_surface_lane.road_surface_id = road_surface.id and road_surface_lane.ref_surface_id <> 0").
		Group("road_section.id,road_section.road_group_id ,road_section.number, road_section.name_origin_th , road_section.name_destination_th")
	if isOwnerData && !isAllData {
		query = query.Where("road_section.ref_depot_code = ?", depotCode)
	}
	if !isOwnerData && !isAllData {
		query = query.Where("road_section.ref_depot_code = ?", "no_data")
	}
	if err := query.Order("road_section.id ASC ").Find(&roadSection).Error; err != nil {
		return nil, err
	}

	return roadSection, nil
}

func (rr *refRepository) GetRefDivisionInitList(isAllData, isOwnerData bool, depotCode string) ([]models.RefDivisionList, error) {
	var division []models.RefDivisionList

	query := rr.conn
	query = query.Select("ref_division.id, ref_division.division_code, ref_division.name , concat('ref_division_code-', division_code) as owner_code_key").
		Joins("JOIN road_section on road_section.ref_division_code = ref_division.division_code").
		Joins("JOIN road on road.road_section_id = road_section.id and road.is_active = true").
		Group("ref_division.id, ref_division.division_code, ref_division.name")

	query = query.Preload("Districts", func(db *gorm.DB) *gorm.DB {
		return db.Select("ref_district.*, concat('ref_district_code-', district_code) as owner_code_key").
			Joins("JOIN road_section on road_section.ref_district_code = ref_district.district_code").
			Joins("JOIN road on road.road_section_id = road_section.id and road.is_active = true").
			Group("ref_district.id, ref_district.district_code, ref_district.name").
			Order("ref_district.id ASC")
	})
	query = query.Preload("Districts.Depots", func(db *gorm.DB) *gorm.DB {
		db = db.Select("*, concat('ref_depot_code-',depot_code) as owner_code_key")
		if isOwnerData && !isAllData {
			db = db.Where("depot_code = ?", depotCode)
		}

		if !isOwnerData && !isAllData {
			db = db.Where("depot_code = ?", "no_data")
		}
		return db
	})
	if err := query.Order("ref_division.id ASC").Find(&division).Error; err != nil {
		return nil, err
	}
	return division, nil
}

func (rr *refRepository) GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error) {
	var refCriteriaMethod []models.RefCriteriaMethod
	if err := rr.conn.Where("is_deleted = ?", false).Find(&refCriteriaMethod).Error; err != nil {
		return refCriteriaMethod, err
	}
	return refCriteriaMethod, nil
}

func (rr *refRepository) GetRefUserOwner() ([]responses.RefUserOwner, error) {

	var refUserOwner []models.RefUserOwner
	if err := rr.conn.Table("ref_user_owner").Find(&refUserOwner).Error; err != nil {
		return []responses.RefUserOwner{}, err
	}

	var res []responses.RefUserOwner
	for _, item := range refUserOwner {
		var data responses.RefUserOwner
		data.ID = item.Id
		data.Name = item.Name
		var refDepots []responses.RefDepot
		if item.Id == 3 {
			if err := rr.conn.Table("ref_depot").Where("district_code = ?", "261").Find(&refDepots).Error; err != nil {
				return []responses.RefUserOwner{}, err
			}
			copier.Copy(&refDepots, refDepots)
			data.RefDepot = refDepots
		} else {
			data.RefDepot = []string{}
		}
		res = append(res, data)
	}
	return res, nil
}
