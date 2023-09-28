package telegram

import (
	"MovieBot/internal/pkg/clients/kinopoisk"
	"MovieBot/internal/pkg/clients/telegram"
	"fmt"
	"strconv"
	"unsafe"
)

const (
	msgHelp  = "There will be help message"
	msgHello = "Hi there! \n\n" + msgHelp

	msgCanNotFindMovie = "I can not find this movie =("
	msgNoSavedMovies   = "You have no saved movies"
	msgSaved           = "Saved!"
	msgAlreadyExists   = "You have already saved this movie in your list"

	msgSorry = "Sorry =("
)

func (p *Processor) movieArrayMessageByTitle(movies []kinopoisk.MovieByTitle) string {
	result := ""
	for i, movie := range movies {
		result += fmt.Sprintf("%d. ", i+1)
		result += p.movieMessageByTitle(movie)
	}

	return result
}

func (p *Processor) movieMessageByTitle(movie kinopoisk.MovieByTitle) string {
	result := movie.Title

	if movie.Year != 0 {
		result += ", " + strconv.Itoa(movie.Year)
	}
	result += "\n"

	if movie.Rating != 0 {
		result += fmt.Sprintf("IMDb: %.2f\n", movie.Rating)
	}

	if movie.Description != "" {
		result += fmt.Sprintf("Description: \n%s\n", movie.Description)
	}

	if movie.Length != 0 {
		result += fmt.Sprintf("Length: %d minutes\n", movie.Length)
	}

	return result
}

func (p *Processor) movieMessageByID(movie kinopoisk.MovieByID) string {
	message := fmt.Sprintf("%s, %d\nIMDb: %.2f\nDescription:\n%s\nLength: %d minutes", movie.Title, movie.Year, movie.Rating.IMDB, movie.Description, movie.MovieLength)

	return message
}

func (p *Processor) makeButton(text string, data string) (telegram.InlineKeyboardButton, error) {
	if unsafe.Sizeof(data) > 64 {
		return telegram.InlineKeyboardButton{}, ErrTooBigData
	}

	return telegram.InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
	}, nil
}
