#!/usr/bin/env bash

set -eo pipefail

if [ -z "$BASE_URL" ]
then
    echo "Please export BASE_URL to continue"
    exit 1
fi

if [ -z "$CLIENT_ID" ]
then
    echo "Please export CLIENT_ID to continue"
    exit 1
fi

if [ -z "$GROUP_ID" ]
then
    echo "Please export GROUP_ID to continue"
    exit 1
fi

if [ -z "$PATIENT_ID" ]
then
    echo "Please export PATIENT_ID to continue"
    exit 1
fi

if [ -z "$PRIVATE_KEY" ]
then
    echo "Please export PRIVATE_KEY to continue"
    exit 1
fi


NODE_OPTIONS="--max-old-space-size=8192" \
BASE_URL=${BASE_URL} \
CLIENT_ID=${CLIENT_ID} \
GROUP_ID=${GROUP_ID} \
PATIENT_ID=${PATIENT_ID} \
PRIVATE_KEY=${PRIVATE_KEY} \

node index.js --pattern "testSuite/**/!(authorization.test.js)"
