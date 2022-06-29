package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/softkot/locker/api"
	"github.com/softkot/locker/client"
	"gitlab.com/skllzz/multiop"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

func init() {

}

func WithNamedLock(ctx context.Context, l api.Locker, name string, block func(context.Context) error) error {
	if name == "" {
		return block(ctx)
	}
	myctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := l.Lock(myctx, name); err != nil {
		return err
	} else {
		return block(myctx)
	}
}

func WithNamedTryLock(ctx context.Context, l api.Locker, name string, block func(context.Context) error) error {
	if name == "" {
		return block(ctx)
	}
	myctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if ok, err := l.TryLock(myctx, name); err != nil {
		return err
	} else {
		if ok {
			return block(myctx)
		} else {
			return nil
		}
	}
}

func TestLockWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	if l, e := client.NewLocker("localhost:8443", os.Getenv("API_KEY"), true); e != nil {
		t.Fatal(e)
	} else {
		runner := multiop.NewParallelRunner(ctx, 16)
		for i := 0; i < 50000; i++ {
			if err := runner.Submit(func(ctx context.Context) error {
				e := WithNamedTryLock(ctx, l, uuid.NewString(), func(ctx context.Context) error {
					t.Logf("lock started")
					select {
					case <-ctx.Done():
					case <-time.After(time.Millisecond * 1):
					}
					t.Logf("lock completed")
					return ctx.Err()
				})
				return e
			}); err != nil {
				t.Fatal(err)
			}
		}
		if err := runner.Await(); err != nil {
			t.Fatal(err)
		}
	}
	cancel()
}

func TestRemoteLocker(t *testing.T) {
	if l, e := client.NewLocker("localhost:8443", "", true); e != nil {
		t.Fatal(e)
	} else {
		var seed = time.Now().Unix()
		//seed=1625211875
		t.Logf("seed %v", seed)
		rand.Seed(seed)
		wg := &sync.WaitGroup{}
		values := make([]int, 50)
		for i := 0; i < 5000; i++ {
			wg.Add(1)
			e := i % len(values)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				<-time.After((time.Duration(rand.Intn(50))) * time.Millisecond)
				id := fmt.Sprintf("n-%v", e)
				if err := l.Lock(ctx, id); err == nil {
					v := rand.Intn(100)
					<-time.After((time.Duration(rand.Intn(5))) * time.Millisecond)
					c := values[e]
					<-time.After((time.Duration(rand.Intn(5))) * time.Millisecond)
					values[e] = c + v
					<-time.After((time.Duration(rand.Intn(5))) * time.Millisecond)
					values[e] = c
				} else {
					t.Log(err)
				}
				cancel()
			}()
		}
		wg.Wait()
		t.Logf("%v", values)
		for _, v := range values {
			if v != 0 {
				t.Fail()
			}
		}
	}
}

func TestLocalLocker(t *testing.T) {
	var seed = time.Now().Unix()
	//seed=1625211875
	t.Logf("seed %v", seed)
	rand.Seed(seed)
	l := NewLocker()
	wg := &sync.WaitGroup{}
	values := make([]int, 10)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		e := i % len(values)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			<-time.After((time.Duration(rand.Intn(50)) + 10) * time.Millisecond)
			id := fmt.Sprintf("n-%v", e)
			if l.Lock(ctx, id); true {
				v := rand.Intn(100)
				<-time.After((time.Duration(rand.Intn(5))) * time.Millisecond)
				c := values[e]
				<-time.After((time.Duration(rand.Intn(5))) * time.Millisecond)
				values[e] = c + v
				<-time.After((time.Duration(rand.Intn(5))) * time.Millisecond)
				values[e] = c
			}
			cancel()
		}()
	}
	wg.Wait()
	t.Logf("%v %v", values, l)
	for _, v := range values {
		if v != 0 {
			t.Fail()
		}
	}

}

func TestLocalTryLocker(t *testing.T) {
	var seed = time.Now().Unix()
	//seed=1625211875
	t.Logf("seed %v", seed)
	rand.Seed(seed)
	l := NewLocker()
	wg := &sync.WaitGroup{}
	values := make([]int, 10)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		e := i % len(values)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			<-time.After((time.Duration(rand.Intn(50)) + 10) * time.Millisecond)
			id := fmt.Sprintf("n-%v", e)
			if ok, _ := l.TryLock(ctx, id); ok {
				values[e]++
				<-time.After(5 * time.Second)
			}
			cancel()
		}()
	}
	wg.Wait()
	t.Logf("%v %v", values, l)
	for _, v := range values {
		if v != 1 {
			t.Fail()
		}
	}
}
