package identitysvc

import (
	"context"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/grpcfx"
	pb "github.com/eser/go-service/pkg/proto-go/broadcast"
)

type BroadcastService struct {
	pb.UnimplementedChannelServiceServer
	pb.UnimplementedMessageServiceServer

	logger *slog.Logger
}

func RegisterGrpcService(container di.Container, grpcService grpcfx.GrpcService, logger *slog.Logger) {
	bs := NewBroadcastService(logger)

	grpcService.RegisterService(&pb.ChannelService_ServiceDesc, bs)
	grpcService.RegisterService(&pb.MessageService_ServiceDesc, bs)
}

func NewBroadcastService(logger *slog.Logger) *BroadcastService {
	return &BroadcastService{logger: logger} //nolint:exhaustruct
}

func (s *BroadcastService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.Channel, error) {
	// uow, err := dbConn.CreateScope()
	// if err != nil {
	// 	return nil, err
	// }
	// defer uow.Close()

	// scopeCtx := uow.Context()
	// channelRepo := NewChannelRepository(uow.Scope())

	// channelRecord, err := channelRepo.GetById(scopeCtx, req.GetId())
	// if err != nil {
	// 	return nil, err
	// }

	channel := &pb.Channel{
		Id:   "123",          // channelRecord.Id,
		Name: "Test Channel", // channelRecord.Name,
	}

	return channel, nil
}

func (s *BroadcastService) List(ctx context.Context, req *pb.ListRequest) (*pb.Channels, error) {
	// Implementation here
	return nil, nil //nolint:nilnil
}

func (s *BroadcastService) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	s.logger.Info(
		"Send",
		slog.String("channelId", req.GetChannelId()),
		slog.Any("message", req.GetMessage()),
	)

	return nil, nil //nolint:nilnil
}
