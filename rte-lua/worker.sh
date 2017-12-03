#!/bin/sh

while true; do
  echo "Starting RTE Lua instance"
  /usr/local/bin/rte-lua worker
  echo "RTE Lua instance terminated"
done
