#!/bin/bash
set -e

export HOST=`hostname -i`
cmd=`eval echo "$@"`

echo "Host: $HOST"
echo "Running CMD: $cmd"

{ exec $cmd 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp
