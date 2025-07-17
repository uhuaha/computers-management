package model

type Computer struct {
	ID                   int
	Name                 string
	IPAddress            string
	MACAddress           string
	EmployeeAbbreviation *string
	Description          *string
}
