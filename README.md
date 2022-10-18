# todo-cli
To Do list with gist repository

[![Release](https://github.com/guionardo/todo-cli/actions/workflows/release.yml/badge.svg)](https://github.com/guionardo/todo-cli/actions/workflows/release.yml)
[![CodeQL](https://github.com/guionardo/todo-cli/actions/workflows/codeql.yml/badge.svg)](https://github.com/guionardo/todo-cli/actions/workflows/codeql.yml)
[![Go](https://github.com/guionardo/todo-cli/actions/workflows/go.yml/badge.svg)](https://github.com/guionardo/todo-cli/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/guionardo/todo-cli/branch/main/graph/badge.svg?token=SbUaUBJkzE)](https://codecov.io/gh/guionardo/todo-cli)
[![Release](https://github.com/guionardo/todo-cli/actions/workflows/release.yml/badge.svg?event=release)](https://github.com/guionardo/todo-cli/actions/workflows/release.yml)

## Commands

```bash
❯ go run . --help
NAME:
   todo-cli - A simple todo app with a cli interface and github gist persistence

USAGE:
   todo-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Guionardo Furlan <guionardo@gmail.com>

COMMANDS:
   notify, n  Notify about pending tasks
   help, h    Shows a list of commands or help for one command
   Setup:
     setup, s  Setup todo app
     sync, s   Synchronize local collection with GIST
   Tasks:
     add, a       Add a new todo item
     update, u    Update a todo item
     list, l      List all todo items
     delete, d    Delete a todo item
     complete, c  Complete a todo item
     act, a       Set current timestamp as action for item
     backup, a    Run a backup of the todo list
     init         Outputs the shell initialization script

GLOBAL OPTIONS:
   --data-folder FOLDER  Load configuration from FOLDER (default: "/home/guionardo/.config/todo-cli") [$TODO_CONFIG]
   --debug               Enable debug mode (default: false)
   --help, -h            show help (default: false)
   --version, -v         print the version (default: false)
```
