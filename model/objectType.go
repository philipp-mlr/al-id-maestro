package model

import "fmt"

const (
	Table            int = 1
	Report           int = 3
	Codeunit         int = 5
	Xmlport          int = 6
	Page             int = 8
	Query            int = 9
	Enum             int = 16
	PageExt          int = 14
	TableExt         int = 15
	ReportExt        int = 22
	EnumExt          int = 17
	PermissionSet    int = 20
	PermisisonSetExt int = 21
)

type ObjectType struct {
	ID   uint `gorm:"primarykey" json:"id"`
	Name string
}

type ObjectTypeQuery struct {
	Query string `json:"query" validate:"required"`
}

func (o ObjectTypeQuery) GetResults(db *DB) ([]ObjectType, error) {
	if db == nil {
		return nil, fmt.Errorf("DB is nil")
	}

	objectTypes := []ObjectType{}

	tx := db.Database.Where("name LIKE ?", "%"+o.Query+"%").Find(&objectTypes)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return objectTypes, nil
}

func NewObjectType(name string) *ObjectType {
	return &ObjectType{
		Name: name,
		ID:   uint(MapNameToInt(name)),
	}
}

func MapNameToInt(name string) int {
	switch name {
	case "table":
		return Table
	case "report":
		return Report
	case "codeunit":
		return Codeunit
	case "xmlport":
		return Xmlport
	case "page":
		return Page
	case "query":
		return Query
	case "enum":
		return Enum
	case "pageextension":
		return PageExt
	case "tableextension":
		return TableExt
	case "reportextension":
		return ReportExt
	case "enumextension":
		return EnumExt
	case "permissionset":
		return PermissionSet
	case "permissionsetextension":
		return PermisisonSetExt
	}
	return 0
}
