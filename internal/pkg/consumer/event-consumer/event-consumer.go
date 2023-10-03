package event_consumer

import (
	"MovieBot/internal/pkg/events"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.UpdateFetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.UpdateFetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	log.Println("bot has started working")
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	wg := sync.WaitGroup{}

	for _, event := range events {
		log.Printf("got new event: %s", event.Text)
		ev := event
		go func() {
			wg.Add(1)
			defer wg.Done()

			if err := c.processor.Process(ev); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}
		}()
	}
	wg.Wait()

	return nil
}
