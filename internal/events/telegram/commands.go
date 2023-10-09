package telegram

import (
	"MovieBot/internal/clients/telegram"
	"MovieBot/internal/events"
	"MovieBot/internal/events/telegram/messages"
	"MovieBot/internal/logger"
	"MovieBot/internal/storage"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
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
	logger.Info("got new command", zap.String("text", text), zap.String("from", username))
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
		if len(parts) != 2 {
			return p.tg.SendMessage(chatID, messages.MsgIncorrectCommand)
		}
		movieNum, err := strconv.Atoi(parts[1])
		if err != nil || movieNum < 1 {
			return p.tg.SendMessage(chatID, messages.MsgIncorrectCommand)
		}
		return p.showMovie(chatID, username, movieNum)
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

func (p *Processor) showMovie(chatID int, username string, movieNum int) error {
	movie, err := p.storage.GetNMovie(username, movieNum)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, messages.MsgNoSavedMovies)
	}

	return p.sendMovie(chatID, movie, movieNum)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "can't do command: send random")
		}
	}()

	movieNum := rand.Int()
	movie, err := p.storage.GetNMovie(username, movieNum)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, messages.MsgNoSavedMovies)
	}

	return p.sendMovie(chatID, movie, movieNum)
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

func (p *Processor) sendMovie(chatID int, movie events.Movie, n int) error {
	buttonNextData := fmt.Sprintf("%s;%d", getNextButton, n)
	buttonWatchData := fmt.Sprintf("%s;%d", watchItButton, movie.ID)
	buttons := make([]telegram.InlineKeyboardButton, 2)
	buttons[0], _ = messages.MakeButton("Next", buttonNextData)
	buttons[1], _ = messages.MakeButton("Watch it!", buttonWatchData)

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
