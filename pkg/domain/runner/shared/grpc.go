package shared

import (
	"context"
	"github.com/MR5356/aurora/pkg/domain/runner/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type TaskGRPCPlugin struct {
	plugin.Plugin
	Impl ITask
}

func (p *TaskGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterITaskServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *TaskGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewITaskClient(c)}, nil
}

type GRPCClient struct {
	client proto.ITaskClient
}

type GRPCServer struct {
	Impl ITask
	proto.UnimplementedITaskServer
}

func (c *GRPCClient) GetInfo() *proto.TaskInfo {
	taskInfo, _ := c.client.GetInfo(context.Background(), &proto.Empty{})
	return taskInfo
}

func (s *GRPCServer) GetInfo(ctx context.Context, req *proto.Empty) (*proto.TaskInfo, error) {
	return s.Impl.GetInfo(), nil
}

func (c *GRPCClient) GetParams() *proto.TaskParams {
	params, _ := c.client.GetParams(context.Background(), &proto.Empty{})
	return params
}

func (s *GRPCServer) GetParams(ctx context.Context, req *proto.Empty) (*proto.TaskParams, error) {
	return s.Impl.GetParams(), nil
}

func (c *GRPCClient) SetParams(params *proto.TaskParams) {
	_, _ = c.client.SetParams(context.Background(), params)
}

func (s *GRPCServer) SetParams(ctx context.Context, params *proto.TaskParams) (*proto.Empty, error) {
	s.Impl.SetParams(params)
	return &proto.Empty{}, nil
}

func (c *GRPCClient) Start() error {
	_, err := c.client.Start(context.Background(), &proto.Empty{})
	return err
}

func (s *GRPCServer) Start(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.Start()
}

func (c *GRPCClient) Stop() error {
	_, err := c.client.Stop(context.Background(), &proto.Empty{})
	return err
}

func (s *GRPCServer) Stop(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.Stop()
}

func (c *GRPCClient) Pause() error {
	_, err := c.client.Pause(context.Background(), &proto.Empty{})
	return err
}

func (s *GRPCServer) Pause(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.Pause()
}

func (c *GRPCClient) Resume() error {
	_, err := c.client.Resume(context.Background(), &proto.Empty{})
	return err
}

func (s *GRPCServer) Resume(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.Resume()
}
