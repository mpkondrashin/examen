/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

list.go

List of tasks
*/
package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"sandboxer/pkg/logging"
)

type TaskListInterface interface {
	NewTask(path string) ID
}

type TaskList struct {
	mx sync.RWMutex
	//changeMX sync.Mutex
	changed    chan struct{}
	Tasks      map[ID]*Task
	TasksCount ID
}

func NewList() *TaskList {
	l := &TaskList{
		Tasks:   make(map[ID]*Task),
		changed: make(chan struct{}, 1000),
	}
	logging.Debugf("%p XXX List Lock (in New List)", l)
	//l.changeMX.Lock()
	l.changed <- struct{}{}
	return l
}

func (l *TaskList) Updated() {
	if len(l.changed) > 0 {
		return
	}
	l.changed <- struct{}{}
}

func (l *TaskList) Length() int {
	return len(l.Tasks)
}

func (l *TaskList) Changes() chan struct{} {
	return l.changed
}

var ErrAlreadyExists = errors.New("task already exist")

func (l *TaskList) NewTask(path string) (ID, error) {
	defer l.lockUnlock()()
	for _, tsk := range l.Tasks {
		if path == tsk.Path {
			logging.Debugf("NewTask. Same path: %s", path)
			tsk.SubmitTime = time.Now()
			l.Updated()
			return 0, fmt.Errorf("%s: %w", path, ErrAlreadyExists)
		}
	}
	logging.Debugf("NewTask %d, %s", l.TasksCount, path)
	tsk := NewTask(l.TasksCount, path)
	l.Tasks[tsk.Number] = tsk
	l.Updated()
	l.TasksCount++
	return tsk.Number, nil
}

func (l *TaskList) DelByID(id ID) {
	defer l.lockUnlock()() //mx.Lock()
	logging.Debugf("DelByID, id = %d, len = %d", id, len(l.Tasks))
	delete(l.Tasks, id)
	l.Updated()
}

func (l *TaskList) Get(num ID) *Task {
	defer l.lockUnlock()()
	return l.Tasks[num]
}

func (l *TaskList) Task(num ID, callback func(tsk *Task) error) error {
	//defer l.lockUnlock()()
	tsk := l.Tasks[num]
	if tsk == nil {
		return fmt.Errorf("missing task #%d", num)
	}
	//defer tsk.lockUnlock()()
	return callback(tsk)
}

func (l *TaskList) lockUnlock() func() {
	logging.Debugf("Lock %p", l)
	l.mx.Lock()
	return l.unlock // func() {} //
}
func (l *TaskList) unlock() {
	logging.Debugf("Unlock %p", l)
	l.mx.Unlock()
}

func (l *TaskList) GetIDs() []ID {
	keys := make([]ID, len(l.Tasks))
	logging.Debugf("keys len = %d", len(l.Tasks))
	i := 0
	for k := range l.Tasks {
		keys[i] = k
		i++
	}
	return keys
}

func (l *TaskList) Process(callback func([]ID)) {
	defer l.lockUnlock()()
	keys := l.GetIDs()
	sort.Slice(keys, func(i, j int) bool {
		return l.Tasks[keys[i]].SubmitTime.After(l.Tasks[keys[j]].SubmitTime)
	})
	logging.Debugf("slice: %v", keys)
	callback(keys)
}

func (l *TaskList) Save(filePath string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
