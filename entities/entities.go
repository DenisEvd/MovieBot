package entities

import "time"

type Movie struct {
	Title       string
	Genre       string
	Rating      float32
	Description string
	Date        time.Time
}
