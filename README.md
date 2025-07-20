# computers-management
This repository provides a REST API for managing the computers of employees of a fictious company. The logic of this service 
is just an MVP. Possible points of improvements are highlighted in a separate section below.

# Registerd routes
These routes can be reached at 'http://localhost:8080':

POST /computers
GET /computers/{computerID}
GET /computers
PUT /computers/{computerID}
GET /employees/{employee}/computers
DELETE /computers/{computerID}

# How to run
Execute `docker compose up` (if you have Docker compose v2 installed) or `docker-compose up` (if you use v1 of Docker compose) to fire up the database (migrations are run implicitly).
Then, start the server by executing `go run cmd/main.go` in the project's root directory.

# Improvements for making this service more robust
- Define OpenAPI specs for documenting the routes and their parameters and request and response bodies.
- Have all handler methods covered by unit tests. For brevity and example only handler.AddComputer() is covered.

- A repository.GetAll() should have pagination implemented or a hard limit for requested resources is set on DB level.
- The computers table should have created_at and updated_at columns to be able to track dates of creation and update.
- Currently the delete repository method executes a hard delete of the given resource. Providing a deleted_at column and executing an UPDATE on the resource to be deleted leads to a soft delete which would keep the resource in the DB.

- Check for context timeout and cancellation by using context.WithTimeout() etc.: One would need to pass a context through all levels and methods - down to the repository level - to make sure the overall call doesn't take longer than required. As implemented at the moment the DB could stall and, thus, the DB call from the repository layer would stall as well as long as the DB stalls.
