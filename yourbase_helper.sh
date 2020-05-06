#!/bin/sh

mkdir -p build
yarn install
yarn build
cd build
git init
git remote add origin https://glacials:${GITHUB_TOKEN}@github.com/glacials/twos.dev
git checkout -b gh-pages
git add .
git -c 'user.name=YourBase' -c 'user.email=ben@yourbase.io' commit -am 'Automatic deploy by YourBase'
git push origin gh-pages --force
cd ../
