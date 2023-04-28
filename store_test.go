package citygraph

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/geomodulus/citygraph/pb"
)

func TestWriteBodyText(t *testing.T) {
	aID := NewID()
	aUUID := UUID(aID)
	fakeGraph := &FakeGraphClient{}
	store := &Store{fakeGraph}

	module := &Module{
		ID: aID.String(),
	}
	q, _ := module.VertexQuery()
	if err := store.WriteBodyText(context.Background(), q, "body text"); err != nil {
		t.Errorf("store.WriteBodyText() returned err: %v", err)
	}

	avq := NewSpecificVertexQuery(aUUID)
	wantReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     NewVertexPropertyQuery(avq, PropertyNameBodyText),
		Value: StringVal("body text"),
	}}
	if diff := cmp.Diff(wantReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteBodyText() sent graph req diff:\n%s\n", diff)
	}
}

func TestWriteJS(t *testing.T) {
	aID := NewID()
	aUUID := UUID(aID)
	fakeGraph := &FakeGraphClient{}
	store := &Store{fakeGraph}

	module := &Module{
		ID: aID.String(),
	}
	q, _ := module.VertexQuery()
	if err := store.WriteJS(context.Background(), q, "js function"); err != nil {
		t.Errorf("store.WriteJS() returned err: %v", err)
	}

	avq := NewSpecificVertexQuery(aUUID)
	wantReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     NewVertexPropertyQuery(avq, PropertyNameJSFunc),
		Value: StringVal("js function"),
	}}
	if diff := cmp.Diff(wantReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteJS() sent graph req diff:\n%s\n", diff)
	}
}

func TestWriteModule(t *testing.T) {
	aID := NewID()
	aUUID := UUID(aID)
	fakeGraph := &FakeGraphClient{}
	store := &Store{fakeGraph}

	module := &Module{
		ID:       aID.String(),
		Name:     "Headline **with emphasis** here",
		Headline: "Headline <em>with emphasis</em> here",
		Camera: map[string]interface{}{
			"sm": map[string]interface{}{
				"zoom": 10.5,
			},
		},
		Creators:     []string{"John Dole", "Jane Dole"},
		PubDate:      "2022-06-14",
		LastUpdated:  "2022-09-01T10:30:00-05:00",
		Categories:   []string{"Open Data", "User-Generated"},
		CodeCredit:   "some programmer",
		Description:  "some description",
		FeatureImage: "some url",
		Teaser: map[string]interface{}{
			"type": "FeatureCollection",
		},
	}
	if err := store.WriteModule(context.Background(), module); err != nil {
		t.Errorf("store.WriteModule() returned err: %v", err)
	}

	slugID, _ := module.SlugID()

	creatorsBytes, _ := json.Marshal(module.Creators)
	cameraBytes, _ := json.Marshal(module.Camera)
	catBytes, _ := json.Marshal(module.Categories)
	creditBytes, _ := json.Marshal(module.CodeCredit)
	descBytes, _ := json.Marshal(module.Description)
	featImgBytes, _ := json.Marshal(module.FeatureImage)
	teaserBytes, _ := json.Marshal(module.Teaser)

	avq := NewSpecificVertexQuery(aUUID)
	wantCreateVertexReqs := []*pb.Vertex{{
		Id: aUUID,
		T:  moduleType,
	}}
	if diff := cmp.Diff(wantCreateVertexReqs, fakeGraph.CreateVertexReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteModule() sent create vertex graph req diff:\n%s\n", diff)
	}

	wantSetVertexPropertiesReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     NewVertexPropertyQuery(avq, "display_name"),
		Value: StringVal(module.Name),
	}, {
		Q:     NewVertexPropertyQuery(avq, "headline_html"),
		Value: StringVal(module.Headline),
	}, {
		Q:     NewVertexPropertyQuery(avq, "slug_id"),
		Value: StringVal(slugID),
	}, {
		Q:     NewVertexPropertyQuery(avq, "slug_title"),
		Value: StringVal(module.SlugTitle()),
	}, {
		Q:     NewVertexPropertyQuery(avq, "creators"),
		Value: Json(creatorsBytes),
	}, {
		Q:     NewVertexPropertyQuery(avq, "camera"),
		Value: Json(cameraBytes),
	}, {
		Q:     NewVertexPropertyQuery(avq, "published_on"),
		Value: StringVal(module.PubDate),
	}, {
		Q:     NewVertexPropertyQuery(avq, PropertyNameUpdatedAt),
		Value: StringVal(module.LastUpdated),
	}, {
		Q:     NewVertexPropertyQuery(avq, "categories"),
		Value: Json(catBytes),
	}, {
		Q:     NewVertexPropertyQuery(avq, "code_credit"),
		Value: Json(creditBytes),
	}, {
		Q:     NewVertexPropertyQuery(avq, "h2"),
		Value: Json(descBytes),
	}, {
		Q:     NewVertexPropertyQuery(avq, "img_url"),
		Value: Json(featImgBytes),
	}, {
		Q:     NewVertexPropertyQuery(avq, "teaser"),
		Value: Json(teaserBytes),
	}}
	if diff := cmp.Diff(wantSetVertexPropertiesReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteModule() sent set vertex property graph req diff:\n%s\n", diff)
	}
}
