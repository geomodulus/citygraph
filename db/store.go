package db

import (
	"context"
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
