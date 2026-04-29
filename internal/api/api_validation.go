package api

import (
	"fmt"
	"unicode/utf8"
)

func (req *createEmployeeRequest) validate() error {
	const minAge = 18
	if req.Age < minAge {
		return fmt.Errorf("Age must be greater than %d", minAge)
	}
	const maxAge = 100
	if req.Age > maxAge {
		return fmt.Errorf("Age must be less than or equal to %d", maxAge)
	}

	if len(req.Name) == 0 {
		return fmt.Errorf("Name is required")
	}
	const minNameLength = 2
	if utf8.RuneCountInString(req.Name) < minNameLength {
		return fmt.Errorf("Name must be at lest %d characters long", minNameLength)
	}

	if len(req.Position) == 0 {
		return fmt.Errorf("Position is required")
	}

	return nil
}
