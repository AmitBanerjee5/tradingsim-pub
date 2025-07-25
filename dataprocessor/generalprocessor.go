package dataprocessor

import (
	"log"
)

type ActionType uint8

const (
	Hold  = 1
	Watch = 0
)

const (
	UpFactor    = 2
	DownFactor  = 1
)

// GeneralProcessor encapsulates the state and logic for trading a single stock.
type GeneralProcessor struct {
	// General Stuff
	Symbol            string
	InitialFund       float64
	CurrentFund       float64
	SmootheningFactor float64
	RiskFactor        float64

	// Current Price
	CurrentPrice     float64
	CurrentTimestamp uint64

	// Variables
	Action ActionType

	// Initial Price
	InitialPrice     float64
	InitialTimestamp uint64

	// Last Price
	LastPrice     float64
	LastTimestamp uint64

	// Peak
	PeakPrice     float64
	PeakTimestamp uint64

	// Peak
	ValleyPrice     float64
	ValleyTimestamp uint64

	// Buy
	BuyPrice     float64
	BuyTimestamp uint64
	BuyAmount    float64

}

// Print Values
// PrintValues logs the current state of the processor for debugging.
func (gp *GeneralProcessor) PrintValues(action string) {
	log.Println(action, gp.Symbol, gp.InitialFund, gp.CurrentFund, gp.SmootheningFactor, gp.RiskFactor, gp.CurrentPrice, gp.CurrentTimestamp, gp.Action, gp.LastPrice, gp.LastTimestamp, gp.PeakPrice, gp.PeakTimestamp, gp.ValleyPrice, gp.ValleyTimestamp, gp.BuyPrice, gp.BuyTimestamp, gp.BuyAmount)
}

// Constructor
// NewGeneralProcessor creates and initializes a new GeneralProcessor.
func NewGeneralProcessor(symbol string, initialFund float64, smootheningFactor float64, riskFactor float64) (gp *GeneralProcessor) {
	gp = &GeneralProcessor{
		Symbol:            symbol,
		InitialFund:       initialFund,
		CurrentFund:       initialFund,
		SmootheningFactor: smootheningFactor,
		RiskFactor:        riskFactor,
		CurrentPrice:      0.0,
		CurrentTimestamp:  0,
		Action:            Watch,
		InitialPrice:      0.0,
		InitialTimestamp:  0,
		LastPrice:         0.0,
		LastTimestamp:     0,
		PeakPrice:         0.0,
		PeakTimestamp:     0,
		ValleyPrice:       0.0,
		ValleyTimestamp:   0,
		BuyPrice:          0.0,
		BuyTimestamp:      0,
		BuyAmount:         0.0,
	}
	//gp.printHeaders()
	return gp
}

// Set Current Fund
// setCurrentFund updates the current fund based on the current price and buy amount.
func (gp *GeneralProcessor) setCurrentFund() {
	if gp.Action == Hold {
		gp.CurrentFund = gp.CurrentPrice * gp.BuyAmount
	}
}
