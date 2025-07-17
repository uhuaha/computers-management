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
	GetComputersByEmployeeID(w http.ResponseWriter, r *http.Request)
	DeleteComputerFromEmployee(w http.ResponseWriter, r *http.Request)
}

func New(handler Handler) *mux.Router {
	router := mux.NewRouter()

	// - The system administrator wants to be able to add a new computer to an employee
	// • The system administrator wants to be informed when an employee is assigned 3 or more computers
	router.HandleFunc("/computers", handler.AddComputer).Methods("POST")

	// • The system administrator wants to be able to get the data of a single computer
	router.HandleFunc("/computers/{computerID}", handler.GetComputerByID).Methods("GET")

	// • The system administrator wants to be able to get all computers
	router.HandleFunc("/computers", handler.GetAllComputers).Methods("GET")

	// • The system administrator wants to be able to assign a computer to another employee
	router.HandleFunc("/computers/{computerID}", handler.UpdateComputer).Methods("PUT")

	// • The system administrator wants to be able to get all assigned computers for an employee
	router.HandleFunc("/employees/{employeeID}/computers", handler.GetComputersByEmployeeID).Methods("GET")

	// • The system administrator wants to be able to remove a computer from an employee
	router.HandleFunc("/employees/{employeeID}/computers/{computerID}", handler.DeleteComputerFromEmployee).Methods("DELETE")

	return router
}
