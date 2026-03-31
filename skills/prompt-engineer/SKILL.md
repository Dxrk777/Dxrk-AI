---
name: prompt-engineer
description: >
  Build and optimize system prompts using patterns from the world's best AI models.
  Use when creating system prompts, agent instructions, or optimizing AI behavior.
author: dxrk
argument-hint: "[task]"
allowed-tools: Read, Grep, Glob
---

# Prompt Engineer

## Context
You have access to 128 system prompts from the world's leading AI models:
- **Anthropic** (29 prompts): Claude's safety rules, tool use, identity
- **OpenAI** (58 prompts): ChatGPT, GPT-4, plugins, vision, code interpreter
- **Google** (17 prompts): Gemini, Bard, search integration
- **xAI** (8 prompts): Grok personas, safety instructions
- **Perplexity** (2 prompts): Search-augmented generation
- **Misc** (13 prompts): Various AI systems

All stored in `prompts/system-prompts-library/`.

## Task
$ARGUMENTS

## How to Build System Prompts

### 1. Research
Read relevant prompts from the library:
- Search for similar use cases: `grep -r "keyword" prompts/system-prompts-library/`
- Study how top models handle: identity, safety, tools, refusal, output format

### 2. Structure (from best practices across all models)
```
## Identity
Who the AI is, what it does, its personality.

## Capabilities
What tools/functions it can use.

## Rules
Safety rules, content policies, behavioral constraints.

## Output Format
How to structure responses (markdown, JSON, etc.)

## Examples
Few-shot examples of good behavior.

## Error Handling
What to do when things go wrong.
```

### 3. Key Patterns Found in Library
- **Safety layers**: Most models have 3+ layers of safety rules
- **Tool use**: Structured function calling with error handling
- **Refusal behavior**: Polite refusal with explanation, not just "I can't do that"
- **Output formatting**: Consistent markdown/JSON structure
- **Context window**: Explicit instructions about what to remember
- **Persona consistency**: Single voice throughout all interactions

### 4. Anti-Patterns (from leaked prompts)
- Overly long prompts that reduce effectiveness
- Contradictory rules between sections
- Missing error handling for edge cases
- No fallback behavior when tools fail

## Rules
- Always base prompts on patterns from the library
- Keep prompts under 2000 tokens when possible
- Test prompts with edge cases before deploying
- Version your prompts for A/B testing
