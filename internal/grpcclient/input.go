package grpcclient

import (
	"github.com/oleg5896/ai-common/logger"
	proto "github.com/oleg5896/ai-proto/gateway"
	"google.golang.org/grpc"
)

func NewInputClient(addr string) (proto.GatewayClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error("Не удалось подключиться к gRPC", err)
		return nil, nil, err
	}
	return proto.NewGatewayClient(conn), conn, nil
}
