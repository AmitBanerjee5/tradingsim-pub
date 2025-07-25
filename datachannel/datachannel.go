package datachannel

import (
	"sync"
)

// StockData represents a single stock tick data point.
type StockData struct {
	Symbol    string
	Price     float64
	Timestamp uint64
}

// Datachannel implements a pub/sub channel for StockData.
type Datachannel struct {
	subscribers []chan *StockData
	mu          sync.RWMutex
	closed      bool
}

// NewStockData creates a new StockData instance.
func NewStockData(symbol string, price float64, timestamp uint64) *StockData {
	return &StockData{symbol, price, timestamp}
}

// NewDatachannel creates a new Datachannel instance.
func NewDatachannel() *Datachannel {
	return &Datachannel{
		mu: sync.RWMutex{},
	}
}

// Subscribe registers a new subscriber and returns a channel for receiving StockData.
func (s *Datachannel) Subscribe() <-chan *StockData {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	r := make(chan *StockData, 10000)

	s.subscribers = append(s.subscribers, r)

	return r
}

// Publish sends a StockData value to all subscribers.
func (s *Datachannel) Publish(value *StockData) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return
	}

	for _, ch := range s.subscribers {
		ch <- value
	}
}

// Close closes all subscriber channels and marks the Datachannel as closed.
func (s *Datachannel) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	for _, ch := range s.subscribers {
		close(ch)
	}

	s.closed = true
}
