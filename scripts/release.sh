#!/usr/bin/env bash

package_name=$(git rev-parse --show-toplevel | xargs basename)
version=$(git describe --abbrev=0 --tags 2>/dev/null || git rev-parse HEAD 2>/dev/null)

set -e

[ ! -d dist/ ] && mkdir dist/
if [ ! -z "$1" ]; then
	cat <<EOF > dist/asset.yaml
---
type: Asset
api_version: core/v2
metadata:
  name: davidpenn/${package_name}
  labels:
  annotations:
    com.github.davidpenn.${package_name}.version: ${version}
spec:
  builds:
EOF
fi

arch=amd64
targets=("darwin" "linux" "windows")
for os in ${targets[@]}; do
	rm -rf bin/
	GOOS=${os} make

	filename=${package_name}_${version}_${os}_${arch}.tar.gz
	[ ! -f dist/$filename ] && tar -czvf dist/$filename ./bin/

	if [ ! -z "$1" ]; then
		cat <<EOF >> dist/asset.yaml
  - filters:
    - entity.system.os == '${os}'
    - entity.system.arch == '${arch}'
    sha512: $(shasum -a 512 dist/${filename} | cut -d ' ' -f1)
    url: ${1%/}/${filename}
EOF
	fi

done
