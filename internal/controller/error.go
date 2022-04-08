package controller

import "fmt"

// NoActionFoundError occurs when no action found.
type NoActionFoundError struct {
	Name string
}

func (e NoActionFoundError) Error() string {
	return fmt.Sprintf("no action found for %s", e.Name)
}
