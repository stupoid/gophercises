package cmd

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const BucketName = "tasks"
const incompleteTasksKey = "incomplete"
const completedTasksKey = "completed"

func getTasks(key string) ([]Task, error) {
	var s []Task
	v, err := Get("/task/task.db", []byte(BucketName), []byte(key))
	if err != nil {
		return s, err
	}
	if v != nil {
		err = json.Unmarshal(v, &s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func putTasks(key string, tasks []Task) error {
	buf, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return Put("/task/task.db", []byte(BucketName), []byte(key), buf)
}

func putTask(key string, t Task) error {
	tasks, err := getTasks(key)
	if err != nil {
		return err
	}
	tasks = append(tasks, t)
	return putTasks(key, tasks)
}

func PopTask(key string, index int) (Task, error) {
	var t Task
	tasks, err := getTasks(key)
	if err != nil {
		return t, err
	}
	if index >= len(tasks) {
		return t, fmt.Errorf("index out of range [%d] with length %d", index, len(tasks))
	}
	t = tasks[index]
	tasks = RemoveIndex(tasks, index)
	err = putTasks(key, tasks)
	if err != nil {
		return t, err
	}
	return t, nil
}

func AddTask(text string) error {
	t := Task{
		Text:      text,
		CreatedAt: time.Now(),
	}
	return putTask(incompleteTasksKey, t)
}

func DoTask(id int) (Task, error) {
	t, err := PopTask(incompleteTasksKey, id)
	if err != nil {
		return t, err
	}
	t.UpdatedAt = time.Now()
	err = putTask(completedTasksKey, t)
	if err != nil {
		return t, err
	}
	return t, nil
}
