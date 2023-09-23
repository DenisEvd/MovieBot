package telegram

import (
	"MovieBot/internal/lib"
	"MovieBot/internal/pkg/storage"
	"github.com/pkg/errors"
	"log"
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
		return p.saveMovie(text, chatID, username)
	}
}

func (p *Processor) saveMovie(text string, chatID int, username string) error {
	movies, err := p.kp.FindMovieByTitle(text, 3)
	if err != nil {
		return errors.Wrap(err, "can't take movies from API")
	}

	if len(movies) == 0 {
		return p.tg.SendMessage(chatID, msgCanNotFindMovie)
	}
	var max float32 = 0
	ind := 0
	for i, movie := range movies {
		if movie.Rating > max {
			ind = i
			max = movie.Rating
		}
	}

	isExists, err := p.storage.IsExists(username, movies[ind].ID)
	if err != nil {
		log.Println("here")
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	p.storage.AddMovie(username, &movies[ind])

	return p.tg.SendMessage(chatID, msgSaved)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		if err != nil {
			err = lib.Wrap("can't do command: send random", err)
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

	if movie.Poster.URL != "" {
		if err = p.tg.SendPhoto(chatID, p.movieMessage(movie), movie.Poster.URL); err != nil {
			return err
		}
	} else {
		if err = p.tg.SendMessage(chatID, p.movieMessage(movie)); err != nil {
			return err
		}
	}

	return p.storage.Remove(username, movieID)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
