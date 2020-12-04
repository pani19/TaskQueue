package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var terminateSignal chan bool
var mutex sync.Mutex
var queue []Task // Make a queue.

type Task struct {
	Id           string
	IsCompleted  bool
	Status       string    // untouched, completed, failed, timeout time.Time // when was the task created
	CreationTime time.Time // when was the task created
	TaskData     string    // field containing data about the task
}

func main() {
	fmt.Println(strings.ToUpper("Main Func"))
	go adder()

	go excutor()
	go cleaner()
	<-terminateSignal
}

func adder() {
	id := 1
	ticker := time.NewTicker(time.Second * time.Duration(1))
	for {
		select {
		case <-ticker.C:
			if len(queue) <= 10 {
				var newTask Task
				newTask.CreationTime = time.Now()
				newTask.Id = fmt.Sprintf("%d", id)
				newTask.IsCompleted = false
				newTask.Status = "untouched"
				newTask.TaskData = "Task Id #" + fmt.Sprintf("%d", id) + " has been assigned to do some important work."
				queue = append(queue, newTask)
				fmt.Println("New task added - ", newTask.Id)
				id++
			}
		}
	}
}

func excutor() {
	ticker := time.NewTicker(time.Second * time.Duration(5))
	for {
		select {
		case <-ticker.C:
			for index := range queue {
				if len(queue) == index {
					break
				}
				if len(queue) < index {
					break
				}
				if queue[index].IsCompleted == false {
					randNo := rand.Intn(10)

					if randNo%3 == 0 {
						queue[index].IsCompleted = true
						queue[index].Status = "completed"
						fmt.Println("task completed - ", queue[index].Id)
					} else if randNo%5 == 0 {
						queue[index].IsCompleted = true
						queue[index].Status = "failed"
						fmt.Println("task failed - ", queue[index].Id)
					}
				}
			}
		}
	}
}

func cleaner() {
	ticker := time.NewTicker(time.Second * time.Duration(20))
	for {
		select {
		case <-ticker.C:

			for index := range queue {

				if len(queue) == index {
					break
				}
				if len(queue) < index {
					break
				}
				if queue[index].IsCompleted == true {
					//remove task from the queue
					fmt.Println("task is going to clean - ", queue[index].Id)
					queue[index] = queue[len(queue)-1]
					queue[len(queue)-1] = Task{}
					queue = queue[:len(queue)-1]

				} else {
					t1 := queue[index].CreationTime
					t2 := time.Now()
					diff := t2.Sub(t1).Seconds()
					if diff >= 4 { //remove task from queue
						queue[index].IsCompleted = true
						queue[index].Status = "timeout"
						fmt.Println("task timeout - ", queue[index].Id)
					} else {
						//add this in queue again
					}
				}
			}
		}
	}
}
