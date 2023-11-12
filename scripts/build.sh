#!/bin/bash

if [ "$1" = "clean" ]
then
    echo "clean out directory"
    rm ./out/* 2>/dev/null
    exit 0
fi

echo "build client and server"
go build -o $GONTP_OUT_DIR/gontp-client $GONTP_HOME/cmd/gontp-client
go build -o $GONTP_OUT_DIR/gontp-server $GONTP_HOME/cmd/gontp-server
echo "finish"