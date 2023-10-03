package telegram

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events/telegram/messages"
	"MovieBot/internal/pkg/storage"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"sort"
	"strings"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	RndCmd   = "/rnd"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.suggestMovie(text, chatID)
	}
}

func (p *Processor) suggestMovie(text string, chatID int) error {
	movies, err := p.kp.FetchMoviesByTitle(text)
	if err != nil {
		return errors.Wrap(err, "can't take movies from API")
	}

	if len(movies) == 0 {
		return p.tg.SendMessage(chatID, messages.MsgCanNotFindMovie)
	}

	sort.Slice(movies, func(i, j int) bool { return movies[i].Rating > movies[j].Rating })

	requestID, err := p.storage.AddRequest(text)
	if err != nil {
		return errors.Wrap(err, "saving request")
	}

	buttonDataNo := fmt.Sprintf("%s;%d", findMoreButton, requestID)
	buttonDataYes := fmt.Sprintf("%s;%d;%d", saveButton, movies[0].ID, requestID)
	buttons := make([]telegram.InlineKeyboardButton, 2)
	buttons[0], _ = messages.MakeButton("No", buttonDataNo)
	buttons[1], _ = messages.MakeButton("Yes", buttonDataYes)

	messageText := messages.MovieMessage(movies[0])

	if movies[0].Poster == "" {
		return p.tg.SendMessageWithInlineKeyboard(chatID, messageText, buttons)
	}

	return p.tg.SendPhotoWithInlineKeyboard(chatID, messageText, movies[0].Poster, buttons)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "can't do command: send random")
		}
	}()

	movieID, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, messages.MsgNoSavedMovies)
	}

	movie, err := p.kp.FetchMovieById(movieID)
	if err != nil {
		return err
	}

	buttonData := fmt.Sprintf("%s;%d", watchItButton, movieID)
	buttons := make([]telegram.InlineKeyboardButton, 2)
	buttons[0], _ = messages.MakeButton("Next", getNextButton)
	buttons[1], _ = messages.MakeButton("Watch it!", buttonData)

	if movie.Poster != "" {
		if err = p.tg.SendPhotoWithInlineKeyboard(chatID, messages.MovieMessage(movie), movie.Poster, buttons); err != nil {
			return err
		}
	} else {
		if err = p.tg.SendMessageWithInlineKeyboard(chatID, messages.MovieMessage(movie), buttons); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, messages.MsgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, messages.MsgHello)
}
