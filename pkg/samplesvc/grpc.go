package samplesvc

import (
	"context"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/grpcfx"
	pb "github.com/eser/go-service/pkg/proto-go/broadcast"
)

type BroadcastService struct {
	pb.UnimplementedChannelServiceServer
	pb.UnimplementedMessageServiceServer

	logger       *slog.Logger
	dataProvider datafx.DataProvider
}

func RegisterGrpcService(container di.Container, grpcService grpcfx.GrpcService, logger *slog.Logger, dataProvider datafx.DataProvider) { //nolint:lll
	bs := NewBroadcastService(logger, dataProvider)

	grpcService.RegisterService(&pb.ChannelService_ServiceDesc, bs)
	grpcService.RegisterService(&pb.MessageService_ServiceDesc, bs)
}

func NewBroadcastService(logger *slog.Logger, dataProvider datafx.DataProvider) *BroadcastService {
	return &BroadcastService{logger: logger, dataProvider: dataProvider} //nolint:exhaustruct
}

func (s *BroadcastService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.Channel, error) {
	channel := &pb.Channel{
		Id:   "123",
		Name: "Test Channel",
	}

	return channel, nil
}

func (s *BroadcastService) List(ctx context.Context, req *pb.ListRequest) (*pb.Channels, error) {
	scope := s.dataProvider.GetDefault().Connection

	channels, err := NewChannelService(scope).List(ctx)
	if err != nil {
		return nil, err
	}

	newChannels := make([]*pb.Channel, len(channels))
	for i, channel := range channels {
		newChannels[i] = &pb.Channel{
			Id:   channel.Id,
			Name: channel.Name.String,
		}
	}

	return &pb.Channels{Channels: newChannels}, nil
}

func (s *BroadcastService) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	s.logger.Info(
		"Send",
		slog.String("channelId", req.GetChannelId()),
		slog.Any("message", req.GetMessage()),
	)

	return nil, nil //nolint:nilnil
}
