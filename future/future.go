package future

import (
	"context"

	"github.com/salesforce/gentler/try"
)

type Future[T any] struct {
	result    try.Try[T]
	wait      chan bool
	completed bool
}

func New[T any](f func() (T, error)) *Future[T] {
	return NewT(func() try.Try[T] {
		return try.New(f())
	})
}

func NewT[T any](f func() try.Try[T]) *Future[T] {
	fut := &Future[T]{
		wait: make(chan bool),
	}
	go func() {
		defer close(fut.wait)
		fut.result = f()
		fut.wait <- true
		fut.completed = true
	}()
	return fut
}

func Successful[T any](v T) *Future[T] {
	return New[T](func() (T, error) {
		return v, nil
	})
}

func Failed[T any](err error) *Future[T] {
	return New[T](func() (T, error) {
		var v T
		return v, err
	})
}

func (f *Future[T]) Get() try.Try[T] {
	return f.GetWithContext(context.Background())
}

func (f *Future[T]) GetWithContext(ctx context.Context) try.Try[T] {
	if f.completed {
		return f.result
	}

	select {
	case <-f.wait:
		return f.result
	case <-ctx.Done():
		return try.Failure[T](ctx.Err())
	}
}

func Map[U any, V any](fut *Future[U], fn func(U) try.Try[V]) *Future[V] {
	return NewT(func() try.Try[V] {
		val, err := fut.Get().Unpack()
		if err != nil {
			return try.Failure[V](err)
		}
		return fn(val)
	})
}

func FlatMap[U any, V any](fut *Future[U], fn func(U) *Future[V]) *Future[V] {
	return NewT(func() try.Try[V] {
		val, err := fut.Get().Unpack()
		if err != nil {
			return try.Failure[V](err)
		}
		return fn(val).Get()
	})
}
