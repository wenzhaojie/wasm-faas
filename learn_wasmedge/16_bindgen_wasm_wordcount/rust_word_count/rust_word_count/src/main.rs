use std::collections::HashMap;

fn main() {
    let text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. \
                Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. \
                Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi \
                ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit \
                in voluptate velit esse cillum dolore eu fugiat nulla pariatur. \
                Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia \
                deserunt mollit anim id est laborum.";

    let num_threads = 4;
    let chunks = split_text_into_chunks(text, num_threads);

    let results: Vec<String> = chunks
        .into_iter()
        .map(|chunk| count_words(chunk))
        .collect();

    let final_result = merge_results(results);

    // Deserialize the final result string into HashMap<String, usize>
    let final_result_map: HashMap<String, usize> = serde_json::from_str(&final_result).unwrap();

    println!("{:?}", final_result_map);
}

fn split_text_into_chunks(text: &str, num_chunks: usize) -> Vec<String> {
    let mut chunks = Vec::with_capacity(num_chunks);
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
        chunks.push(chunk);
    }

    chunks
}

fn count_words(chunk: String) -> String {
    let mut word_count = HashMap::new();
    for word in chunk.split_whitespace() {
        *word_count.entry(word.to_string()).or_insert(0) += 1;
    }
    serde_json::to_string(&word_count).unwrap()
}

fn merge_results(results: Vec<String>) -> String {
    let mut merged_result = HashMap::new();
    for result in results {
        let word_count: HashMap<String, usize> = serde_json::from_str(&result).unwrap();
        for (word, count) in word_count {
            *merged_result.entry(word).or_insert(0) += count;
        }
    }
    serde_json::to_string(&merged_result).unwrap()
}
