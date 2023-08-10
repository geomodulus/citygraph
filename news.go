package citygraph

import "github.com/geomodulus/citygraph/pb"

// News Types
var (
	// Vertex types
	NewsPublisher    = pb.Identifier{Value: "news-publisher"}
	NewsArticle      = pb.Identifier{Value: "news-article"}
	NewsWireArticle  = pb.Identifier{Value: "news-wire-article"}
	NewsWireBulletin = pb.Identifier{Value: "news-wire-bulletin"}
	NewsAuthor       = pb.Identifier{Value: "news-author"}
	NewsCategory     = pb.Identifier{Value: "news-category"}
	NewsImage        = pb.Identifier{Value: "news-image"}
	NewsGeoJSON      = pb.Identifier{Value: "news-geojson"}

	// Edge types
	Published      = pb.Identifier{Value: "published"}       // author or publisher -> publish -> item.
	PublishedBy    = pb.Identifier{Value: "published-by"}    // item -> published by -> author or publisher
	PublishesAbout = pb.Identifier{Value: "publishes-about"} // item or publisher -> about -> a category.
	ItemAbout      = pb.Identifier{Value: "item-about"}      // item or publisher -> about -> a category.
	IllustratedBy  = pb.Identifier{Value: "illustrated-by"}
	CoversTopic    = pb.Identifier{Value: "covers-topic"}
	Featuring      = pb.Identifier{Value: "featuring"}
)

