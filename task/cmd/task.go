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

func getTasks(path, bucket, key string) ([]Task, error) {
	var s []Task
	v, err := Get(path, []byte(bucket), []byte(key))
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

func putTasks(path, bucket, key string, tasks []Task) error {
	buf, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return Put(path, []byte(bucket), []byte(key), buf)
}

func putTask(path, bucket, key string, t Task) error {
	tasks, err := getTasks(path, bucket, key)
	if err != nil {
		return err
	}
	tasks = append(tasks, t)
	return putTasks(path, bucket, key, tasks)
}

func PopTask(path, bucket, key string, index int) (Task, error) {
	var t Task
	tasks, err := getTasks(path, bucket, key)
	if err != nil {
		return t, err
	}
	if index >= len(tasks) {
		return t, fmt.Errorf("index out of range [%d] with length %d", index, len(tasks))
	}
	t = tasks[index]
	tasks = RemoveIndex(tasks, index)
	err = putTasks(path, bucket, key, tasks)
	if err != nil {
		return t, err
	}
	return t, nil
}
