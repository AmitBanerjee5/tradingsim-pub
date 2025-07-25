#/bin/bash

# These variables should be in accordance with the configuration
# @configprocessor/sampleconfig.json
export DOCKER_INFLUXDB_INIT_MODE=setup
export DOCKER_INFLUXDB_INIT_USERNAME=tradingsim
export DOCKER_INFLUXDB_INIT_PASSWORD=tradingsim
export DOCKER_INFLUXDB_INIT_ORG=amibaner-org
export DOCKER_INFLUXDB_INIT_BUCKET=amibaner-trading
export DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=s7Vs1yBfmcLVLCg3NQbokmNsFiZjwGH0mEECcWpDkqki0yXxIcud6aQqGt9bGw75euKZ-Ki0d5Nldr4fv4_Rzw==

# Start InfluxDB
nohup influxd &>influx.out &

sleep 10

echo "Privsioning setup"
influx setup \
  --username ${DOCKER_INFLUXDB_INIT_USERNAME} \
  --password ${DOCKER_INFLUXDB_INIT_PASSWORD} \
  --token ${DOCKER_INFLUXDB_INIT_ADMIN_TOKEN} \
  --org ${DOCKER_INFLUXDB_INIT_ORG} \
  --bucket ${DOCKER_INFLUXDB_INIT_BUCKET} \
  --force

sleep 10

echo "Provisioning Dashboard"
influx apply \
  --file /root/stock_data.json \
  --org ${DOCKER_INFLUXDB_INIT_ORG} \
  --token ${DOCKER_INFLUXDB_INIT_ADMIN_TOKEN} \
  --force yes

echo "Running Trading Simulator"
cat /app/sampleconfig.json | /app/tradingsim
