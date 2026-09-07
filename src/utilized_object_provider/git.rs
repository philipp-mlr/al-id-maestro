use std::{
    collections::{BTreeMap, HashMap},
    path::Path,
};

use axum::extract::path;
use git2::{Error, Repository, RepositoryInitOptions};
use regex::Regex;
use walkdir::WalkDir;

use crate::{id_state::IDState, object_tree::ObjectTree, object_type::ObjectType};

// here will the git provider live
//
pub struct GitUtilizedObjectProvider;

impl super::UtilizedObjectProvider for GitUtilizedObjectProvider {
    fn generate_id_tree(&self) -> ObjectTree {
        let url = "https://github.com/fond-of/bc-apps";
        let path = Path::new("/home/philipp/repos/al-id-maestro/repos");

        let mut callbacks = git2::RemoteCallbacks::new();
        callbacks.credentials(|_url, username_from_url, _allowed_types| {
            git2::Cred::userpass_plaintext("philipp-mlr", "")
        });

        callbacks.update_tips(|tip, oid1, oid2| {
            println!("Something is updating: {}", tip);
            true
        });

        let mut fo = git2::FetchOptions::new();
        fo.remote_callbacks(callbacks);

        let mut builder = git2::build::RepoBuilder::new();
        builder.fetch_options(fo);

        let repo = Self::open_repo(&mut builder, path, url).unwrap();

        let branches = repo.branches(Some(git2::BranchType::Remote)).unwrap();

        for branch_result in branches {
            let (branch, _btype) = branch_result.unwrap();
            branch.get().target().unwrap();

            repo.set_head(branch.name().unwrap().unwrap());

            let mut cb = git2::build::CheckoutBuilder::new();
            repo.checkout_head(Some(cb.force()));

            let re = Regex::new(r#"^(\w+) (\d{1,6}) "?"?([^"]*)?"?"#).unwrap();
            for entry in WalkDir::new(path) {
                let entry = entry.unwrap();
                if !entry.file_type().is_file() {
                    continue;
                }
                let path = entry.path();
                let text = std::fs::read_to_string(path).unwrap();
                for m in re.find_iter(&text) {
                    println!("{}:{}: {:?}", path.display(), m.start(), m.as_str());
                }
            }

            println!("checked out to {}", branch.name().unwrap().unwrap());
        }

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

impl GitUtilizedObjectProvider {
    fn open_repo(
        builder: &mut git2::build::RepoBuilder,
        path: &Path,
        url: &str,
    ) -> Result<Repository, git2::Error> {
        let name = url.strip_prefix("https://").unwrap().replace("/", "-");

        match builder.clone(url, path) {
            Ok(v) => Ok::<git2::Repository, git2::Error>(v),
            Err(e) => {
                if e.code() == git2::ErrorCode::Exists && e.class() == git2::ErrorClass::Invalid {
                    println!(
                        "A repository already exists at {}. Opening instead",
                        path.to_str().unwrap()
                    );
                    return Repository::open(path);
                }
                return Err(e);
            }
        };
        Err(git2::Error::new(
            git2::ErrorCode::GenericError,
            git2::ErrorClass::None,
            "Unexpected error occurred",
        ))
    }
}
