package models

import "time"

// Todo ...
type RoadSurface struct {
	Id                        int       `json:"id"`
	RoadId                    int       `json:"road_id"`
	Year                      int       `json:"year"`
	KmStart                   float64   `json:"km_start"`
	KmEnd                     float64   `json:"km_end"`
	WidthSurface              float64   `json:"width_surface"`
	WidthShoulderLeft         float64   `json:"width_shoulder_left"`
	WidthShoulderRight        float64   `json:"width_shoulder_right"`
	TheGeom                   string    `json:"the_geom"`
	Revision                  int       `json:"revision"`
	Status                    string    `json:"status"`
	IdParent                  int       `json:"id_parent"`
	RejectReason              string    `json:"reject_reason"`
	ThicknessSurface          float64   `json:"thickness_surface"`
	ThicknessSurfaceConcrete  float64   `json:"thickness_surface_concrete" gorm:"column:thickness_surface_concrete"`
	RefSurfaceIdShoulderLeft  int       `json:"ref_surface_id_shoulder_left"`
	RefSurfaceIdShoulderRight int       `json:"ref_surface_id_shoulder_right"`
	ThicknessBase             *float64  `json:"thickness_base"`
	RefMaterialBaseId         int       `json:"ref_material_base_id"`
	ThicknessSubbase          *float64  `json:"thickness_subbase"`
	RefMaterialSubbaseId      int       `json:"ref_material_subbase_id"`
	ThicknessSubgrade         *float64  `json:"thickness_subgrade"`
	RefMaterialSubgradeId     int       `json:"ref_material_subgrade_id"`
	HashData                  string    `json:"hash_data"`
	SurfaceCrosssectionCode   int       `json:"surface_crosssection_code"`
	RefStructureID            int       `json:"ref_structure_id"`
	AreaSurface               float64   `json:"area_surface"`
	CreatedBy                 int       `json:"created_by"`
	CreatedDate               time.Time `json:"created_date"`
	UpdatedBy                 int       `json:"updated_by"`
	UpdatedDate               time.Time `json:"updated_date"`
	ThicknessConcreteSlab     *float64  `json:"thickness_concrete_slab"`
	SurfaceGroup              int       `json:"surface_group"`
}

type RoadSurfacePointer struct {
	Id                        int       `json:"id"`
	RoadId                    int       `json:"road_id"`
	Year                      int       `json:"year"`
	KmStart                   float64   `json:"km_start"`
	KmEnd                     float64   `json:"km_end"`
	WidthSurface              float64   `json:"width_surface"`
	WidthShoulderLeft         float64   `json:"width_shoulder_left"`
	WidthShoulderRight        float64   `json:"width_shoulder_right"`
	TheGeom                   string    `json:"the_geom"`
	Revision                  int       `json:"revision"`
	Status                    string    `json:"status"`
	IdParent                  int       `json:"id_parent"`
	RejectReason              string    `json:"reject_reason"`
	ThicknessSurface          float64   `json:"thickness_surface"`
	ThicknessSurfaceConcrete  *float64  `json:"thickness_surface_concrete"`
	RefSurfaceIdShoulderLeft  int       `json:"ref_surface_id_shoulder_left"`
	RefSurfaceIdShoulderRight int       `json:"ref_surface_id_shoulder_right"`
	ThicknessBase             *float64  `json:"thickness_base"`
	RefMaterialBaseId         *int      `json:"ref_material_base_id"`
	ThicknessSubbase          *float64  `json:"thickness_subbase"`
	RefMaterialSubbaseId      *int      `json:"ref_material_subbase_id"`
	ThicknessSubgrade         *float64  `json:"thickness_subgrade"`
	RefMaterialSubgradeId     *int      `json:"ref_material_subgrade_id"`
	HashData                  string    `json:"hash_data"`
	SurfaceCrosssectionCode   int       `json:"surface_crosssection_code"`
	RefStructureID            int       `json:"ref_structure_id"`
	AreaSurface               float64   `json:"area_surface"`
	CreatedBy                 int       `json:"created_by"`
	CreatedDate               time.Time `json:"created_date"`
	UpdatedBy                 int       `json:"updated_by"`
	UpdatedDate               time.Time `json:"updated_date"`
	ThicknessConcreteSlab     *float64  `json:"thickness_concrete_slab"`
	SurfaceGroup              int       `json:"surface_group"`
}

type RoadSurfaceData struct {
	RoadSurface
	RoadSurfaceLane RoadSurfaceLane `json:"road_surface_lane" gorm:"ForeignKey:RoadSurfaceId;AssociationForeignKey:Id"`
}

type RoadSurfaceData2 struct {
	RoadSurface
	RoadSurfaceLane []RoadSurfaceLane `json:"road_surface_lane" gorm:"ForeignKey:RoadSurfaceId;AssociationForeignKey:Id"`
}

type RoadSurfaceGroupType struct {
	SurfaceGroupType string `json:"surface_group_type"`
}

type RoadSurfacePrepareData struct {
	RoadSurface
	RoadSurfaceLane        []RoadSurfaceLanePrePareData `json:"road_surface_lane" gorm:"ForeignKey:RoadSurfaceId;AssociationForeignKey:ID"`
	RefStructureSurface    RefStructureSurface          `json:"ref_structure_surface" gorm:"ForeignKey:SurfaceCrosssectionCode;AssociationForeignKey:ID"`
	RefSurfaceShoulderLeft RefSurface                   `json:"ref_surface_shoulder_left" gorm:"ForeignKey:RefSurfaceIdShoulderLeft;AssociationForeignKey:ID"`
	RefMaterialBase        RefMaterialBase              `json:"ref_material_base" gorm:"ForeignKey:RefMaterialBaseId;AssociationForeignKey:ID"`
	RefMaterialSubbase     RefMaterialSubbase           `json:"ref_material_subbase" gorm:"ForeignKey:RefMaterialSubbaseId;AssociationForeignKey:ID"`
	RefMaterialSubgrade    RefMaterialSubgrade          `json:"ref_material_subgrade" gorm:"ForeignKey:RefMaterialSubgradeId;AssociationForeignKey:ID"`
}

type RoadSurfaceForPreload struct {
	Id     int `json:"id"`
	RoadId int `json:"road_id"`
}

type RoadSurfaceForCount struct {
	Id     int    `json:"id"`
	RoadId int    `json:"road_id"`
	Status string `json:"status"`
}

type RoadSurfaceIcon struct {
	ID        int    `json:"id"`
	RoadId    int    `json:"-"`
	Name      string `json:"name"`
	ColorCode string `json:"color_code"`
}

type RefSurfaceRoad struct {
	Id           int          `json:"-"`
	RoadId       int          `json:"road_id"`
	RefSurfaceId IntDataArray `json:"ref_Surface_Id" gorm:"type:integer[]"`
}

// TableName use to specific table
func (b *RoadSurface) TableName() string {
	return "road_surface"
}

func (b *RoadSurfaceForCount) TableName() string {
	return "road_surface"
}

func (b *RoadSurfaceForPreload) TableName() string {
	return "road_surface"
}

func (b *RoadSurfaceData) TableName() string {
	return "road_surface"
}

func (b *RoadSurfaceIcon) TableName() string {
	return "road_surface"
}

func (b *RefSurfaceRoad) TableName() string {
	return "road_surface"
}
