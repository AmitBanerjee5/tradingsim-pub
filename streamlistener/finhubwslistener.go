package streamlistener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	dc "main/datachannel"
	"sync"

	"github.com/gorilla/websocket"
)

// FinhubWSListener listens to Finhub WebSocket for real-time stock trades.
type FinhubWSListener struct {
	Token   string
	Url     string
	Symbols []string
}

type finhubStockDatum struct {
	Symbol     string   `json:"s"`
	Price      float64  `json:"p"`
	Timestamp  uint64   `json:"t"`
	Conditions []string `json:"c" omitempty:"true"`
}

type finhubResponse struct {
	Type      string             `json:"type"`
	StockData []finhubStockDatum `json:"data"`
}

// ListenFinhub subscribes to Finhub WebSocket trades and publishes them to Datachannels.
func (st *FinhubWSListener) ListenFinhub(p map[string]*dc.Datachannel, wg *sync.WaitGroup, ctx context.Context) {
	wg.Add(1)
	defer wg.Done()

	dialStr := fmt.Sprintf("%s?token=%s", st.Url, st.Token)
	w, _, err := websocket.DefaultDialer.Dial(dialStr, nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	for _, s := range st.Symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		w.WriteMessage(websocket.TextMessage, msg)
	}

	var msg finhubResponse
	for {
		select {
		case <-ctx.Done():
			log.Println("Closing Listener")
			return
		default:
			err := w.ReadJSON(&msg)
			if err != nil {
				panic(err)
			}
			for _, data := range msg.StockData {
				p[data.Symbol].Publish(dc.NewStockData(data.Symbol, data.Price, data.Timestamp))
			}
		}
	}
}
