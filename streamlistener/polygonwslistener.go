package streamlistener

import (
	"context"
	"log"
	dc "main/datachannel"
	"sync"

	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)

// PolygonWSListener listens to Polygon WebSocket for real-time stock trades.
type PolygonWSListener struct {
	Token   string
	Symbols []string
}

// ListenPolygon subscribes to Polygon WebSocket trades and publishes them to Datachannels.
func (st *PolygonWSListener) ListenPolygon(p map[string]*dc.Datachannel, wg *sync.WaitGroup, ctx context.Context) {
	wg.Add(1)
	defer wg.Done()

	c, err := polygonws.New(polygonws.Config{
		APIKey: st.Token,
		Feed:   polygonws.Delayed,
		Market: polygonws.Stocks,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer c.Close()

	for _, s := range st.Symbols {
		_ = c.Subscribe(polygonws.StocksTrades, s)
	}

	if err := c.Connect(); err != nil {
		log.Fatalf("%v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Closing Listener")
			return
		case <-c.Error():
			return
		case out, more := <-c.Output():
			if !more {
				return
			}
			switch trade := out.(type) {
			case models.EquityTrade:
				p[trade.Symbol].Publish(dc.NewStockData(trade.Symbol, trade.Price, uint64(trade.Timestamp)))
				//log.Println("%v\n", trade.Price)
			}
		}
	}
}
