package models

import "github.com/xuri/excelize/v2"

type Icon struct {
	SheetName  string
	J          interface{}
	Alp        rune
	Ext        string
	Row        string
	FolderPath string
	F          *excelize.File
	Merge      int
}
