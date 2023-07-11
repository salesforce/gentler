package try

type Try[T any] struct {
	success T
	failure error
}

func New[T any](value T, err error) Try[T] {
	return Try[T]{
		success: value,
		failure: err,
	}
}

func Success[T any](value T) Try[T] {
	return Try[T]{
		success: value,
		failure: nil,
	}
}

func Failure[T any](err error) Try[T] {
	return Try[T]{
		failure: err,
	}
}

func (t Try[T]) Success() T {
	return t.success
}

func (t Try[T]) Failure() error {
	return t.failure
}

func (t Try[T]) Unpack() (T, error) {
	return t.success, t.failure
}
