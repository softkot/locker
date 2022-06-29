package server

import (
	"context"
	"github.com/softkot/locker/api"
	"github.com/softkot/locker/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
)

var apiKeyHeader = "API_KEY"
var staticApiKey = ""

func init() {
	if value := os.Getenv(apiKeyHeader); value != "" {
		staticApiKey = value
	}
}
func checkAuth(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	apiKeys := md.Get(apiKeyHeader)
	var apiKey string
	if len(apiKeys) > 0 {
		apiKey = apiKeys[0]
	}
	if apiKey != staticApiKey {
		return status.Error(codes.Unauthenticated, "GRPC api key is invalid")
	}
	return nil
}
func CheckUnaryApiKeyAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := checkAuth(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func CheckStreamApiKeyAuth(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := checkAuth(ss.Context()); err != nil {
		return err
	}
	return handler(srv, ss)
}

type lockerServer struct {
	locker api.Locker
}

func (l *lockerServer) Lock(entity *client.Entity, server client.Locker_LockServer) error {
	ctx := server.Context()
	if err := l.locker.Lock(ctx, entity.Name); err != nil {
		return err
	}
	if err := server.Send(entity); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (l *lockerServer) TryLock(entity *client.Entity, server client.Locker_TryLockServer) error {
	ctx := server.Context()
	if ok, err := l.locker.TryLock(ctx, entity.Name); err != nil {
		return err
	} else {
		if ok {
			if err := server.Send(entity); err != nil {
				return err
			}
			<-ctx.Done()
		}
	}
	return nil
}

func NewLockerServer() client.LockerServer {
	return &lockerServer{locker: api.NewLocker()}
}
