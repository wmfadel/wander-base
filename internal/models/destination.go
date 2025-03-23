package models

import (
	"github.com/twpayne/go-geom"
)

type Destination struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Location    *geom.Point `json:"location"` // PostGIS POINT (longitude, latitude)
	Description string      `json:"description,omitempty"`
}
