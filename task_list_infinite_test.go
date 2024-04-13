package tasklist

import (
	"slices"
	"testing"
	"time"
)

func TestTaskListInfinite(t *testing.T) {
	taskList := NewTaskListInfinite[string, string](3, 100*time.Millisecond)
	results := make(map[string]string)
	errors := make(map[string]string)
	var timeouts []string
	taskList.OnResult = func(id string, s string) {
		results[id] = s
	}
	taskList.OnError = func(id string, s string) {
		errors[id] = s
	}
	taskList.OnTimeout = func(id string) {
		timeouts = append(timeouts, id)
	}
	taskList.Start()
	taskList.Add("id-1", successTask(time.Millisecond*50))
	taskList.Add("id-2", successTask(time.Millisecond*20))
	taskList.Add("id-3", failTask(time.Millisecond*90))
	taskList.Add("id-4", failTask(time.Millisecond*150))
	taskList.Finish()
	if results["id-1"] != "Success!" {
		t.Errorf("\"id-1\" should be 'Success!'")
	}
	if results["id-2"] != "Success!" {
		t.Errorf("\"id-2\" should be 'Success!'")
	}
	if errors["id-3"] != "Fail!" {
		t.Errorf("\"id-3\" should be 'Fail!'")
	}
	if !slices.Contains(timeouts, "id-4") {
		t.Errorf("timeouts should contain \"id-4\"")
	}
}
