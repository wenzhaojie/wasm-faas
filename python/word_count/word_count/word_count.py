import json
from collections import defaultdict
import re


def word_count(input_string):
    # Create a dictionary to store word counts
    counts = defaultdict(int)

    # Normalize and filter words
    words = re.findall(r'\b\w+\b', input_string.lower())

    # Count each word
    for word in words:
        counts[word] += 1

    return counts


if __name__ == '__main__':
    # Example usage
    input_text = "Hello world! hello."
    print(word_count(input_text))
