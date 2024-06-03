use std::collections::HashMap;

fn main() {
    println!("Starting main function...");

    let text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. \
                Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. \
                Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi \
                ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit \
                in voluptate velit esse cillum dolore eu fugiat nulla pariatur. \
                Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia \
                deserunt mollit anim id est laborum.".to_string();

    let num_threads = 4;
    let chunks = split_text_into_chunks(text, num_threads);
    println!("Text split into chunks: {:?}", chunks);

    // count_words

    let merge_results = merge_results(chunks);
    println!("Results merged: {}", merge_results);

    // Deserialize the final result string into HashMap<String, i32>
    let final_result_map: HashMap<String, i32> = serde_json::from_str(&merge_results).unwrap();

    println!("Final result map: {:?}", final_result_map);
    println!("Exiting main function...");
}

fn split_text_into_chunks(text: String, num_chunks: usize) -> String {
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


fn count_words(chunk: &str) -> String {
    println!("Starting count_words function...");
    let mut word_count = HashMap::new();
    for word in chunk.split_whitespace() {
        *word_count.entry(word.to_string()).or_insert(0) += 1;
    }
    let serialized_word_count = serde_json::to_string(&word_count).unwrap();
    serialized_word_count
}


fn merge_results(results: String) -> String {
    println!("Starting merge_results function...");
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
