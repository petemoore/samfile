#!/bin/bash -eu
cd "$(dirname "${0}")"/..
NEW_VERSION="${1}"
git grep -l 'github\.com/petemoore/samfile/v[0-9][0-9]*' | while read file; do
  cat "${file}" > x
  cat x | sed 's/\(github\.com\/petemoore\/samfile\/v\)[0-9][0-9]*/\1'"${NEW_VERSION%%.*}/g" > "${file}"
  git add "${file}"
  rm x
done
