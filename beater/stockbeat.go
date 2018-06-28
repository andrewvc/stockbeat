package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/andrewvc/stockbeat/config"
	"github.com/andrewvc/stockbeat/fetcher"
)

type Stockbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

var logger = logp.NewLogger("stockbeat.fetcher")

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Stockbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Stockbeat) Run(b *beat.Beat) error {
	logp.Info("stockbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		quotes, err := fetcher.RetrieveQuotes(bt.config.Symbols)

		if err != nil {
			logger.Error(err)
			continue
		}

		for _, quote := range *quotes {
			event := beat.Event{
				Timestamp: time.Unix(0, quote.LatestUpdate * int64(time.Millisecond)),
				Fields: common.MapStr{
					"price": quote.LatestPrice,
					"companyName": quote.CompanyName,
					"latestVolume": quote.LatestVolume,
					"primaryExchange": quote.PrimaryExchange,
					"symbol": quote.Symbol,
					"sector": quote.Sector,
				},
			}
			bt.client.Publish(event)
			logger.Info(fmt.Sprintf("Event sent for %s", quote.Symbol))
		}

		counter++
	}
}

func (bt *Stockbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
