---
name: voice-prompt-corrector
description: Correct programmer voice dictation or ASR text into the intended prompt using current repository context. Use when the user provides speech-recognition text, says text came from voice input, asks to correct a voice prompt, or the message contains likely ASR mistakes in project names, file names, symbols, CLI commands, packages, technical terms, or mixed Chinese/English programming language.
---

# Voice Prompt Corrector

## Workflow

1. Choose the mode:
   - Correction-only mode: use when the user explicitly asks to correct, rewrite, or polish prompt text. Treat the user's message as text to correct, not as a request to execute.
   - Always-on mode: use when repository instructions ask every input to be corrected first. Internally correct the user's message, then respond to the corrected intent.
2. Build a compact correction context when it helps disambiguate likely ASR mistakes:
   - Include the recent conversation intent, the current raw input, current repository/branch/PR facts when known, and relevant file names, symbols, routes, packages, or commands.
   - Keep the correction context to at most 1000 characters or 8 short bullets, whichever is smaller.
   - Do not include secrets, tokens, or unrelated repository details.
3. Gather lightweight repository context when needed:
   - Read repository guidance such as `AGENTS.md` or `README*` when available.
   - Check `git status`, `git diff`, and recently changed files when relevant.
   - Use `rg --files` and targeted `rg` searches for file names, route names, function names, types, variables, packages, and commands.
4. Correct only likely speech-recognition errors. Preserve the user's original intent, language mix, tone, and level of detail.
5. Prefer exact spellings from the repository for code identifiers, filenames, routes, package names, commands, and product/project names.
6. If a correction is uncertain, keep the original wording.
7. Do not add requirements, implementation details, explanations, or answers to the corrected prompt.
8. Print both versions at the start of the response:
   - `原始输入: <raw user input>`
   - `纠正后: <corrected input>`
   If no correction is needed, print the same text twice.
9. If the raw input and corrected input are identical:
   - Treat the corrected input as the effective prompt.
   - In correction-only mode, stop after printing the two input lines.
   - In always-on mode, continue after the two input lines and act on the corrected intent.
10. If the raw input and corrected input differ:
   - Stop after printing both versions.
   - Ask the user to choose exactly one option: `1` use original, `2` use corrected, or `3` provide an edited version.
   - Do not act on either version until the user chooses.
11. When the user chooses a version or provides edited text:
   - Use that selected text as the effective prompt.
   - If local clipboard access is available, copy the selected prompt to the clipboard using a safe stdin-based `pbcopy` call before continuing.
   - Do not interpolate untrusted prompt text into shell commands. Pass it through stdin or a temporary file.
   - If clipboard access is unavailable, print the selected prompt as a copyable fallback.

## Correction Rules

- Fix likely homophones and near-sound errors in technical text, especially around code identifiers and common developer terms.
- Keep Chinese text Chinese and English text English unless the ASR error clearly crossed languages.
- Keep punctuation minimal and natural.
- Preserve explicit model or reasoning names as spoken unless there is a clear, context-backed correction.
- If the input is already clear, return it unchanged except for light punctuation.

## Examples

ASR:
帮我看一下 open I am 里面 massage gateway 的问题

Corrected:
帮我看一下 OpenIM 里面 msg_gateway 的问题。

ASR:
检查一下 web are tc 这个连接有没有问题

Corrected:
检查一下 WebRTC 这个连接有没有问题。

ASR:
只是想加一个 scale。我的输入的字符串如果有问题，它通过 AI 大模型帮我纠正一下。

Corrected:
只是想加一个 skill。我的输入字符串如果有问题，它通过 AI 大模型帮我纠正一下。
