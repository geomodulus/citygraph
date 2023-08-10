package citygraph

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/geomodulus/citygraph/pb"
)

var (
	PlaceType = pb.Identifier{Value: "poi"}

	HasOtherLocation = &pb.Identifier{Value: "has-other-location"}
)

type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

// Restaurant contains details specific to restaurant POIs.
type Restaurant struct {
	Style         string `json:"style,omitempty"`
	OpenTable     string `json:"open_table,omitempty"`
	Ritual        string `json:"ritual,omitempty"`
	SkipTheDishes string `json:"skip_the_dishes,omitempty"`
	Custom        string `json:"custom,omitempty"`
}

// Bar contains details specific to bar POIs.
type Bar struct {
	Style string `json:"style,omitempty"`
}

// Cafe contains details specific to cafe POIs.
type Cafe struct {
	Style string `json:"style,omitempty"`
}

type PlaceLocation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Status is all OK if missing/empty otherwise "Closed" and "Temporarily closed".
	Status string `json:"status,omitempty"`
	// PhoneNumber should be formatted like this: "416-902-7200".
	PhoneNumber string `json:"phone_number"`
	// Address holds the street address of the location in the form "123 Sesame St W".
	Address  string   `json:"street_address"`
	City     string   `json:"city,omitempty"`
	Location LngLat   `json:"location"`
	Images   []*Image `json:"images,omitempty"`
	// Restaurant is optional, will be present only for some restaurant-type locations.
	Restaurant *Restaurant `json:"restaurant"`
	// Bar is optional, will be present only for some bar-type locations.
	Bar *Bar `json:"bar"`
	// Cafe is optional, will be present only for some cafe-type locations.
	Cafe *Cafe `json:"cafe"`
}

func (p *PlaceLocation) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (p *PlaceLocation) SlugID() (string, error) {
	id, err := p.UUID()
	if err != nil {
		return "", err
	}
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if _, err := b.WriteString(base64.URLEncoding.EncodeToString(idBytes)[:22]); err != nil {
		return "", fmt.Errorf("base64 encode: %v", err)
	}
	return b.String(), err
}

func (p *PlaceLocation) VertexQuery() (*pb.VertexQuery, error) {
	id, err := p.UUID()
	if err != nil {
		return nil, err
	}
	return NewSpecificVertexQuery(UUID(id)), nil
}

type Place struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Desc    string    `json:"description"`
	AddedAt time.Time `json:"added_at"`
	// PhoneNumber should be formatted like this: "416-902-7200".
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"street_address"`
	// Status is all OK if missing/empty otherwise "Closed" and "Temporarily closed".
	Status string `json:"status,omitempty"`
	// Established records when the business was established.
	Established string `json:"established"`
	// Instagram is the business or location's Instagram username.
	Instagram string `json:"instagram"`
	// Twitter is the business or location's Twitter username.
	Twitter   string `json:"twitter"`
	ForFeed   bool   `json:"for_feed"`
	URL       string `json:"url"`
	ShortName string `json:"short_name"`
	Type      string `json:"type"`
	// Restaurant is optional, will be present only for some restaurant-type locations.
	Restaurant *Restaurant `json:"restaurant"`
	// Bar is optional, will be present only for some bar-type locations.
	Bar *Bar `json:"bar"`
	// Cafe is optional, will be present only for some cafe-type locations.
	Cafe      *Cafe            `json:"cafe"`
	Sprite    string           `json:"sprite"`
	Location  LngLat           `json:"location"`
	Images    []*Image         `json:"images,omitempty"`
	Locations []*PlaceLocation `json:"locations"`
}

func (p *Place) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (p *Place) SlugTitle() string {
	useName := p.Name
	if p.ShortName != "" {
		useName = p.ShortName
	}
	titleParts := strings.Split(useName, " ")
	if len(titleParts) > 5 {
		titleParts = titleParts[:5]
	}
	cutAndMashed := strings.ToLower(strings.Join(titleParts, "-"))
	sanitized := url.PathEscape(punctRE.ReplaceAllString(cutAndMashed, ""))

	return sanitized
}

func (p *Place) VertexQuery() (*pb.VertexQuery, error) {
	id, err := p.UUID()
	if err != nil {
		return nil, err
	}
	return NewSpecificVertexQuery(UUID(id)), nil
}
