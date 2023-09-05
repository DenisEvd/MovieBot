package lib

import "fmt"

func Wrap(description string, err error) error {
	return fmt.Errorf("%s: %w", description, err)
}
