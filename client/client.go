package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io"
	"time"
)

type Locker struct {
	client LockerClient
}

type apiKeyAuth struct {
	apiKey string
}

func (t *apiKeyAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	res := make(map[string]string)
	res["API_KEY"] = t.apiKey
	return res, nil
}

func (t *apiKeyAuth) RequireTransportSecurity() bool {
	return true
}

func NewLocker(apiEndpoint string, apiKey string, insecure bool) (*Locker, error) {
	pool, _ := x509.SystemCertPool()
	transportCreds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: insecure,
		RootCAs:            pool,
	})
	if lockerConnection, err := grpc.Dial(
		apiEndpoint,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			Timeout:             time.Second * 15,
			PermitWithoutStream: true,
		}),
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithPerRPCCredentials(&apiKeyAuth{apiKey: apiKey}),
	); err != nil {
		return nil, err
	} else {
		return &Locker{
			client: NewLockerClient(lockerConnection),
		}, nil
	}
}

func (l *Locker) Lock(ctx context.Context, name string) error {
	if strm, err := l.client.Lock(ctx, &Entity{Name: name}); err != nil {
		return err
	} else {
		if _, err := strm.Recv(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Locker) TryLock(ctx context.Context, name string) (bool, error) {
	if strm, err := l.client.TryLock(ctx, &Entity{Name: name}); err != nil {
		return false, err
	} else {
		if _, err := strm.Recv(); err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
	}
	return true, nil
}
