use crate::{
    object_tree::ObjectTree, object_type::ObjectType,
    utilized_object_provider::UtilizedObjectProvider,
};

pub struct ObjectManager {
    id_tree: ObjectTree,
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
