package main

import (
	"crypto/tls"
	_ "embed"
	"github.com/softkot/locker/client"
	"github.com/softkot/locker/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net/http"
	"os"
	"time"
)

// go get -u google.golang.org/protobuf/cmd/protoc-gen-go
// go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
// openssl req -x509 -nodes -days 36500 -newkey rsa:2048 -subj "/CN=locker.backend" -keyout server.key -out server.crt
//go:generate protoc -I proto --go_opt=paths=source_relative  --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative --go_out=client --go-grpc_out=client locker.proto

//go:embed server.crt
var serverCrt []byte

//go:embed server.key
var serverKey []byte

func main() {
	if cert, err := tls.X509KeyPair(serverCrt, serverKey); err != nil {
		log.Fatal(err)
	} else {
		tlsConfig := &tls.Config{
			Certificates:             []tls.Certificate{cert},
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
		}
		grpcServer := grpc.NewServer(
			grpc.WriteBufferSize(4096),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{MinTime: 10 * time.Second, PermitWithoutStream: true}),
			grpc.KeepaliveParams(keepalive.ServerParameters{Timeout: 10 * time.Second}),
			grpc.StreamInterceptor(server.CheckStreamApiKeyAuth),
			grpc.UnaryInterceptor(server.CheckUnaryApiKeyAuth),
		)
		client.RegisterLockerServer(grpcServer, server.NewLockerServer())
		handler := func(resp http.ResponseWriter, req *http.Request) {
			execFile, _ := os.Executable()
			resp.Header().Set("Cache-Control", "private, max-age=3600")
			resp.Header().Set("Micro-Service", execFile)
			if req.RequestURI == "/" {
				http.Error(resp, http.StatusText(http.StatusOK), http.StatusOK)
				return
			}
			grpcServer.ServeHTTP(resp, req)
			return
		}
		h2s := &http2.Server{}
		httpsSrv := &http.Server{
			Addr:        "0.0.0.0:8443",
			Handler:     h2c.NewHandler(http.HandlerFunc(handler), h2s),
			IdleTimeout: 120 * time.Second,
			TLSConfig:   tlsConfig,
		}
		log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
	}

}
