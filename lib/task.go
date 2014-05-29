package lib

import (
	"github.com/Unknwon/com"
	"labix.org/v2/mgo/bson"
)

type Task struct {
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `json:"name" validate:"min=1"`
	Done bool          `json:"done"`
}

type TaskJSON struct {
	Task *Task `json:"task"`
}

func NewTask() *Task {
	return &Task{}
}

func (task *Task) Clone() *Task {
	cloneTask := *task
	return &cloneTask
}

func sanitize(value string) string {
	return com.Trim(com.HtmlEncode(value))
}

func (task *Task) Sanitize() *Task {
	task.Name = sanitize(task.Name)
	return task
}

func (task *Task) ToJSON() *TaskJSON {
	return &TaskJSON{task}
}

type Tasks []*Task

type TasksJSON struct {
	Tasks []*Task `json:"tasks"`
}

func (ts Tasks) ToJSON() *TasksJSON {
	return &TasksJSON{ts}
}
