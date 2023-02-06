#!/bin/bash

echo "VERCEL_GIT_COMMIT_REF: $VERCEL_GIT_COMMIT_REF"

$(git diff HEAD^ HEAD --quiet .)
DIFF=$?

if [[ "$VERCEL_GIT_COMMIT_REF" == "master" ]] ; then
    echo "✅  Master branch updated. Building..."
  exit 1;
elif [[ $DIFF -eq 1 ]] ; then
    echo "✅  Changes detected. Building..."
    exit 1;
else
  echo "🛑  Not building - no changes detected and not on master"
  exit 0;
fi
