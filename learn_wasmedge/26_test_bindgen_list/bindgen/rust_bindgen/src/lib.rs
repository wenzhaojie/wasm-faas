#[allow(unused_imports)]
#[allow(dead_code)]
use wasmedge_bindgen::*;
use wasmedge_bindgen_macro::*;


fn main() {
    let text = String::from("Hello, world! This is a test string.");
    let num_parts = 5;

    let parts = split_text(text, num_parts);

    for (index, part) in parts.iter().enumerate() {
        println!("Part {}: '{}'", index + 1, part);
    }
}

#[wasmedge_bindgen] // 使用 wasmedge_bindgen 宏
pub fn split_text(text: String, num: usize) -> Vec<String> {
    let part_size = text.chars().count() / num;
    let mut extra = text.chars().count() % num;

    let mut result: Vec<String> = Vec::new();
    let mut chars = text.chars();

    for _ in 0..num {
        let mut part = String::new();
        for _ in 0..part_size {
            if let Some(c) = chars.next() {
                part.push(c);
            }
        }

        if extra > 0 {
            if let Some(c) = chars.next() {
                part.push(c);
            }
            extra -= 1;
        }

        result.push(part);
    }

    result
}
