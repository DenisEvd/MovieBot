package tg_fetcher

import "MovieBot/internal/events"

type UpdateFetcher interface {
	Fetch(limit int) ([]events.Event, error)
}
