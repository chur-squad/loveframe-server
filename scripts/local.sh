#!/bin/bash

set -e -E -u -o pipefail

cd `dirname $0` && cd ..
CURRENT=`pwd`

function start
{
    set_env
    go build -gcflags='all=-N -l' && ./loveframe-server
}

function set_env
{
  source "$CURRENT"/scripts/local_env.sh
}

CMD=$1
shift
$CMD $*
