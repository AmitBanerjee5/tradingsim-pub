# tradingsim

## Overview

`tradingsim` is a Go-based trading simulation platform for processing real-time and historical stock tick data, executing algorithmic trading strategies, and storing results in InfluxDB. It supports multiple data sources, including Polygon and Finhub WebSockets, as well as CSV file-based historical data.

## Features
- Real-time data streaming from Polygon and Finhub WebSocket APIs
- Historical data replay from CSV files
- Algorithmic trading strategies (spike trading, step-up trading)
- Configurable stocks, funds, and strategy parameters via JSON
- Results and metrics stored in InfluxDB
- Modular and extensible architecture

## Project Structure
- `main.go` — Application entry point
- `configprocessor/` — Configuration objects and parsing
- `datachannel/` — Pub/sub channels for stock data
- `dataprocessor/` — Trading strategy implementations
- `streamlistener/` — Listeners for Polygon/Finhub WebSockets and CSV file reader
- `docker/` — Docker helper scripts
- `influxdbdashboard/` — InfluxDB dashboard JSON

## Getting Started

### Prerequisites
- Go 1.24+ (not needed for docker/containerized setup)
- InfluxDB instance (not needed for docker/containerized setup)
- Polygon and/or Finhub API credentials
- Docker (for docker/containerized setup)

### Installation (Local)
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd tradingsim
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Configure your simulation:
   - Edit `configprocessor/sampleconfig.json` with your stock symbols, funds, API tokens, and InfluxDB details.

### Running the Simulation (Local)
```bash
cat configprocessor/sampleconfig.json | go run main.go
```

## Using the Makefile

The project provides a Makefile for common development tasks:

- Build the binary:
  ```bash
  make build
  ```
- Run the binary:
  ```bash
  make run
  ```
- Format code:
  ```bash
  make fmt
  ```
- Lint code:
  ```bash
  make lint
  ```
- Clean up build artifacts:
  ```bash
  make clean
  ```
- Build and run the Docker image/container:
  ```bash
  make docker-run
  ```
- Clean and run the Docker container/image:
  ```bash
  make docker-clean
  ```

## Running with Docker

The project includes a Dockerfile that builds the Go binary and runs it alongside InfluxDB 2.7.11. The container:
- Sets up InfluxDB with user/password `tradingsim:tradingsim`
- Creates a bucket named `amibaner-trading`
- Loads a dashboard from `influxdbdashboard/stock_data.json`
- Runs the trading simulator

To build and run with Docker:

```bash
make docker-run
make docker-clean
```

Or manually:

```bash
docker build -t tradingsim .
docker run -i --rm -p 8086:8086 tradingsim &
```

Or using the built-in scripts:
```bash
./docker/createdocker.sh
./docker/removedocker.sh
```

The container will start InfluxDB and the trading simulator. You can access InfluxDB at [http://localhost:8086](http://localhost:8086) with the credentials `tradingsim` / `tradingsim`. Once logged in checkout the dashboard named `Stock Data`. You may want to set an autorefresh cycle.

## Configuration
See `configprocessor/sampleconfig.json` for an example configuration file. Key sections include:
- `stocktick.websocket.finhub` — Finhub WebSocket URL and token
- `stocktick.websocket.polygon` — Polygon WebSocket token
- `stocktick.filereader.polygon` — Path to historical CSV data
- `datarepo.influxdb` — InfluxDB connection details
- `stocks` — List of stocks to simulate, with initial fund and strategy parameters

## License
MIT License. See `LICENSE` for details.
