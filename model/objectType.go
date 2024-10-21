package model

import "strings"

type ObjectType string

const (
	Table                  ObjectType = "Table"
	TableExtention         ObjectType = "TableExtension"
	Page                   ObjectType = "Page"
	PageExtention          ObjectType = "PageExtension"
	Report                 ObjectType = "Report"
	ReportExtention        ObjectType = "ReportExtension"
	Enum                   ObjectType = "Enum"
	EnumExtention          ObjectType = "EnumExtension"
	PermissionSet          ObjectType = "PermissionSet"
	PermisisonSetExtention ObjectType = "PermissionSetExtension"
	Codeunit               ObjectType = "Codeunit"
	XMLPort                ObjectType = "XMLPort"
	Query                  ObjectType = "Query"
	MenuSuite              ObjectType = "MenuSuite"
)

func MapObjectType(objectType string) ObjectType {
	switch strings.ToLower(objectType) {
	case strings.ToLower(string(Table)):
		return Table
	case strings.ToLower(string(TableExtention)):
		return TableExtention
	case strings.ToLower(string(Page)):
		return Page
	case strings.ToLower(string(PageExtention)):
		return PageExtention
	case strings.ToLower(string(Report)):
		return Report
	case strings.ToLower(string(ReportExtention)):
		return ReportExtention
	case strings.ToLower(string(Enum)):
		return Enum
	case strings.ToLower(string(EnumExtention)):
		return EnumExtention
	case strings.ToLower(string(PermissionSet)):
		return PermissionSet
	case strings.ToLower(string(PermisisonSetExtention)):
		return PermisisonSetExtention
	case strings.ToLower(string(Codeunit)):
		return Codeunit
	case strings.ToLower(string(XMLPort)):
		return XMLPort
	case strings.ToLower(string(Query)):
		return Query
	case strings.ToLower(string(MenuSuite)):
		return MenuSuite
	default:
		return ""
	}
}

func GetObjectTypes() []ObjectType {
	return []ObjectType{
		Table,
		TableExtention,
		Page,
		PageExtention,
		Report,
		ReportExtention,
		Enum,
		EnumExtention,
		PermissionSet,
		PermisisonSetExtention,
		Codeunit,
		XMLPort,
		Query,
		MenuSuite,
	}
}
