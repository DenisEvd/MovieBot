package events

type MovieFetcher interface {
	FetchMovieById(id int) (Movie, error)
	FetchMoviesByTitle(title string) ([]Movie, error)
}

type UpdateFetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	CallbackQuery
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Movie struct {
	ID          int
	Title       string
	Year        int
	Description string
	Poster      string
	Rating      float32
	Length      int
}
