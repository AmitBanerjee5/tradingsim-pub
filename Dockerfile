# Dockerfile for tradingsim with Go 1.24.3 and InfluxDB 2.7.11

# Build stage: build the Go binary
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o tradingsim main.go

# Final stage: run the app and InfluxDB
FROM  influxdb:2.7.11-alpine

# Copy the executable
WORKDIR /app

# Copy the Go binary
COPY --from=builder /app/tradingsim .
RUN chmod +x /app/tradingsim
COPY configprocessor/sampleconfig.json /app/sampleconfig.json

# Copy Influx Dashboard config
COPY influxdbdashboard/stock_data.json /root/stock_data.json

# Copy Docker Startup Script
COPY docker/dockerstartup.sh /root/dockerstartup.sh
RUN chmod +x /root/dockerstartup.sh

CMD ["/root/dockerstartup.sh"]
