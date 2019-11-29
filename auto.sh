#!/bin/bash


echo "start compiler"
go build
echo "start compiler ok"

sleep 1

echo "kill & cp & run"

killall -9 gofileserver

nohup ./gofileserver >gofileserver.log 2>&1 &

ps -aux | grep gofileserver

echo "kill & cp & run ok"
