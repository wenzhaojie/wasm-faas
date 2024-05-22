use std::collections::HashMap;

#[no_mangle]
pub fn split_text_into_chunks(text: &str, num_chunks: usize) -> String {
    let mut chunks = String::new();
    let words: Vec<&str> = text.split_whitespace().collect();
    let chunk_size = words.len() / num_chunks;

    for i in 0..num_chunks {
        let start = i * chunk_size;
        let end = if i == num_chunks - 1 {
            words.len()
        } else {
            (i + 1) * chunk_size
        };

        let chunk = words[start..end].join(" ");
        chunks.push_str(&count_words(&chunk));
        chunks.push('\n');
    }
    chunks
}

#[no_mangle]
pub fn count_words(chunk: &str) -> String {
    let mut word_count = HashMap::new();
    for word in chunk.split_whitespace() {
        *word_count.entry(word.to_string()).or_insert(0) += 1;
    }
    let serialized_word_count = serde_json::to_string(&word_count).unwrap();
    serialized_word_count
}

#[no_mangle]
pub fn merge_results(results: String) -> String {
    let mut merged_result = HashMap::new();
    let chunks: Vec<&str> = results.trim().split('\n').collect(); // 去除首尾空白字符，然后分割
    for chunk in chunks {
        let word_count: HashMap<String, i32> = serde_json::from_str(chunk).unwrap();
        for (word, count) in word_count {
            *merged_result.entry(word).or_insert(0) += count;
        }
    }
    let serialized_merged_result = serde_json::to_string(&merged_result).unwrap();
    serialized_merged_result
}
