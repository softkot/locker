package server

import (
	"context"
	"fmt"
	"github.com/softkot/locker/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"sync"
)

type nameHolder struct {
	name  string
	usage int
	locks int
	cond  *sync.Cond
}

type lockersHolder struct {
	createHolders int
	deleteHolders int
	locker        *sync.Mutex
	waiters       map[string]*nameHolder
}

func (h *nameHolder) String() string {
	return fmt.Sprintf("%v", h.usage)
}

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

func (l *lockersHolder) String() string {
	l.locker.Lock()
	defer l.locker.Unlock()
	return fmt.Sprintf("Locker created: %v , deleted: %v waiters %v", l.createHolders, l.deleteHolders, l.waiters)
}

func NewLocker() *lockersHolder {
	return &lockersHolder{
		locker:  &sync.Mutex{},
		waiters: make(map[string]*nameHolder),
	}
}
func (l *lockersHolder) getHolder(name string) *nameHolder {
	var holder *nameHolder
	if h, ok := l.waiters[name]; ok {
		holder = h
	} else {
		l.createHolders++
		holder = &nameHolder{
			name:  name,
			usage: 0,
			cond:  sync.NewCond(l.locker),
		}
		l.waiters[name] = holder
	}
	return holder
}

// Lock / - захватывает имя и владеет им до завершения контекста
func (l *lockersHolder) Lock(ctx context.Context, name string) {
	l.locker.Lock()
	defer l.locker.Unlock()
	var holder = l.getHolder(name)
	defer l.unlock(ctx, holder)
	holder.usage++
	for holder.locks > 0 {
		holder.cond.Wait()
	}
	holder.locks++
}

// TryLock / - пробует захватить имя и в случае успеха владеет им до завершения контекста
func (l *lockersHolder) TryLock(ctx context.Context, name string) bool {
	l.locker.Lock()
	defer l.locker.Unlock()
	var holder = l.getHolder(name)
	if holder.locks > 0 {
		return false
	}
	holder.usage++
	holder.locks++
	l.unlock(ctx, holder)
	return true
}

func (l *lockersHolder) unlock(ctx context.Context, h *nameHolder) {
	go func() {
		<-ctx.Done()
		l.locker.Lock()
		defer l.locker.Unlock()
		h.locks--
		h.usage--
		if h.usage == 0 {
			l.deleteHolders++
			delete(l.waiters, h.name)
		}
		h.cond.Signal()
	}()
}

type lockerServer struct {
	locker *lockersHolder
}

func (l *lockerServer) Lock(entity *client.Entity, server client.Locker_LockServer) error {
	ctx := server.Context()
	l.locker.Lock(ctx, entity.Name)
	if err := server.Send(entity); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (l *lockerServer) TryLock(entity *client.Entity, server client.Locker_TryLockServer) error {
	ctx := server.Context()
	if l.locker.TryLock(ctx, entity.Name) {
		if err := server.Send(entity); err != nil {
			return err
		}
		<-ctx.Done()
	}
	return nil
}

func NewLockerServer() client.LockerServer {
	return &lockerServer{locker: NewLocker()}
}
