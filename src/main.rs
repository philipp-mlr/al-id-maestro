mod object_manager;
mod object_type;

fn main() {}

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
// fn initalize_id_range(start: u32, end: u32) -> BTreeMap<u32, bool> {
//     let mut id_range: BTreeMap<u32, bool> = BTreeMap::new();
//     for i in start..end + 1 {
//         id_range.insert(i, false);
//     }
//     id_range
// }
//
// fn print_range(id_range: &BTreeMap<u32, bool>) {
//     for (key, value) in id_range {
//         println!("{key} {value}");
//     }
// }
//
// async fn run() {
//     let app = Router::new().route("/new/{object_type}", get(new_object));
//
//     let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
//
//     axum::serve(listener, app).await.unwrap();
// }
//
