package router

import (
	"net/http"

	"todoapp/repository"

	"github.com/gorilla/mux"
)

func Router() {
	router := mux.NewRouter()

	router.HandleFunc("/todo", repository.GetAllTodos).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo/{id}", repository.GetTodoById).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo", repository.CreateTodo).Methods("POST", "OPTIONS")
	router.HandleFunc("/todo/{id}", repository.UpdateTodo).Methods("PUT", "OPTIONS")
	router.HandleFunc("/todo/{id}", repository.DeleteTodo).Methods("DELETE", "OPTIONS")
	http.Handle("/", router)
}
