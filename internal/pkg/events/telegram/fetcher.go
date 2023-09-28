package telegram

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events"
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

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "can't get events")
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, p.event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) event(upd telegram.Update) events.Event {
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

func fetchType(upd telegram.Update) events.Type {
	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}

	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}

func fetchText(upd telegram.Update) string {
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	if upd.Message != nil {
		return upd.Message.Text
	}

	return ""
}
