#!/bin/bash

if [ $# -eq 0 ]
  then
    echo "No arguments supplied, please pass number of peer"
    exit 1
fi

docker attach sdcc_project_peer_"$1"
