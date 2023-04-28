package citygraph

import (
	"context"
	//	"fmt"
	//	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/geomodulus/citygraph/pb"
)

const bufSize = 1024 * 1024

type fakeGraphServer struct {
	pb.UnimplementedIndraDBServer

	getVerticesReqs             []*pb.VertexQuery
	getVerticesResps            [][]*pb.Vertex
	createVertexReqs            []*pb.Vertex
	createVertexFromTypeReqs    []*pb.Identifier
	createVertexFromTypeResps   []*pb.Uuid
	deleteVerticesReqs          []*pb.VertexQuery
	getVertexPropertiesReqs     []*pb.VertexPropertyQuery
	getVertexPropertiesNames    []string
	getVertexPropertiesResps    [][]*pb.VertexProperty
	setVertexPropertiesReqs     []*pb.SetVertexPropertiesRequest
	getAllVertexPropertiesReqs  []*pb.VertexQuery
	getAllVertexPropertiesResps [][]*pb.VertexProperties

	getEdgesReqs              []*pb.EdgeQuery
	getEdgesResps             [][]*pb.Edge
	deleteEdgesReqs           []*pb.EdgeQuery
	getEdgePropertiesReqs     []*pb.EdgePropertyQuery
	getEdgePropertiesNames    []string
	getEdgePropertiesResps    [][]*pb.EdgeProperty
	getAllEdgePropertiesReqs  []*pb.EdgeQuery
	getAllEdgePropertiesResps [][]*pb.EdgeProperties
}

func (f *fakeGraphServer) GetVertices(q *pb.VertexQuery, stream pb.IndraDB_GetVerticesServer) error {
	f.getVerticesReqs = append(f.getVerticesReqs, q)
	vertices, remaining := f.getVerticesResps[0], f.getVerticesResps[1:]
	f.getVerticesResps = remaining

	for _, vertex := range vertices {
		if err := stream.Send(vertex); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeGraphServer) CreateVertex(ctx context.Context, vtx *pb.Vertex) (*pb.CreateResponse, error) {
	f.createVertexReqs = append(f.createVertexReqs, vtx)
	return &pb.CreateResponse{Created: true}, nil
}

func (f *fakeGraphServer) CreateVertexFromType(ctx context.Context, id *pb.Identifier) (*pb.Uuid, error) {
	f.createVertexFromTypeReqs = append(f.createVertexFromTypeReqs, id)
	r, remaining := f.createVertexFromTypeResps[0], f.createVertexFromTypeResps[1:]
	f.createVertexFromTypeResps = remaining

	return r, nil
}

func (f *fakeGraphServer) DeleteVertices(ctx context.Context, q *pb.VertexQuery) (*emptypb.Empty, error) {
	f.deleteVerticesReqs = append(f.deleteVerticesReqs, q)
	return &emptypb.Empty{}, nil
}

func (f *fakeGraphServer) GetVertexProperties(q *pb.VertexPropertyQuery, stream pb.IndraDB_GetVertexPropertiesServer) error {
	f.getVertexPropertiesReqs = append(f.getVertexPropertiesReqs, q)
	vertexProperties, remaining := f.getVertexPropertiesResps[0], f.getVertexPropertiesResps[1:]
	f.getVertexPropertiesResps = remaining

	for _, props := range vertexProperties {
		if err := stream.Send(props); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeGraphServer) SetVertexProperties(ctx context.Context, req *pb.SetVertexPropertiesRequest) (*emptypb.Empty, error) {
	f.setVertexPropertiesReqs = append(f.setVertexPropertiesReqs, req)

	return &emptypb.Empty{}, nil
}

func (f *fakeGraphServer) GetAllVertexProperties(q *pb.VertexQuery, stream pb.IndraDB_GetAllVertexPropertiesServer) error {
	f.getAllVertexPropertiesReqs = append(f.getAllVertexPropertiesReqs, q)
	vertexProperties, remaining := f.getAllVertexPropertiesResps[0], f.getAllVertexPropertiesResps[1:]
	f.getAllVertexPropertiesResps = remaining

	for _, props := range vertexProperties {
		if err := stream.Send(props); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeGraphServer) GetEdges(q *pb.EdgeQuery, stream pb.IndraDB_GetEdgesServer) error {
	f.getEdgesReqs = append(f.getEdgesReqs, q)
	edges, remaining := f.getEdgesResps[0], f.getEdgesResps[1:]
	f.getEdgesResps = remaining

	for _, edge := range edges {
		if err := stream.Send(edge); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeGraphServer) DeleteEdges(ctx context.Context, q *pb.EdgeQuery) (*emptypb.Empty, error) {
	f.deleteEdgesReqs = append(f.deleteEdgesReqs, q)
	return &emptypb.Empty{}, nil
}

func (f *fakeGraphServer) GetEdgeProperties(q *pb.EdgePropertyQuery, stream pb.IndraDB_GetEdgePropertiesServer) error {
	f.getEdgePropertiesReqs = append(f.getEdgePropertiesReqs, q)
	edgeProperties, remaining := f.getEdgePropertiesResps[0], f.getEdgePropertiesResps[1:]
	f.getEdgePropertiesResps = remaining

	for _, props := range edgeProperties {
		if err := stream.Send(props); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeGraphServer) GetAllEdgeProperties(q *pb.EdgeQuery, stream pb.IndraDB_GetAllEdgePropertiesServer) error {
	f.getAllEdgePropertiesReqs = append(f.getAllEdgePropertiesReqs, q)
	edgeProperties, remaining := f.getAllEdgePropertiesResps[0], f.getAllEdgePropertiesResps[1:]
	f.getAllEdgePropertiesResps = remaining

	for _, props := range edgeProperties {
		if err := stream.Send(props); err != nil {
			return err
		}
	}
	return nil
}

func newServer(lis *bufconn.Listener, fakeServer *fakeGraphServer) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterIndraDBServer(s, fakeServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	return s
}

func bufDialer(lis *bufconn.Listener) func(ctx context.Context, addr string) (net.Conn, error) {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestGetVertices(t *testing.T) {
	wantVertices := []*pb.Vertex{
		{Id: &pb.Uuid{Value: []byte("vertex-a")}, T: &CityPrimary},
		{Id: &pb.Uuid{Value: []byte("vertex-b")}, T: &CityPrimary},
	}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
		getVerticesResps:           [][]*pb.Vertex{wantVertices},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificVertexQuery([]*pb.Uuid{{Value: []byte("vertex-a")}, {Value: []byte("vertex-b")}}...)
	gotVertices, err := client.GetVertices(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]*pb.VertexQuery{query}, fakeServer.getVerticesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("GetVertices() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantVertices, gotVertices, protocmp.Transform()); diff != "" {
		t.Errorf("GetVertices() returned diff:\n%s", diff)
	}
}

func TestCreateVertex(t *testing.T) {
	wantUUID := &pb.Uuid{Value: []byte("some-uuid-value")}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	if err = client.CreateVertex(ctx, wantUUID, &CityPrimary); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]*pb.Vertex{{Id: wantUUID, T: &CityPrimary}}, fakeServer.createVertexReqs, protocmp.Transform()); diff != "" {
		t.Errorf("CreateVertex() sent req diff:\n%s\n", diff)
	}
}

func TestCreateVertexFromType(t *testing.T) {
	wantUUID := &pb.Uuid{Value: []byte("some-uuid-value")}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
		createVertexFromTypeResps:  []*pb.Uuid{wantUUID},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	gotUUID, err := client.CreateVertexFromType(ctx, &CityPrimary)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]*pb.Identifier{&CityPrimary}, fakeServer.createVertexFromTypeReqs, protocmp.Transform()); diff != "" {
		t.Errorf("CreateVertexFromType() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantUUID, gotUUID, protocmp.Transform()); diff != "" {
		t.Errorf("CreateVertexFromType() returned diff:\n%s", diff)
	}
}

func TestDeleteVertices(t *testing.T) {
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificVertexQuery([]*pb.Uuid{{Value: []byte("vertex-a")}, {Value: []byte("vertex-b")}}...)
	if err := client.DeleteVertices(ctx, query); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]*pb.VertexQuery{query}, fakeServer.deleteVerticesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("DeleteVertices() sent req diff:\n%s\n", diff)
	}
}

func TestGetVertexProperties(t *testing.T) {
	wantVertexProperties := []*pb.VertexProperty{
		{Id: &pb.Uuid{Value: []byte("vertex-a")}, Value: &pb.Json{Value: "a value"}},
		{Id: &pb.Uuid{Value: []byte("vertex-b")}, Value: &pb.Json{Value: "b value"}},
	}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
		getVertexPropertiesResps:   [][]*pb.VertexProperty{wantVertexProperties},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificVertexQuery([]*pb.Uuid{{Value: []byte("vertex-a")}, {Value: []byte("vertex-b")}}...)
	gotVertexProperties, err := client.GetVertexProperties(ctx, query, "some-property")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(
		[]*pb.VertexPropertyQuery{{Inner: query, Name: &pb.Identifier{Value: "some-property"}}},
		fakeServer.getVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("GetVertexProperties() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantVertexProperties, gotVertexProperties, protocmp.Transform()); diff != "" {
		t.Errorf("GetVertexProperties() returned diff:\n%s", diff)
	}
}

func TestSetVertexProperties(t *testing.T) {
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificVertexQuery([]*pb.Uuid{
		{Value: []byte("vertex-a")},
		{Value: []byte("vertex-b")},
	}...)
	if err := client.SetVertexProperties(ctx, query, "some-property", "some-value"); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(
		[]*pb.SetVertexPropertiesRequest{{
			Q: &pb.VertexPropertyQuery{
				Inner: NewSpecificVertexQuery([]*pb.Uuid{
					{Value: []byte("vertex-a")},
					{Value: []byte("vertex-b")},
				}...),
				Name: &pb.Identifier{Value: "some-property"},
			},
			Value: &pb.Json{Value: `"some-value"`},
		}},
		fakeServer.setVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("SetVertexProperties() sent req diff:\n%s\n", diff)
	}
}

func TestGetAllVertexProperties(t *testing.T) {
	wantVertexProperties := []*pb.VertexProperties{{
		Vertex: &pb.Vertex{
			Id: &pb.Uuid{Value: []byte("vertex-1")},
			T:  &pb.Identifier{Value: "some-type"},
		},
		Props: []*pb.NamedProperty{
			{Name: &pb.Identifier{Value: "property-a"}, Value: &pb.Json{Value: `"a value"`}},
			{Name: &pb.Identifier{Value: "property-b"}, Value: &pb.Json{Value: `"b value"`}},
		},
	}, {
		Vertex: &pb.Vertex{
			Id: &pb.Uuid{Value: []byte("vertex-2")},
			T:  &pb.Identifier{Value: "some-other-type"},
		},
		Props: []*pb.NamedProperty{
			{Name: &pb.Identifier{Value: "property-c"}, Value: &pb.Json{Value: `"c value"`}},
		},
	}}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer:  pb.UnimplementedIndraDBServer{},
		getAllVertexPropertiesResps: [][]*pb.VertexProperties{wantVertexProperties},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificVertexQuery([]*pb.Uuid{{Value: []byte("vertex-1")}, {Value: []byte("vertex-2")}}...)
	gotVertexProperties, err := client.GetAllVertexProperties(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(
		[]*pb.VertexQuery{query},
		fakeServer.getAllVertexPropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("GetAllVertexProperties() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantVertexProperties, gotVertexProperties, protocmp.Transform()); diff != "" {
		t.Errorf("GetAllVertexProperties() returned diff:\n%s", diff)
	}
}

func TestGetEdges(t *testing.T) {
	now := time.Now()
	keyA := &pb.EdgeKey{
		OutboundId: &pb.Uuid{Value: []byte("vertex-a")},
		T:          &pb.Identifier{Value: "some-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-b")},
	}
	keyB := &pb.EdgeKey{
		OutboundId: &pb.Uuid{Value: []byte("vertex-b")},
		T:          &pb.Identifier{Value: "some-other-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-a")},
	}
	wantEdges := []*pb.Edge{
		{Key: keyA, CreatedDatetime: timestamppb.New(now)},
		{Key: keyB, CreatedDatetime: timestamppb.New(now.Add(-5 * time.Minute))},
	}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
		getEdgesResps:              [][]*pb.Edge{wantEdges},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificEdgeQuery(keyA, keyB)
	gotEdges, err := client.GetEdges(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]*pb.EdgeQuery{query}, fakeServer.getEdgesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("GetEdges() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantEdges, gotEdges, protocmp.Transform()); diff != "" {
		t.Errorf("GetEdges() returned diff:\n%s", diff)
	}
}

func TestDeleteEdges(t *testing.T) {
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificEdgeQuery([]*pb.EdgeKey{{
		OutboundId: &pb.Uuid{Value: []byte("vertex-a")},
		T:          &pb.Identifier{Value: "some-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-b")},
	}, {
		OutboundId: &pb.Uuid{Value: []byte("vertex-b")},
		T:          &pb.Identifier{Value: "some-other-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-a")},
	}}...)
	if err := client.DeleteEdges(ctx, query); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]*pb.EdgeQuery{query}, fakeServer.deleteEdgesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("DeleteEdges() sent req diff:\n%s\n", diff)
	}
}

func TestGetEdgeProperties(t *testing.T) {
	keyA, keyB := &pb.EdgeKey{
		OutboundId: &pb.Uuid{Value: []byte("vertex-a")},
		T:          &pb.Identifier{Value: "some-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-c")},
	}, &pb.EdgeKey{
		OutboundId: &pb.Uuid{Value: []byte("vertex-b")},
		T:          &pb.Identifier{Value: "some-other-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-d")},
	}
	wantEdgeProperties := []*pb.EdgeProperty{{
		Key:   keyA,
		Value: &pb.Json{Value: "a value"},
	}, {
		Key:   keyB,
		Value: &pb.Json{Value: "b value"},
	}}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
		getEdgePropertiesResps:     [][]*pb.EdgeProperty{wantEdgeProperties},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificEdgeQuery(keyA, keyB)
	gotEdgeProperties, err := client.GetEdgeProperties(ctx, query, "some-property")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(
		[]*pb.EdgePropertyQuery{{Inner: query, Name: &pb.Identifier{Value: "some-property"}}},
		fakeServer.getEdgePropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("GetEdgeProperties() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantEdgeProperties, gotEdgeProperties, protocmp.Transform()); diff != "" {
		t.Errorf("GetEdgeProperties() returned diff:\n%s", diff)
	}
}

func TestGetAllEdgeProperties(t *testing.T) {
	now := time.Now()
	keyA := &pb.EdgeKey{
		OutboundId: &pb.Uuid{Value: []byte("vertex-a")},
		T:          &pb.Identifier{Value: "some-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-b")},
	}
	keyB := &pb.EdgeKey{
		OutboundId: &pb.Uuid{Value: []byte("vertex-b")},
		T:          &pb.Identifier{Value: "some-other-edge-type"},
		InboundId:  &pb.Uuid{Value: []byte("vertex-a")},
	}
	wantEdgeProperties := []*pb.EdgeProperties{{
		Edge: &pb.Edge{Key: keyA, CreatedDatetime: timestamppb.New(now)},
		Props: []*pb.NamedProperty{
			{Name: &pb.Identifier{Value: "property-a"}, Value: &pb.Json{Value: `"a value"`}},
			{Name: &pb.Identifier{Value: "property-b"}, Value: &pb.Json{Value: `"b value"`}},
		},
	}, {
		Edge: &pb.Edge{Key: keyB, CreatedDatetime: timestamppb.New(now.Add(-5 * time.Minute))},
		Props: []*pb.NamedProperty{
			{Name: &pb.Identifier{Value: "property-c"}, Value: &pb.Json{Value: `"c value"`}},
		},
	}}
	fakeServer := &fakeGraphServer{
		UnimplementedIndraDBServer: pb.UnimplementedIndraDBServer{},
		getAllEdgePropertiesResps:  [][]*pb.EdgeProperties{wantEdgeProperties},
	}

	lis := bufconn.Listen(bufSize)
	s := newServer(lis, fakeServer)
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewClient(conn)

	query := NewSpecificEdgeQuery(keyA, keyB)
	gotEdgeProperties, err := client.GetAllEdgeProperties(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(
		[]*pb.EdgeQuery{query},
		fakeServer.getAllEdgePropertiesReqs, protocmp.Transform()); diff != "" {
		t.Errorf("GetAllEdgeProperties() sent req diff:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantEdgeProperties, gotEdgeProperties, protocmp.Transform()); diff != "" {
		t.Errorf("GetAllEdgeProperties() returned diff:\n%s", diff)
	}
}
