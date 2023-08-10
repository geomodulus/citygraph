// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FeedProducerClient is the client API for FeedProducer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedProducerClient interface {
	// AddAnnouncement is the method used to provided announcement to Feed Producer.
	AddAnnouncement(ctx context.Context, in *AddAnnouncementRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// AddContent is the method used to provide content items to Feed Producer.
	AddContent(ctx context.Context, in *AddContentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// RemoveContent is the method used to tell Feed Producer to forget an item.
	RemoveContent(ctx context.Context, in *RemoveContentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ReadLatest returns the latest items from the feed.
	ReadLatest(ctx context.Context, in *ReadLatestRequest, opts ...grpc.CallOption) (*ReadLatestResponse, error)
	// ListActiveReleases returns a list of active releases only.
	ListActiveReleases(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (FeedProducer_ListActiveReleasesClient, error)
	// ListAllReleases returns a list of every release, active and waiting.
	ListAllReleases(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (FeedProducer_ListAllReleasesClient, error)
	// QueueItem schedules an item for immediate release.
	QueueItem(ctx context.Context, in *QueueItemRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type feedProducerClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedProducerClient(cc grpc.ClientConnInterface) FeedProducerClient {
	return &feedProducerClient{cc}
}

func (c *feedProducerClient) AddAnnouncement(ctx context.Context, in *AddAnnouncementRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/feed_producer.FeedProducer/AddAnnouncement", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedProducerClient) AddContent(ctx context.Context, in *AddContentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/feed_producer.FeedProducer/AddContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedProducerClient) RemoveContent(ctx context.Context, in *RemoveContentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/feed_producer.FeedProducer/RemoveContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedProducerClient) ReadLatest(ctx context.Context, in *ReadLatestRequest, opts ...grpc.CallOption) (*ReadLatestResponse, error) {
	out := new(ReadLatestResponse)
	err := c.cc.Invoke(ctx, "/feed_producer.FeedProducer/ReadLatest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedProducerClient) ListActiveReleases(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (FeedProducer_ListActiveReleasesClient, error) {
	stream, err := c.cc.NewStream(ctx, &FeedProducer_ServiceDesc.Streams[0], "/feed_producer.FeedProducer/ListActiveReleases", opts...)
	if err != nil {
		return nil, err
	}
	x := &feedProducerListActiveReleasesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FeedProducer_ListActiveReleasesClient interface {
	Recv() (*ReleaseItem, error)
	grpc.ClientStream
}

type feedProducerListActiveReleasesClient struct {
	grpc.ClientStream
}

func (x *feedProducerListActiveReleasesClient) Recv() (*ReleaseItem, error) {
	m := new(ReleaseItem)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *feedProducerClient) ListAllReleases(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (FeedProducer_ListAllReleasesClient, error) {
	stream, err := c.cc.NewStream(ctx, &FeedProducer_ServiceDesc.Streams[1], "/feed_producer.FeedProducer/ListAllReleases", opts...)
	if err != nil {
		return nil, err
	}
	x := &feedProducerListAllReleasesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FeedProducer_ListAllReleasesClient interface {
	Recv() (*ReleaseItem, error)
	grpc.ClientStream
}

type feedProducerListAllReleasesClient struct {
	grpc.ClientStream
}

func (x *feedProducerListAllReleasesClient) Recv() (*ReleaseItem, error) {
	m := new(ReleaseItem)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *feedProducerClient) QueueItem(ctx context.Context, in *QueueItemRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/feed_producer.FeedProducer/QueueItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedProducerServer is the server API for FeedProducer service.
// All implementations must embed UnimplementedFeedProducerServer
// for forward compatibility
type FeedProducerServer interface {
	// AddAnnouncement is the method used to provided announcement to Feed Producer.
	AddAnnouncement(context.Context, *AddAnnouncementRequest) (*emptypb.Empty, error)
	// AddContent is the method used to provide content items to Feed Producer.
	AddContent(context.Context, *AddContentRequest) (*emptypb.Empty, error)
	// RemoveContent is the method used to tell Feed Producer to forget an item.
	RemoveContent(context.Context, *RemoveContentRequest) (*emptypb.Empty, error)
	// ReadLatest returns the latest items from the feed.
	ReadLatest(context.Context, *ReadLatestRequest) (*ReadLatestResponse, error)
	// ListActiveReleases returns a list of active releases only.
	ListActiveReleases(*emptypb.Empty, FeedProducer_ListActiveReleasesServer) error
	// ListAllReleases returns a list of every release, active and waiting.
	ListAllReleases(*emptypb.Empty, FeedProducer_ListAllReleasesServer) error
	// QueueItem schedules an item for immediate release.
	QueueItem(context.Context, *QueueItemRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedFeedProducerServer()
}

// UnimplementedFeedProducerServer must be embedded to have forward compatible implementations.
type UnimplementedFeedProducerServer struct {
}

func (UnimplementedFeedProducerServer) AddAnnouncement(context.Context, *AddAnnouncementRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAnnouncement not implemented")
}
func (UnimplementedFeedProducerServer) AddContent(context.Context, *AddContentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddContent not implemented")
}
func (UnimplementedFeedProducerServer) RemoveContent(context.Context, *RemoveContentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveContent not implemented")
}
func (UnimplementedFeedProducerServer) ReadLatest(context.Context, *ReadLatestRequest) (*ReadLatestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadLatest not implemented")
}
func (UnimplementedFeedProducerServer) ListActiveReleases(*emptypb.Empty, FeedProducer_ListActiveReleasesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListActiveReleases not implemented")
}
func (UnimplementedFeedProducerServer) ListAllReleases(*emptypb.Empty, FeedProducer_ListAllReleasesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListAllReleases not implemented")
}
func (UnimplementedFeedProducerServer) QueueItem(context.Context, *QueueItemRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueueItem not implemented")
}
func (UnimplementedFeedProducerServer) mustEmbedUnimplementedFeedProducerServer() {}

// UnsafeFeedProducerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedProducerServer will
// result in compilation errors.
type UnsafeFeedProducerServer interface {
	mustEmbedUnimplementedFeedProducerServer()
}

func RegisterFeedProducerServer(s grpc.ServiceRegistrar, srv FeedProducerServer) {
	s.RegisterService(&FeedProducer_ServiceDesc, srv)
}

func _FeedProducer_AddAnnouncement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAnnouncementRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedProducerServer).AddAnnouncement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/feed_producer.FeedProducer/AddAnnouncement",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedProducerServer).AddAnnouncement(ctx, req.(*AddAnnouncementRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedProducer_AddContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedProducerServer).AddContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/feed_producer.FeedProducer/AddContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedProducerServer).AddContent(ctx, req.(*AddContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedProducer_RemoveContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedProducerServer).RemoveContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/feed_producer.FeedProducer/RemoveContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedProducerServer).RemoveContent(ctx, req.(*RemoveContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedProducer_ReadLatest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadLatestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedProducerServer).ReadLatest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/feed_producer.FeedProducer/ReadLatest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedProducerServer).ReadLatest(ctx, req.(*ReadLatestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedProducer_ListActiveReleases_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FeedProducerServer).ListActiveReleases(m, &feedProducerListActiveReleasesServer{stream})
}

type FeedProducer_ListActiveReleasesServer interface {
	Send(*ReleaseItem) error
	grpc.ServerStream
}

type feedProducerListActiveReleasesServer struct {
	grpc.ServerStream
}

func (x *feedProducerListActiveReleasesServer) Send(m *ReleaseItem) error {
	return x.ServerStream.SendMsg(m)
}

func _FeedProducer_ListAllReleases_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FeedProducerServer).ListAllReleases(m, &feedProducerListAllReleasesServer{stream})
}

type FeedProducer_ListAllReleasesServer interface {
	Send(*ReleaseItem) error
	grpc.ServerStream
}

type feedProducerListAllReleasesServer struct {
	grpc.ServerStream
}

func (x *feedProducerListAllReleasesServer) Send(m *ReleaseItem) error {
	return x.ServerStream.SendMsg(m)
}

func _FeedProducer_QueueItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedProducerServer).QueueItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/feed_producer.FeedProducer/QueueItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedProducerServer).QueueItem(ctx, req.(*QueueItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedProducer_ServiceDesc is the grpc.ServiceDesc for FeedProducer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedProducer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "feed_producer.FeedProducer",
	HandlerType: (*FeedProducerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddAnnouncement",
			Handler:    _FeedProducer_AddAnnouncement_Handler,
		},
		{
			MethodName: "AddContent",
			Handler:    _FeedProducer_AddContent_Handler,
		},
		{
			MethodName: "RemoveContent",
			Handler:    _FeedProducer_RemoveContent_Handler,
		},
		{
			MethodName: "ReadLatest",
			Handler:    _FeedProducer_ReadLatest_Handler,
		},
		{
			MethodName: "QueueItem",
			Handler:    _FeedProducer_QueueItem_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListActiveReleases",
			Handler:       _FeedProducer_ListActiveReleases_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListAllReleases",
			Handler:       _FeedProducer_ListAllReleases_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}