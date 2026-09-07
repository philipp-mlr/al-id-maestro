use crate::{
    object_manager::ObjectManager,
    object_type::ObjectType,
    utilized_object_provider::{git::GitUtilizedObjectProvider, mock::MockUtilizedObjectProvider},
};

mod id_state;
mod object_manager;
mod object_tree;
mod object_type;
mod utilized_object_provider;

fn main() {
    let mut object_manager = ObjectManager::new(GitUtilizedObjectProvider);
    new(&mut object_manager);
    new(&mut object_manager);
    new(&mut object_manager);
    new(&mut object_manager);
}

fn new(object_manager: &mut ObjectManager) {
    match object_manager.get_free_id_for_type(object_type::ObjectType::Codeunit) {
        Ok(i) => println!("got id {} for object {}", i, ObjectType::Codeunit),
        Err(e) => eprintln!("{}", e),
    }
}

// use axum::{
//     Router,
//     extract::{Path, Query},
//     http::{HeaderMap, Response, StatusCode},
//     response::Json,
//     routing::{get, post},
// };
// use serde_json::{Value, json};

// #[tokio::main]
// async fn main() {
//     let mut object_ranges = HashMap::new();
//
//     let start = 50000u32;
//     let end = 99999u32;
//
//     object_ranges.insert(Type::Codeunit, initalize_id_range(start, end));
//
//     match object_rangbes.get(&Type::Codeunit) {
//         Some(range) => {
//             println!("Yes!");
//             print_range(range);
//         }
//         None => println!("Something strange happened"),
//     }
//
//     println!(
//         "{}",
//         object_ranges
//             .get(&Type::Codeunit)
//             .unwrap()
//             .get(&50000)
//             .unwrap()
//     );
//
//     run().await;
// }
//
//
//
// async fn run() {
//     let app = Router::new().route("/new/{object_type}", get(new_object));
//
//     let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
//
//     axum::serve(listener, app).await.unwrap();
// }
//
