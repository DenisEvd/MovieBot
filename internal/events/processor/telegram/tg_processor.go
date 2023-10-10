package telegram

import (
	"MovieBot/internal/clients/telegram"
	"MovieBot/internal/events"
	"MovieBot/internal/events/movie_fetcher"
	"MovieBot/internal/events/tg_fetcher"
	"MovieBot/internal/storage"
	"github.com/pkg/errors"
)

type TgProcessor struct {
	tg      telegram.TgClient
	kp      movie_fetcher.MovieFetcher
	storage storage.Storage
}

const (
	processingMsgError      = "can't process message"
	processingCallbackError = "can't process callback query"
)

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func NewTgProcessor(client telegram.TgClient, kp movie_fetcher.MovieFetcher, storage storage.Storage) *TgProcessor {
	return &TgProcessor{
		tg:      client,
		kp:      kp,
		storage: storage,
	}
}

func (p *TgProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processCallback(event)
	default:
		return errors.Wrap(ErrUnknownEventType, processingMsgError)
	}
}

func (p *TgProcessor) processMessage(event events.Event) error {
	meta, err := messageMeta(event)
	if err != nil {
		return errors.Wrap(err, processingMsgError)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return errors.Wrap(err, processingMsgError)
	}

	return nil
}

func (p *TgProcessor) processCallback(event events.Event) error {
	meta, err := callbackMeta(event)
	if err != nil {
		return errors.Wrap(err, processingCallbackError)
	}

	if err := p.doButton(meta.CallbackID, meta.ChatID, meta.MessageID, event.Text, meta.Username); err != nil {
		return errors.Wrap(err, processingCallbackError)
	}

	return nil
}

func messageMeta(event events.Event) (tg_fetcher.MessageMeta, error) {
	res, ok := event.Meta.(tg_fetcher.MessageMeta)
	if !ok {
		return tg_fetcher.MessageMeta{}, errors.Wrap(ErrUnknownMetaType, "can't get meta")
	}

	return res, nil
}

func callbackMeta(event events.Event) (tg_fetcher.CallbackMeta, error) {
	res, ok := event.Meta.(tg_fetcher.CallbackMeta)
	if !ok {
		return tg_fetcher.CallbackMeta{}, errors.Wrap(ErrUnknownMetaType, "can't get meta")
	}

	return res, nil
}
