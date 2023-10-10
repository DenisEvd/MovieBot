package telegram

import (
	"MovieBot/internal/clients/telegram"
	"MovieBot/internal/events"
	"MovieBot/internal/events/processor/messages"
	"MovieBot/internal/logger"
	"MovieBot/internal/storage"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
var ErrDoNotHaveRequestId = errors.New("don't have request id")

func (p *TgProcessor) doButton(callbackID string, chatID int, messageID int, data string, username string) error {
	logger.Info("got new callback query", zap.String("from", username), zap.String("data", data))

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
		return p.showNextMovie(callbackID, chatID, messageID, username, parts[1])
	case watchItButton:
		return p.watchThisMovie(callbackID, chatID, messageID, username, parts[1])
	default:
		return ErrUnknownDataType
	}
}

func (p *TgProcessor) cancelSearch(callbackID string, chatID int, messageID int) error {
	err := p.tg.AnswerCallbackQuery(callbackID, messages.MsgOkay)
	if err != nil {
		return errors.Wrap(err, "canceling search")
	}

	return p.editMessage(chatID, messageID)
}

func (p *TgProcessor) showMoreMovies(callbackID string, chatID int, messageID int, requestID string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, messages.MsgOkay)
	if err != nil {
		return err
	}

	var request string
	if requestID == "" {
		return ErrDoNotHaveRequestId
	}
	id, _ := strconv.Atoi(requestID)
	request, err = p.storage.DeleteRequest(id)
	if errors.Is(err, storage.ErrNoRequest) {
		if err = p.editMessage(chatID, messageID); err != nil {
			return errors.Wrap(err, "error show more movies")
		}

		return p.tg.SendMessage(chatID, messages.MsgTryAgain)
	}

	if err != nil {
		return err
	}

	movies, err := p.kp.FetchMoviesByTitle(request)
	if err != nil {
		return errors.Wrap(err, "error show more movies")
	}

	if err := p.editMessage(chatID, messageID); err != nil {
		return errors.Wrap(err, "error show more movies")
	}

	if len(movies) == 0 {
		return p.tg.SendMessage(chatID, messages.MsgCanNotFindMovie)
	}
	sort.Slice(movies, func(i, j int) bool { return movies[i].Rating > movies[j].Rating })

	buttons, err := p.makeMoviesButtons(movies)
	if err != nil {
		return err
	}

	messageText := messages.MovieArrayMessage(movies)

	return p.tg.SendMessageWithInlineKeyboard(chatID, messages.HeaderOfMoreMovieList+messageText, buttons)
}

func (p *TgProcessor) saveMovie(callbackID string, chatID int, messageID int, data string, username string, requestID string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, messages.MsgSaved)
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
	isExists, err := p.storage.IsExistRecord(username, movieID)
	if err != nil {
		return err
	}

	if isExists {
		if err := p.editMessage(chatID, messageID); err != nil {
			return errors.Wrap(err, "edit message")
		}

		isWatched, err := p.storage.IsWatched(username, movieID)
		if err != nil {
			return errors.Wrap(err, "check is watched")
		}

		if isWatched {
			return p.tg.SendMessage(chatID, messages.MsgAlreadyWatched)
		}

		return p.tg.SendMessage(chatID, messages.MsgAlreadyExists)
	}

	movie, err := p.kp.FetchMovieById(movieID)
	if err != nil {
		return err
	}

	err = p.storage.AddMovie(username, movie)
	if err != nil {
		return errors.Wrap(err, "saving movie")
	}

	return p.editMessage(chatID, messageID)
}

func (p *TgProcessor) showNextMovie(callbackID string, chatID int, messageID int, username string, n string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, messages.MsgOkay)
	if err != nil {
		return err
	}

	err = p.editMessage(chatID, messageID)
	if err != nil {
		return err
	}

	movieNum, err := strconv.Atoi(n)
	if err != nil {
		return err
	}

	return p.showMovie(chatID, username, movieNum+1)
}

func (p *TgProcessor) watchThisMovie(callbackID string, chatID int, messageID int, username string, movieID string) error {
	err := p.tg.AnswerCallbackQuery(callbackID, messages.MsgEnjoyWatching)
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(movieID)

	err = p.storage.Watch(username, id)
	if err != nil {
		return errors.Wrap(err, "removing movie")
	}

	return p.tg.EditMessageReplyMarkup(chatID, messageID)
}

func (p *TgProcessor) editMessage(chatID int, messageID int) error {
	success, err := p.tg.DeleteMessage(chatID, messageID)
	if err != nil {
		return errors.Wrap(err, "error showing more movies")
	}

	if !success {
		if err := p.tg.EditMessageReplyMarkup(chatID, messageID); err != nil {
			return errors.Wrap(err, "error showing more movies")
		}
	}

	return nil
}

func (p *TgProcessor) makeMoviesButtons(movies []events.Movie) ([]telegram.InlineKeyboardButton, error) {
	buttons := make([]telegram.InlineKeyboardButton, 0, len(movies)+1)
	for i, movie := range movies {
		buttonData := fmt.Sprintf("%s;%d", saveButton, movie.ID)
		button, err := messages.MakeButton(strconv.Itoa(i+1), buttonData)
		if err != nil {
			return []telegram.InlineKeyboardButton{}, err
		}
		buttons = append(buttons, button)
	}
	cancelButton, _ := messages.MakeButton(messages.ButtonCancel, canselButton)
	buttons = append(buttons, cancelButton)

	return buttons, nil
}
