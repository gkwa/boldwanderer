#!/usr/bin/env bash

set -e

tmp=$(mktemp -d boldwanderer.XXXXX)

if [ -z "${tmp+x}" ] || [ -z "$tmp" ]; then
    echo "Error: \$tmp is not set or is an empty string."
    exit 1
fi

{
    rg --files . \
        | grep -v $tmp/filelist.txt \
        | grep -vE 'boldwanderer$' \
        | grep -v README.org \
        | grep -v make_txtar.sh \
        | grep -v go.sum \
        | grep -v go.mod \
        | grep -v Makefile \
        | grep -v cmd/main.go \
        | grep -v logger.go \
        # | grep -v boldwanderer.go \

} | tee $tmp/filelist.txt
tar -cf $tmp/boldwanderer.tar -T $tmp/filelist.txt
mkdir -p $tmp/boldwanderer
tar xf $tmp/boldwanderer.tar -C $tmp/boldwanderer
rg --files $tmp/boldwanderer
txtar-c $tmp/boldwanderer | pbcopy

rm -rf $tmp
