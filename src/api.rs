use std::str::FromStr;

use axum::{
    extract::Path,
    http::{HeaderMap, StatusCode},
};

use crate::{
    object_manager::{self, ObjectManager},
    object_type::ObjectType,
};

struct Api {
    object_manager: ObjectManager,
}

impl Api {
    pub fn new(object_manager: ObjectManager) -> Api {
        Api {
            object_manager: object_manager,
        }
    }

    pub async fn new_object(
        &mut self,
        Path(object_type): Path<String>,
        headers: HeaderMap,
    ) -> Result<String, StatusCode> {
        let api_key = headers.get("X-API-KEY");

        match api_key {
            Some(_v) => println!("yes"),
            None => return Err(StatusCode::UNAUTHORIZED),
        };

        let object_type = match ObjectType::from_str(object_type.as_str()) {
            Ok(v) => v,
            Err(e) => {
                eprintln!("Received incorrect object type: {:?}", e);
                ObjectType::None
            }
        };

        if object_type == ObjectType::None {
            return Err(StatusCode::BAD_REQUEST);
        }

        self.object_manager.get_free_id_for_type(object_type);

        Ok("420".to_string())
    }
}
