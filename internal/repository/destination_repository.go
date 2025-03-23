package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/twpayne/go-geom"
	"github.com/wmfadel/wander-base/internal/models"
)

type DestinationRepository struct {
	db *sql.DB
}

func NewDestinationRepository(db *sql.DB) *DestinationRepository {
	return &DestinationRepository{db: db}
}

func (r *DestinationRepository) Save(destination *models.Destination) error {
	query := `
        INSERT INTO destinations (name, location, description)
        VALUES ($1, ST_GeomFromText($2, 4326), $3) RETURNING id`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare save destination query: %w", err)
	}
	defer stmt.Close()

	wkt := fmt.Sprintf("POINT(%f %f)", destination.Location.X(), destination.Location.Y()) // X=long, Y=lat
	err = stmt.QueryRow(destination.Name, wkt, destination.Description).Scan(&destination.ID)
	if err != nil {
		return fmt.Errorf("failed to save destination: %w", err)
	}
	return nil
}

func (r *DestinationRepository) GetByID(id int64) (*models.Destination, error) {
	query := `
        SELECT id, name, ST_AsText(location), description 
        FROM destinations WHERE id = $1`
	row := r.db.QueryRow(query, id)

	dest := &models.Destination{}
	var locationWKT string
	err := row.Scan(&dest.ID, &dest.Name, &locationWKT, &dest.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get destination %d: %w", id, err)
	}

	if locationWKT != "" {
		point, err := geomFromWKT(locationWKT)
		if err != nil {
			return nil, fmt.Errorf("failed to parse location: %w", err)
		}
		dest.Location = point
	}

	return dest, nil
}

func (r *DestinationRepository) GetAll() ([]models.Destination, error) {
	query := `
        SELECT id, name, ST_AsText(location), description 
        FROM destinations`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query destinations: %w", err)
	}
	defer rows.Close()

	var destinations []models.Destination
	for rows.Next() {
		var dest models.Destination
		var locationWKT string
		err := rows.Scan(&dest.ID, &dest.Name, &locationWKT, &dest.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan destination: %w", err)
		}
		if locationWKT != "" {
			point, err := geomFromWKT(locationWKT)
			if err != nil {
				return nil, fmt.Errorf("failed to parse location: %w", err)
			}
			dest.Location = point
		}
		destinations = append(destinations, dest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating destinations: %w", err)
	}
	return destinations, nil
}

// Helper function to parse WKT
// Updated geomFromWKT to parse WKT manually
func geomFromWKT(wkt string) (*geom.Point, error) {
	// Expecting wkt like "POINT(-119.5383 37.8651)"
	if !strings.HasPrefix(wkt, "POINT(") || !strings.HasSuffix(wkt, ")") {
		return nil, fmt.Errorf("invalid WKT format: %s", wkt)
	}

	// Extract coordinates from "POINT(long lat)"
	coordsStr := strings.TrimPrefix(wkt, "POINT(")
	coordsStr = strings.TrimSuffix(coordsStr, ")")
	coords := strings.Split(coordsStr, " ")

	if len(coords) != 2 {
		return nil, fmt.Errorf("invalid number of coordinates in WKT: %s", wkt)
	}

	longitude, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse longitude: %w", err)
	}

	latitude, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse latitude: %w", err)
	}

	// Create geom.Point with longitude (X) and latitude (Y)
	return geom.NewPointFlat(geom.XY, []float64{longitude, latitude}), nil
}
