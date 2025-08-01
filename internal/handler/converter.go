package handler

import (
	"uhuaha/computers-management/internal/model"
)

func convertComputerModelToDTO(computer model.Computer) GetComputerByIDResponse {
	return GetComputerByIDResponse{
		ID:                   computer.ID,
		Name:                 computer.Name,
		IPAddress:            computer.IPAddress,
		MACAddress:           computer.MACAddress,
		EmployeeAbbreviation: computer.EmployeeAbbreviation,
		Description:          computer.Description,
	}
}

func convertComputerModelsToDTOs(computers []model.Computer) GetComputersResponse {
	computerDTOs := make([]GetComputerByIDResponse, len(computers))

	for i, computer := range computers {
		computerDTOs[i] = GetComputerByIDResponse{
			ID:                   computer.ID,
			Name:                 computer.Name,
			IPAddress:            computer.IPAddress,
			MACAddress:           computer.MACAddress,
			EmployeeAbbreviation: computer.EmployeeAbbreviation,
			Description:          computer.Description,
		}
	}

	return GetComputersResponse{
		Computers: computerDTOs,
	}
}
