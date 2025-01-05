package samplesvc

import (
	"context"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/grpcfx"
	pb "github.com/eser/go-service/pkg/proto-go/broadcast"
	"github.com/eser/go-service/pkg/samplesvc/adapters/storage"
	"github.com/eser/go-service/pkg/samplesvc/business/channel"
)

type BroadcastService struct {
	pb.UnimplementedChannelServiceServer
	pb.UnimplementedMessageServiceServer

	logger       *slog.Logger
	dataRegistry *datafx.Registry
}

func RegisterGrpcService(container di.Container, grpcService grpcfx.GrpcService, logger *slog.Logger, dataRegistry *datafx.Registry) { //nolint:lll
	bs := NewBroadcastService(logger, dataRegistry)

	grpcService.RegisterService(&pb.ChannelService_ServiceDesc, bs)
	grpcService.RegisterService(&pb.MessageService_ServiceDesc, bs)
}

func NewBroadcastService(logger *slog.Logger, dataRegistry *datafx.Registry) *BroadcastService {
	return &BroadcastService{logger: logger, dataRegistry: dataRegistry} //nolint:exhaustruct
}

func (s *BroadcastService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.Channel, error) {
	dataSource := s.dataRegistry.GetDefaultSql()
	repo := storage.NewChannelRepository(dataSource)
	service := channel.NewService(repo)

	channel, err := service.GetById(ctx, req.GetId())
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	result := &pb.Channel{
		Id:   channel.Id,
		Name: channel.Name.String,
	}

	return result, nil
}

func (s *BroadcastService) List(ctx context.Context, req *pb.ListRequest) (*pb.Channels, error) {
	dataSource := s.dataRegistry.GetDefaultSql()
	repo := storage.NewChannelRepository(dataSource)
	service := channel.NewService(repo)

	channels, err := service.List(ctx)
	if err != nil {
		return nil, err //nolint:wrapcheck
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
