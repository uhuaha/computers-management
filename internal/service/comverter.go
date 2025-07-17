package service

import (
	"database/sql"
	"uhuaha/computers-management/internal/db/postgres/dbo"
	"uhuaha/computers-management/internal/model"
)

func convertComputerModelToDBO(c model.Computer) dbo.Computer {
	return dbo.Computer{
		ID:                   c.ID,
		Name:                 c.Name,
		IPAddress:            c.IPAddress,
		MACAddress:           c.MACAddress,
		EmployeeAbbreviation: stringToNullString(c.EmployeeAbbreviation),
		Description:          stringToNullString(c.Description),
	}
}

func stringToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	
	return sql.NullString{Valid: false}
}

func convertComputerDBOToModel(dbo dbo.Computer) model.Computer {
	return model.Computer{
		ID:                   dbo.ID,
		Name:                 dbo.Name,
		IPAddress:            dbo.IPAddress,
		MACAddress:           dbo.MACAddress,
		EmployeeAbbreviation: nullStringToPointer(dbo.EmployeeAbbreviation),
		Description:          nullStringToPointer(dbo.Description),
	}
}

func nullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}

	return nil
}
