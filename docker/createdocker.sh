#/bin/bash
docker build -t tradingsim .
docker run -i --rm -p 8086:8086 tradingsim &
