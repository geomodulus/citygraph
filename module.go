package citygraph

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/geomodulus/citygraph/pb"
)

var ModuleType = &pb.Identifier{Value: "module"}

type Module struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"display_name"`
	Headline     string                 `json:"headline_html,omitempty"`
	Description  string                 `json:"desc"`
	FeatureImage string                 `json:"img_url"`
	Format       string                 `json:"format"`
	Categories   []string               `json:"categories"`
	Creators     []string               `json:"creators"`
	Camera       map[string]interface{} `json:"camera,omitempty"`
	PubDate      string                 `json:"pub_date"`
	LastUpdated  string                 `json:"last_updated,omitempty"`
	CodeCredit   string                 `json:"code_credit"`
	Teaser       map[string]interface{} `json:"teaser,omitempty"`
}

func (a *Module) UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (a *Module) SlugID() (string, error) {
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

func (a *Module) SlugTitle() string {
	titleParts := strings.Split(a.Name, " ")
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

func (a *Module) VertexQuery() (*pb.VertexQuery, error) {
	id, err := a.UUID()
	if err != nil {
		return nil, err
	}
	return NewSpecificVertexQuery(UUID(id)), nil
}
