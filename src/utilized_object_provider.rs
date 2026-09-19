use crate::object_tree::ObjectTree;

pub(crate) mod git;
pub(crate) mod mock;

pub trait UtilizedObjectProvider {
    fn generate_id_tree(&self) -> ObjectTree;
}
