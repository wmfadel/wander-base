package requests

import "time"

type PatchEvent struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	DateTime    *time.Time `json:"dateTime"`
}

func (pe PatchEvent) IsEmpty() bool {
	if pe.Name == nil && pe.Description == nil && pe.Location == nil && pe.DateTime == nil {
		return true
	}
	return false
}
