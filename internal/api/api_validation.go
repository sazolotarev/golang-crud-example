package api

import (
	"fmt"
	"unicode/utf8"
)

const (
	minAge        = 18
	maxAge        = 100
	minNameLength = 2
)

func (req *createEmployeeRequest) validate() error {
	if req.Age < minAge {
		return fmt.Errorf("Age must be greater than %d", minAge)
	}
	if req.Age > maxAge {
		return fmt.Errorf("Age must be less than or equal to %d", maxAge)
	}
	if len(req.Name) == 0 {
		return fmt.Errorf("Name is required")
	}
	if utf8.RuneCountInString(req.Name) < minNameLength {
		return fmt.Errorf("Name must be at least %d characters long", minNameLength)
	}
	if len(req.Position) == 0 {
		return fmt.Errorf("Position is required")
	}
	return nil
}
