import os
import re

base_dir = r"C:\Users\fzxak\Desktop\flowus_pages"

def format_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Normalize line endings
    content = content.replace('\r\n', '\n')
    
    # Process line by line to remove continuous spaces and tabs, preserving leading spaces/tabs for indentation
    lines = content.split('\n')
    new_lines = []
    for line in lines:
        if not line.strip():
            new_lines.append('')
            continue
            
        # Extract leading spaces and tabs
        match = re.match(r'^([\s\t]*)', line)
        leading = match.group(1) if match else ''
        rest = line[len(leading):]
        
        # Replace continuous spaces and tabs in the rest of the line with a single space
        rest = re.sub(r'[ \t]+', ' ', rest)
        new_lines.append(leading + rest)
        
    content = '\n'.join(new_lines)
    
    # Replace continuous newlines (3 or more) with exactly 2 newlines (\n\n)
    # Dify splits by \n\n, so we want to make sure paragraphs are separated by strictly \n\n.
    content = re.sub(r'\n{3,}', '\n\n', content)
    
    # Save the formatted file back
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

count = 0
for root, dirs, files in os.walk(base_dir):
    for file in files:
        if file.endswith('.md'):
            try:
                format_file(os.path.join(root, file))
                count += 1
            except Exception as e:
                print(f"Error formatting {file}: {e}")

print(f"Processed {count} files.")
