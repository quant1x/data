#!/bin/bash
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/kline-win-amd64.exe' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/kline-mac-amd64' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/realtime-mac-amd64' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/quant-win-amd64.exe' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/xdxr-win-amd64.exe' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/xdxr-mac-amd64' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/kline-win-amd64.exe' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/kline' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/kline-mac-amd64' --prune-empty --tag-name-filter cat -- --all
git filter-branch --force --index-filter 'git rm -rf --cached --ignore-unmatch bin/quant-mac-amd64' --prune-empty --tag-name-filter cat -- --all
