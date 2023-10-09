package messages

import (
	"MovieBot/internal/clients/telegram"
	"MovieBot/internal/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseMessage_ShouldParseMovieToString(t *testing.T) {
	movie := events.Movie{
		Title:       "Film",
		Year:        2005,
		Description: "Desc",
		Rating:      5.6,
		Length:      123,
	}

	msg := MovieMessage(movie)

	assert.Equal(t, msg, "Film, 2005\nIMDb: 5.60\nDescription:\nDesc\nLength: 123 minutes\n")
}

func Test_ParseMessage_ShouldParseMovieWithoutFields(t *testing.T) {
	movie := events.Movie{
		Title:  "Film",
		Year:   2005,
		Rating: 5.6,
		Length: 0,
	}

	msg := MovieMessage(movie)

	assert.Equal(t, msg, "Film, 2005\nIMDb: 5.60\n")
}

func Test_ParseMessage_ShouldParseMovieWithoutTitle(t *testing.T) {
	movie := events.Movie{
		Title:       "",
		Year:        2005,
		Description: "Desc",
		Rating:      5.6,
		Length:      123,
	}

	msg := MovieMessage(movie)

	assert.Equal(t, msg, "")
}

func Test_ParseArrayMessages_ShouldParseMovies(t *testing.T) {
	movies := []events.Movie{{
		Title:       "Film1",
		Year:        2005,
		Description: "Desc",
		Rating:      5.6,
		Length:      123,
	},
		{
			Title:  "Film2",
			Rating: 5.6,
			Length: 123,
		}}

	msg := MovieArrayMessage(movies)

	assert.Equal(t, msg, "1. Film1, 2005\nIMDb: 5.60\nDescription:\nDesc\nLength: 123 minutes\n2. Film2\nIMDb: 5.60\nLength: 123 minutes\n")
}

func Test_ParseArrayMessages_ShouldParseEmptyList(t *testing.T) {
	var movies []events.Movie

	msg := MovieArrayMessage(movies)

	assert.Equal(t, msg, "")
}

func Test_MakeButton_ShouldMakeButton(t *testing.T) {
	button, err := MakeButton("Press", "1543")

	assert.NoError(t, err)
	assert.Equal(t, button, telegram.InlineKeyboardButton{Text: "Press", CallbackData: "1543"})
}

func Test_MakeButton_ShouldReturnErr(t *testing.T) {
	button, err := MakeButton("Press", "15432572358235672657634576345783475638756348756345763457864385765")

	assert.Error(t, err)
	assert.Equal(t, err, ErrTooBigData)
	assert.Equal(t, button, telegram.InlineKeyboardButton{Text: "", CallbackData: ""})
}
