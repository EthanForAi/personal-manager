# Paper Personal Manager Prototype

This directory contains a static PaperCSS design prototype for the personal-manager frontend.

Open `index.html` directly in a browser:

```sh
open /Users/wenxu/.codex/worktrees/bab5/personal-manager/design/paper-personal-manager/index.html
```

The prototype uses PaperCSS from the official unpkg CDN and does not require a build step.

Static screenshots are included:

- `screenshots/desktop.png`
- `screenshots/mobile.png`

## Scope

- Visual design first, backend integration later.
- Current backend fields are represented: `userid`, `name`, `email`, and `phone`.
- Current API routes are represented: `/create`, `/read`, `/update`, `/delete`, and `/check`.
- Prototype interactions use local draft data and mirror the backend's main validation messages.

## Implementation Notes

The next coding step can reuse the same files and replace the prototype submit handler with `fetch()` calls:

```js
await fetch(endpoint, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(payload),
});
```

Codex can turn this prototype into production frontend code directly. If a Figma-style artifact is needed first, the practical path is to export screenshots from this prototype or use the Figma plugin to recreate the same layout as frames.
