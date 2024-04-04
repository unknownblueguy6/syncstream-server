package room

import "github.com/google/uuid"

func (e *Event) IsValid(id uuid.UUID) bool {
	switch {
	case !(e.Type > ZERO && e.Type <= MESSAGE):
		return false
	case id != e.SourceID:
		return false
	// case e.Data != nil:
	// 	keys := []string{}
	// 	for k, _ := range e.Data {
	// 		keys = append(keys, k)
	// 	}
	// 	slices.SortFunc[]()

	default:
		return true
	}
}
