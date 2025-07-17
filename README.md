# computers-management
This repository provides a REST API for managing the computers of employees of a fictious company. The logic of this service 
is just an MVP. Possible points of improvements are highlighted in a separate section below.

# Registerd routes
These routes can be reached at 'http://localhost:8080':

POST /computers
GET /computers/{computerID}
GET /computers
PUT /computers/{computerID}
GET /employees/{employeeID}/computers
GET /employees/{employeeID}/computers/{computerID}

# How to run the server
Execute `go run cmd/main.go` in the project's root directory.

# Possible improvements
- Define OpenAPI specs for documenting the routes and their parameters and request and response bodies.
- A repository.GetAll() should have pagination implemented or a hard limit for requested resources is set on DB level.
