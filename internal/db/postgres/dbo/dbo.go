package dbo

import "database/sql"

type Computer struct {
	ID                   int            `db:"id"`
	Name                 string         `db:"name"`
	IPAddress            string         `db:"ip_address"`
	MACAddress           string         `db:"mac_address"`
	EmployeeAbbreviation sql.NullString `db:"employee_abbreviation"`
	Description          sql.NullString `db:"description"`
}
