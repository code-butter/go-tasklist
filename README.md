Task List for Go
================

Churn through asynchronous tasks easily with Task List with support for worker limits and timeouts. 


Usage
=====

* Create an instance of TaskList with Success and Error types
* Specify number of concurrent workers
* Specify timeout
* Add functions that return the relevant Success or Error values

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