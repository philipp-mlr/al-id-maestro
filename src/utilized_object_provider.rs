pub(crate) mod mock;

use crate::object_manager::IDTree;

pub trait UtilizedObjectProvider {
    fn generate_id_tree(&self) -> IDTree;
}
