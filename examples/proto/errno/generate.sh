#!/bin/sh

echo 'Generating errno'
ERRORS=$(find ./ -type f -name '*.proto')
for ERROR in $ERRORS; do
  echo $ERROR
  protoc \
  -I. -I../..\
  --gofast_out=. \
  --gofast_opt paths=source_relative \
  --wzo-errno_out=. \
  --wzo-errno_opt epk=github.com/wyy-go/wzo/core/errors \
  --wzo-errno_opt paths=source_relative \
  $ERROR
done
