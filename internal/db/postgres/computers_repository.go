package postgres

import (
	"uhuaha/computers-management/internal/db/postgres/dbo"
)

type DBConn interface{}

type Repository struct {
	dbConn DBConn
}

func NewRepository(dbConn DBConn) *Repository {
	return &Repository{
		dbConn: dbConn,
	}
}

func (r *Repository) AddComputer(computer dbo.Computer) (int, error) {
	computerID := 0

	return computerID, nil
}

func (r *Repository) GetComputer(computerID int) (dbo.Computer, error) {
	return dbo.Computer{}, nil
}

func (r *Repository) GetAllComputers() ([]dbo.Computer, error) {
	return []dbo.Computer{}, nil
}

func (r *Repository) UpdateComputer(computerID int, data dbo.Computer) error {
	return nil
}

func (r *Repository) GetComputersByEmployee(employeeID int) ([]dbo.Computer, error) {
	return []dbo.Computer{}, nil
}

func (r *Repository) DeleteComputerFromEmployee(computerID, employeeID int) error {
	return nil
}
