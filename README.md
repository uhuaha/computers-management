# computers-management
This repository provides a REST API for managing the computers of employees of a fictious company. The logic of this service 
is just an MVP. Possible points of improvements are highlighted in a separate section below.

## Registered routes
These routes can be reached at `http://localhost:8081`:

- `POST /computers`
- `GET /computers/{computerID}`
- `GET /computers`
- `PUT /computers/{computerID}`
- `GET /employees/{employee}/computers`
- `DELETE /computers/{computerID}`

## How to run
Execute `docker compose up` (if you have Docker compose v2 installed) or `docker-compose up` (if you use v1 of Docker compose) to fire up the database (migrations are run implicitly) and the notify service.
Then, start the server by executing `go run cmd/main.go` in the project's root directory.

## How to test
Import the provided Postman collection and test the endpoints once the docker containers and the server are running. Execute `make test` in order to run all unit tests.

## Possible areas of improvement
- Define OpenAPI specs for documenting the routes and their parameters as well as their request and response bodies.
- Have all handler methods covered by unit tests. For brevity and example only handler.AddComputer() is covered.
- A repository.GetAll() should have pagination (LIMIT and OFFSET clauses) implemented or have a hard limit for requested resources on DB level (see: LIMIT clause).
- The computers table should have created_at and updated_at columns to be able to track dates of creation and update.
- Currently the delete repository method executes a hard delete of the given resource. Providing a deleted_at column and executing an UPDATE on the resource to be deleted leads to a soft delete which would keep the resource in the DB.
- Use context.WithTimeout() in the handler functions: It's recommended to pass a context through all levels and methods - down to the repository level - to make sure the overall call doesn't take longer than a maximum amount of time because the context times out after a specified time period.
- Catching, logging and recovering from panics that happen anywhere in the code is recommended. One could e.g. wrap the mux.Router with a recover middleware that would prevent the server from crashing silently.
- We would need many more unit tests and end-to-end tests. E2e testing could e.g. spin up an SQLight database for in-memory storage during testing.
- The storage connection's configuration could be read from environment variables specified in the docker-compose file and read by the concrete connection implementation - instead of providing a hard-coded connection string inside the code.
- The linting can be extended by proving an own golang-ci.yml file that has all necessary linters enabled.