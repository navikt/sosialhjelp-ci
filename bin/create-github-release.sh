#!/usr/bin/env bash

usage() {
    echo "$0 -p PROJECT -c COMMIT -t TAG -m MESSAGE -a AUTH_TOKEN"
}

while getopts 'a:c:hm:p:t:' option; do
    case "${option}" in
        a) AUTH_TOKEN=${OPTARG};;
        c) COMMIT=${OPTARG};;
        h) usage; exit;;
        m) MESSAGE=${OPTARG};;
        p) PROJECT=${OPTARG};;
        t) TAG=${OPTARG};;
    esac
done

# FIXME: Better error handling
if [[ -z ${AUTH_TOKEN} || -z ${COMMIT} || -z ${MESSAGE} || -z ${PROJECT} || -z ${TAG} ]]; then
    usage; exit 1;
fi

RESPONSE=$(curl -sIX GET "https://api.github.com/repos/navikt/$PROJECT/releases/tags/$TAG" \
    -H "Accept: application/vnd.github.v3.full+json" \
    -H "Authorization: token $AUTH_TOKEN")

if [[ ${RESPONSE} == *"Status: 200 OK"* ]]; then
    exit # The release already exists
fi

curl -so /dev/null -X POST "https://api.github.com/repos/navikt/$PROJECT/git/tags" \
    -H "Accept: application/vnd.github.v3.full+json" \
    -H "Authorization: token $AUTH_TOKEN" \
    -H "Content-Type: application/json" \
    --data '{"tag": "'${TAG}'", "message": "'${TAG}'", "object": "'${COMMIT}'", "type": "commit"}'

RESPONSE=$(curl -sX POST "https://api.github.com/repos/navikt/$PROJECT/releases" \
    -H "Accept: application/vnd.github.v3.full+json" \
    -H "Authorization: token $AUTH_TOKEN" \
    -H "Content-Type: application/json" \
    --data '{"tag_name": "'${TAG}'", "target_commitish": "'${COMMIT}'", "name": "'${TAG}'", "body": "'"$MESSAGE"'"}')

if [[ ${RESPONSE} == *"\"id\""* ]]; then
    echo ${RESPONSE} | jq ".id"
fi
