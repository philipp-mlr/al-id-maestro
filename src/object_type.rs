use std::{
    fmt::{self, Display},
    str::FromStr,
};

#[derive(Debug, Hash, Eq, PartialEq)]
pub enum ObjectType {
    None,
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

impl Display for ObjectType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "{}",
            match self {
                ObjectType::None => "None",
                ObjectType::Codeunit => "Codeunit",
                ObjectType::Table => "Table",
                ObjectType::TableExtension => "TableExtension",
                ObjectType::Page => "Page",
                ObjectType::PageExtension => "PageExtension",
                ObjectType::Report => "Report",
                ObjectType::ReportExtension => "ReportExtension",
                ObjectType::Enum => "Enum",
                ObjectType::EnumExtension => "EnumExtension",
                ObjectType::PermissionSet => "PermissionSet",
                ObjectType::PermissionSetExtension => "PermissionSetExtension",
                ObjectType::MenuSuite => "MenuSuite",
                ObjectType::Query => "Query",
                ObjectType::XMLPort => "XMLPort",
            }
        )
    }
}
