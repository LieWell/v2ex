#! /bin/bash

set -ex

cd /root/v2ex
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o v2ex_server
chmod +x ./v2ex_server
nohup ./v2ex_server -c ./config.yaml > /dev/null 2>&1 &
exit 0