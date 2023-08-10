package graphtest

import (
	"context"
	"errors"
	"sync"

	"github.com/geomodulus/citygraph"
	"github.com/geomodulus/citygraph/pb"
)

type FakeGraphClient struct {
	sync.Mutex

	CreateVertexReqs            []*pb.Vertex
	CreateVertexFromTypeReqs    []*pb.Identifier
	CreateVertexFromTypeResps   []*pb.Uuid
	GetVerticesReqs             []*pb.VertexQuery
	GetVerticesResps            [][]*pb.Vertex
	GetVertexPropertiesReqs     []*pb.VertexPropertyQuery
	GetVertexPropertiesResps    [][]*pb.VertexProperty
	SetVertexPropertiesReqs     []*pb.SetVertexPropertiesRequest
	GetAllVertexPropertiesReqs  []*pb.VertexQuery
	GetAllVertexPropertiesResps [][]*pb.VertexProperties

	CreateEdgeReqs            []*pb.EdgeKey
	DeleteEdgesReqs           []*pb.EdgeQuery
	GetEdgesReqs              []*pb.EdgeQuery
	GetEdgesResps             [][]*pb.Edge
	GetEdgePropertiesReqs     []*pb.EdgePropertyQuery
	GetEdgePropertiesResps    [][]*pb.EdgeProperty
	SetEdgePropertiesReqs     []*pb.SetEdgePropertiesRequest
	GetAllEdgePropertiesReqs  []*pb.EdgeQuery
	GetAllEdgePropertiesResps [][]*pb.EdgeProperties
}

func (f *FakeGraphClient) CreateVertex(ctx context.Context, id *pb.Uuid, t *pb.Identifier) error {
	f.Lock()
	defer f.Unlock()

	f.CreateVertexReqs = append(f.CreateVertexReqs, &pb.Vertex{Id: id, T: t})
	return nil
}

func (f *FakeGraphClient) CreateVertexFromType(ctx context.Context, t *pb.Identifier) (*pb.Uuid, error) {
	f.Lock()
	defer f.Unlock()

	if len(f.CreateVertexFromTypeResps) == 0 {
		return nil, errors.New("CreateVertexFromType: fake has no response to return")
	}
	f.CreateVertexFromTypeReqs = append(f.CreateVertexFromTypeReqs, t)
	id, remaining := f.CreateVertexFromTypeResps[0], f.CreateVertexFromTypeResps[1:]
	f.CreateVertexFromTypeResps = remaining
	return id, nil
}

func (f *FakeGraphClient) GetVertices(ctx context.Context, query *pb.VertexQuery) ([]*pb.Vertex, error) {
	f.Lock()
	defer f.Unlock()

	if len(f.GetVerticesResps) == 0 {
		return nil, errors.New("GetVertices: fake has no response to return")
	}
	f.GetVerticesReqs = append(f.GetVerticesReqs, query)
	vtxs, remaining := f.GetVerticesResps[0], f.GetVerticesResps[1:]
	f.GetVerticesResps = remaining
	return vtxs, nil
}

func (c *FakeGraphClient) DeleteVertices(ctx context.Context, query *pb.VertexQuery) error {
	return nil
}

func (f *FakeGraphClient) GetVertexProperties(ctx context.Context, query *pb.VertexQuery, name string) ([]*pb.VertexProperty, error) {
	f.Lock()
	defer f.Unlock()

	if len(f.GetVertexPropertiesResps) == 0 {
		return nil, errors.New("GetVertexProperties: fake has no response to return")
	}
	f.GetVertexPropertiesReqs = append(f.GetVertexPropertiesReqs, &pb.VertexPropertyQuery{Inner: query, Name: &pb.Identifier{Value: name}})
	vps, remaining := f.GetVertexPropertiesResps[0], f.GetVertexPropertiesResps[1:]
	f.GetVertexPropertiesResps = remaining
	return vps, nil
}

func (f *FakeGraphClient) SetVertexProperties(ctx context.Context, query *pb.VertexQuery, name string, jsonValue interface{}) error {
	f.Lock()
	defer f.Unlock()

	req, err := citygraph.SetVertexPropertiesRequest(query, name, jsonValue)
	if err != nil {
		return err
	}
	f.SetVertexPropertiesReqs = append(f.SetVertexPropertiesReqs, req)
	return nil
}

func (f *FakeGraphClient) GetAllVertexProperties(ctx context.Context, query *pb.VertexQuery) ([]*pb.VertexProperties, error) {
	f.Lock()
	defer f.Unlock()

	if len(f.GetAllVertexPropertiesResps) == 0 {
		return nil, errors.New("GetAllVertexProperties: fake has no response to return")
	}
	f.GetAllVertexPropertiesReqs = append(f.GetAllVertexPropertiesReqs, query)
	vps, remaining := f.GetAllVertexPropertiesResps[0], f.GetAllVertexPropertiesResps[1:]
	f.GetAllVertexPropertiesResps = remaining
	return vps, nil
}

func (f *FakeGraphClient) CreateEdge(ctx context.Context, outbound *pb.Uuid, t *pb.Identifier, inbound *pb.Uuid) error {
	f.Lock()
	defer f.Unlock()

	f.CreateEdgeReqs = append(f.CreateEdgeReqs, &pb.EdgeKey{OutboundId: outbound, T: t, InboundId: inbound})
	return nil
}

func (f *FakeGraphClient) DeleteEdges(ctx context.Context, query *pb.EdgeQuery) error {
	f.Lock()
	defer f.Unlock()

	f.DeleteEdgesReqs = append(f.DeleteEdgesReqs, query)
	return nil
}

func (f *FakeGraphClient) GetEdges(ctx context.Context, query *pb.EdgeQuery) ([]*pb.Edge, error) {
	f.Lock()
	defer f.Unlock()

	if len(f.GetEdgesResps) == 0 {
		return nil, errors.New("GetEdges: fake has no response to return")
	}
	f.GetEdgesReqs = append(f.GetEdgesReqs, query)
	edges, remaining := f.GetEdgesResps[0], f.GetEdgesResps[1:]
	f.GetEdgesResps = remaining
	return edges, nil
}

func (f *FakeGraphClient) GetEdgeProperties(ctx context.Context, query *pb.EdgeQuery, name string) ([]*pb.EdgeProperty, error) {
	f.Lock()
	defer f.Unlock()

	if len(f.GetEdgePropertiesResps) == 0 {
		return nil, errors.New("GetEdgeProperties: fake has no response to return")
	}
	f.GetEdgePropertiesReqs = append(f.GetEdgePropertiesReqs, &pb.EdgePropertyQuery{Inner: query, Name: &pb.Identifier{Value: name}})
	eps, remaining := f.GetEdgePropertiesResps[0], f.GetEdgePropertiesResps[1:]
	f.GetEdgePropertiesResps = remaining
	return eps, nil
}

func (f *FakeGraphClient) SetEdgeProperties(ctx context.Context, query *pb.EdgeQuery, name string, jsonValue interface{}) error {
	f.Lock()
	defer f.Unlock()

	req, err := citygraph.SetEdgePropertiesRequest(query, name, jsonValue)
	if err != nil {
		return err
	}
	f.SetEdgePropertiesReqs = append(f.SetEdgePropertiesReqs, req)
	return nil
}

func (f *FakeGraphClient) GetAllEdgeProperties(ctx context.Context, query *pb.EdgeQuery) ([]*pb.EdgeProperties, error) {
	f.Lock()
	defer f.Unlock()

	f.GetAllEdgePropertiesReqs = append(f.GetAllEdgePropertiesReqs, query)
	eps, remaining := f.GetAllEdgePropertiesResps[0], f.GetAllEdgePropertiesResps[1:]
	f.GetAllEdgePropertiesResps = remaining
	return eps, nil
}

func (f *FakeGraphClient) NewBulkSender(ctx context.Context) (citygraph.BulkSender, error) {
	return nil, nil
}
