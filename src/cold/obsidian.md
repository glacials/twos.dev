---
date: 2023-07-09
filename: obsidian.html
type: post
---

# Obsidian Plugins

I run many Obsidian plugins and they're constantly changing, but below are the ones that
have stuck around.

- [Copy as HTML](https://github.com/jenningsb2/copy-as-html)
  - I treat this as “copy with formatting” so I can copy e.g. bulleted lists from my notes and paste them into an email to send to others; without it, the copied text is unparsed markdown
- [Editor Syntax Highlight](https://github.com/deathau/cm-editor-syntax-highlight-obsidian)
  - Syntax highlighting for fenced code blocks, just like GitHub
- [Excalidraw](https://github.com/zsviczian/obsidian-excalidraw-plugin)
  - Unbelievably well-integrated Excalidraw support; embed drawings in notes but also let specific elements in drawings link to notes or dynamically inherit content from notes, and a lot more
- [Global Search and Replace](https://github.com/MahmoudFawzyKhalil/obsidian-global-search-and-replace)
  - Does what it says on the tin; I expect this will become a first-party feature soon
- [Linter](https://github.com/platers/obsidian-linter)
  - Automatically apply Markdown linting on save
  - Cleans up inconsistencies like # line breaks between paragraphs, heading capitalization, re-index ordered lists, etc.
- [Numerals](https://github.com/gtg922r/obsidian-numerals)
  - Adds a code block type for math, including calculating and rendering answers for each line
  - Basically Jupyter Notebooks for math; easy to show work and avoid mistakes since the answers are calculated on render
- [QuickShare](https://github.com/mcndt/obsidian-quickshare)
  - Publish a note to an auto-expiring shareable URL
- [Smart Links](https://github.com/kemayo/obsidian-smart-links)
  - Allows specifying patterns to automatically link
  - I use it so I can write e.g. OPSFLOW-123 to link to a Jira ticket
- I use the **Daily notes** core plugin, which adds a button/command to create a note for today with a predesignated place and title (e.g. `Daily/YYYY-MM-DD.md`). These community plugins extend that:
  - [Daily Note Outline](https://github.com/iiz00/obsidian-daily-note-outline)
    - A single pane that shows an outline of _all_ daily notes in reverse chronological order
    - I have this persistently in my right sidebar
  - [Daily Notes Opener](https://github.com/reorx/obsidian-daily-notes-opener)
    - Auto open today’s note when opening Obsidian
  - [Obsidian TODO](https://github.com/larslockefeer/obsidian-plugin-todo)
    - Collects TODOs from all notes (daily or not) and shows them in one pane
    - This is also persistently in my right sidebar
    - These are markdown todos, which Obsidian natively supports (e.g. `- [ ] Go shopping`); for a more powerful experience look at [Tasks](https://github.com/obsidian-tasks-group/obsidian-tasks)
  - [Natural Language Dates](https://github.com/argenos/nldates-obsidian)
    - Type e.g. `@yesterday` or `@tomorrow` to link to the daily note for those days
  - [Review](https://github.com/ryanjamurphy/review-obsidian)
    - Mark the current spot (any note, any header) as needing review at whatever date you want (e.g. tomorrow); it will be added to that day's daily note
    - For a different workflow see [Reminder](https://github.com/uphy/obsidian-reminder)
