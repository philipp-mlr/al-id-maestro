use std::{
    collections::{BTreeMap, HashMap},
    path::Path,
};

use git2::RepositoryInitOptions;

use crate::{id_state::IDState, object_tree::ObjectTree, object_type::ObjectType};

// here will the git provider live
//
pub struct GitUtilizedObjectProvider;

impl super::UtilizedObjectProvider for GitUtilizedObjectProvider {
    fn generate_id_tree(&self) -> ObjectTree {
        let url = "https://github.com/microsoft/BCapps";

        let mut callbacks = git2::RemoteCallbacks::new();
        callbacks.credentials(|_url, username_from_url, _allowed_types| {
            git2::Cred::userpass_plaintext("philipp-mlr", "")
        });

        let mut fo = git2::FetchOptions::new();
        fo.remote_callbacks(callbacks);

        let mut builder = git2::build::RepoBuilder::new();
        builder.fetch_options(fo);

        let repo = builder
            .clone(url, Path::new("/home/philipp/repos/al-id-maestro/repos"))
            .unwrap();

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
