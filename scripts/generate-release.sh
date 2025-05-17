#!/bin/bash

BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH_NAME" != "master" ]; then
  TAG_FORMAT="${BRANCH_NAME}-v${version}"
else
  TAG_FORMAT="v${version}"
fi
npx semantic-release --tag-format "${TAG_FORMAT}"
