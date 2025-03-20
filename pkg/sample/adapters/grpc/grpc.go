package sample

import (
	"context"
	"log/slog"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/di"
	"github.com/eser/ajan/grpcfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/go-service/pkg/sample/adapters/grpc/generated"
	"github.com/eser/go-service/pkg/sample/adapters/storage"
	"github.com/eser/go-service/pkg/sample/business/channels"
)

type BroadcastService struct {
	generated.UnimplementedChannelServiceServer
	generated.UnimplementedMessageServiceServer

	logger       *logfx.Logger
	dataRegistry *datafx.Registry
}

func RegisterGrpcService(container di.Container, grpcService grpcfx.GrpcService, logger *logfx.Logger, dataRegistry *datafx.Registry) { //nolint:lll
	bs := NewBroadcastService(logger, dataRegistry)

	grpcService.RegisterService(&generated.ChannelService_ServiceDesc, bs)
	grpcService.RegisterService(&generated.MessageService_ServiceDesc, bs)
}

func NewBroadcastService(logger *logfx.Logger, dataRegistry *datafx.Registry) *BroadcastService {
	return &BroadcastService{logger: logger, dataRegistry: dataRegistry} //nolint:exhaustruct
}

func (s *BroadcastService) GetById(ctx context.Context, req *generated.GetByIdRequest) (*generated.Channel, error) {
	queries, err := storage.NewFromDefault(s.dataRegistry)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	service := channels.NewService(queries)

	channel, err := service.GetById(ctx, req.GetId())
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	result := &generated.Channel{
		Id:   channel.Id,
		Name: channel.Name.String,
	}

	return result, nil
}

func (s *BroadcastService) List(ctx context.Context, req *generated.ListRequest) (*generated.Channels, error) {
	queries, err := storage.NewFromDefault(s.dataRegistry)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	service := channels.NewService(queries)

	channels, err := service.List(ctx)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	newChannels := make([]*generated.Channel, len(channels))
	for i, channel := range channels {
		newChannels[i] = &generated.Channel{
			Id:   channel.Id,
			Name: channel.Name.String,
		}
	}

	return &generated.Channels{Channels: newChannels}, nil
}

func (s *BroadcastService) Send(ctx context.Context, req *generated.SendRequest) (*generated.SendResponse, error) {
	s.logger.Info(
		"Send",
		slog.String("channelId", req.GetChannelId()),
		slog.Any("message", req.GetMessage()),
	)

	return nil, nil //nolint:nilnil
}
