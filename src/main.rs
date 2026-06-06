use std::collections::{BTreeMap, HashMap};

struct Object {
    id: u32,
    used: bool,
}

#[derive(Debug, Hash, Eq, PartialEq)]
enum Type {
    Codeunit,
    Table,
}

fn main() {
    let mut objects = HashMap::new();
    let tree: BTreeMap<u32, bool> = BTreeMap::new();
    objects.insert(Type::Codeunit, tree);
    test(1);
}

fn test(i: i32) -> i32 {
    assert!(i != 0);
    i
}
