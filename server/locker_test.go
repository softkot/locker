package server

import (
	"context"
	"fmt"
	"github.com/softkot/locker/client"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {

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
	l := newLocker()
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
	l := newLocker()
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
			if l.TryLock(ctx, id) {
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
