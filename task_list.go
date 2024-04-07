package task_list

import "time"

type TaskList[R any, E any] struct {
	Results       map[string]R
	Errors        map[string]E
	Timeouts      []string
	TaskFunctions []taskFunction[R, E]
	Workers       uint
	Timeout       time.Duration
}

func NewTaskList[R any, E any](workers uint, timeout time.Duration) *TaskList[R, E] {
	return &TaskList[R, E]{Workers: workers, Timeout: timeout}
}

func (tm *TaskList[R, E]) Add(id string, fn func() (R, E)) {
	tm.TaskFunctions = append(tm.TaskFunctions, taskFunction[R, E]{id, fn})
}

func (tm *TaskList[R, E]) Work() {
	total := len(tm.TaskFunctions)
	next := make(chan bool, total)
	var working uint = 0
	done := 0
	tm.Results = make(map[string]R)
	tm.Errors = make(map[string]E)
	tm.Timeouts = make([]string, 0)
	for i := 0; i < total; i++ {
		index := i
		nextFunction := func() {
			id := tm.TaskFunctions[index].Id
			doneChannel := make(chan taskResult[R, E])
			performFunction := func() {
				result, err := tm.TaskFunctions[index].Func()
				doneChannel <- taskResult[R, E]{Result: result, Error: err}
			}
			go performFunction()
			timeout := time.After(tm.Timeout)
			select {
			case result := <-doneChannel:
				if !isZeroValue(result.Result) {
					tm.Results[id] = result.Result
				}
				if !isZeroValue(result.Error) {
					tm.Errors[id] = result.Error
				}
			case <-timeout:
				tm.Timeouts = append(tm.Timeouts, id)
			}
			next <- true
		}
		go nextFunction()
		working++
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
	}
	for done < total {
		select {
		case <-next:
			done++
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
