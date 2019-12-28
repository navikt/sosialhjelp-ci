#!/usr/bin/env sh

if [ "${INPUT_ASSERT_MASTER}" = "true" ]; then
  git clone -qb "${INPUT_TAG}" "https://github.com/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}"
  if [ "$(git merge-base origin/master HEAD)" != "$(git rev-parse HEAD)" ]; then
    echo "${INPUT_TAG} has not been merged into master" 1>&2
    exit 1
  fi
else
  # If no assert_master, only single branch is needed
  git clone -qb "${INPUT_TAG}" --single-branch "https://github.com/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}"
fi

# Add .version to json file if "input_vars" is set
if [ -f "${INPUT_VARS}" ]; then
  TEMP_FILE="/tmp/input_vars.json"
  jq ".version = \"${INPUT_TAG}\"" "${INPUT_VARS}" > ${TEMP_FILE}
  if [ -s ${TEMP_FILE} ]; then
    mv ${TEMP_FILE} "${INPUT_VARS}"
  fi
fi
