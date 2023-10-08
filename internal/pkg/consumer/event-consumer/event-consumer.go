package event_consumer

import (
	"MovieBot/internal/pkg/events"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Consumer struct {
	logger    *zap.Logger
	fetcher   events.UpdateFetcher
	processor events.Processor
	batchSize int
}

func New(logger *zap.Logger, fetcher events.UpdateFetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		logger:    logger,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	c.logger.Info("bot has started working")
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.logger.Error("error get events", zap.Error(err))
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			c.logger.Error("error handle events batch", zap.Error(err))
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	wg := sync.WaitGroup{}

	for _, event := range events {
		c.logger.Info("got new event", zap.String("event text", event.Text))
		ev := event
		go func() {
			wg.Add(1)
			defer wg.Done()

			if err := c.processor.Process(ev); err != nil {
				c.logger.Error("error handling event", zap.Error(err))
			}
		}()
	}
	wg.Wait()

	return nil
}
