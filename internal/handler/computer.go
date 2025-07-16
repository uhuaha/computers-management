package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ComputerMgmtHandler struct {
	// TODO: service
}

func New() *ComputerMgmtHandler {
	return &ComputerMgmtHandler{}
}

func (c *ComputerMgmtHandler) CreateComputer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Computer created")); err != nil {
		fmt.Println("Failed to write response: ", err)
	}
}

func (c *ComputerMgmtHandler) GetComputerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramComputerID := vars["computerID"]
	_, err := strconv.Atoi(paramComputerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("The URL parameter computerID must be an integer.")); err != nil {
			fmt.Println("Failed to write response: ", err)
		}

		return
	}

	if _, err := w.Write([]byte("Get computer with ID: " + paramComputerID)); err != nil {
		fmt.Println("Failed to write response: ", err)
	}
}

func (c *ComputerMgmtHandler) GetAllComputers(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("These are all computers that are currently registered...")); err != nil {
		fmt.Println("Failed to write response: ", err)
	}
}

func (c *ComputerMgmtHandler) UpdateComputer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramComputerID := vars["computerID"]
	_, err := strconv.Atoi(paramComputerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("The URL parameter computerID must be an integer.")); err != nil {
			fmt.Println("Failed to write response: ", err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
	if _, err := w.Write([]byte("Update computer with ID: " + paramComputerID)); err != nil {
		fmt.Println("Failed to write response: ", err)
	}
}

func (c *ComputerMgmtHandler) GetComputersByEmployeeID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramEmployeeID := vars["employeeID"]
	_, err := strconv.Atoi(paramEmployeeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("The URL parameter employeeID must be an integer.")); err != nil {
			fmt.Println("Failed to write response: ", err)
		}

		return
	}

	if _, err := w.Write([]byte("These are all computers that are currently registered for employee with ID: " + paramEmployeeID)); err != nil {
		fmt.Println("Failed to write response: ", err)
	}
}

func (c *ComputerMgmtHandler) DeleteComputerFromEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramEmployeeID := vars["employeeID"]
	_, err := strconv.Atoi(paramEmployeeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("The URL parameter employeeID must be an integer.")); err != nil {
			fmt.Println("Failed to write response: ", err)
		}

		return
	}

	paramComputerID := vars["computerID"]
	_, err = strconv.Atoi(paramComputerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("The URL parameter computerID must be an integer.")); err != nil {
			fmt.Println("Failed to write response: ", err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
	if _, err := w.Write([]byte("Deleted computer with ID " + paramComputerID + " for employee with ID: " + paramEmployeeID)); err != nil {
		fmt.Println("Failed to write response: ", err)
	}
}
