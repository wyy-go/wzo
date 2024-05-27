#!/bin/sh

echo 'Generating errno'
ERRORS=$(find ./ -type f -name '*.proto')
for ERROR in $ERRORS; do
  echo $ERROR
  protoc \
  -I. -I../.. -I../../../proto\
  --gofast_out=. \
  --gofast_opt paths=source_relative \
  --wzo-errno_out=. \
  --wzo-errno_opt paths=source_relative \
  $ERROR
done
