TaskList for Go
================

Churn through asynchronous tasks easily with TaskList with support for worker limits and timeouts. 


Usage
=====

* Create an instance of `TaskList` or `TaskListInfinite` with Success and Error types. Specify concurrent workers and
  worker timeout.
* `TaskList` is designed to run through the supplied functions and then finish, with all results in fields `Results`, 
  `Errors`, `Timeouts`. Add functions via `Add` and then call `Work` to work through all the jobs. It returns when
  finished.
* `TaskListInfinite` is designed to run indefinitely. Specify the callbacks by assigning to `OnResult`, `OnError`, and 
  `OnTimeout`. Call `Start` and add functions via `Add`. Call `Pause` and `Unpause` to halt starting new tasks, 
  temporarily, or call `Finish` to work through the remaining tasks and then return.
* Add functions that return the relevant Success or Error values
  * Success or error values depend on the zero value of the type. 


TaskList Example
----------------
```go
package main

import (
	"fmt"
	tasklist "github.com/code-butter/go-tasklist"
	"io"
	"net/http"
	"time"
)

func main() {
    websites := []string{"google.com", "github.com", "yahoo.com", "microsoft.com"}
    taskList := tasklist.NewTaskList[string, string](2, 20*time.Second)
    for i, website := range websites {
        taskList.Add(fmt.Sprint(i), func() (string, string) {
            response, err := http.Get(fmt.Sprintf("https://%s", website))
            if err != nil {
                return "", fmt.Sprintf("There was an error getting the website: %s", err)
            }
            defer response.Body.Close()
            body, err := io.ReadAll(response.Body)
            if err != nil {
                return "", fmt.Sprintf("There was an error reading the HTML body: %s", err)
            }
            return string(body), ""
        })
    }
    taskList.Work()

    // Successful tasks are in taskList.Results["id"]
    // Errored tasks are in taskList.Error["id"]
    // Timed out IDs are in a slice at taskList.Timeout
}

```

Please see tests for more examples. 