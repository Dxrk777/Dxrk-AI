---
name: codebase-visualizer
description: >
  Generate an interactive collapsible tree visualization of the codebase. Use when exploring a new repo, understanding project structure, or identifying large files.
allowed-tools: Bash(python *)
argument-hint: "[directory]"
---

# Codebase Visualizer

Generate an interactive HTML tree view that shows the project's file structure with collapsible directories.

## Usage

Run the visualization script:
```bash
python ${CLAUDE_SKILL_DIR}/scripts/visualize.py $ARGUMENTS
```

This creates `dxrk-codebase-map.html` in the current directory.

## What It Shows
- **Collapsible directories**: Click to expand/collapse
- **File sizes**: Displayed next to each file
- **Colors**: Different colors for file types (.go, .md, .json, etc.)
- **Summary**: File count, directory count, total size, file type breakdown

## After Generation
Open `dxrk-codebase-map.html` in a browser to explore.
