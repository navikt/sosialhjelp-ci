#!/usr/bin/env bash

# Create a newline escaped message from last commit
if [[ $(git rev-parse --is-inside-work-tree 2> /dev/null) == "true" ]]; then
    git log --format=%B -n 1 | sed ':a;N;$!ba;s/\n/\\n/g'
fi
