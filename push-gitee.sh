#!/usr/bin/env bash

set -e

git remote set-url origin https://gitee.com/quant1x/data.git
git checkout master
git push
git push --tags
git remote -vv