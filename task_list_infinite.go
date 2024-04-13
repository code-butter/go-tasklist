package tasklist

import "time"

type TaskListInfinite[R any, E any] struct {
	OnResult          func(string, R)
	OnError           func(string, E)
	OnTimeout         func(string)
	TaskFunctions     []taskFunction[R, E]
	Workers           uint
	Timeout           time.Duration
	Working           bool
	Paused            bool
	Finishing         bool
	FinishedFunctions int
}

func NewTaskListInfinite[R any, E any](workers uint, timeout time.Duration) *TaskListInfinite[R, E] {
	return &TaskListInfinite[R, E]{Workers: workers, Timeout: timeout}
}

func (tm *TaskListInfinite[R, E]) Add(id string, fn func() (R, E)) *string {
	if tm.Finishing {
		err := "Cannot add new task, finishing up"
		return &err
	}
	tm.TaskFunctions = append(tm.TaskFunctions, taskFunction[R, E]{id, fn})
	return nil
}

func (tm *TaskListInfinite[R, E]) Start() {
	if tm.Working {
		return
	}
	tm.Working = true
	go tm.doWork()
}

func (tm *TaskListInfinite[R, E]) doWork() {
	next := make(chan bool, tm.Workers)
	var working uint = 0
	done := 0
	currentIndex := 0
	for tm.Working {
		index := currentIndex
		select {
		case <-next:
			working--
			done++
		default:
			if working >= tm.Workers {
				<-next
				working--
				done++
			}
		}
		tm.FinishedFunctions = done
		if tm.Paused || index > len(tm.TaskFunctions)-1 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		nextFunction := func() {
			performTaskFunction(tm.Timeout, tm.TaskFunctions[index], func(id string, result *taskResult[R, E]) {
				if result == nil && tm.OnTimeout != nil {
					tm.OnTimeout(id)
				} else if result != nil {
					if !isZeroValue(result.Result) && tm.OnResult != nil {
						tm.OnResult(id, result.Result)
					}
					if !isZeroValue(result.Error) && tm.OnError != nil {
						tm.OnError(id, result.Error)
					}
				}
				next <- true
			})
		}
		go nextFunction()
		working++
		currentIndex++
	}
}

func (tm *TaskListInfinite[R, E]) Pause() {
	tm.Paused = true
}

func (tm *TaskListInfinite[R, E]) Resume() {
	tm.Paused = false
}

func (tm *TaskListInfinite[R, E]) Finish() {
	for tm.Working {
		if tm.FinishedFunctions >= len(tm.TaskFunctions) {
			tm.Finishing = true
			tm.Working = false
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
