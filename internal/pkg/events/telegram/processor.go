package telegram

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events"
	"MovieBot/internal/pkg/storage"
	"github.com/pkg/errors"
)

type Processor struct {
	tg      *telegram.Client
	kp      events.MovieFetcher
	storage *storage.Storage
}

const (
	processingMsgError      = "can't process message"
	processingCallbackError = "can't process callback query"
)

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func NewProcessor(client *telegram.Client, kp events.MovieFetcher, storage *storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		kp:      kp,
		storage: storage,
	}
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processCallback(event)
	default:
		return errors.Wrap(ErrUnknownEventType, processingMsgError)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := messageMeta(event)
	if err != nil {
		return errors.Wrap(err, processingMsgError)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return errors.Wrap(err, processingMsgError)
	}

	return nil
}

func (p *Processor) processCallback(event events.Event) error {
	meta, err := callbackMeta(event)
	if err != nil {
		return errors.Wrap(err, processingCallbackError)
	}

	if err := p.doButton(meta.CallbackID, meta.ChatID, meta.MessageID, event.Text, meta.Username); err != nil {
		return errors.Wrap(err, processingCallbackError)
	}

	return nil
}

func messageMeta(event events.Event) (MessageMeta, error) {
	res, ok := event.Meta.(MessageMeta)
	if !ok {
		return MessageMeta{}, errors.Wrap(ErrUnknownMetaType, "can't get meta")
	}

	return res, nil
}

func callbackMeta(event events.Event) (CallbackMeta, error) {
	res, ok := event.Meta.(CallbackMeta)
	if !ok {
		return CallbackMeta{}, errors.Wrap(ErrUnknownMetaType, "can't get meta")
	}

	return res, nil
}
