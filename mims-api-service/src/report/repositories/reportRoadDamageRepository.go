package repositories

import (
	"gitlab.com/mims-api-service/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
)

// ////////////// NEW MIMS ////////////////
func (r *Repository) GetRoadInfoForRoadDamage(roadID int) (*models.RoadReportInfo, error) {
	var data models.RoadReportInfo
	err := r.conn.
		Table("road").
		Select(`
        road_group.name AS road_group_name,
        road_info.name AS road_name,
        road.road_code,
        road_info.km_start,
        road_info.km_end,
		road_info.road_color_code
    	`).
		Joins("JOIN road_group ON road_group.id = road.road_group_id").
		Joins("JOIN road_info ON road_info.road_id = road.id").
		Where("road_id = ?", roadID).
		Find(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository) GetReportDamageForRoadDamage(year, roadID int) (*models.DataReportDamage, []models.DataRoadDamage, []models.PositionDamage, error) {

	var data models.DataReportDamageModel
	err := r.conn.
		Table("road").
		Select(`
        road_group.name AS road_group_name,
        road_info.name AS road_name,
        road.road_code,
        road_info.km_start,
        road_info.km_end
    	`).
		Joins("JOIN road_group ON road_group.id = road.road_group_id").
		Joins("JOIN road_info ON road_info.road_id = road.id").
		Where("road.id = ?", roadID). // Fix: change "road_id" to "road.id"
		Find(&data).Error
	if err != nil {
		return nil, nil, nil, err
	}

	var detail []models.DataRoadDamageModel

	subquery := r.conn.
		Table("road_damage").
		Select(`
		lane_no, 
		COUNT(*) AS count_lanes`).
		Where("year = ? AND status = ? AND road_id = ?", year, "A", roadID).
		Group("lane_no").
		Having("COUNT(*) = 1")

	err = r.conn.
		Table("road_damage rd").
		Select(`
		rd.lane_no, 
		rd.ac_icrack, 
		rd.ac_ucrack, 
		rd.ac_ravelling, 
		rd.ac_patching, 
		rd.ac_pothole_area, 
		rd.ac_bleeding, 
		rd.ac_pothole_count,
		rd.cc_transverse_crack, 
		rd.cc_non_transverse_crack, 
		rd.cc_spalling,
		rd.cc_corner_break, 
		rd.cc_joint_seal_damage, 
		rd.cc_patching,
		rd.cc_scaling
		`).
		Joins("JOIN (?) subquery ON rd.lane_no = subquery.lane_no", subquery).
		Where("rd.year = ? AND rd.status = ? AND rd.road_id = ?", year, "A", roadID).
		Order("lane_no ASC").
		Find(&detail).Error

	if err != nil {
		return nil, nil, nil, err
	}

	var position []models.PositionDamage

	err = r.conn.
		Table("road_damage").
		Select(`
		road_damage.lane_no,
        road_damage_m.km,
        road_damage_m.img_filepath,
        road_damage_m.ac_icrack, 
		road_damage_m.ac_ucrack, 
		road_damage_m.ac_ravelling, 
		road_damage_m.ac_patching, 
		road_damage_m.ac_pothole_area, 
		road_damage_m.ac_bleeding, 
		road_damage_m.ac_pothole_count,
		road_damage_m.cc_transverse_crack, 
		road_damage_m.cc_non_transverse_crack, 
		road_damage_m.cc_spalling,
		road_damage_m.cc_corner_break, 
		road_damage_m.cc_joint_seal_damage, 
		road_damage_m.cc_patching,
		road_damage_m.cc_scaling
    	`).
		Joins(`JOIN road_damage_range ON road_damage_range.road_damage_id = road_damage.id`).
		Joins(`JOIN road_damage_m ON road_damage_m.road_damage_range_id = road_damage_range.id`).
		Where("road_damage.year = ? AND road_damage.status = ? AND road_damage.road_id = ?", year, "A", roadID).
		Order("road_damage.lane_no ASC, road_damage_m.km ASC").
		Find(&position).Error

	if err != nil {
		return nil, nil, nil, err
	}

	var result models.DataReportDamage

	var result2 []models.DataRoadDamage

	copier.Copy(&result, &data)
	copier.Copy(&result2, &detail)

	result.Detail = result2
	return &result, result2, position, nil
}

func (r *Repository) GetRoadSectionByIDForRoadDamage(id int) (models.RoadSection, error) {
	result := models.RoadSection{}
	err := r.conn.Where("id = ?", id).First(&result).Error
	if err != nil {
		return models.RoadSection{}, err
	}
	return result, nil
}

func (r *Repository) GetRoadFromSectionIDForRoadDamage(SectionID int) ([]int, error) {
	var roadIDs []int
	err := r.conn.Model(&models.Road{}).Where("road_section_id = ?", SectionID).Pluck("id", &roadIDs).Error
	if err != nil {
		return nil, err
	}
	return roadIDs, nil
}

func (r *Repository) GetRoadGroupByIDForRoadDamage(id int) (models.RoadGroup, error) {
	result := models.RoadGroup{}
	err := r.conn.Where("id = ?", id).First(&result).Error
	if err != nil {
		return models.RoadGroup{}, err
	}
	return result, nil
}
