#!/bin/bash

TODO_CLI_CHECK_TIME=0

__todo_cli_ps1() {
    local __todo_cli_elapsed_seconds=$((SECONDS - TODO_CLI_CHECK_TIME))
    if [[ "$TODO_CLI_CHECK_TIME" -eq "0" || "$__todo_cli_elapsed_seconds" -gt "300" ]]; then
        TODO_CLI_CHECK_TIME=$SECONDS
        todo-cli notify
    fi
}

PROMPT_COMMAND="$PROMPT_COMMAND __todo_cli_ps1"
