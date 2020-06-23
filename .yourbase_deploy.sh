#!/bin/sh

GITHUB_USERNAME="glacials"      # Must be the same user $GITHUB_TOKEN token belongs to
GITHUB_ORGANIZATION="glacials"  # Organization that holds the repository
GITHUB_REPOSITORY="twos.dev"    # Repository name
BUILD_DIR="dist"                # Set this to where your application gets built into

# This is the author and message information YourBase will use on each deploy
# commit. These commits will only appear on gh-pages, not on master, and will
# be overwritten by each new deploy.
COMMITTER_NAME="YourBase"
COMMITTER_EMAIL="ben@yourbase.io"
COMMIT_MESSAGE="Automatic deploy by YourBase"

mkdir -p ${BUILD_DIR}

# Replace these lines with any steps necessary to get your application built.
yarn install
yarn build

# You shouldn't need to change anything below here.
cd ${BUILD_DIR}
git init
git remote add origin https://${GITHUB_USERNAME}:${GITHUB_TOKEN}@github.com/${GITHUB_ORGANIZATION}/${GITHUB_REPOSITORY}
git checkout -b gh-pages
git add .
git -c "user.name=${COMMITTER_NAME}" -c "user.email=${COMMITTER_EMAIL}" commit --allow-empty -am "${COMMIT_MESSAGE}"
git push origin gh-pages --force
