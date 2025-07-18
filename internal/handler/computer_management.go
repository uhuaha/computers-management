package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"uhuaha/computers-management/internal/model"

	"github.com/bdlm/log"
	"github.com/gorilla/mux"
)

type ComputerMgmtService interface {
	AddComputer(computer model.Computer) (int, error)
	GetComputer(computerID int) (model.Computer, error)
	GetAllComputers() ([]model.Computer, error)
	UpdateComputer(computerID int, data model.Computer) error
	GetComputersByEmployee(employeeID int) ([]model.Computer, error)
	DeleteComputerFromEmployee(computerID, empoyeeID int) error
}

type ComputerMgmtHandler struct {
	computerMgmtService ComputerMgmtService
}

func New(service ComputerMgmtService) *ComputerMgmtHandler {
	return &ComputerMgmtHandler{
		computerMgmtService: service,
	}
}

func (c *ComputerMgmtHandler) AddComputer(w http.ResponseWriter, r *http.Request) {
	var data AddComputerRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error("failed to decode the request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	computer := model.Computer{
		Name:                 data.Name,
		IPAddress:            data.IPAddress,
		MACAddress:           data.MACAddress,
		EmployeeAbbreviation: data.EmployeeAbbreviation,
		Description:          data.Description,
	}

	computerID, err := c.computerMgmtService.AddComputer(computer)
	if err != nil {
		log.Error("failed to add computer: " + err.Error())
		http.Error(w, "Failed to add computer", http.StatusInternalServerError)
		return
	}

	response := AddComputerResponse{
		ID: computerID,
	}

	res, err := json.Marshal(response)
	if err != nil {
		log.Error("failed to encode response: " + err.Error())
		w.WriteHeader(http.StatusCreated)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(res); err != nil {
		log.Error("failed to write response body: " + err.Error())
	}
}

func (c *ComputerMgmtHandler) GetComputerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramComputerID := vars["computerID"]
	computerID, err := strconv.Atoi(paramComputerID)
	if err != nil {
		log.Error("failed to parse URL parameter computerID: " + err.Error())
		http.Error(w, "Invalid URL parameter", http.StatusBadRequest)
		return
	}

	computer, err := c.computerMgmtService.GetComputer(computerID)
	if err != nil {
		log.Error("failed to get computer by ID: " + err.Error())
		http.Error(w, "Failed to get computer by ID", http.StatusInternalServerError)
		return
	}

	response := convertComputerModelToDTO(computer)

	res, err := json.Marshal(response)
	if err != nil {
		log.Error("failed to encode response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		log.Error("failed to write response body: " + err.Error())
	}
}

func (c *ComputerMgmtHandler) GetAllComputers(w http.ResponseWriter, r *http.Request) {
	computers, err := c.computerMgmtService.GetAllComputers()
	if err != nil {
		log.Error("failed to get all computers: " + err.Error())
		http.Error(w, "Failed to get all computers", http.StatusInternalServerError)
		return
	}

	response := convertComputerModelsToDTOs(computers)

	res, err := json.Marshal(response)
	if err != nil {
		log.Error("failed to encode response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		log.Error("failed to write response body: " + err.Error())
	}
}

func (c *ComputerMgmtHandler) UpdateComputer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramComputerID := vars["computerID"]
	computerID, err := strconv.Atoi(paramComputerID)
	if err != nil {
		log.Error("failed to parse URL parameter computerID: " + err.Error())
		http.Error(w, "Invalid URL parameter", http.StatusBadRequest)
		return
	}

	var data UpdateComputerRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error("failed to decode the request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	computer := model.Computer{
		Name:                 data.Name,
		IPAddress:            data.IPAddress,
		MACAddress:           data.MACAddress,
		EmployeeAbbreviation: data.EmployeeAbbreviation,
		Description:          data.Description,
	}

	if err := c.computerMgmtService.UpdateComputer(computerID, computer); err != nil {
		log.Error("failed to update computer: " + err.Error())
		http.Error(w, "Failed to update computer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ComputerMgmtHandler) GetComputersByEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramEmployeeID := vars["employeeID"]
	employeeID, err := strconv.Atoi(paramEmployeeID)
	if err != nil {
		log.Error("failed to parse URL parameter employeeID: " + err.Error())
		http.Error(w, "Invalid URL parameter", http.StatusBadRequest)
		return
	}

	computers, err := c.computerMgmtService.GetComputersByEmployee(employeeID)
	if err != nil {
		log.Error("failed to get computers by employee: " + err.Error())
		http.Error(w, "Failed to get computers by employee", http.StatusInternalServerError)
		return
	}

	response := convertComputerModelsToDTOs(computers)

	res, err := json.Marshal(response)
	if err != nil {
		log.Error("failed to encode response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		log.Error("failed to write response body: " + err.Error())
	}
}

func (c *ComputerMgmtHandler) DeleteComputerFromEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paramEmployeeID := vars["employeeID"]
	employeeID, err := strconv.Atoi(paramEmployeeID)
	if err != nil {
		log.Error("failed to parse URL parameter employeeID: " + err.Error())
		http.Error(w, "Invalid URL parameter", http.StatusBadRequest)
		return
	}

	paramComputerID := vars["computerID"]
	computerID, err := strconv.Atoi(paramComputerID)
	if err != nil {
		log.Error("failed to parse URL parameter computerID: " + err.Error())
		http.Error(w, "Invalid URL parameter", http.StatusBadRequest)
		return
	}

	if err := c.computerMgmtService.DeleteComputerFromEmployee(computerID, employeeID); err != nil {
		log.Error("failed to delete computer from employee: " + err.Error())
		http.Error(w, "Failed to delete computer from employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
