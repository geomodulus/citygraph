package citygraph

import (
	"github.com/geomodulus/citygraph/pb"
)

var (
	// CityPrimary identifies a city that is the biggest in its geographic area
	// like Toronto, London or Stockholm.
	CityPrimary = pb.Identifier{Value: "city-primary"}

	// CityIntersection identifies an intersection inside a city.
	CityIntersection = pb.Identifier{Value: "city-intersection"}

	// CityNeighbourhood identifies an neighbourhood inside a city.
	CityNeighbourhood = pb.Identifier{Value: "city-neighbourhood"}

	// CityMayorOffice identifies the mayor's office of a city but not any
	// particular mayor individually.
	CityMayorOffice = pb.Identifier{Value: "city-mayor-office"}

	// CityWard identifies a Ward, a typical city constituency.
	CityWard = pb.Identifier{Value: "city-ward"}

	// A waterfront recreation location that may or may not feature sand.
	Beach = pb.Identifier{Value: "beach"}
)

// Municipal anchor vertices
var (
	Toronto = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6bc0-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &CityPrimary,
	}
)
var (
	//HasOffice points from a municipality to offices it maintains, like mayor.
	HasOffice = pb.Identifier{Value: "has-office"}

	//HasConstituency points from a municipality to it's electoral constituencies.
	HasConstituency = pb.Identifier{Value: "has-constituency"}
)

// Toronto Beaches
var (
	MarieCurtisParkEastBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6ed8-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	SunnysideBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6ec1-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	HanlansPointBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6ea4-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	GibraltarPointBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6e96-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	CentreIslandBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6e78-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	WardsIslandBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6e50-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	CherryBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6e44-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	WoodbineBeaches = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6e2d-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	KewBalmyBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6e18-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
	BluffersParkBeach = &pb.Vertex{
		Id: &pb.Uuid{
			Value: MustParseUUIDBytes("ae5c6db6-ffc1-11eb-b867-244bfe5bf61a"),
		},
		T: &Beach,
	}
)
