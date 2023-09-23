package telegram

import (
	"MovieBot/internal/pkg/clients/kinopoisk"
	"fmt"
)

const (
	msgHelp  = "There will be help message"
	msgHello = "Hi there! \n\n" + msgHelp

	msgCanNotFindMovie = "I can not find this movie =("
	msgUnknownCommand  = "Unknown command"
	msgNoSavedMovies   = "You have no saved movies"
	msgSaved           = "Saved!"
	msgAlreadyExists   = "You have already saved this movie in your list"
)

func (p *Processor) movieMessage(movie kinopoisk.Movie) string {
	message := fmt.Sprintf("%s, %d\nIMDb: %.2f\nDescription:\n%s\nLength: %d minutes", movie.Title, movie.Year, movie.Rating.IMDB, movie.Description, movie.MovieLength)

	return message
}
