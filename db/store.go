package db

import (
	"context"
	"fmt"
	"sort"

	"github.com/geomodulus/citygraph"
	"github.com/geomodulus/citygraph/pb"

	"github.com/google/uuid"
)

type Store struct {
	citygraph.GraphClient
}

func (s *Store) WriteBodyText(ctx context.Context, q *pb.VertexQuery, body string) error {
	return s.SetVertexProperties(ctx, q, citygraph.PropertyNameBodyText, body)
}

func (s *Store) WriteJS(ctx context.Context, q *pb.VertexQuery, fn string) error {
	return s.SetVertexProperties(ctx, q, citygraph.PropertyNameJSFunc, fn)
}

func (s *Store) WriteTeaserGeoJSON(ctx context.Context, q *pb.VertexQuery, teaser map[string]interface{}) error {
	return s.SetVertexProperties(ctx, q, "teaser", teaser)
}

func (s *Store) WriteTeaserJS(ctx context.Context, q *pb.VertexQuery, fn string) error {
	return s.SetVertexProperties(ctx, q, "teaser_function", fn)
}

func (s *Store) WriteArticle(ctx context.Context, article *citygraph.Article) error {
	id, err := article.UUID()
	if err != nil {
		return err
	}
	if err := s.CreateVertex(ctx, citygraph.UUID(id), citygraph.ArticleType); err != nil {
		return nil
	}
	q, err := article.VertexQuery()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameDisplayName, article.Name); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "past_names", article.PastNames); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "headline_html", article.Headline); err != nil {
		return err
	}
	slugID, err := article.SlugID()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "slug_id", slugID); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "slug_title", article.SlugTitle()); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "slug_titles", article.AllSlugTitles()); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "creators", article.Authors); err != nil {
		return err
	}
	if article.Camera != nil {
		if err := s.SetVertexProperties(ctx, q, "camera", article.Camera); err != nil {
			return err
		}
	}
	if err := s.SetVertexProperties(ctx, q, "published_on", article.PubDate); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameUpdatedAt, article.LastUpdated); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "categories", article.Categories); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "code_credit", article.CodeCredit); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "h2", article.Description); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameImgURL, article.FeatureImage); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "pitch", article.Pitch); err != nil {
		return err
	}
	if len(article.Teaser) > 0 {
		if err := s.SetVertexProperties(ctx, q, "teaser", article.Teaser); err != nil {
			return err
		}
	}
	if err := s.SetVertexProperties(ctx, q, "format", article.Format); err != nil {
		return err
	}

	for _, relatedIDStr := range article.Related {
		relatedID, err := uuid.Parse(relatedIDStr)
		if err != nil {
			continue
		}
		if err := s.CreateEdge(ctx, citygraph.UUID(relatedID), &citygraph.IsRelated, citygraph.UUID(id)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) WriteArticleGeoJSONDataset(ctx context.Context, aUUID uuid.UUID, dataset *citygraph.GeoJSONDataset) error {
	id, err := dataset.UUID()
	if err != nil {
		return err
	}
	if err := s.CreateVertex(ctx, citygraph.UUID(id), &citygraph.NewsGeoJSON); err != nil {
		return err
	}

	q, err := dataset.VertexQuery()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameSource, dataset.Source); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "sources", dataset.Sources); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameGeoJSONURL, dataset.URL); err != nil {
		return err
	}
	if err := s.CreateEdge(ctx, citygraph.UUID(aUUID), &citygraph.IllustratedBy, citygraph.UUID(id)); err != nil {
		return err
	}
	return nil
}

func (s *Store) WriteArticleGeoJSONDatasetWithData(ctx context.Context, aUUID uuid.UUID, dataset *citygraph.GeoJSONDataset, data interface{}) error {
	id, err := dataset.UUID()
	if err != nil {
		return err
	}
	if err := s.CreateVertex(ctx, citygraph.UUID(id), &citygraph.NewsGeoJSON); err != nil {
		return err
	}
	q, err := dataset.VertexQuery()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameSource, dataset.Source); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "sources", dataset.Sources); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameGeoJSONFeature, data); err != nil {
		return err
	}
	if err := s.CreateEdge(ctx, citygraph.UUID(aUUID), &citygraph.IllustratedBy, citygraph.UUID(id)); err != nil {
		return err
	}
	return nil
}

func (s *Store) WriteModule(ctx context.Context, module *citygraph.Module) error {
	id, err := module.UUID()
	if err != nil {
		return err
	}
	if err := s.CreateVertex(ctx, citygraph.UUID(id), citygraph.ModuleType); err != nil {
		return nil
	}
	q, err := module.VertexQuery()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameDisplayName, module.Name); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "headline_html", module.Headline); err != nil {
		return err
	}
	slugID, err := module.SlugID()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "slug_id", slugID); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "slug_title", module.SlugTitle()); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "creators", module.Creators); err != nil {
		return err
	}
	if module.Camera != nil {
		if err := s.SetVertexProperties(ctx, q, "camera", module.Camera); err != nil {
			return err
		}
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameFormat, module.Format); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "published_on", module.PubDate); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameUpdatedAt, module.LastUpdated); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "categories", module.Categories); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "code_credit", module.CodeCredit); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "h2", module.Description); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, citygraph.PropertyNameImgURL, module.FeatureImage); err != nil {
		return err
	}
	if len(module.Teaser) > 0 {
		if err := s.SetVertexProperties(ctx, q, "teaser", module.Teaser); err != nil {
			return err
		}
	}

	return nil
}

type ArticleListing struct {
	ID string `json:"id"`
	// Name is the article headline, no HTML tags.
	Name string `json:"name"`
	// Headline is an <h1> tag containing the headline.
	Headline   string   `json:"headline"`
	Subhead    string   `json:"subhead"`
	SlugPath   string   `json:"slug_path"`
	PubDate    string   `json:"pubdate"`
	UpdatedAt  string   `json:"updated_at"`
	Authors    []string `json:"authors"`
	Categories []string `json:"categories"`
	ImgURL     string   `json:"img_url"`
}

func (s *Store) WriteArticleListings(ctx context.Context, _ *Store, articles []*citygraph.Article) error {
	latest := []*ArticleListing{}
	for _, article := range articles {
		if article.PubDate == "" || !article.IsLive {
			continue
		}

		articlePath, err := article.Path()
		if err != nil {
			return fmt.Errorf("failed to get article path: %w", err)
		}
		latest = append(latest, &ArticleListing{
			ID:         article.ID,
			Name:       article.Name,
			Headline:   article.Headline,
			Subhead:    article.Description,
			SlugPath:   articlePath,
			PubDate:    article.PubDate,
			UpdatedAt:  article.LastUpdated,
			Authors:    article.Authors,
			Categories: article.Categories,
			ImgURL:     article.FeatureImage,
		})
	}

	sort.Slice(latest, func(i, j int) bool {
		// Case where both items have UpdatedAt.
		if latest[i].UpdatedAt != "" && latest[j].UpdatedAt != "" {
			return latest[i].UpdatedAt > latest[j].UpdatedAt
		}

		// Case where only i item has UpdatedAt.
		if latest[i].UpdatedAt != "" && latest[j].UpdatedAt == "" {
			return latest[i].UpdatedAt > latest[j].PubDate
		}

		// Case where only j item has UpdatedAt.
		if latest[i].UpdatedAt == "" && latest[j].UpdatedAt != "" {
			return latest[i].PubDate > latest[j].UpdatedAt
		}

		// Case where neither item has UpdatedAt, compare PubDate.
		return latest[i].PubDate > latest[j].PubDate
	})

	if err := s.SetVertexProperties(
		ctx,
		citygraph.NewSpecificVertexQuery(citygraph.Torontoverse.Id),
		"latest_articles_public",
		latest,
	); err != nil {
		return err
	}

	return nil
}
