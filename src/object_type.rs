#[derive(Debug, Hash, Eq, PartialEq)]
pub enum ObjectType {
    Codeunit,
    Table,
    TableExtension,
    Page,
    PageExtension,
    Report,
    ReportExtension,
    Enum,
    EnumExtension,
    PermissionSet,
    PermissionSetExtension,
    MenuSuite,
    Query,
    XMLPort,
}

use std::str::FromStr;

impl FromStr for ObjectType {
    type Err = String;

    fn from_str(input: &str) -> Result<ObjectType, Self::Err> {
        match input {
            "Codeunit" => Ok(ObjectType::Codeunit),
            "Table" => Ok(ObjectType::Table),
            "TableExtension" => Ok(ObjectType::TableExtension),
            "Page" => Ok(ObjectType::Page),
            "PageExtension" => Ok(ObjectType::PageExtension),
            "Report" => Ok(ObjectType::Report),
            "ReportExtension" => Ok(ObjectType::ReportExtension),
            "Enum" => Ok(ObjectType::Enum),
            "EnumExtension" => Ok(ObjectType::EnumExtension),
            "PermissionSet" => Ok(ObjectType::PermissionSet),
            "PermissionSetExtension" => Ok(ObjectType::PermissionSetExtension),
            "MenuSuite" => Ok(ObjectType::MenuSuite),
            "Query" => Ok(ObjectType::Query),
            "XMLPort" => Ok(ObjectType::XMLPort),
            _ => Err("unknown object type".to_string()),
        }
    }
}
