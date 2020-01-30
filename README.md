# sensu-plugins

## Files
 * bin/check

## Usage

**mysql-replication-status** example:
```sh
check mysql-replication-status --host 127.0.0.1 --user root --password secret
```

## Installation

```sh
scripts/release.sh https://url-you-will-serve-tarballs/
# upload dist/*.tar.gz => url-you-will-serve-tarball
sensuctl create -f dist/asset.yaml
```
