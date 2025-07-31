// Package postgres provides utilities and wrappers for interacting with PostgreSQL databases.
package postgres

import (
	"database/sql"
	"fmt"
	"uhuaha/computers-management/internal/db/postgres/dbo"
	"uhuaha/computers-management/internal/errors"
)

type Repository struct {
	dbConn *sql.DB
}

func NewRepository(dbConn *sql.DB) *Repository {
	return &Repository{
		dbConn: dbConn,
	}
}

// AddComputer inserts a new computer into the database and returns its generated ID,
// or an error if the insertion fails.
func (r *Repository) AddComputer(computer dbo.Computer) (int, error) {
	query := `
		INSERT INTO computers (name, ip_address, mac_address, employee_abbreviation, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`

	stmt, err := r.dbConn.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare insert statement: %w", err)
	}

	defer stmt.Close()

	var computerID int
	err = stmt.QueryRow(computer.Name, computer.IPAddress, computer.MACAddress, computer.EmployeeAbbreviation, computer.Description).Scan(&computerID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert computer: %w", err)
	}

	return computerID, nil
}

// GetComputer retrieves a computer by its ID from the database.
// It returns the computer or an error if the record is not found or the query fails.
func (r *Repository) GetComputer(computerID int) (dbo.Computer, error) {
	stmt, err := r.dbConn.Prepare(`SELECT * FROM computers WHERE id = $1;`)
	if err != nil {
		return dbo.Computer{}, fmt.Errorf("failed to prepare select statement: %w", err)
	}

	defer stmt.Close()

	var computerDBO dbo.Computer
	err = stmt.QueryRow(computerID).Scan(
		&computerDBO.ID,
		&computerDBO.Name,
		&computerDBO.IPAddress,
		&computerDBO.MACAddress,
		&computerDBO.EmployeeAbbreviation,
		&computerDBO.Description,
	)
	if err == sql.ErrNoRows {
		return dbo.Computer{}, errors.NewNotFound("computer not found")
	} else if err != nil {
		return dbo.Computer{}, fmt.Errorf("failed to query computer: %w", err)
	}

	return computerDBO, nil
}

// GetAllComputers retrieves all computers from the database. The number of records returned is limited to 100.
// It returns a list of computers or an error if the query fails.
func (r *Repository) GetAllComputers() ([]dbo.Computer, error) {
	stmt, err := r.dbConn.Prepare(`SELECT * FROM computers LIMIT 100;`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare select statement: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to query computers: %w", err)
	}
	defer rows.Close()

	var computerDBOs []dbo.Computer

	for rows.Next() {
		var c dbo.Computer

		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.IPAddress,
			&c.MACAddress,
			&c.EmployeeAbbreviation,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		computerDBOs = append(computerDBOs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return computerDBOs, nil
}

// UpdateComputer updates an existing computer's details in the database. It returns an error if the update fails.
func (r *Repository) UpdateComputer(computerID int, data dbo.Computer) error {
	stmt, err := r.dbConn.Prepare(`
		UPDATE computers 
		SET name = $1, ip_address = $2, mac_address = $3, employee_abbreviation = $4, description = $5
		WHERE id = $6;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(data.Name, data.IPAddress, data.MACAddress, data.EmployeeAbbreviation, data.Description, computerID)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	return nil
}

// GetComputersByEmployee retrieves all computers associated with a specific employee abbreviation.
// It returns a list of computers or an error if the query fails.
func (r *Repository) GetComputersByEmployee(employee string) ([]dbo.Computer, error) {
	stmt, err := r.dbConn.Prepare(`SELECT * FROM computers WHERE employee_abbreviation = $1 LIMIT 100;`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare select statement: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(employee)
	if err != nil {
		return nil, fmt.Errorf("failed to query computers for an employee: %w", err)
	}
	defer rows.Close()

	var computerDBOs []dbo.Computer

	for rows.Next() {
		var c dbo.Computer

		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.IPAddress,
			&c.MACAddress,
			&c.EmployeeAbbreviation,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		computerDBOs = append(computerDBOs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return computerDBOs, nil
}

// DeleteComputer removes a computer from the database by its ID.
// It returns an error if the deletion fails.
func (r *Repository) DeleteComputer(computerID int) error {
	stmt, err := r.dbConn.Prepare(`DELETE FROM computers WHERE id = $1;`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(computerID)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	return nil
}
