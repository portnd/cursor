package models

// type RoadSurfacePreload struct {
// 	RoadSurfaceLane
// 	RefSurface	RefSurface					` gorm:"ForeignKey:ID;references:RefSurfaceId"`
// 	RoadSurface	[]RoadSurfaceGroup `gorm:"ForeignKey:RoadSurfaceId;references:Id"`
// }

type RoadSurfacePreload struct {
	RoadSurface
	//Road                          Road                             ` gorm:"ForeignKey:Id;references:RoadId"`
	RoadInfo                      RoadInfo                         `gorm:"ForeignKey:RoadId;references:RoadId"`
	RoadSurfaceLane               []RoadSurfaceLanePrePareDataLane `gorm:"ForeignKey:RoadSurfaceId;references:Id"`
	RefSurfaceIdShoulderLeftData  SurfaceShoulder                  ` gorm:"ForeignKey:ID;references:RefSurfaceIdShoulderRight"`
	RefSurfaceIdShoulderRightData SurfaceShoulder                  ` gorm:"ForeignKey:ID;references:RefSurfaceIdShoulderRight"`
	MaterialBase                  RefMaterialBase                  ` gorm:"ForeignKey:ID;references:RefMaterialBaseId"`
	MaterialSubbase               RefMaterialSubbase               ` gorm:"ForeignKey:ID;references:RefMaterialSubbaseId"`
	MaterialSubgrade              RefMaterialSubgrade              ` gorm:"ForeignKey:ID;references:RefMaterialSubgradeId"`
	// RefDataStatus RefDataStatus				` gorm:"ForeignKey:ID;references:Status"`
}

// type RoadRefDirectionPreload struct {
// 	Road
// 	RefDirection     RefDirection       `json:"direction" gorm:"ForeignKey:RefDirectionId;references:ID"`
// }
