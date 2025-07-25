package main

import (
	"context"
	"io"
	"log"
	cp "main/configprocessor"
	dc "main/datachannel"
	dp "main/dataprocessor"
	st "main/streamlistener"
	"os"
	"os/signal"
	"sync"
	"syscall"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Main entry point for the trading simulation application.
// Reads configuration, sets up data channels, processors, and listeners, and manages shutdown.
func main() {

	symbols := []string{}
	cs, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading: %v", err)
	}
	config := cp.ParseConfig(&cs)

	// Create new InfluxDB client
	iClient := influxdb2.NewClient(config.DataRepo.InfluxDB.Url, config.DataRepo.InfluxDB.Token)
	defer iClient.Close()

	// Signal Handler
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stock Traders
	dcs := make(map[string]*dc.Datachannel)
	wgs := make(map[string]*sync.WaitGroup)
	for _, s := range config.Stocks {
		dcs[s.Symbol] = dc.NewDatachannel()
		wgs[s.Symbol] = &sync.WaitGroup{}
		symbols = append(symbols, s.Symbol)

		// Initiate Stock Data Processor
		stockProcessor := dp.NewGeneralProcessor(s.Symbol, s.Fund, s.SmootheningFactor, s.RiskFactor)
		go stockProcessor.StepupTrade(dcs[s.Symbol], wgs[s.Symbol], &iClient, config, ctx)
	}

	// Stock Listeners
	lWg := sync.WaitGroup{}

	// -- Few methods to read tick data -
	// Some are real time, some are historic
	// Only one should be activated

	// Method 1: Finhub Websocket (Second aggregate os something, not all raw data)
	//finhubListener := &st.FinhubWSListener{Token: config.StockTick.Websocket.Finhub.Token, Url: config.StockTick.Websocket.Finhub.Url, Symbols: symbols}
	//go finhubListener.ListenFinhub(dcs, &lWg, ctx)

	// Method 2: Polygon Websocket (More granular real raw Tick data, 15 min delayed)
	polygonListener := &st.PolygonWSListener{Token: config.StockTick.Websocket.Polygon.Token, Symbols: symbols}
	go polygonListener.ListenPolygon(dcs, &lWg, ctx)

	// Method 3: Polygon Historical data from File
	// File can be downloaded by using aws s3 secret or from WebUI
	//polygonReader := &st.PolygonReader{FileLocation: config.StockTick.FileReader.Polygon.FileLocation}
	//go polygonReader.ReadFile(dcs, &lWg, ctx)

	// Waiting for Shutdown
	sig := <-sigs
	log.Println()
	log.Println(sig)
	cancel()

	// Waiting for the listener be done
	lWg.Wait()

	// Waiting for the trading be settled
	for _, s := range symbols {
		wgs[s].Wait()
		dcs[s].Close()
	}
}
