package main

import (
	"log"
	"net/http"

	h "github.com/donnpebe/todoapi/lib"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
)

func main() {
	log.Println("Starting Server...")
	var err error

	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./public")))
	api := r.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/tasks", h.ErrorHandler(h.IndexTasksHandler)).Methods("GET")
	v1.HandleFunc("/tasks", h.ErrorHandler(h.CreateTaskHandler)).Methods("POST")
	v1.HandleFunc("/tasks/{id}", h.ErrorHandler(h.ShowTaskHandler)).Methods("GET")
	v1.HandleFunc("/tasks/{id}", h.ErrorHandler(h.UpdateTaskHandler)).Methods("PUT")
	v1.HandleFunc("/tasks/{id}/done", h.ErrorHandler(h.DoneTaskHandler)).Methods("PUT")
	v1.HandleFunc("/tasks/{id}/undone", h.ErrorHandler(h.UndoneTaskHandler)).Methods("PUT")
	v1.HandleFunc("/tasks/{id}", h.ErrorHandler(h.DeleteTaskHandler)).Methods("DELETE")
	http.Handle("/", r)

	log.Println("Starting mongo db session...")
	h.Session, err = mgo.Dial("localhost")
	defer h.Session.Close()
	if err != nil {
		panic(err)
	}

	h.Session.SetMode(mgo.Monotonic, true)
	h.Collection = h.Session.DB("Todo").C("tasks")

	log.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
