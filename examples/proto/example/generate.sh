#!/bin/sh

echo 'Generating api'
PROTOS=$(find ./ -type f -name '*.proto')

for PROTO in $PROTOS; do
  echo $PROTO
  protoc \
    -I. -I$(dirname $PROTO) \
    -I../../third_party \
    --gofast_out=. \
    --gofast_opt paths=source_relative \
    --wzo-gin_out=. \
    --wzo-gin_opt paths=source_relative \
    --wzo-gin_opt allow_empty_patch_body=true \
    $PROTO
done

for PROTO in $PROTOS; do
  echo $PROTO
  protoc \
    -I. -I$(dirname $PROTO) \
    -I../../third_party \
    --gofast_out=. \
    --gofast_opt paths=source_relative \
    --rpcx_out=. \
    --rpcx_opt paths=source_relative \
    $PROTO
done