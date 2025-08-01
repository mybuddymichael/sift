# Sift

Prioritize your tasks from Things.app, in the terminal.

![Screenshot of the program in use, in light mode, showing several tasks prioritized, completed, or yet to be prioritized](https://r1vysk5peykhs5gu.public.blob.vercel-storage.com/sift-light-2-n48jrtBYV9W9scSrQWmAv1n947NiEH.png)
![Screenshot of the program in use, in dark mode, showing several tasks prioritized, completed, or yet to be prioritized](https://r1vysk5peykhs5gu.public.blob.vercel-storage.com/sift-dark-2-3mEHc52mvAixv4pmD9ffJuFqhlXies.png)

## How to use

1. Install with Homebrew: `brew install mybuddymichael/tap/sift-things`
2. Run the command: `sift`
3. Use the arrow keys to start prioritizing tasks.
4. Reset all priorities with `ctrl+r`.
5. Quit with `ctrl+c`.

### Options

- `--refresh-interval <seconds>`: Set the refresh interval for getting updates from Things.app (default: 3 seconds)

## How it works

- Sift requires Things.app to be installed and running on your Mac.
- It displays tasks in the Today list, and will poll Things for updates every 3
seconds by default (configurable with `--refresh-interval`).
- Sift does not write any data to Things. It only stores parent-child
relationships between tasks.
- Priorities persist across Sift and Things restarts.

## The sorting method

> [!NOTE]
> This is an in-the-weeds description of the sorting method.

- When a task is chosen, the one *not* chosen is updated to note that its parent is the task that was chosen.
  - In this way, we create a tree of tasks, where each task tracks its parent.
- In order to pick the two tasks being compared, we gather all of the tasks and assign them levels.
  - Tasks with no parents are at the highest level, their children are at the next level, and so on.
  - We choose tasks to compare by finding the highest level with multiple tasks.
- If a level only has one task, and all its ancestors are the only tasks at each of their levels, then we know the task is fully prioritized.
- The technical term for this structure is a directed acyclic graph (DAG).
- There will never be any cycles in the DAG, which means that we assume that if task C is a child of task B, and task B is a child of task A, then task C must be lower in priority than task A.

## Tech stack

- Go
- [Bubbletea](https://github.com/charmbracelet/bubbletea) (TUI framework)

## Prior art

- [Todournament](https://github.com/alltom/todournament) by [Tom Lieber](https://github.com/alltom)
