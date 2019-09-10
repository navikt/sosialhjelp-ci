#!/usr/bin/env bash

usage() {
    echo "$0 -p PROJECT -i RELEASE_ID -n ASSET_NAME -f FILE_PATH -a AUTH_TOKEN [-m MIME_TYPE]"
}

while getopts 'a:f:hi:n:p:' option; do
    case "${option}" in
        a) AUTH_TOKEN=${OPTARG};;
        f) FILE_PATH=${OPTARG};;
        h) usage; exit;;
        i) RELEASE_ID=${OPTARG};;
        m) MIME_TYPE=${OPTARG};;
        n) ASSET_NAME=${OPTARG};;
        p) PROJECT=${OPTARG};;
    esac
done

# FIXME: Better error handling
if [[ -z ${AUTH_TOKEN} || -z ${FILE_PATH} || -z ${RELEASE_ID} || -z ${ASSET_NAME} || -z ${PROJECT} ]]; then
    usage; exit 1;
fi

if [[ -z ${MIME_TYPE} ]]; then
    MIME_TYPE="application/octet-stream"
fi

curl -X POST "https://api.github.com/repos/navikt/$PROJECT/releases/$RELEASE_ID/assets?name=$ASSET_NAME" \
    -H "Accept: application/vnd.github.v3.full+json" \
    -H "Authorization: token $AUTH_TOKEN" \
    -H "Content-Type: application/tar+gzip" \
    -T "$FILE_PATH"
