# Sift

## Plan

We are going to build a prioritizer app.

In this app, we will gather all the tasks from Things.app, and then the tool
will help prioritize them.

The app will do this by comparing two tasks at a time, and asking the user
which one is more important.

Then the user says that one task is more important, that task will become the
parent of the other task.

In this way, we will build an directed acyclic graph (DAG) of the tasks.

When the tasks are first submitted, we won't know which one is more important.
However, we will know the first priority task once it becomes the root of all
the other tasks, and there are no other single tasks that are disconnected from
the root.

We can continue doing this process until we have a single chain of tasks.

If we have three tasks, we know that we have three minus one comparisons to
make to determine the top priority. For a given level in the tree, we will
always have n - 1 comparisons to make to determine the next priority, where n
is the number of tasks in the level.

If we query Things and find new tasks after having already compared other
tasks, the new tasks will be added as independent roots, and therefore will be
compared to the task or tasks at the same level, that is, at the root level.

In this way, a new task can immediately become the root of the tree, if it is
more important that the existing root.

Tasks are fully prioritized when they are the only task at their level, and
every ancestor task above them is also the only task at their level. Tasks are
not fully prioritized when they are not the only task at their level, or if any
ancestors above them are not the only task at their level. That is true at any
level, including the root level. When a task is fully prioritized, it should be
styled differently, showing that it will no longer be compared to other tasks.

The list of tasks should be sorted by level. Within a level, the order of tasks
does not matter.

## Commands

- **Dev**: `go run .` (run app)
- **Build**: `go build -o /dev/null`
- **Lint**: `golangci-lint run`
- **Format**: `gofumpt -w .`
- **Test**: `go test`

## Tech stack

- Golang
- Bubbletea (for the UI)
- Lipgloss (for styling text and elements)

## Implementation

- Implement features in the smallest possible amount of code.
- Always attempt to use existing code to solve a problem.
- If existing code is not available, think deeply about how we might refactor
it to make it more resusable.
- Keep files as small and focused as possible.
- Keep functions as small and focused as possible.
- Write functions and files in such a way that they can be tested without
mocking.
- Always use a test-driven-development (TDD) flow.
- Always write tests first before implementing any functionality.
- Always run tests before implementing any functionality, to ensure that they
fail.
- When finished with work, always format code, then lint, then test, then
build, in this order.

## Misc

- Confirm you've read this document by starting responses with "I have read and
understand the instructions."
