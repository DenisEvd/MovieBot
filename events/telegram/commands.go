package telegram

import (
	"MovieBot/entities"
	"MovieBot/lib"
	"MovieBot/storage"
	"errors"
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

func (p *Processor) saveMovie(text string, chatID int, username string) (err error) {
	defer func() { err = lib.Wrap("can't do command: save page", err) }()

	movie := &entities.Movie{
		Title: text,
	}

	isExists, err := p.storage.IsExists(username, movie)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	p.storage.AddMovie(username, movie)

	return p.tg.SendMessage(chatID, msgSaved)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = lib.Wrap("can't do command: send random", err) }()

	movie, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, msgNoSavedMovies)
	}

	if err = p.tg.SendMessage(chatID, movie.Title); err != nil {
		return err
	}

	return p.storage.Remove(username, movie)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
