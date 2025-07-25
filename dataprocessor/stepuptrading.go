package dataprocessor

import (
	"context"
	"log"
	cf "main/configprocessor"
	dc "main/datachannel"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Processor
// StepupTrade executes a step-up trading strategy on incoming stock data.
func (gp *GeneralProcessor) StepupTrade(ch *dc.Datachannel, wg *sync.WaitGroup, iClient *influxdb2.Client, config *cf.TSConfig, ctx context.Context) {

	wg.Add(1)
	defer wg.Done()
	s := ch.Subscribe()

	// Influx
	iWriteAPI := (*iClient).WriteAPIBlocking(config.DataRepo.InfluxDB.Org, config.DataRepo.InfluxDB.Bucket)

	for {
		select {
		case <-ctx.Done():
			// Implement force Sale
			log.Println("Shutting down trading\n")
			return
		case stockdata, ok := <-s:
			if !ok {
				log.Println("Couldn't retrieve stock data from channel, exiting\n")
				break
			}

			// Initial data
			if gp.LastPrice == 0.0 {
				gp.LastPrice = stockdata.Price
				gp.LastTimestamp = stockdata.Timestamp
				gp.InitialPrice = stockdata.Price
				gp.InitialTimestamp = stockdata.Timestamp
				gp.CurrentPrice = stockdata.Price
				gp.CurrentTimestamp = stockdata.Timestamp
				gp.BuyPrice = stockdata.Price
				gp.BuyTimestamp = stockdata.Timestamp
				gp.setPeak()
				gp.setValley()
				continue
			}

			gp.LastPrice = gp.CurrentPrice
			gp.LastTimestamp = gp.CurrentTimestamp
			gp.CurrentPrice = stockdata.Price
			gp.CurrentTimestamp = stockdata.Timestamp

			// State Machine
			if gp.CurrentPrice >= gp.LastPrice {
				if gp.CurrentPrice > gp.PeakPrice {
					gp.setPeak()
				}
			} else {
				if gp.CurrentPrice < gp.ValleyPrice {
					gp.setValley()
				}
			}
			gp.trySell()
			gp.tryBuy()

			// Common for all Conditions
			gp.setCurrentFund()
			p := influxdb2.NewPoint("stockdata",
				map[string]string{"Symbol": gp.Symbol},
				map[string]interface{}{"IniFnd": gp.InitialFund,
					"CrtFnd": gp.CurrentFund, "CrtPrc": gp.CurrentPrice,
					"CrtTsp": gp.CurrentTimestamp, "Action": gp.Action,
					"PckPrc": gp.PeakPrice, "PckTmsp": gp.PeakTimestamp,
					"ValPrc": gp.ValleyPrice, "ValTmsp": gp.ValleyTimestamp,
					"BuyPrc": gp.BuyPrice, "BuyTmsp": gp.BuyTimestamp,
					"BuyAmt": gp.BuyAmount, "NoActnUpr": (gp.BuyPrice + gp.RiskFactor*UpFactor),
					"NoActnLwr": (gp.PeakPrice - gp.RiskFactor*DownFactor),
					"PrftPctg":  ((gp.CurrentFund - gp.InitialFund) / gp.InitialFund) * 100,
					"ChgPctg":   ((gp.CurrentPrice - gp.InitialPrice) / gp.InitialPrice) * 100},
				time.Now())
			iWriteAPI.WritePoint(context.Background(), p)
			iWriteAPI.Flush(context.Background())
		}
	}
}

// Check if the stock can be bought
// tryBuy checks if a buy condition is met and executes a buy.
func (gp *GeneralProcessor) tryBuy() {
	if gp.Action == Watch && (gp.CurrentPrice > (gp.ValleyPrice+gp.RiskFactor)) {
		gp.Action = Hold
		gp.BuyPrice = gp.CurrentPrice
		gp.BuyTimestamp = gp.CurrentTimestamp
		gp.BuyAmount = gp.CurrentFund / gp.CurrentPrice
		gp.setPeak()
		gp.setValley()

		//Print
		gp.PrintValues("BUY")
	}
}

// Check if the stock can be bought
// trySell checks if a sell condition is met and executes a sell.
func (gp *GeneralProcessor) trySell() {
	if gp.Action == Hold && (gp.BuyPrice != gp.CurrentPrice) &&
		((gp.CurrentPrice >= (gp.BuyPrice+gp.RiskFactor*UpFactor) && gp.CurrentPrice <= (gp.PeakPrice-gp.SmootheningFactor)) ||
			gp.CurrentPrice < (gp.PeakPrice-gp.RiskFactor*DownFactor)) {
		gp.CurrentFund = gp.CurrentPrice * gp.BuyAmount
		gp.Action = Watch
		gp.setPeak()
		gp.setValley()

		//Print
		gp.PrintValues("SELL")
	}
}

// Set Peak
// setPeak sets the peak price and timestamp.
func (gp *GeneralProcessor) setPeak() {
	gp.PeakPrice = gp.CurrentPrice
	gp.PeakTimestamp = gp.CurrentTimestamp
}

// Set Valley
// setValley sets the valley price and timestamp.
func (gp *GeneralProcessor) setValley() {
	gp.ValleyPrice = gp.CurrentPrice
	gp.ValleyTimestamp = gp.CurrentTimestamp
}
