package skipgo

type Iter[T any] interface {
	Next() (T, bool)
}
