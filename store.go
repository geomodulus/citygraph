package citygraph

import (
	"context"
	"github.com/geomodulus/citygraph/pb"
	//"github.com/geomodulus/witness/zones"
)

type Store struct {
	GraphClient
}

func (s *Store) WriteBodyText(ctx context.Context, q *pb.VertexQuery, body string) error {
	return s.SetVertexProperties(ctx, q, PropertyNameBodyText, body)
}

func (s *Store) WriteJS(ctx context.Context, q *pb.VertexQuery, fn string) error {
	return s.SetVertexProperties(ctx, q, PropertyNameJSFunc, fn)
}

func (s *Store) WriteModule(ctx context.Context, module *Module) error {
	id, err := module.UUID()
	if err != nil {
		return err
	}
	if err := s.CreateVertex(ctx, UUID(id), moduleType); err != nil {
		return nil
	}
	q, err := module.VertexQuery()
	if err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, PropertyNameDisplayName, module.Name); err != nil {
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
	if err := s.SetVertexProperties(ctx, q, PropertyNameFormat, module.Format); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, "published_on", module.PubDate); err != nil {
		return err
	}
	if err := s.SetVertexProperties(ctx, q, PropertyNameUpdatedAt, module.LastUpdated); err != nil {
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
	if err := s.SetVertexProperties(ctx, q, PropertyNameImgURL, module.FeatureImage); err != nil {
		return err
	}
	if len(module.Teaser) > 0 {
		if err := s.SetVertexProperties(ctx, q, "teaser", module.Teaser); err != nil {
			return err
		}
	}

	return nil
}
