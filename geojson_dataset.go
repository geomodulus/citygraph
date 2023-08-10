package citygraph

import (
	"github.com/google/uuid"

	"github.com/geomodulus/citygraph/pb"
)

type Source struct {
	Title        string `json:"title"`
	Provider     string `json:"provider"`
	URL          string `json:"url"`
	DateAcquired string `json:"date_acquired"`
	Description  string `json:"description"`
	License      string `json:"license"`
	LicenseURL   string `json:"license_url"`
}

type GeoJSONDataset struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	Render  string    `json:"render,omitempty"`
	Source  *Source   `json:"source,omitempty"`
	Sources []*Source `json:"sources,omitempty"`
}

func (d *GeoJSONDataset) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (d *GeoJSONDataset) VertexQuery() (*pb.VertexQuery, error) {
	id, err := d.UUID()
	if err != nil {
		return nil, err
	}
	return NewSpecificVertexQuery(UUID(id)), nil
}
