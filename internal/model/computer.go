// Package model defines the core domain models for the computer management system.
package model

// Computer represents a computer in the management system, including
// its network information and optional employee assignment.
type Computer struct {
	ID                   int
	Name                 string
	IPAddress            string
	MACAddress           string
	EmployeeAbbreviation *string
	Description          *string
}
