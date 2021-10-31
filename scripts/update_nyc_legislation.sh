#!/bin/bash
set -e

mkdir -p build
go build -o build ./sync_legistar

if [ "$(git config --get user.email)" == "" ]; then
    git config --global user.email "jehiah@gmail.com"
    git config --global user.name "Jehiah Czebotar"
fi

if [ ! -e ../nyc_legislation ]; then
    git clone https://github.com/jehiah/nyc_legislation.git ../nyc_legislation
    pushd ../nyc_legislation >/dev/null
    git remote rm origin
    git remote add origin https://jehiah:$GH_TOKEN@github.com/jehiah/nyc_legislation.git
    popd >/dev/null
fi

./build/sync_legistar --target-dir=../nyc_legislation


pushd ../nyc_legislation
git status
popd