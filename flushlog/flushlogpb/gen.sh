#!/usr/bin/env bash

GOGOPROTO_ROOT="$(go list -mod=readonly -f '{{ .Dir }}' -m github.com/gogo/protobuf)"
GOGOPROTO_PATH="${GOGOPROTO_ROOT}:${GOGOPROTO_ROOT}/protobuf"

protoc --gogofast_out=:. -I=. \
  -I="${GOGOPROTO_PATH}" \
  ./*.proto

sed -i.bak -E 's/import _ \"gogoproto\"//g' ./flushlog.pb.go
sed -i.bak -E 's/import _ \"google\/protobuf\"//g' ./flushlog.pb.go
sed -i.bak -E 's/\t_ \"google\/protobuf\"//g' -- ./flushlog.pb.go
rm -f ./*.bak
goimports -w ./*.pb.go
