#!/usr/bin/env python

# This script will parse all the MD files stored under test/test plans
# and build a json structure with all the information found in those
# markdown files. Nesting the entries based on the Headers (H1, H2, H3).
# USAGE: python test2json.py

import os
import re
import json

def parse_markdown(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        content = file.read()

    data = {}
    current_test_suite = None
    current_test_case = None

    lines = content.split('\n')
    for line in lines:
        if line.startswith('# '):
            data['name'] = line[2:].strip()
            data['Test Suite'] = []
        elif line.startswith('**UC ID**: '):
            data['ID'] = line[len('**UC ID**: '):].strip()
        elif line.startswith('**Description**: '):
            data['Description'] = line[len('**Description**: '):].strip()
        elif line.startswith('## '):
            if 'Review' in line:
                data['Test Suite'].append(current_test_suite)
            else:
                current_test_suite = {'Test Cases': []}
        elif line.startswith('**ID**: '):
            if current_test_suite:
                current_test_suite['ID'] = line[len('**ID**: '):].strip()
        # We don't want to add a test case entry when finding ### Test Cases
        elif line.startswith('### ') and not line.startswith('### Test Cases'):
            current_test_case = {
                'name': line[4:].strip(),
                'Test Type': None,
                'Test Description': None
            }
        elif line.startswith('**Test Type**: '):
            if current_test_case:
                if ',' in line:
                    current_test_case['Test Type'] = []
                    for test_type in line[len('**Test Type**: '):].strip().split(","):
                        current_test_case['Test Type'].append(test_type.strip())
                else:
                    current_test_case['Test Type'] = line[len('**Test Type**: '):].strip()
        elif line.startswith('**Test Description**: '):
            if current_test_case:
                current_test_case['Test Description'] = line[len('**Test Description**: '):].strip()
                current_test_suite['Test Cases'].append(current_test_case)
    return data

def process_files():
    json_data = []

    for file_name in os.listdir():
        if file_name.endswith('.md'):
            file_path = os.path.join(os.getcwd(), file_name)
            json_data.append(parse_markdown(file_path))

    return json_data

if __name__ == "__main__":
    result = process_files()

    # Output the result as a JSON string
    print(json.dumps(result, indent=4))