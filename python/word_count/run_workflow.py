import json
from collections import defaultdict
import re
from multiprocessing import Pool
from word_count.word_count import word_count


def split_text(input_text, n_parts):
    words = input_text.split()
    total_words = len(words)
    part_size = total_words // n_parts
    parts = []
    current_index = 0

    for i in range(n_parts):
        # Start index for this part
        start = current_index

        # Calculate the end index for this part
        # We add part_size, but make adjustments if it's the last part
        if i == n_parts - 1:
            end = total_words
        else:
            # Try to make the parts as even as possible
            end = start + part_size
            # If not the last part, adjust end to prevent splitting inside a word
            # This makes sure we handle cases where words are not evenly divisible by n_parts
            while end < total_words and i != n_parts - 1:
                # Extend end index if it does not exceed total_words and the next char is not a space
                if (end + 1 < total_words) and (input_text[end] != ' ' and input_text[end + 1] != ' '):
                    end += 1
                else:
                    break

        # Append the current part to the parts list
        parts.append(' '.join(words[start:end]))

        # Update current_index to new start for next part
        current_index = end

    return parts


def combine_counts(count_dicts):
    final_counts = defaultdict(int)
    for d in count_dicts:
        for word, count in d.items():
            final_counts[word] += count
    return json.dumps(final_counts, indent=4)

def process_text(part):
    return word_count(part)

def run_workflow(input_text):
    # Degree of parallelism
    parallelism = 4  # You can adjust this based on your system capabilities

    # Split text into parts that won't break words
    parts = split_text(input_text, parallelism)

    # Create a pool of workers to process each part
    with Pool(parallelism) as pool:
        results = pool.map(process_text, parts)

    # Combine results from all parts
    combined_result = combine_counts(results)

    # Output the combined results
    print(combined_result)


if __name__ == '__main__':
    # 从 input.txt 文件中读取文本
    with open('input.txt', 'r', encoding='utf-8') as f:
        input_text = f.read()
    run_workflow(input_text)
