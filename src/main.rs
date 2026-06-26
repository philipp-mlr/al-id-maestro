use std::collections::{BTreeMap, HashMap};

use axum::{Router, routing::post};

#[derive(Debug, Hash, Eq, PartialEq)]
enum Type {
    Codeunit,
    // Table,
    // TableExtension,
    // Page,
    // PageExtension,
    // Report,
    // ReportExtension,
    // Enum,
    // EnumExtension,
    // PermissionSet,
    // PermissionSetExtension,
    // MenuSuite,
    // Query,
    // XMLPort,
    // Unknown,
}

fn main() {
    // let mut object_ranges = HashMap::new();
    //
    // let start = 50000u32;
    // let end = 99999u32;
    //
    // object_ranges.insert(Type::Codeunit, initalize_id_range(start, end));
    //
    // match object_ranges.get(&Type::Codeunit) {
    //     Some(range) => {
    //         println!("Yes!");
    //         print_range(range);
    //     }
    //     None => println!("Something strange happened"),
    // }

    run();
}

#[tokio::main]
async fn run() {
    let app = Router::new().route("/new", post(new_handler));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();

    axum::serve(listener, app).await.unwrap();
}

async fn new_handler() {
    print!("Got post");
}

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
