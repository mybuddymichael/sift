# Sift

Prioritize your tasks from Things.app, in the terminal.

![Screenshot of the program in use, showing several tasks prioritized, completed, or yet to be prioritized](https://r1vysk5peykhs5gu.public.blob.vercel-storage.com/sift-light-PItyiEXdxlxcuS7R5xj87Ir1b14mFN.png)

## How to use

1. Install with Homebrew: `brew install mybuddymichael/tap/sift-things`
2. Run the command: `sift`
3. Use the arrow keys to start prioritizing tasks.
4. Reset all priorities with `ctrl+r`.
5. Quit with `ctrl+c`.

## How it works

- Sift requires Things.app to be installed and running on your Mac.
- It displays tasks in the Today list, and will poll Things for updates every 5
seconds.
- Sift does not write any data to Things. It only stores parent-child
relationships between tasks.
- Priorities persist across Sift and Things restarts.

## Tech stack

- Go
- [Bubbletea](https://github.com/charmbracelet/bubbletea) (TUI framework)

## Prior art

- [Todournament](https://github.com/alltom/todournament) by [Tom Lieber](https://github.com/alltom)
