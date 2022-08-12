package pkg

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	mu sync.Mutex

	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due,
	}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)
	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStore) DeleteAllTask() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

func (ts *TaskStore) GetAllTask() []Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	allTasks := make([]Task, 0, len(ts.tasks))
	for _, elem := range ts.tasks {
		allTasks = append(allTasks, elem)
	}
	return allTasks
}

func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var tasks []Task
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue
			}
		}
	}
	return tasks
}

func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
