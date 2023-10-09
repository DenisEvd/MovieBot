package event_consumer

import (
	"MovieBot/internal/events"
	"MovieBot/internal/events/processor"
	"MovieBot/internal/events/tg_fetcher"
	"MovieBot/internal/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   tg_fetcher.UpdateFetcher
	processor processor.Processor
	batchSize int
}

func New(fetcher tg_fetcher.UpdateFetcher, processor processor.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	logger.Info("bot has started working")

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			logger.Error("error get events", zap.Error(err))
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			logger.Error("error handle events batch", zap.Error(err))
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	wg := sync.WaitGroup{}

	for _, event := range events {
		logger.Info("got new event", zap.String("event text", event.Text))
		ev := event
		go func() {
			wg.Add(1)
			defer wg.Done()

			if err := c.processor.Process(ev); err != nil {
				logger.Error("error handling event", zap.Error(err))
			}
		}()
	}
	wg.Wait()

	return nil
}
