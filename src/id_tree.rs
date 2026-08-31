use std::collections::{BTreeMap, HashMap};

use crate::{id_state::IDState, object_type::ObjectType};

pub struct IDTree {
    ids: HashMap<ObjectType, BTreeMap<u32, IDState>>,
}

impl IDTree {
    pub fn new(map: HashMap<ObjectType, BTreeMap<u32, IDState>>) -> IDTree {
        IDTree { ids: map }
    }

    pub fn request_id(&mut self, object_type: ObjectType) -> Result<u32, String> {
        let id = self.find_free_id(&object_type)?;
        self.issue_id(&object_type, &id);
        Ok(id)
    }

    fn find_free_id(&self, object_type: &ObjectType) -> Result<u32, String> {
        let tree = self.ids.get(&object_type).unwrap();

        for (i, b) in tree.iter() {
            if *b == IDState::Free {
                return Ok(*i);
            }
        }

        Err("No free id".to_string())
    }

    fn issue_id(&mut self, object_type: &ObjectType, id: &u32) {
        self.set_id_state(object_type, id, IDState::Issued);
    }

    fn utilize_id(&mut self, object_type: &ObjectType, id: &u32) {
        self.set_id_state(object_type, id, IDState::Utilized);
    }

    fn set_id_state(&mut self, object_type: &ObjectType, id: &u32, id_state: IDState) {
        *self
            .ids
            .get_mut(&object_type)
            .unwrap()
            .get_mut(&id)
            .unwrap() = id_state;
    }

    fn get_id_state(&self, object_type: &ObjectType, id: &u32) -> Option<&IDState> {
        self.ids.get(&object_type).unwrap().get(&id)
    }
}
