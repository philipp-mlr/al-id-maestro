// use thiserror::Error;
//
// #[derive(Debug, Error)]
// pub enum Error {
//     #[error("item not found")]
//     NotFound,
//     #[error("item already exists")]
//     AlreadyExists,
//     #[error("invalid input: {0}")]
//     InvalidInput(String),
//     #[error(transparent)]
//     Io(#[from] std::io::Error),
// }
//

use std::collections::{BTreeMap, HashMap};

use crate::object_type::ObjectType;

pub struct ObjectManager {
    object_ranges: HashMap<ObjectType, BTreeMap<u32, bool>>,
}

impl ObjectManager {
    pub fn get_free_id(&self, object_type: ObjectType) -> Result<u32, String> {
        let range = self.object_ranges.get(&object_type).unwrap();

        for (i, b) in range.iter() {
            if !*b {
                return Ok(*i);
            }
        }

        Err("No free id".to_string())
    }
}
