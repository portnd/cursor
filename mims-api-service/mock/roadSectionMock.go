package mockdata

import (
	"github.com/lib/pq"
	"gitlab.com/mims-api-service/models"
)

var RoadSectionMock = models.RoadSection{
	Id:                1,
	RoadGroupId:       1,
	Number:            "0001",
	NameOriginTH:      "Road Origin Name (Thai)",
	NameDestinationTH: "Destination Name (Thai)",
	NameOriginEn:      "Origin Name (English)",
	NameDestinationEn: "Destination Name (English)",
	KmStart:           0.0,
	KmEnd:             10.5,
	Distance:          10.5,
	ProvinceCode:      pq.StringArray{"10"},
	Province:          pq.StringArray{"กรุงเทพมหานคร"},
	RefDivisionCode:   "410",
	RefDivision:       RefDivision,
	RefDistrictCode:   "419",
	RefDistrict:       RefDistrict,
	RefDepotCode:      "41904",
	RefDepot:          RefDepotMock,
}

var RoadSectionNilDestMock = models.RoadSection{
	Id:                1,
	RoadGroupId:       1,
	Number:            "0001",
	NameOriginTH:      "Road Origin Name (Thai)",
	NameDestinationTH: "",
	NameOriginEn:      "Origin Name (English)",
	NameDestinationEn: "",
	KmStart:           0.0,
	KmEnd:             10.5,
	Distance:          10.5,
	ProvinceCode:      pq.StringArray{"10"},
	Province:          pq.StringArray{"กรุงเทพมหานคร"},
	RefDivisionCode:   "410",
	RefDivision:       RefDivision,
	RefDistrictCode:   "419",
	RefDistrict:       RefDistrict,
	RefDepotCode:      "41904",
	RefDepot:          RefDepotMock,
}

var RoadSectionByIdMock = models.RoadSectionById{
	RoadSection: RoadSectionMock,
	RoadGroup:   RoadGroupMock,
}

var RoadSectionDataMock = models.RoadSectionData{
	RoadSection: RoadSectionMock,
	Roads:       []models.RoadData{RoadDataMock},
}
