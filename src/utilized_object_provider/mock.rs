use std::collections::{BTreeMap, HashMap};

use crate::{id_state::IDState, object_tree::ObjectTree, object_type::ObjectType};

pub struct MockUtilizedObjectProvider;

impl super::UtilizedObjectProvider for MockUtilizedObjectProvider {
    fn generate_id_tree(&self) -> ObjectTree {
        let start = 50000u32;
        let end = 99999u32;

        let mut id_range: BTreeMap<u32, IDState> = BTreeMap::new();
        for i in start..end + 1 {
            id_range.insert(i, IDState::Free);
        }

        let mut map = HashMap::new();
        map.insert(ObjectType::Codeunit, id_range);
        ObjectTree::new(map)
    }
}
