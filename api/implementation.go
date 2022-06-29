package api

import (
	"context"
	"fmt"
	"sync"
)

type Locker interface {
	// Lock / - захватывает имя и владеет им до завершения контекста
	Lock(ctx context.Context, name string) error
	// TryLock / - пробует захватить имя и в случае успеха владеет им до завершения контекста
	TryLock(ctx context.Context, name string) (bool, error)
}

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

func (l *lockersHolder) String() string {
	l.locker.Lock()
	defer l.locker.Unlock()
	return fmt.Sprintf("localLocker created: %v , deleted: %v waiters %v", l.createHolders, l.deleteHolders, l.waiters)
}

func NewLocker() Locker {
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
func (l *lockersHolder) Lock(ctx context.Context, name string) error {
	l.locker.Lock()
	defer l.locker.Unlock()
	var holder = l.getHolder(name)
	defer l.unlock(ctx, holder)
	holder.usage++
	for holder.locks > 0 && ctx.Err() == nil {
		holder.cond.Wait()
	}
	holder.locks++
	return nil
}

// TryLock / - пробует захватить имя и в случае успеха владеет им до завершения контекста
func (l *lockersHolder) TryLock(ctx context.Context, name string) (bool, error) {
	l.locker.Lock()
	defer l.locker.Unlock()
	var holder = l.getHolder(name)
	if holder.locks > 0 {
		return false, nil
	}
	holder.usage++
	holder.locks++
	l.unlock(ctx, holder)
	return true, nil
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
