package citygraph

import (
	"context"
	"encoding/json"
	"io"
	"math"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/geomodulus/citygraph/pb"
)

type GraphClient interface {
	CreateVertex(context.Context, *pb.Uuid, *pb.Identifier) error
	CreateVertexFromType(context.Context, *pb.Identifier) (*pb.Uuid, error)
	DeleteVertices(context.Context, *pb.VertexQuery) error
	GetVertexProperties(context.Context, *pb.VertexQuery, string) ([]*pb.VertexProperty, error)
	SetVertexProperties(context.Context, *pb.VertexQuery, string, interface{}) error
	GetVertices(context.Context, *pb.VertexQuery) ([]*pb.Vertex, error)
	GetAllVertexProperties(context.Context, *pb.VertexQuery) ([]*pb.VertexProperties, error)
	DeleteEdges(context.Context, *pb.EdgeQuery) error
	CreateEdge(context.Context, *pb.Uuid, *pb.Identifier, *pb.Uuid) error
	GetEdges(context.Context, *pb.EdgeQuery) ([]*pb.Edge, error)
	GetEdgeProperties(context.Context, *pb.EdgeQuery, string) ([]*pb.EdgeProperty, error)
	SetEdgeProperties(context.Context, *pb.EdgeQuery, string, interface{}) error
	GetAllEdgeProperties(context.Context, *pb.EdgeQuery) ([]*pb.EdgeProperties, error)

	NewBulkSender(ctx context.Context) (BulkSender, error)
}

// NewSpecificVertexQuery returns a new query for the provided UUIDs.
func NewSpecificVertexQuery(ids ...*pb.Uuid) *pb.VertexQuery {
	return &pb.VertexQuery{
		Query: &pb.VertexQuery_Specific{
			Specific: &pb.SpecificVertexQuery{Ids: ids},
		},
	}
}

// NewSpecificEdgeQuery returns a new query for the provided edge keys.
func NewSpecificEdgeQuery(keys ...*pb.EdgeKey) *pb.EdgeQuery {
	return &pb.EdgeQuery{
		Query: &pb.EdgeQuery_Specific{
			Specific: &pb.SpecificEdgeQuery{Keys: keys},
		},
	}
}

// NewVertexPropertyQuery generates a vertex property query.
func NewVertexPropertyQuery(inner *pb.VertexQuery, name string) *pb.VertexPropertyQuery {
	return &pb.VertexPropertyQuery{
		Inner: inner,
		Name:  &pb.Identifier{Value: name},
	}
}

// NewEdgePropertyQuery generates an edge property query.
func NewEdgePropertyQuery(inner *pb.EdgeQuery, name string) *pb.EdgePropertyQuery {
	return &pb.EdgePropertyQuery{
		Inner: inner,
		Name:  &pb.Identifier{Value: name},
	}
}

// NewPipeEdgeQuery returns a new query for the provided edge keys.
func NewPipeEdgeQuery(inner *pb.VertexQuery, dir pb.EdgeDirection, t *pb.Identifier) *pb.EdgeQuery {
	return &pb.EdgeQuery{
		Query: &pb.EdgeQuery_Pipe{
			Pipe: &pb.PipeEdgeQuery{
				Inner:     inner,
				Direction: dir,
				T:         t,
				Limit:     math.MaxInt32,
			},
		},
	}
}

// NewPipeEdgeQueryLimit returns a new query for the provided edge keys.
func NewPipeEdgeQueryLimit(inner *pb.VertexQuery, dir pb.EdgeDirection, t *pb.Identifier, limit int) *pb.EdgeQuery {
	return &pb.EdgeQuery{
		Query: &pb.EdgeQuery_Pipe{
			Pipe: &pb.PipeEdgeQuery{
				Inner:     inner,
				Direction: dir,
				T:         t,
				Limit:     uint32(limit),
			},
		},
	}
}

// NewPipeEdgeQueryLimitOffset returns a new query for the provided edge keys.
func NewPipeEdgeQueryLimitOffset(inner *pb.VertexQuery, dir pb.EdgeDirection, t *pb.Identifier, limit int, low *timestamppb.Timestamp) *pb.EdgeQuery {
	return &pb.EdgeQuery{
		Query: &pb.EdgeQuery_Pipe{
			Pipe: &pb.PipeEdgeQuery{
				Inner:     inner,
				Direction: dir,
				T:         t,
				Low:       low,
				Limit:     uint32(limit),
			},
		},
	}
}

// NewPipeVertexQuery returns a new query for the provided edge keys.
func NewPipeVertexQuery(inner *pb.EdgeQuery, dir pb.EdgeDirection, t *pb.Identifier) *pb.VertexQuery {
	return &pb.VertexQuery{
		Query: &pb.VertexQuery_Pipe{
			Pipe: &pb.PipeVertexQuery{
				Inner:     inner,
				Direction: dir,
				T:         t,
			},
		},
	}
}

// SetVertexPropertiesRequest encodes the supplied value as JSON and returns a
// request to set that value for the supplied VertexQuery and property name.
func SetVertexPropertiesRequest(inner *pb.VertexQuery, name string, jsonValue interface{}) (*pb.SetVertexPropertiesRequest, error) {
	jsonStr, err := json.Marshal(jsonValue)
	if err != nil {
		return nil, err
	}
	return &pb.SetVertexPropertiesRequest{
		Q:     &pb.VertexPropertyQuery{Inner: inner, Name: &pb.Identifier{Value: name}},
		Value: &pb.Json{Value: string(jsonStr)},
	}, nil
}

// SetEdgePropertiesRequest encodes the supplied value as JSON and returns a
// request to set that value for the supplied EdgeQuery and property name.
func SetEdgePropertiesRequest(inner *pb.EdgeQuery, name string, jsonValue interface{}) (*pb.SetEdgePropertiesRequest, error) {
	jsonStr, err := json.Marshal(jsonValue)
	if err != nil {
		return nil, err
	}
	return &pb.SetEdgePropertiesRequest{
		Q:     &pb.EdgePropertyQuery{Inner: inner, Name: &pb.Identifier{Value: name}},
		Value: &pb.Json{Value: string(jsonStr)},
	}, nil
}

type Client struct {
	graph   pb.IndraDBClient
	counter uint32
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{pb.NewIndraDBClient(conn), 0}
}

type BulkSender interface {
	Send(*pb.BulkInsertItem) error
	CloseAndRecv() (*emptypb.Empty, error)
}

func (c *Client) NewBulkSender(ctx context.Context) (BulkSender, error) {
	return c.graph.BulkInsert(ctx)
}

func (c *Client) GetVertices(ctx context.Context, query *pb.VertexQuery) ([]*pb.Vertex, error) {
	stream, err := c.graph.GetVertices(ctx, query)
	if err != nil {
		return nil, err
	}
	var res []*pb.Vertex
	for {
		vtx, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, vtx)
	}
	return res, nil
}

func (c *Client) CreateVertex(ctx context.Context, id *pb.Uuid, t *pb.Identifier) error {
	_, err := c.graph.CreateVertex(ctx, &pb.Vertex{Id: id, T: t})
	return err
}

func (c *Client) CreateVertexFromType(ctx context.Context, t *pb.Identifier) (*pb.Uuid, error) {
	return c.graph.CreateVertexFromType(ctx, t)
}

func (c *Client) DeleteVertices(ctx context.Context, query *pb.VertexQuery) error {
	_, err := c.graph.DeleteVertices(ctx, query)
	return err
}

func (c *Client) GetVertexProperties(ctx context.Context, query *pb.VertexQuery, name string) ([]*pb.VertexProperty, error) {
	stream, err := c.graph.GetVertexProperties(ctx, NewVertexPropertyQuery(query, name))
	if err != nil {
		return nil, err
	}
	var res []*pb.VertexProperty
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (c *Client) SetVertexProperties(ctx context.Context, query *pb.VertexQuery, name string, jsonValue interface{}) error {
	req, err := SetVertexPropertiesRequest(query, name, jsonValue)
	if err != nil {
		return err
	}
	_, err = c.graph.SetVertexProperties(ctx, req)
	return err
}

func (c *Client) GetAllVertexProperties(ctx context.Context, query *pb.VertexQuery) ([]*pb.VertexProperties, error) {
	stream, err := c.graph.GetAllVertexProperties(ctx, query)
	if err != nil {
		return nil, err
	}
	var res []*pb.VertexProperties
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (c *Client) GetEdges(ctx context.Context, query *pb.EdgeQuery) ([]*pb.Edge, error) {
	stream, err := c.graph.GetEdges(ctx, query)
	if err != nil {
		return nil, err
	}
	var res []*pb.Edge
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (c *Client) CreateEdge(ctx context.Context, outbound *pb.Uuid, t *pb.Identifier, inbound *pb.Uuid) error {
	_, err := c.graph.CreateEdge(ctx, &pb.EdgeKey{OutboundId: outbound, T: t, InboundId: inbound})
	return err
}

func (c *Client) DeleteEdges(ctx context.Context, query *pb.EdgeQuery) error {
	_, err := c.graph.DeleteEdges(ctx, query)
	return err
}

func (c *Client) GetEdgeProperties(ctx context.Context, query *pb.EdgeQuery, name string) ([]*pb.EdgeProperty, error) {
	stream, err := c.graph.GetEdgeProperties(ctx, NewEdgePropertyQuery(query, name))
	if err != nil {
		return nil, err
	}
	var res []*pb.EdgeProperty
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (c *Client) SetEdgeProperties(ctx context.Context, query *pb.EdgeQuery, name string, jsonValue interface{}) error {
	req, err := SetEdgePropertiesRequest(query, name, jsonValue)
	if err != nil {
		return err
	}
	_, err = c.graph.SetEdgeProperties(ctx, req)
	return err
}

func (c *Client) GetAllEdgeProperties(ctx context.Context, query *pb.EdgeQuery) ([]*pb.EdgeProperties, error) {
	stream, err := c.graph.GetAllEdgeProperties(ctx, query)
	if err != nil {
		return nil, err
	}
	var res []*pb.EdgeProperties
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
