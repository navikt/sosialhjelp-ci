#!/usr/bin/env bash

usage() {
    echo "$0 -p PROJECT -c COMMIT -t TAG -m MESSAGE -a AUTH_TOKEN"
}

while getopts 'a:f:hi:n:p:' option; do
    case "${option}" in
        a) AUTH_TOKEN=${OPTARG};;
        f) FILE_NAME=${OPTARG};;
        h) usage; exit;;
        i) RELEASE_ID=${OPTARG};;
        n) ASSET_NAME=${OPTARG};;
        p) PROJECT=${OPTARG};;
    esac
done

# FIXME: Better error handling
if [[ -z ${AUTH_TOKEN} || -z ${FILE_NAME} || -z ${RELEASE_ID} || -z ${ASSET_NAME} || -z ${PROJECT} ]]; then
    usage; exit 1;
fi

curl -X POST "https://api.github.com/repos/navikt/$PROJECT/releases/$RELEASE_ID/assets?name=$ASSET_NAME" \
    -H "Accept: application/vnd.github.v3.full+json" \
    -H "Authorization: token $AUTH_TOKEN" \
    -H "Content-Type: application/tar+gzip" \
    -T "$FILE_NAME"
