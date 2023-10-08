package telegram

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events"
	"MovieBot/internal/pkg/events/telegram/messages"
	"MovieBot/internal/pkg/storage"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	RndCmd   = "/rnd"
	AllCmd   = "/all"
	ShowCmd  = "/show"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	p.logger.Info("got new command", zap.String("text", text), zap.String("from", username))
	text = strings.TrimSpace(text)
	parts := strings.Split(text, " ")

	switch parts[0] {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case AllCmd:
		return p.sendAll(chatID, username)
	case ShowCmd:
		return p.showMovie(chatID, username, parts[1])
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

func (p *Processor) showMovie(chatID int, username string, text string) error {
	movieNum, err := strconv.Atoi(text)
	if err != nil || movieNum < 1 {
		p.logger.Error("error convert string to int or invalid string", zap.String("data", text))
		return p.tg.SendMessage(chatID, messages.MsgIncorrectCommand)
	}

	movies, err := p.storage.GetAll(username)
	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, messages.MsgNoSavedMovies)
	}

	if err != nil {
		return errors.Wrap(err, "error send all")
	}

	movie := movies[movieNum%len(movies)]

	return p.sendMovie(chatID, movie)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "can't do command: send random")
		}
	}()

	movie, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, messages.MsgNoSavedMovies)
	}

	return p.sendMovie(chatID, movie)
}

func (p *Processor) sendAll(chatID int, username string) error {
	movies, err := p.storage.GetAll(username)

	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, messages.MsgNoSavedMovies)
	}

	if err != nil {
		return errors.Wrap(err, "error send all")
	}

	message := messages.MovieArrayMessage(movies)
	return p.tg.SendMessage(chatID, message)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, messages.MsgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, messages.MsgHello)
}

func (p *Processor) sendMovie(chatID int, movie events.Movie) error {
	buttonData := fmt.Sprintf("%s;%d", watchItButton, movie.ID)
	buttons := make([]telegram.InlineKeyboardButton, 2)
	buttons[0], _ = messages.MakeButton("Next", getNextButton)
	buttons[1], _ = messages.MakeButton("Watch it!", buttonData)

	if movie.Poster != "" {
		if err := p.tg.SendPhotoWithInlineKeyboard(chatID, messages.MovieMessage(movie), movie.Poster, buttons); err != nil {
			return err
		}
	} else {
		if err := p.tg.SendMessageWithInlineKeyboard(chatID, messages.MovieMessage(movie), buttons); err != nil {
			return err
		}
	}

	return nil
}
