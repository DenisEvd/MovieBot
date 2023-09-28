package telegram

import (
	"MovieBot/internal/pkg/clients/telegram"
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
	movies, err := p.kp.FindMovieByTitle(text, 4)
	if err != nil {
		return errors.Wrap(err, "can't take movies from API")
	}

	if len(movies) == 0 {
		return p.tg.SendMessage(chatID, msgCanNotFindMovie)
	}

	sort.Slice(movies, func(i, j int) bool { return movies[i].Rating > movies[j].Rating })

	id, err := p.storage.AddRequest(text)
	if err != nil {
		return errors.Wrap(err, "saving request")
	}

	buttonDataNo := fmt.Sprintf("%s;%d", findMoreButton, id)
	buttonDataYes := fmt.Sprintf("%s;%d", saveButton, movies[0].ID)
	buttons := make([]telegram.InlineKeyboardButton, 2)
	buttons[0], _ = p.makeButton("No", buttonDataNo)
	buttons[1], _ = p.makeButton("Yes", buttonDataYes)

	return p.tg.SendPhotoWithInlineKeyboard(chatID, p.movieMessageByTitle(movies[0]), movies[0].Poster, buttons)
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
		return p.tg.SendMessage(chatID, msgNoSavedMovies)
	}

	movie, err := p.kp.GetMovieByID(movieID)
	if err != nil {
		return err
	}

	buttonData := fmt.Sprintf("%s;%d", watchItButton, movieID)
	buttons := make([]telegram.InlineKeyboardButton, 2)
	buttons[0], _ = p.makeButton("Next", getNextButton)
	buttons[1], _ = p.makeButton("Watch it!", buttonData)

	if movie.Poster.URL != "" {
		if err = p.tg.SendPhotoWithInlineKeyboard(chatID, p.movieMessageByID(movie), movie.Poster.URL, buttons); err != nil {
			return err
		}
	} else {
		if err = p.tg.SendMessageWithInlineKeyboard(chatID, p.movieMessageByID(movie), buttons); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
