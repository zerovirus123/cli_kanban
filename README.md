# CLI_Kanban

![example](https://github.com/zerovirus123/cli_kanban/blob/master/public/demo.png?raw=true)

A simple command line interface based Kanban board built using BubbleTea.
BubbleTea is a framework for building terminal apps. 

Tasks can be moved to the next column. The app also supports task addition and deletion.
Persistence is achieved with the `storage.json` file.

# Usage

To run the program, type `go run .`.

Use the `left/right` or `A/D` keys to move between lists. 

Use the up/down buttons to move between tasks.

Press `Enter` to move a selected task to the next column.

Use the `n` button to open a task form.

Use the `x` button to delete the focused task.

Press `q` to quit the kanban board.