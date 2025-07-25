#/bin/bash
DID=`docker ps -a | grep tradingsim | head -1 | awk '{print $1}'`
docker exec -it $DID /bin/bash
