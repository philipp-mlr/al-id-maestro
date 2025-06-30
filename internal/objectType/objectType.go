package objectType

import "strings"

type Type string

const (
	Table                  Type = "Table"
	TableExtention         Type = "TableExtension"
	Page                   Type = "Page"
	PageExtention          Type = "PageExtension"
	Report                 Type = "Report"
	ReportExtention        Type = "ReportExtension"
	Enum                   Type = "Enum"
	EnumExtention          Type = "EnumExtension"
	PermissionSet          Type = "PermissionSet"
	PermisisonSetExtention Type = "PermissionSetExtension"
	Codeunit               Type = "Codeunit"
	XMLPort                Type = "XMLPort"
	Query                  Type = "Query"
	MenuSuite              Type = "MenuSuite"
	Unknown                Type = "Unknown"
)

func MapObjectType(objectType string) Type {
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
		return Unknown
	}
}

func GetObjectTypes() []Type {
	return []Type{
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
