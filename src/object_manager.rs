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

use std::{
    collections::{BTreeMap, HashMap, hash_map},
    hash::Hash,
};

use crate::{
    id_state::IDState,
    object_type::{self, ObjectType},
    utilized_object_provider::UtilizedObjectProvider,
};

pub struct ObjectManager {
    id_tree: IDTree,
}

impl ObjectManager {
    pub fn new(range_provider: impl UtilizedObjectProvider) -> ObjectManager {
        ObjectManager {
            id_tree: range_provider.generate_id_tree(),
        }
    }

    pub fn get_free_id_for_type(&mut self, object_type: ObjectType) -> Result<u32, String> {
        self.id_tree.request_id(object_type)
    }
}
