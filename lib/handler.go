package lib

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v1"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	Session    *mgo.Session
	Collection *mgo.Collection
)

type DeletedJSON struct {
	Deleted bool          `json:"deleted"`
	Id      bson.ObjectId `json:"id"`
}

type ErrorsJSON struct {
	Errors map[string][]string `json:"errors"`
}

type ErrorJSON struct {
	Errors string `json:"errors"`
}

func check(err interface{}) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func stringify(errs map[string][]error) map[string][]string {
	clearErrors := make(map[string][]string)
	for k, d := range errs {
		var errorText []string
		for _, err := range d {
			errorText = append(errorText, err.Error())
		}
		clearErrors[k] = errorText
	}
	return clearErrors
}

func ErrorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var (
				j   []byte
				err error
			)

			if e := recover(); e != nil {
				if ex, ok := e.(error); ok {
					j, err = json.Marshal(&ErrorJSON{ex.Error()})
					check(err)
				} else if ex, ok := e.(map[string][]error); ok {
					clearErrors := stringify(ex)
					j, err = json.Marshal(&ErrorsJSON{clearErrors})
					check(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(422)
				w.Write(j)
			}
		}()
		fn(w, r)
	}
}

func IndexTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks Tasks
	iter := Collection.Find(nil).Iter()
	task := NewTask()
	for iter.Next(task) {
		tasks = append(tasks, task.Clone())
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(tasks.ToJSON())
	check(err)

	w.Write(j)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := NewTask()
	err := json.NewDecoder(r.Body).Decode(task)
	check(err)

	valid, errs := validator.Validate(task)
	if !valid {
		check(errs)
	}

	objId := bson.NewObjectId()
	task.Id = objId
	task.Done = false

	err = Collection.Insert(task)
	check(err)

	log.Printf("Inserted new task %s with name %s", task.Id, task.Name)

	j, err := json.Marshal(task.ToJSON())
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func ShowTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	valid := bson.IsObjectIdHex(vars["id"])
	if !valid {
		check(errors.New("invalid id"))
	}
	id := bson.ObjectIdHex(vars["id"])
	task := NewTask()
	err := Collection.FindId(id).One(task)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(task.ToJSON())
	check(err)

	w.Write(j)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	valid := bson.IsObjectIdHex(vars["id"])
	if !valid {
		check(errors.New("invalid id"))
	}
	id := bson.ObjectIdHex(vars["id"])
	taskParams := NewTask()

	err := json.NewDecoder(r.Body).Decode(taskParams)
	check(err)

	valid, errs := validator.Validate(taskParams)
	if !valid {
		check(errs)
	}

	task := NewTask()
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"name": taskParams.Name}},
		ReturnNew: true,
	}
	_, err = Collection.FindId(id).Apply(change, task)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(task.ToJSON())
	check(err)

	w.Write(j)
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	valid := bson.IsObjectIdHex(vars["id"])
	if !valid {
		check(errors.New("invalid id"))
	}
	id := bson.ObjectIdHex(vars["id"])
	task := NewTask()
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"done": true}},
		ReturnNew: true,
	}
	_, err := Collection.FindId(id).Apply(change, task)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(task.ToJSON())
	check(err)

	w.Write(j)
}

func UndoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	valid := bson.IsObjectIdHex(vars["id"])
	if !valid {
		check(errors.New("invalid id"))
	}
	id := bson.ObjectIdHex(vars["id"])
	task := NewTask()
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"done": false}},
		ReturnNew: true,
	}
	_, err := Collection.FindId(id).Apply(change, task)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(task.ToJSON())
	check(err)

	w.Write(j)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	valid := bson.IsObjectIdHex(vars["id"])
	if !valid {
		check(errors.New("invalid id"))
	}
	id := bson.ObjectIdHex(vars["id"])
	err := Collection.RemoveId(id)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(&DeletedJSON{true, id})
	check(err)

	w.Write(j)
}
