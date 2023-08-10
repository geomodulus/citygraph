// Package citygraph contains code for interacting with the Torontoverse.
package citygraph

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"github.com/geomodulus/citygraph/pb"
)

var (
	// PropertyNameDisplayName identifies the name for this vertex that should
	// be displayed when thnis vertex is viewed. If this data is meaningful, it
	// should be duplicated in another property with a more descriptive name.
	PropertyNameDisplayName = "display_name"

	PropertyNameUpdatedAt = "updated_at"
)

var (
	PropertyNameBodyText = "body"
	// PropertyNameBodyTextFormat identifies the format for the module's display, default is
	// "content-flex" but also possible in "content-map".
	PropertyNameFormat = "format"
	// PropertyNameImgURL identifies a property containing a URL pointing at an
	// image that should be used for promotion, eg. OpenGraph, feed entries, etc.
	PropertyNameImgURL = "img_url"
	PropertyNameJSFunc = "javascript_function"
)

// LngLat represents a geographic coordinate in the map. Encodes directly into
// a Mapbox-compatible json object.
type LngLat struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// Json creates a JSON protobuf object value from the provided json byte array.
func Json(b []byte) *pb.Json {
	return &pb.Json{Value: string(b)}
}

// IntVal creates a JSON protobuf value from a provided int.
func IntVal(num int) *pb.Json {
	return &pb.Json{Value: fmt.Sprintf(`%d`, num)}
}

// FloatVal creates a JSON protobuf value from a provided int.
func FloatVal(num float64) *pb.Json {
	return &pb.Json{Value: fmt.Sprintf(`%f`, num)}
}

// LngLatVal creates a JSON protobuf value from a provided LngLat.
func LngLatVal(loc LngLat) *pb.Json {
	return &pb.Json{Value: fmt.Sprintf(`{"lng":%g,"lat":%g}`, loc.Lng, loc.Lat)}
}

// StringVal creates a JSON protobuf value from the provided string.
func StringVal(str string) *pb.Json {
	jsonStr, err := json.Marshal(str)
	if err != nil {
		panic(err)
	}
	return &pb.Json{Value: string(jsonStr)}
}

// NewID() returns a v1 UUID (or panics, which seems reasonable). The v1 UUID
// is used here for it's convenient ordering.
func NewID() uuid.UUID {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id
}

// UUID wraps a UUID (preferably v1 for ordering) in a graph client proto.
func UUID(id uuid.UUID) *pb.Uuid {
	b, err := id.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return &pb.Uuid{Value: b}
}

func MustParseUUIDBytes(id string) []byte {
	b, err := uuid.MustParse(id).MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

// punctRE is a regular expression used to remove punctution from slug titles.
var punctRE = regexp.MustCompile(`([^\w\s-_\.~])`)
