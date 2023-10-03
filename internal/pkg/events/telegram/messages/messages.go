package messages

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events"
	"errors"
	"fmt"
)

const (
	MsgHelp  = "There will be help message"
	MsgHello = "Hi there! \n\n" + MsgHelp

	MsgCanNotFindMovie = "I can not find this movie =("
	MsgNoSavedMovies   = "You have no saved movies"
	MsgSaved           = "Saved!"
	MsgAlreadyExists   = "You have already saved this movie in your list"

	MsgSorry = "Sorry =("
)

var ErrTooBigData = errors.New("data is bigger then 64 bytes")

func MovieArrayMessage(movies []events.Movie) string {
	result := ""
	for i, movie := range movies {
		result += fmt.Sprintf("%d. ", i+1)
		result += MovieMessage(movie)
	}

	return result
}

func MovieMessage(movie events.Movie) string {
	result := movie.Title
	if result == "" {
		return ""
	}

	if movie.Year != 0 {
		result += fmt.Sprintf(", %d", movie.Year)
	}
	result += "\n"

	if movie.Rating != 0 {
		result += fmt.Sprintf("IMDb: %.2f\n", movie.Rating)
	}

	if movie.Description != "" {
		result += fmt.Sprintf("Description:\n%s\n", movie.Description)
	}

	if movie.Length != 0 {
		result += fmt.Sprintf("Length: %d minutes\n", movie.Length)
	}

	return result
}

func MakeButton(text string, data string) (telegram.InlineKeyboardButton, error) {
	if len([]byte(data)) > 64 {
		return telegram.InlineKeyboardButton{}, ErrTooBigData
	}

	return telegram.InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
	}, nil
}
