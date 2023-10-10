package telegram

import (
	"MovieBot/internal/clients/telegram"
	mockTg "MovieBot/internal/clients/telegram/mock"
	"MovieBot/internal/events"
	mockMovieFetcher "MovieBot/internal/events/movie_fetcher/mock"
	"MovieBot/internal/events/processor/messages"
	"MovieBot/internal/events/tg_fetcher"
	"MovieBot/internal/storage"
	mockRepo "MovieBot/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Test_Process_Start(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "/start",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().SendMessage(105, messages.MsgHello).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_Help(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "/help",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().SendMessage(105, messages.MsgHelp).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_SuggestMovie_ShouldSendMovie(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	movies := []events.Movie{{
		ID:          13,
		Title:       "Movie",
		Description: "Desc",
		Poster:      "url",
	}}

	event := events.Event{
		Type: events.Message,
		Text: "Movie",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	message := messages.MovieMessage(movies[0])

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	movieFetcher.EXPECT().FetchMoviesByTitle("Movie").Return(movies, nil)
	rep.EXPECT().AddRequest("Movie").Return(1, nil)
	tg.EXPECT().SendPhotoWithInlineKeyboard(105, message, "url", []telegram.InlineKeyboardButton{
		{Text: messages.ButtonNo, CallbackData: findMoreButton + ";1"},
		{Text: messages.ButtonYes, CallbackData: saveButton + ";13" + ";1"},
	})

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_SuggestMovie_ShouldSendCanNotFindMessage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "Movie",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	movieFetcher.EXPECT().FetchMoviesByTitle("Movie").Return([]events.Movie{}, nil)
	tg.EXPECT().SendMessage(105, messages.MsgCanNotFindMovie)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_PickN_ShouldSendMovieWithPoster(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "/show 5",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	rep.EXPECT().GetNMovie("user", 5).Return(events.Movie{
		ID:     15,
		Title:  "Movie",
		Year:   2005,
		Poster: "url",
		Rating: 5.7,
	}, nil)
	tg.EXPECT().SendPhotoWithInlineKeyboard(105, messages.MovieMessage(events.Movie{
		ID:     15,
		Title:  "Movie",
		Year:   2005,
		Rating: 5.7,
	}), "url", []telegram.InlineKeyboardButton{
		{Text: "Next", CallbackData: "next;5"},
		{Text: "Watch it!", CallbackData: "watch;15"},
		{Text: messages.ButtonCancel, CallbackData: "no"},
	}).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_PickN_ShouldSendMovieWithoutPoster(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "/show 5",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	rep.EXPECT().GetNMovie("user", 5).Return(events.Movie{
		ID:     15,
		Title:  "Movie",
		Year:   2005,
		Rating: 5.7,
	}, nil)
	tg.EXPECT().SendMessageWithInlineKeyboard(105, messages.MovieMessage(events.Movie{
		ID:     15,
		Title:  "Movie",
		Year:   2005,
		Rating: 5.7,
	}), []telegram.InlineKeyboardButton{
		{Text: "Next", CallbackData: "next;5"},
		{Text: "Watch it!", CallbackData: "watch;15"},
		{Text: "Cancel", CallbackData: "no"},
	}).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_GetAll_ShouldSendList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "/all",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	movies := []events.Movie{{
		ID:     15,
		Title:  "Movie",
		Year:   2005,
		Poster: "url",
		Rating: 5.7,
	},
		{
			ID:     7,
			Title:  "Movie 1",
			Year:   2006,
			Rating: 6.1,
		}}

	rep.EXPECT().GetAll("user").Return(movies, nil)

	tg.EXPECT().SendMessage(105, messages.HeaderOfMovieList+messages.MovieArrayMessage(movies)).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_GetAll_ShouldSendNoMoviesMessage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.Message,
		Text: "/all",
		Meta: tg_fetcher.MessageMeta{
			ChatID:   105,
			Username: "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	rep.EXPECT().GetAll("user").Return([]events.Movie{}, storage.ErrNoSavedMovies)

	tg.EXPECT().SendMessage(105, messages.MsgNoSavedMovies)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_editMessage_ShouldDelete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.editMessage(5, 10)

	assert.NoError(t, err)
}

func Test_Process_editMessage_ShouldEdit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().DeleteMessage(5, 10).Return(false, nil)
	tg.EXPECT().EditMessageReplyMarkup(5, 10).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.editMessage(5, 10)

	assert.NoError(t, err)
}

func Test_Process_editMessage_ShouldReturnError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().DeleteMessage(5, 10).Return(false, nil)
	tg.EXPECT().EditMessageReplyMarkup(5, 10).Return(errors.New("err"))

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.editMessage(5, 10)

	assert.Error(t, err)
}

func Test_Process_CallbackCancel(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.CallbackQuery,
		Text: "no",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgOkay).Return(nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackShowMoreMovies_ShouldSendMoviesList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.CallbackQuery,
		Text: findMoreButton + ";5",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	movies := []events.Movie{
		{
			ID:     22,
			Title:  "Serial",
			Rating: 5.2,
		},
		{
			ID:    11,
			Title: "Movie",
			Year:  2007,
		},
	}

	sort.Slice(movies, func(i, j int) bool { return movies[i].Rating > movies[j].Rating })

	message := messages.MovieArrayMessage(movies)

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgOkay).Return(nil)
	rep.EXPECT().DeleteRequest(5).Return("Movie", nil)
	movieFetcher.EXPECT().FetchMoviesByTitle("Movie").Return(movies, nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)
	tg.EXPECT().SendMessageWithInlineKeyboard(5, messages.HeaderOfMoreMovieList+message, []telegram.InlineKeyboardButton{
		{Text: "1", CallbackData: saveButton + ";22"},
		{Text: "2", CallbackData: saveButton + ";11"},
		{Text: messages.ButtonCancel, CallbackData: canselButton},
	})

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackShowMoreMovies_ShouldSendMessageCanNotFindMovies(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.CallbackQuery,
		Text: findMoreButton + ";5",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgOkay).Return(nil)
	rep.EXPECT().DeleteRequest(5).Return("Movie", nil)
	movieFetcher.EXPECT().FetchMoviesByTitle("Movie").Return([]events.Movie{}, nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)
	tg.EXPECT().SendMessage(5, messages.MsgCanNotFindMovie)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackSaveMovie_ShouldSaveMovieAndDeleteRequest(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	movie := events.Movie{
		ID:     22,
		Title:  "Serial",
		Rating: 5.2,
	}

	event := events.Event{
		Type: events.CallbackQuery,
		Text: saveButton + ";5" + ";1",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgSaved).Return(nil)
	rep.EXPECT().DeleteRequest(1).Return("", nil)
	rep.EXPECT().IsExistRecord("user", 5).Return(false, nil)
	movieFetcher.EXPECT().FetchMovieById(5).Return(movie, nil)
	rep.EXPECT().AddMovie("user", movie).Return(nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackSaveMovie_ShouldSendMessageIsExists(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.CallbackQuery,
		Text: saveButton + ";5",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgSaved).Return(nil)
	rep.EXPECT().IsExistRecord("user", 5).Return(true, nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)
	rep.EXPECT().IsWatched("user", 5).Return(false, nil)
	tg.EXPECT().SendMessage(5, messages.MsgAlreadyExists)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackSaveMovie_ShouldSendMessageIsWatched(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.CallbackQuery,
		Text: saveButton + ";5",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgSaved).Return(nil)
	rep.EXPECT().IsExistRecord("user", 5).Return(true, nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)
	rep.EXPECT().IsWatched("user", 5).Return(true, nil)
	tg.EXPECT().SendMessage(5, messages.MsgAlreadyWatched)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackShowNextMovie_ShouldSendMovie(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	movie := events.Movie{
		ID:    2,
		Title: "Movie",
		Year:  1998,
	}

	event := events.Event{
		Type: events.CallbackQuery,
		Text: getNextButton + ";5",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	message := messages.MovieMessage(movie)

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgOkay).Return(nil)
	tg.EXPECT().DeleteMessage(5, 10).Return(true, nil)
	rep.EXPECT().GetNMovie("user", 6).Return(movie, nil)
	tg.EXPECT().SendMessageWithInlineKeyboard(5, message, []telegram.InlineKeyboardButton{
		{Text: messages.ButtonNext, CallbackData: getNextButton + ";6"},
		{Text: messages.ButtonWatch, CallbackData: watchItButton + ";2"},
		{Text: messages.ButtonCancel, CallbackData: canselButton},
	})

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}

func Test_Process_CallbackWatchIt_ShouldMarkMovieWatched(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	event := events.Event{
		Type: events.CallbackQuery,
		Text: watchItButton + ";5",
		Meta: tg_fetcher.CallbackMeta{
			CallbackID: "123",
			ChatID:     5,
			MessageID:  10,
			Username:   "user",
		},
	}

	rep := mockRepo.NewMockStorage(ctl)
	tg := mockTg.NewMockTgClient(ctl)
	movieFetcher := mockMovieFetcher.NewMockMovieFetcher(ctl)

	tg.EXPECT().AnswerCallbackQuery("123", messages.MsgEnjoyWatching).Return(nil)
	rep.EXPECT().Watch("user", 5).Return(nil)
	tg.EXPECT().EditMessageReplyMarkup(5, 10).Return(nil)

	processor := NewTgProcessor(tg, movieFetcher, rep)
	err := processor.Process(event)
	assert.NoError(t, err)
}
