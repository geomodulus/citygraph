package db

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/geomodulus/citygraph"
	"github.com/geomodulus/citygraph/graphtest"
	"github.com/geomodulus/citygraph/pb"
)

func TestWriteBodyText(t *testing.T) {
	aID := citygraph.NewID()
	aUUID := citygraph.UUID(aID)
	fakeGraph := &graphtest.FakeGraphClient{}
	store := &Store{fakeGraph}

	module := &citygraph.Module{
		ID: aID.String(),
	}
	q, _ := module.VertexQuery()
	if err := store.WriteBodyText(context.Background(), q, "body text"); err != nil {
		t.Errorf("store.WriteBodyText() returned err: %v", err)
	}

	avq := citygraph.NewSpecificVertexQuery(aUUID)
	wantReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     citygraph.NewVertexPropertyQuery(avq, citygraph.PropertyNameBodyText),
		Value: citygraph.StringVal("body text"),
	}}
	if diff := cmp.Diff(wantReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteBodyText() sent graph req diff:\n%s\n", diff)
	}
}

func TestWriteJS(t *testing.T) {
	aID := citygraph.NewID()
	aUUID := citygraph.UUID(aID)
	fakeGraph := &graphtest.FakeGraphClient{}
	store := &Store{fakeGraph}

	module := &citygraph.Module{
		ID: aID.String(),
	}
	q, _ := module.VertexQuery()
	if err := store.WriteJS(context.Background(), q, "js function"); err != nil {
		t.Errorf("store.WriteJS() returned err: %v", err)
	}

	avq := citygraph.NewSpecificVertexQuery(aUUID)
	wantReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     citygraph.NewVertexPropertyQuery(avq, citygraph.PropertyNameJSFunc),
		Value: citygraph.StringVal("js function"),
	}}
	if diff := cmp.Diff(wantReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteJS() sent graph req diff:\n%s\n", diff)
	}
}

func TestWriteArticle(t *testing.T) {
	aID, bID := citygraph.NewID(), citygraph.NewID()
	aUUID := citygraph.UUID(aID)
	fakeGraph := &graphtest.FakeGraphClient{}
	store := &Store{fakeGraph}

	article := &citygraph.Article{
		ID:        aID.String(),
		Name:      "Headline **with emphasis** here",
		PastNames: []string{"Headline where we forgot emphasis here"},
		Headline:  "Headline <em>with emphasis</em> here",
		Authors:   []string{"Raoul Duke", "Hunter Thompson"},
		Camera: map[string]interface{}{
			"sm": map[string]interface{}{
				"zoom": 10.5,
			},
		},
		PubDate:      "2022-06-14",
		LastUpdated:  "2022-09-01T10:30:00-05:00",
		Categories:   []string{"Unit testing"},
		CodeCredit:   "some programmer",
		Description:  "some description",
		FeatureImage: "some url",
		Pitch:        85,
		Related:      []string{bID.String()},
		Teaser: map[string]interface{}{
			"type": "FeatureCollection",
		},
		Format: "content-fullscreen",
	}
	if err := store.WriteArticle(context.Background(), article); err != nil {
		t.Errorf("store.WriteArticle() returned err: %v", err)
	}

	slugID, _ := article.SlugID()

	pastNamesBytes, _ := json.Marshal(article.PastNames)
	allSlugTitlesBytes, _ := json.Marshal(article.AllSlugTitles())
	authorsBytes, _ := json.Marshal(article.Authors)
	cameraBytes, _ := json.Marshal(article.Camera)
	catBytes, _ := json.Marshal(article.Categories)
	creditBytes, _ := json.Marshal(article.CodeCredit)
	descBytes, _ := json.Marshal(article.Description)
	featImgBytes, _ := json.Marshal(article.FeatureImage)
	pitchBytes, _ := json.Marshal(article.Pitch)
	teaserBytes, _ := json.Marshal(article.Teaser)

	avq := citygraph.NewSpecificVertexQuery(aUUID)
	wantCreateVertexReqs := []*pb.Vertex{{
		Id: aUUID,
		T:  citygraph.ArticleType,
	}}
	if diff := cmp.Diff(wantCreateVertexReqs, fakeGraph.CreateVertexReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteArticle() sent create vertex graph req diff:\n%s\n", diff)
	}

	wantSetVertexPropertiesReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     citygraph.NewVertexPropertyQuery(avq, "display_name"),
		Value: citygraph.StringVal(article.Name),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "past_names"),
		Value: citygraph.Json(pastNamesBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "headline_html"),
		Value: citygraph.StringVal(article.Headline),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "slug_id"),
		Value: citygraph.StringVal(slugID),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "slug_title"),
		Value: citygraph.StringVal(article.SlugTitle()),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "slug_titles"),
		Value: citygraph.Json(allSlugTitlesBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "creators"),
		Value: citygraph.Json(authorsBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "camera"),
		Value: citygraph.Json(cameraBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "published_on"),
		Value: citygraph.StringVal(article.PubDate),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, citygraph.PropertyNameUpdatedAt),
		Value: citygraph.StringVal(article.LastUpdated),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "categories"),
		Value: citygraph.Json(catBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "code_credit"),
		Value: citygraph.Json(creditBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "h2"),
		Value: citygraph.Json(descBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "img_url"),
		Value: citygraph.Json(featImgBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "pitch"),
		Value: citygraph.Json(pitchBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "teaser"),
		Value: citygraph.Json(teaserBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "format"),
		Value: citygraph.StringVal(article.Format),
	}}
	if diff := cmp.Diff(wantSetVertexPropertiesReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteArticle() sent set vertex property graph req diff:\n%s\n", diff)
	}

	wantCreateEdgeReqs := []*pb.EdgeKey{{
		OutboundId: citygraph.UUID(bID),
		T:          &citygraph.IsRelated,
		InboundId:  citygraph.UUID(aID),
	}}
	if diff := cmp.Diff(wantCreateEdgeReqs, fakeGraph.CreateEdgeReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteArticle() sent graph create edge req diff:\n%s\n", diff)
	}
}

//func TestWriteFeedEntryHTML(t *testing.T) {
//	aID := citygraph.NewID()
//	aUUID := citygraph.UUID(aID)
//	fakeGraph := &graphtest.FakeGraphClient{}
//	store := &Store{fakeGraph}
//
//	article := &citygraph.Article{
//		ID: aID.String(),
//	}
//	q, _ := article.VertexQuery()
//	if err := store.WriteFeedEntryHTML(context.Background(), q, "feed entry html"); err != nil {
//		t.Errorf("store.WriteFeedEntryHTML() returned err: %v", err)
//	}
//
//	avq := citygraph.NewSpecificVertexQuery(aUUID)
//	wantReqs := []*pb.SetVertexPropertiesRequest{{
//		Q:     citygraph.NewVertexPropertyQuery(avq, "feed_entry_html"),
//		Value: citygraph.StringVal("feed entry html"),
//	}}
//	if diff := cmp.Diff(wantReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
//		t.Errorf("store.WriteFeedEntryHTML() sent graph req diff:\n%s\n", diff)
//	}
//}

func TestWriteModule(t *testing.T) {
	aID := citygraph.NewID()
	aUUID := citygraph.UUID(aID)
	fakeGraph := &graphtest.FakeGraphClient{}
	store := &Store{fakeGraph}

	module := &citygraph.Module{
		ID:       aID.String(),
		Name:     "Headline **with emphasis** here",
		Headline: "Headline <em>with emphasis</em> here",
		Camera: map[string]interface{}{
			"sm": map[string]interface{}{
				"zoom": 10.5,
			},
		},
		Creators:     []string{"John Dole", "Jane Dole"},
		Format:       "content-map",
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

	avq := citygraph.NewSpecificVertexQuery(aUUID)
	wantCreateVertexReqs := []*pb.Vertex{{
		Id: aUUID,
		T:  citygraph.ModuleType,
	}}
	if diff := cmp.Diff(wantCreateVertexReqs, fakeGraph.CreateVertexReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteModule() sent create vertex graph req diff:\n%s\n", diff)
	}

	wantSetVertexPropertiesReqs := []*pb.SetVertexPropertiesRequest{{
		Q:     citygraph.NewVertexPropertyQuery(avq, "display_name"),
		Value: citygraph.StringVal(module.Name),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "headline_html"),
		Value: citygraph.StringVal(module.Headline),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "slug_id"),
		Value: citygraph.StringVal(slugID),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "slug_title"),
		Value: citygraph.StringVal(module.SlugTitle()),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "creators"),
		Value: citygraph.Json(creatorsBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "camera"),
		Value: citygraph.Json(cameraBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, citygraph.PropertyNameFormat),
		Value: citygraph.StringVal(module.Format),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "published_on"),
		Value: citygraph.StringVal(module.PubDate),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, citygraph.PropertyNameUpdatedAt),
		Value: citygraph.StringVal(module.LastUpdated),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "categories"),
		Value: citygraph.Json(catBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "code_credit"),
		Value: citygraph.Json(creditBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "h2"),
		Value: citygraph.Json(descBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "img_url"),
		Value: citygraph.Json(featImgBytes),
	}, {
		Q:     citygraph.NewVertexPropertyQuery(avq, "teaser"),
		Value: citygraph.Json(teaserBytes),
	}}
	if diff := cmp.Diff(wantSetVertexPropertiesReqs, fakeGraph.SetVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("store.WriteModule() sent set vertex property graph req diff:\n%s\n", diff)
	}
}
