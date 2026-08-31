
pub async fn new_object(
    Path(object_type): Path<String>,
    headers: HeaderMap,
) -> Result<String, StatusCode> {
    let api_key = headers.get("X-API-KEY");

    match api_key {
        Some(_v) => println!("yes"),
        None => return Err(StatusCode::UNAUTHORIZED),
    };

    let type_result = Type::from_str(object_type.as_str());

    let result = match type_result {
        Ok(v) => result,
        Err(e) => eprintln!("Received incorrect object type: {:?}", e);

    }

    println!("{}", object_type);
    Ok("420".to_string())
}
