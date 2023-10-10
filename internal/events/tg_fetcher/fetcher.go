package tg_fetcher

import "MovieBot/internal/events"

type MessageMeta struct {
	ChatID   int
	Username string
}

type CallbackMeta struct {
	CallbackID string
	ChatID     int
	MessageID  int
	Username   string
}

type UpdateFetcher interface {
	Fetch(limit int) ([]events.Event, error)
}
