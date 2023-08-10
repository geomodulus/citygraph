package citygraph

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/geomodulus/citygraph/pb"
)

var (
	ArticleType = &pb.Identifier{Value: "news-contributor-article"}
	IsRelated   = pb.Identifier{Value: "is-related"}
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Article struct {
	LoadedFrom   string                 `json:"-"`
	ID           string                 `json:"id"`
	Name         string                 `json:"display_name"`
	PastNames    []string               `json:"past_names"`
	Headline     string                 `json:"headline_html,omitempty"`
	Description  string                 `json:"h2"`
	FeatureImage string                 `json:"img_url"`
	IsLive       bool                   `json:"is_live"`
	Categories   []string               `json:"categories"`
	Authors      []string               `json:"authors"`
	Pitch        float64                `json:"pitch,omitempty"`
	Camera       map[string]interface{} `json:"camera,omitempty"`
	Promo        string                 `json:"promo,omitempty"`
	PubDate      string                 `json:"pub_date"`
	LastUpdated  string                 `json:"last_updated,omitempty"`
	CodeCredit   string                 `json:"code_credit"`
	PromoAt      time.Time              `json:"promo_start"`
	PromoUntil   time.Time              `json:"promo_until"`
	PromoWait    Duration               `json:"promo_wait"`
	PromoExpires time.Time              `json:"promo_expires,omitempty"`
	Slug         string                 `json:"slug,omitempty"`
	//	Loc             citygraph.LngLat       `json:"loc,omitempty"`
	//	Zoom            float64                `json:"zoom,omitempty"`
	Related         []string               `json:"related,omitempty"`
	GeoJSONDatasets []*GeoJSONDataset      `json:"geojson_datasets,omitempty"`
	Teaser          map[string]interface{} `json:"teaser,omitempty"`
	Format          string                 `json:"format,omitempty"`
}

func (a *Article) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (a *Article) SlugID() (string, error) {
	id, err := a.UUID()
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

func (a *Article) LoadBodyText() (string, error) {
	body, err := os.ReadFile(filepath.Join(a.LoadedFrom, "article.html"))
	if err != nil {
		return "", fmt.Errorf("error reading article body: %v", err)
	}
	return string(body), nil
}

func makeSlugTitle(title string) string {
	titleParts := strings.Split(title, " ")
	if len(titleParts) > 5 {
		titleParts = titleParts[:5]
	}
	cutAndMashed := strings.ToLower(strings.Join(titleParts, "-"))
	sanitized := strings.TrimRight(
		url.PathEscape(punctRE.ReplaceAllString(cutAndMashed, "")),
		".",
	)

	return sanitized
}

func (a *Article) SlugTitle() string {
	return makeSlugTitle(a.Name)
}

// AllSlugTitles includes past slug titles form former headlines that may still be out there as
// valid links.
func (a *Article) AllSlugTitles() (slugs []string) {
	for _, name := range append(a.PastNames, a.Name) {
		slugs = append(slugs, makeSlugTitle(name))
	}
	return
}

func (a *Article) Path() (string, error) {
	slugID, err := a.SlugID()
	if err != nil {
		return "", fmt.Errorf("slug id: %v", err)
	}
	return fmt.Sprintf("/articles/%s/%s", slugID, a.SlugTitle()), nil
}

func (a *Article) VertexQuery() (*pb.VertexQuery, error) {
	id, err := a.UUID()
	if err != nil {
		return nil, err
	}
	return NewSpecificVertexQuery(UUID(id)), nil
}
