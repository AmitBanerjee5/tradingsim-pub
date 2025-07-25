package streamlistener

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	dc "main/datachannel"
	"os"
	"strconv"
	"sync"
	"time"
)

// PolygonReader reads historical stock tick data from a CSV file and publishes to Datachannels.
type PolygonReader struct {
	FileLocation string
}

// ReadFile reads the CSV file and publishes each tick to the appropriate Datachannel.
func (st *PolygonReader) ReadFile(p map[string]*dc.Datachannel, wg *sync.WaitGroup, ctx context.Context) {
	wg.Add(1)
	defer wg.Done()

	// open file
	f, err := os.Open(st.FileLocation)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	for {
		select {
		case <-ctx.Done():
			log.Println("Closing Listener")
			return
		default:
			rec, err := csvReader.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v", err)
			}

			channel, exists := p[rec[0]]
			if exists {
				price, _ := strconv.ParseFloat(rec[6], 64)
				ts, _ := strconv.ParseUint(rec[8], 0, 64)
				channel.Publish(dc.NewStockData(rec[0], price, ts))
				time.Sleep(2 * time.Millisecond)
			}
		}
	}
}
