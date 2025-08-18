// Package service implements the business logic for the computer management system,
// acting as a bridge between HTTP handlers and the storage layer.
package service

import (
	"fmt"
	"uhuaha/computers-management/internal/db/postgres/dbo"
	"uhuaha/computers-management/internal/model"

	"github.com/bdlm/log"
)

type ComputerRepository interface {
	AddComputer(computer dbo.Computer) (int, error)
	GetComputer(computerID int) (dbo.Computer, error)
	GetAllComputers() ([]dbo.Computer, error)
	UpdateComputer(computerID int, data dbo.Computer) error
	GetComputersByEmployee(employee string) ([]dbo.Computer, error)
	DeleteComputer(computerID int) error
}

type MessageSender interface {
	SendMessage(employeeAbbreviation string)
}

type ComputerMgmtService struct {
	repository ComputerRepository
	notifier   MessageSender
}

func NewComputerMgmtService(repo ComputerRepository, notifier MessageSender) *ComputerMgmtService {
	return &ComputerMgmtService{
		repository: repo,
		notifier:   notifier,
	}
}

// AddComputer stores a new computer and returns its generated ID. If there are 3 or more computers assigned to the same employee,
// it sends a notification to the system administrator.
func (s *ComputerMgmtService) AddComputer(computer model.Computer) (int, error) {
	computerDBO := convertComputerModelToDBO(computer)

	computerID, err := s.repository.AddComputer(computerDBO)
	if err != nil {
		return 0, fmt.Errorf("failed to add a computer: %w", err)
	}

	// Notify system administrator if there are 3 or more computers assigned to the given employee.
	if computer.EmployeeAbbreviation != nil {
		computers, err := s.GetComputersByEmployee(*computer.EmployeeAbbreviation)
		if err != nil {
			log.Error("failed to get computers: " + err.Error())
			return 0, fmt.Errorf("failed to get computers for employee %q: %w", *computer.EmployeeAbbreviation, err)
		}

		if len(computers) >= 3 {
			go s.notifier.SendMessage(*computer.EmployeeAbbreviation)
		}
	}

	return computerID, nil
}

// GetComputer retrieves a computer by its ID.
func (s *ComputerMgmtService) GetComputer(computerID int) (model.Computer, error) {
	computerDBO, err := s.repository.GetComputer(computerID)
	if err != nil {
		return model.Computer{}, fmt.Errorf("failed to get computer with ID=%d: %w", computerID, err)
	}

	computer := convertComputerDBOToModel(computerDBO)

	return computer, nil
}

// GetAllComputers returns a list of all computers in the system.
func (s *ComputerMgmtService) GetAllComputers() ([]model.Computer, error) {
	computerDBOs, err := s.repository.GetAllComputers()
	if err != nil {
		return []model.Computer{}, fmt.Errorf("failed to get all computers: %w", err)
	}

	computers := make([]model.Computer, len(computerDBOs))
	for i, dbo := range computerDBOs {
		computers[i] = convertComputerDBOToModel(dbo)
	}

	return computers, nil
}

// UpdateComputer updates the data of an existing computer identified by its ID.
func (s *ComputerMgmtService) UpdateComputer(computerID int, data model.Computer) error {
	data.ID = computerID
	computerDBO := convertComputerModelToDBO(data)

	err := s.repository.UpdateComputer(computerID, computerDBO)
	if err != nil {
		return fmt.Errorf("failed to update the computer with ID=%d: %w", computerID, err)
	}

	// Notify system administrator if there are 3 or more computers assigned to the given employee.
	if data.EmployeeAbbreviation != nil {
		computers, err := s.GetComputersByEmployee(*data.EmployeeAbbreviation)
		if err != nil {
			log.Error("failed to get computers: " + err.Error())
			return fmt.Errorf("failed to get computers for employee %q: %w", *data.EmployeeAbbreviation, err)
		}

		if len(computers) >= 3 {
			go s.notifier.SendMessage(*data.EmployeeAbbreviation)
		}
	}

	return nil
}

// GetComputersByEmployee retrieves all computers assigned to the specified employee.
func (s *ComputerMgmtService) GetComputersByEmployee(employee string) ([]model.Computer, error) {
	computerDBOs, err := s.repository.GetComputersByEmployee(employee)
	if err != nil {
		return []model.Computer{}, fmt.Errorf("failed to get computers for employee %q: %w", employee, err)
	}

	computers := make([]model.Computer, len(computerDBOs))
	for i, dbo := range computerDBOs {
		computers[i] = convertComputerDBOToModel(dbo)
	}

	return computers, nil
}

// DeleteComputer removes a computer by its ID.
func (s *ComputerMgmtService) DeleteComputer(computerID int) error {
	err := s.repository.DeleteComputer(computerID)
	if err != nil {
		return fmt.Errorf("failed to delete computer with ID=%d: %w", computerID, err)
	}

	return nil
}
