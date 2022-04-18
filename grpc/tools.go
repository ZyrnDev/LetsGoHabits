package grpc

import (
	"context"

	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ToolsServer struct {
	proto.UnimplementedToolsServer
}

func (s *ToolsServer) Ping(ctx context.Context, in *empty.Empty) (*timestamp.Timestamp, error) {
	return timestamppb.Now(), nil
	// return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
