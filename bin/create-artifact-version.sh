#!/usr/bin/env sh

if [ "$(git rev-parse --is-inside-work-tree 2> /dev/null)" = "true" ]; then
    GIT_COMMIT_HASH=$(git log -n 1 --pretty=format:'%h')
    GIT_COMMIT_DATE=$(git log -1 --pretty='%ad' --date=format:'%Y%m%d.%H%M')
    echo "1.0_${GIT_COMMIT_DATE}_${GIT_COMMIT_HASH}"
fi
