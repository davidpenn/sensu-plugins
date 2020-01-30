#!/usr/bin/env bash

package_name=$(git rev-parse --show-toplevel | xargs basename)
version=$(git describe --abbrev=0 --tags 2>/dev/null || git rev-parse HEAD 2>/dev/null)

targets=(
	"bin/darwin/amd64/check"
	"bin/linux/amd64/check"
	"bin/windows/amd64/check"
)

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

for t in ${targets[@]}; do
	[ ! -f $t ] && make $t

	parts=(${t///// })
	os=${parts[1]%/}
	arch=${parts[2]%/}
	filename=${package_name}_${version}_${os}_${arch}.tar.gz

	# bsd tar
	[ ! -f dist/$filename ] && tar -czvf dist/$filename -C $(dirname $t) -s /./bin/ .

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
