#!/usr/bin/env python3

import os
import subprocess
import sys

# Directory to search for templates
TEMPLATE_DIR = './Infrastructure'

# File extensions to consider as CloudFormation templates
TEMPLATE_EXTENSIONS = ['.yaml']

# List to store validation errors
failed_templates = []

def is_template_file(file_name):
    """Check if the file has a valid CloudFormation extension."""
    return any(file_name.endswith(ext) for ext in TEMPLATE_EXTENSIONS)

def validate_template_with_cfn_lint(file_path):
    """Run cfn-lint on a given CloudFormation template."""
    try:
        print(f"Linting {file_path} with cfn-lint...")

        # Run cfn-lint on the template file
        result = subprocess.run(
            ['cfn-lint', file_path],
            capture_output=True,
            text=True
        )

        # Check if cfn-lint found any issues
        if result.returncode != 0:
            print(f"Linting failed for {file_path}:\n{result.stdout}\n{result.stderr}")
            failed_templates.append(file_path)
        else:
            print(f"Template {file_path} passed linting")

    except Exception as e:
        print(f"An error occurred while linting {file_path}: {str(e)}")
        failed_templates.append(file_path)

def find_templates():
    """Find all template files recursively within TEMPLATE_DIR."""
    template_files = []
    for root, dirs, files in os.walk(TEMPLATE_DIR):
        for file in files:
            if is_template_file(file):
                template_files.append(os.path.join(root, file))
    return template_files

def main():
    # Find all CloudFormation templates
    templates = find_templates()

    if not templates:
        print("No CloudFormation templates found!")
        sys.exit(1)

    # Validate each template with cfn-lint
    for template in templates:
        validate_template_with_cfn_lint(template)

    # If any validation failed, exit with status code 1
    if failed_templates:
        print(f"\nLinting failed for the following templates: {', '.join(failed_templates)}")
        sys.exit(1)
    else:
        print("\nAll templates passed linting.")

if __name__ == '__main__':
    main()
