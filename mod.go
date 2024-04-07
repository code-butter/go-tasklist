package task_list

import "reflect"

type taskResult[R any, E any] struct {
	Result R
	Error  E
}

type taskFunction[R any, E any] struct {
	Id   string
	Func func() (R, E)
}

func isZeroValue[T any](value T) bool {
	zeroValue := reflect.Zero(reflect.TypeOf(value))
	return reflect.DeepEqual(value, zeroValue.Interface())
}
