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
    --wzo-resty_out=. \
    --wzo-resty_opt paths=source_relative \
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


echo 'Generating api swagger'
protoc \
  -I. \
  -I../../third_party \
  --openapiv2_out docs \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_merge=true \
  --openapiv2_opt merge_file_name=swagger \
  --openapiv2_opt enums_as_ints=true \
  --openapiv2_opt json_names_for_fields=false \
$PROTOS