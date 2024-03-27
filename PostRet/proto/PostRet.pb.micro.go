// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/PostRet.proto

package PostRet

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for PostRet service

func NewPostRetEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for PostRet service

type PostRetService interface {
	PostRet(ctx context.Context, in *PostRetRequest, opts ...client.CallOption) (*PostRetResponse, error)
}

type postRetService struct {
	c    client.Client
	name string
}

func NewPostRetService(name string, c client.Client) PostRetService {
	return &postRetService{
		c:    c,
		name: name,
	}
}

func (c *postRetService) PostRet(ctx context.Context, in *PostRetRequest, opts ...client.CallOption) (*PostRetResponse, error) {
	req := c.c.NewRequest(c.name, "PostRet.PostRet", in)
	out := new(PostRetResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PostRet service

type PostRetHandler interface {
	PostRet(context.Context, *PostRetRequest, *PostRetResponse) error
}

func RegisterPostRetHandler(s server.Server, hdlr PostRetHandler, opts ...server.HandlerOption) error {
	type postRet interface {
		PostRet(ctx context.Context, in *PostRetRequest, out *PostRetResponse) error
	}
	type PostRet struct {
		postRet
	}
	h := &postRetHandler{hdlr}
	return s.Handle(s.NewHandler(&PostRet{h}, opts...))
}

type postRetHandler struct {
	PostRetHandler
}

func (h *postRetHandler) PostRet(ctx context.Context, in *PostRetRequest, out *PostRetResponse) error {
	return h.PostRetHandler.PostRet(ctx, in, out)
}
