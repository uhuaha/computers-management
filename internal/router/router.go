// Package router sets up the HTTP routes for the computer management service
// using the Gorilla Mux router.
package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	AddComputer(w http.ResponseWriter, r *http.Request)
	GetComputerByID(w http.ResponseWriter, r *http.Request)
	GetAllComputers(w http.ResponseWriter, r *http.Request)
	UpdateComputer(w http.ResponseWriter, r *http.Request)
	GetComputersByEmployee(w http.ResponseWriter, r *http.Request)
	DeleteComputer(w http.ResponseWriter, r *http.Request)
}

// New creates and returns a new Gorilla Mux router configured with all
// routes for the computer management service.
func New(handler Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/computers", handler.AddComputer).Methods("POST")
	router.HandleFunc("/computers/{computerID}", handler.GetComputerByID).Methods("GET")
	router.HandleFunc("/computers", handler.GetAllComputers).Methods("GET")
	router.HandleFunc("/computers/{computerID}", handler.UpdateComputer).Methods("PUT")
	router.HandleFunc("/employees/{employee}/computers", handler.GetComputersByEmployee).Methods("GET")
	router.HandleFunc("/computers/{computerID}", handler.DeleteComputer).Methods("DELETE")

	return router
}
