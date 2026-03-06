package models

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type RefDistrict struct {
	Id           int    `json:"id"`
	DistrictCode string `json:"district_code"`
	Name         string `json:"name"`
	NameEn       string `json:"name_en"`
	TheGeom      string `json:"the_geom"`
}

type RefDistrictRes struct {
	Id           int    `json:"id"`
	DistrictCode string `json:"district_code"`
	Name         string `json:"name"`
	NameEn       string `json:"name_en"`
}

type RefDistrictInitData struct {
	Id           int                `json:"id"`
	DistrictCode string             `json:"district_code"`
	Name         string             `json:"name"`
	Depots       []RefDepotInitData `json:"depots" gorm:"foreignKey:DistrictCode;references:DistrictCode"`
}

type RefDistrictInit struct {
	Id           int            `json:"id"`
	DistrictCode string         `json:"district_code"`
	Name         string         `json:"name"`
	OwnerCodeKey string         `json:"owner_code_key"`
	DivisionCode string         `json:"-"`
	Depots       []RefDepotInit `json:"depots" gorm:"foreignKey:DistrictCode;references:DistrictCode"`
}

func (b *RefDistrict) TableName() string {
	return "ref_district"
}

func (b *RefDistrictInitData) TableName() string {
	return "ref_district"
}

func (b *RefDistrictInit) TableName() string {
	return "ref_district"
}

func (b *RefDistrictRes) TableName() string {
	return "ref_district"
}

type IntDataArray []int

func (a *IntDataArray) Scan(src interface{}) error {
	nilArray := []int{}
	if src == "{NULL}" {
		*a = nilArray
	} else {
		trimLeft := strings.TrimLeft(src.(string), "{")
		trimRight := strings.TrimRight(trimLeft, "}")
		arrayStr := strings.Split(trimRight, ",")
		intSlice := make([]int, 0, len(arrayStr))
		for _, valueStr := range arrayStr {
			if valueStr != "NULL" {
				value, _ := strconv.Atoi(valueStr)
				intSlice = append(intSlice, value)
			}
		}
		*a = intSlice
	}
	return nil
}

func (a IntDataArray) Value() (driver.Value, error) {
	return a, nil
}
