package task_list

import (
	"slices"
	"testing"
	"time"
)

func successTask(delay time.Duration) func() (string, string) {
	return func() (string, string) {
		time.Sleep(delay)
		return "Success!", ""
	}
}

func failTask(delay time.Duration) func() (string, string) {
	return func() (string, string) {
		time.Sleep(delay)
		return "", "Fail!"
	}
}

func TestTaskList(t *testing.T) {
	taskList := NewTaskList[string, string](3, 100*time.Millisecond)
	taskList.Add("id-1", successTask(time.Millisecond*50))
	taskList.Add("id-2", successTask(time.Millisecond*20))
	taskList.Add("id-3", failTask(time.Millisecond*90))
	taskList.Add("id-4", failTask(time.Millisecond*150))
	taskList.Work()
	if taskList.Results["id-1"] != "Success!" {
		t.Errorf("taskList.Results[\"id-1\"] should be 'Success!'")
	}
	if taskList.Results["id-2"] != "Success!" {
		t.Errorf("taskList.Results[\"id-2\"] should be 'Success!'")
	}
	if taskList.Errors["id-3"] != "Fail!" {
		t.Errorf("taskList.Errors[\"id-3\"] should be 'Fail!'")
	}
	if !slices.Contains(taskList.Timeouts, "id-4") {
		t.Errorf("taskList.Timeouts should contain \"id-4\"")
	}
}
