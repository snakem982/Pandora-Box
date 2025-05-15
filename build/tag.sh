#!/bin/bash

cd $(dirname $(readlink -f "$0"))
cd ..

# 检查是否提供了版本号参数
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1

# 创建带注释的 Git tag
git tag -a "$VERSION" -m "$VERSION"

# 推送到远程仓库
git push origin "$VERSION"
