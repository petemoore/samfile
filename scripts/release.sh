#!/bin/bash

# This script is used to generate releases of samfile. It should be the only
# way that releases are created. There are two phases, the first is checking
# that the code is in a clean and working state. The second phase is modifying
# files, tagging, commiting and pushing to github.

# exit in case of bad exit code
set -e

OFFICIAL_GIT_REPO='git@github.com:petemoore/samfile'

# step up a directory from this script
cd "$(dirname "${0}")/.."

NEW_VERSION="${1}"

if [ -z "${1}" ]; then
  echo "Please supply version number for release, e.g. ./release.sh 7.2.0" >&2
  exit 1
fi

VALID_FORMAT='^[1-9][0-9]*\.\(0\|[1-9][0-9]*\)\.\(0\|[1-9]\)\([0-9]*alpha[1-9][0-9]*\|[0-9]*\)$'
FORMAT_EXPLANATION='should be "<a>.<b>.<c>" OR "<a>.<b>.<c>alpha<d>" where a>=1, b>=0, c>=0, d>=1 and a,b,c,d are integers, with no leading zeros'

if ! echo "${NEW_VERSION}" | grep -q "${VALID_FORMAT}"; then
  echo "Release version '${NEW_VERSION}' not allowed (${FORMAT_EXPLANATION})" >&2
  exit 65
fi

if [ $(git ls-remote -t "${OFFICIAL_GIT_REPO}" "v${NEW_VERSION}" | wc -l) != '0' ]; then
  echo "Cannot release as version ${NEW_VERSION} since tag v${NEW_VERSION} already exists on ${OFFICIAL_GIT_REPO}" >&2
  exit 66
fi

# Make sure git tag doesn't already exist on remote
if [ "$(git ls-remote -t "${OFFICIAL_GIT_REPO}" "v${NEW_VERSION}" 2>&1 | wc -l | tr -d ' ')" != '0' ]; then
  echo "git tag '${NEW_VERSION}' already exists remotely on ${OFFICIAL_GIT_REPO},"
  echo "or there was an error checking whether it existed:"
  git ls-remote -t "${OFFICIAL_GIT_REPO}" "v${NEW_VERSION}"
  exit 67
fi

# Local changes will not be in the release, so they should be dealt with before
# continuing. git stash can help here! Untracked files can make it into release
# so let's make sure we have none of them either.
modified="$(git status --porcelain)"
if [ -n "$modified" ]; then
  echo "There are changes in the local tree. This probably means"
  echo "you'll do something unintentional. For safety's sake, please"
  echo 'revert or stash them!'
  echo
  git status
  exit 68
fi

# ******** If making a NON-alpha release only **********
# Check that the current HEAD is also the tip of the official repo master
# branch. If the commits match, it does not matter what the local branch
# name is, or even if we have a detached head.
if ! echo "${NEW_VERSION}" | grep -q "alpha"; then
  remoteMasterSha="$(git ls-remote "${OFFICIAL_GIT_REPO}" master | cut -f1)"
  localSha="$(git rev-parse HEAD)"
  if [ "${remoteMasterSha}" != "${localSha}" ]; then
    echo "Locally, you are on commit ${localSha}."
    echo "The remote petemoore/samfile repo master branch is on commit ${remoteMasterSha}."
    echo "Make sure to git push/pull so that they both point to the same commit."
    exit 69
  fi
fi

scripts/refresh_readme.sh
git add README.md
git commit -m "Refreshed README with samfile output from v${NEW_VERSION}"
git tag -s "v${NEW_VERSION}" -m "Release ${NEW_VERSION}"
# only ensure master is updated if it is a non-alpha release
if ! echo "${NEW_VERSION}" | grep -q "alpha"; then
  git push "${OFFICIAL_GIT_REPO}" "+HEAD:refs/heads/master"
  git fetch --all
fi
git push "${OFFICIAL_GIT_REPO}" "+refs/tags/v${NEW_VERSION}:refs/tags/v${NEW_VERSION}"
