package handler

type AddComputerRequest struct {
	Name                 string  `json:"name"`
	IPAddress            string  `json:"ip_address"`
	MACAddress           string  `json:"mac_address"`
	EmployeeAbbreviation *string `json:"employee_abbreviation"`
	Description          *string `json:"description"`
}

type AddComputerResponse struct {
	ID int `json:"id"`
}

type GetComputerByIDResponse struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	IPAddress            string  `json:"ip_address"`
	MACAddress           string  `json:"mac_address"`
	EmployeeAbbreviation *string `json:"employee_abbreviation,omitempty"`
	Description          *string `json:"description,omitempty"`
}

type GetComputersResponse struct {
	Computers []GetComputerByIDResponse `json:"computers"`
}

type UpdateComputerRequest struct {
	Name                 string  `json:"name"`
	IPAddress            string  `json:"ip_address"`
	MACAddress           string  `json:"mac_address"`
	EmployeeAbbreviation *string `json:"employee_abbreviation"`
	Description          *string `json:"description"`
}
