#!/usr/bin/env bash

# Bash completion
if [ -n "$BASH_VERSION" ]; then
    # kubectl completion
    if command -v kubectl &> /dev/null; then
        eval "$(kubectl completion bash)"
    fi

    # talosctl completion
    if command -v talosctl &> /dev/null; then
        eval "$(talosctl completion bash)"
    fi

    # wild completion
    if command -v wild &> /dev/null; then
        eval "$(wild completion bash)"
    fi
fi

source <(wild instance env)
