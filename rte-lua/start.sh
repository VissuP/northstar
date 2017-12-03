#!/bin/bash
echo "Starting $NUM_WORKERS workers"

WORKER=0
while [  $WORKER -lt $NUM_WORKERS ]; do
  echo "Starting worker $WORKER"
  /usr/local/bin/worker.sh &
  let WORKER=WORKER+1
done

while true; do
  echo "Starting management interface"
  /usr/local/bin/rte-lua management
  echo "Management interface was stopped"
done
