use std::collections::HashMap;
use serde_json;
use wasmedge_bindgen_macro::*;
#[allow(unused_imports)]
use wasmedge_bindgen::*;


#[wasmedge_bindgen]
pub fn word_count(input: String) -> String {
    let mut counts = HashMap::new();
    let words = input
        .split_whitespace()
        .map(|word| word.to_lowercase())
        .map(|word| word.chars().filter(|c| c.is_alphanumeric()).collect::<String>());

    for word in words {
        if !word.is_empty() {
            *counts.entry(word).or_insert(0) += 1;
        }
    }
    // 用serde_json::to_string_pretty(&counts).unwrap()可以将HashMap转换为json字符串
    let json_str = serde_json::to_string_pretty(&counts).unwrap();
    json_str
}

