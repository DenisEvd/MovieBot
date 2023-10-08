package messages

import (
	"MovieBot/internal/pkg/clients/telegram"
	"MovieBot/internal/pkg/events"
	"errors"
	"fmt"
	"strings"
)

const (
	MsgHelp  = "There will be help message"
	MsgHello = "Hi there! \n\n" + MsgHelp

	MsgCanNotFindMovie = "I can not find this movie =("
	MsgNoSavedMovies   = "You have no saved movies"
	MsgSaved           = "Saved!"
	MsgAlreadyExists   = "You have already saved this movie in your list"
	MsgAlreadyWatched  = "You have already watched this movie"

	MsgSorry = "Sorry =("
)

var ErrTooBigData = errors.New("data is bigger then 64 bytes")

func MovieArrayMessage(movies []events.Movie) string {
	result := strings.Builder{}
	for i, movie := range movies {
		inf := MovieMessage(movie)
		if inf != "" {
			result.WriteString(fmt.Sprintf("%d. ", i+1))
			result.WriteString(inf)
		}
	}

	return result.String()
}

func MovieMessage(movie events.Movie) string {
	result := strings.Builder{}
	result.WriteString(movie.Title)
	if result.Len() == 0 {
		return ""
	}

	if movie.Year != 0 {
		result.WriteString(fmt.Sprintf(", %d", movie.Year))
	}
	result.WriteString("\n")

	if movie.Rating != 0 {
		result.WriteString(fmt.Sprintf("IMDb: %.2f\n", movie.Rating))
	}

	if movie.Description != "" {
		result.WriteString(fmt.Sprintf("Description:\n%s\n", movie.Description))
	}

	if movie.Length != 0 {
		result.WriteString(fmt.Sprintf("Length: %d minutes\n", movie.Length))
	}

	return result.String()
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
