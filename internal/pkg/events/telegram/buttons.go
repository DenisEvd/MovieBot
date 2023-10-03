package telegram

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	canselButton = "no"

	findMoreButton = "more"
	saveButton     = "save"

	getNextButton = "next"
	watchItButton = "watch"
)

var ErrUnknownDataType = errors.New("unknown data type")
var ErrTooBigData = errors.New("data is bigger then 64 bytes")
var ErrDoNotHaveRequestId = errors.New("don't have request id")

func (p *Processor) doButton(callbackID string, chatID int, messageID int, data string, username string) error {
	log.Printf("got new callback query from '%s'", username)

	parts := strings.Split(data, ";")

	switch parts[0] {
	case canselButton:
		return p.cancelSearch(callbackID, chatID, messageID)
	case findMoreButton:
		return p.showMoreMovies(callbackID, chatID, messageID, parts[1])
	case saveButton:
		requestID := ""
		if len(parts) == 3 {
			requestID = parts[2]
		}
		return p.saveMovie(callbackID, chatID, messageID, parts[1], username, requestID)
	case getNextButton:
		return p.showNextMovie(callbackID, chatID, messageID, username)
	case watchItButton:
		return p.watchThisMovie(callbackID, chatID, messageID, username, parts[1])
	default:
		return ErrUnknownDataType
	}
}

func (p *Processor) cancelSearch(callbackID string, chatID int, messageID int) error {
	err := p.tg.AnswerCallbackQuery(callbackID, msgSorry)
	if err != nil {
		return errors.Wrap(err, "canceling search")
	}

	return p.tg.EditMessageReplyMarkup(chatID, messageID)
}

func (p *Processor) showMoreMovies(callbackID string, chatID int, messageID int, requestID string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, "Ok!")
	if err != nil {
		return err
	}

	var request string
	if requestID == "" {
		return ErrDoNotHaveRequestId
	}
	id, _ := strconv.Atoi(requestID)
	request, err = p.storage.DeleteRequest(id)
	if err != nil {
		return err
	}

	movies, err := p.kp.FetchMoviesByTitle(request)

	sort.Slice(movies, func(i, j int) bool { return movies[i].Rating > movies[j].Rating })

	buttons, err := p.makeMoviesButtons(movies)
	if err != nil {
		return err
	}

	messageText := p.movieArrayMessage(movies)

	err = p.tg.EditMessageReplyMarkup(chatID, messageID)

	return p.tg.SendMessageWithInlineKeyboard(chatID, messageText, buttons)
}

func (p *Processor) saveMovie(callbackID string, chatID int, messageID int, data string, username string, requestID string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, msgSaved)
	if err != nil {
		return err
	}

	if requestID != "" {
		id, _ := strconv.Atoi(requestID)
		_, err := p.storage.DeleteRequest(id)
		if err != nil {
			return err
		}
	}

	movieID, err := strconv.Atoi(data)
	if err != nil {
		return errors.Wrap(err, "can't convert id to int")
	}
	isExists, err := p.storage.IsExists(username, movieID)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	movie, err := p.kp.FetchMovieById(movieID)
	if err != nil {
		return err
	}

	err = p.storage.AddMovie(username, movieID, movie.Title)
	if err != nil {
		return errors.Wrap(err, "saving movie")
	}

	return p.tg.EditMessageReplyMarkup(chatID, messageID)
}

func (p *Processor) showNextMovie(callbackID string, chatID int, messageID int, username string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, "Ok!")
	if err != nil {
		return err
	}

	err = p.tg.DeleteMessage(chatID, messageID)
	if err != nil {
		return err
	}

	return p.sendRandom(chatID, username)
}

func (p *Processor) watchThisMovie(callbackID string, chatID int, messageID int, username string, movieID string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, "Enjoy!")
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(movieID)

	err = p.storage.Remove(username, id)
	if err != nil {
		return errors.Wrap(err, "removing movie")
	}

	return p.tg.EditMessageReplyMarkup(chatID, messageID)
}

func (p *Processor) makeMoviesButtons(movies []events.Movie) ([]telegram.InlineKeyboardButton, error) {
	buttons := make([]telegram.InlineKeyboardButton, 0, len(movies)+1)
	for i, movie := range movies {
		buttonData := fmt.Sprintf("%s;%d", saveButton, movie.ID)
		button, err := p.makeButton(strconv.Itoa(i+1), buttonData)
		if err != nil {
			return []telegram.InlineKeyboardButton{}, err
		}
		buttons = append(buttons, button)
	}
	cancelButton, _ := p.makeButton("Cancel", canselButton)
	buttons = append(buttons, cancelButton)

	return buttons, nil
}
