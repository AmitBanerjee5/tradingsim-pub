#/bin/bash

DID=`docker ps -a | grep tradingsim | head -1 | awk '{print $1}'`

echo "Killing Trading Simulator"
PID=`docker exec -it $DID ps -ef | grep '/app/tradingsim' | grep -v grep | awk '{print $1}'`
docker exec -it $DID kill -2 $PID

echo "Sleeping for 10 sec"
sleep 10

echo "Stopping and removing docker process"
docker stop $DID &>/dev/null
docker rm $DID &>/dev/null

echo "Removing the docker image"
docker image rm tradingsim
