package processor

import "MovieBot/internal/events"

type Processor interface {
	Process(e events.Event) error
}
