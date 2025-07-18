package service

import (
	"fmt"
	"uhuaha/computers-management/internal/db/postgres/dbo"
	"uhuaha/computers-management/internal/model"
)

type ComputerRepository interface {
	AddComputer(computer dbo.Computer) (int, error)
	GetComputer(computerID int) (dbo.Computer, error)
	GetAllComputers() ([]dbo.Computer, error)
	UpdateComputer(computerID int, data dbo.Computer) error
	GetComputersByEmployee(employeeID int) ([]dbo.Computer, error)
	DeleteComputerFromEmployee(computerID, employeeID int) error
}

type ComputerMgmtService struct {
	repository ComputerRepository
}

func NewComputerMgmtService(repo ComputerRepository) *ComputerMgmtService {
	return &ComputerMgmtService{
		repository: repo,
	}
}

func (s *ComputerMgmtService) AddComputer(computer model.Computer) (int, error) {
	// TODO: check if abbreviation has only 3 characters

	computerDBO := convertComputerModelToDBO(computer)

	computerID, err := s.repository.AddComputer(computerDBO)
	if err != nil {
		return 0, fmt.Errorf("failed to add a computer: %w", err)
	}

	// TODO: check if 3 or more computers are assigned -> if yes: inform admin via msg service
	// Call GetComputersByEmployee()...

	return computerID, nil
}

func (s *ComputerMgmtService) GetComputer(computerID int) (model.Computer, error) {
	computerDBO, err := s.repository.GetComputer(computerID)
	if err != nil {
		return model.Computer{}, fmt.Errorf("failed to get computer with ID=%d: %w", computerID, err)
	}

	computer := convertComputerDBOToModel(computerDBO)

	return computer, nil
}

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

func (s *ComputerMgmtService) UpdateComputer(computerID int, data model.Computer) error {
	data.ID = computerID
	dataToBeUpdated := convertComputerModelToDBO(data)

	err := s.repository.UpdateComputer(computerID, dataToBeUpdated)
	if err != nil {
		return fmt.Errorf("failed to update the computer with ID=%d: %w", computerID, err)
	}

	return nil
}

func (s *ComputerMgmtService) GetComputersByEmployee(employeeID int) ([]model.Computer, error) {
	computerDBOs, err := s.repository.GetComputersByEmployee(employeeID)
	if err != nil {
		return []model.Computer{}, fmt.Errorf("failed to get computers for employee with ID=%d: %w", employeeID, err)
	}

	computers := make([]model.Computer, len(computerDBOs))
	for i, dbo := range computerDBOs {
		computers[i] = convertComputerDBOToModel(dbo)
	}

	return computers, nil
}

func (s *ComputerMgmtService) DeleteComputerFromEmployee(computerID, employeeID int) error {
	err := s.repository.DeleteComputerFromEmployee(computerID, employeeID)
	if err != nil {
		return fmt.Errorf("failed to delete computer (ID=%d) from employee (ID=%d): %w", computerID, employeeID, err)
	}

	return nil
}
