package configprocessor

import (
	"encoding/json"
)

// TSConfig holds the top-level configuration for the trading simulator.
type TSConfig struct {
	StockTick StockTickCon    `json:"stocktick"`
	DataRepo  Repo            `json:"datarepo"`
	Stocks    []TradingStocks `json:"stocks"`
}

type StockTickCon struct {
	Websocket  WebsocketCon  `json:"websocket"`
	FileReader FileReaderObj `json:"filereader"`
}

type FileReaderObj struct {
	Polygon PolygonReader `json:"polygon"`
}

type PolygonReader struct {
	FileLocation string `json:"location"`
}

type WebsocketCon struct {
	Finhub  FinhubCon  `json:"finhub"`
	Polygon PolygonCon `json:"polygon"`
}

type FinhubCon struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type PolygonCon struct {
	Token string `json:"token"`
}

type Repo struct {
	InfluxDB InfluxRepo `json:"influxdb"`
}

type InfluxRepo struct {
	Url    string `json:"url"`
	Token  string `json:"token"`
	Bucket string `json:"bucket"`
	Org    string `json:"org"`
}

type TradingStocks struct {
	Symbol            string  `json:"symbol"`
	Fund              float64 `json:"fund"`
	SmootheningFactor float64 `json:"smootheningfactor"`
	RiskFactor        float64 `json:"riskfactor"`
}

// ParseConfig parses a JSON configuration and returns a TSConfig struct.
func ParseConfig(s *[]byte) *TSConfig {
	var sl TSConfig
	err := json.Unmarshal(*s, &sl)
	if err != nil {
		panic(err)
	}
	return &sl
}
