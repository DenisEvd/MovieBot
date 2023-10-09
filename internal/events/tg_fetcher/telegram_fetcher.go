package tg_fetcher

import (
	telegram2 "MovieBot/internal/clients/telegram"
	"MovieBot/internal/events"
	"github.com/pkg/errors"
)

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

type Fetcher struct {
	tg     *telegram2.Client
	offset int
}

func NewFetcher(tg *telegram2.Client) *Fetcher {
	return &Fetcher{
		tg: tg,
	}
}

func (f *Fetcher) Fetch(limit int) ([]events.Event, error) {
	updates, err := f.tg.Updates(f.offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "can't get events")
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, f.event(u))
	}

	f.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (f *Fetcher) event(upd telegram2.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	switch updType {
	case events.Message:
		res.Meta = MessageMeta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	case events.CallbackQuery:
		res.Meta = CallbackMeta{
			CallbackID: upd.CallbackQuery.ID,
			ChatID:     upd.CallbackQuery.Message.Chat.ID,
			MessageID:  upd.CallbackQuery.Message.ID,
			Username:   upd.CallbackQuery.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram2.Update) events.Type {
	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}

	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}

func fetchText(upd telegram2.Update) string {
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	if upd.Message != nil {
		return upd.Message.Text
	}

	return ""
}
