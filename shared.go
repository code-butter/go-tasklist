package tasklist

import (
	"reflect"
	"time"
)

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

func performTaskFunction[R any, E any](timeoutDuration time.Duration, fn taskFunction[R, E], callback func(string, *taskResult[R, E])) {
	doneChannel := make(chan taskResult[R, E])
	performFunction := func() {
		result, err := fn.Func()
		doneChannel <- taskResult[R, E]{Result: result, Error: err}
	}
	go performFunction()
	timeout := time.After(timeoutDuration)
	select {
	case result := <-doneChannel:
		callback(fn.Id, &result)
	case <-timeout:
		callback(fn.Id, nil)
	}
}
