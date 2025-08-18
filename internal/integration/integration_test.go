//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"uhuaha/computers-management/internal/handler"
	"uhuaha/computers-management/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	internal_postgres "uhuaha/computers-management/internal/db/postgres"
)

const (
	testDBName          = "testdb"
	testDBUser          = "testuser"
	testDBPassword      = "testpw"
	testDBHost          = "localhost"
	defaultPostgresPort = "5432"
)

var (
	db                  *sql.DB
	h                   *handler.ComputerMgmtHandler
	wg                  sync.WaitGroup
	notificationPayload []byte
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start PostgreSQL test container
	testContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPassword),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	if err != nil {
		log.Fatalf("failed to start postgres test container: %v", err)
	}

	connectionString, err := testContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	externalPort, err := testContainer.MappedPort(ctx, defaultPostgresPort)
	if err != nil {
		log.Fatal("failed to get externally mapped port: ", err)
	}

	// Create connection to DB
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	defer db.Close()

	// Exectute migrations
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	databaseURL := "postgres://" + testDBUser + ":" + testDBPassword + "@" +
		testDBHost + ":" + externalPort.Port() + "/" + testDBName + "?sslmode=disable"

	migrationDriver, err := migrate.New(
		"file://"+filepath.Join(workingDirectory, "../../migrations"), // Path to migration files
		databaseURL,
	)
	if err != nil {
		log.Fatalf("failed to init migrations: %v", err)
	}

	if err := migrationDriver.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Create a test notifyServer for receiving notifications
	notifyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done() // Mark that notifier is finished

		body, _ := io.ReadAll(r.Body)
		notificationPayload = body
		w.WriteHeader(http.StatusOK)
	}))

	defer notifyServer.Close()

	// Create repository, services and API handler
	repository := internal_postgres.NewRepository(db)
	notifier := service.NewNotifier(notifyServer.URL)
	computerMgmtService := service.NewComputerMgmtService(repository, notifier)
	h = handler.New(computerMgmtService)

	// Run all tests
	exitCode := m.Run()

	// Teardown
	if err := testContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(exitCode)
}

func TestAddAndGetComputersIntegration(t *testing.T) {
	defer truncateTable()

	computersToBeAdded := []map[string]any{
		{
			"name":        "TestPC-01",
			"ip_address":  "192.168.1.100",
			"mac_address": "AA:BB:CC:DD:EE:F1",
		},
		{
			"name":                  "TestPC-02",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F2",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #1 for employee EMP",
		},
		{
			"name":                  "TestPC-03",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F3",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #2 for employee EMP",
		},
	}

	t.Run("Add all computers to the DB", func(t *testing.T) {
		for _, computer := range computersToBeAdded {
			resp, err := addComputer(computer)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)
		}
	})

	t.Run("Get a computer by its ID returns 200", func(t *testing.T) {
		computerID := 1
		resp := getComputerByID(computerID)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var computer handler.GetComputerByIDResponse
		err = json.Unmarshal(body, &computer)
		require.NoError(t, err)

		assert.Equal(t, computersToBeAdded[0]["name"], computer.Name)
		assert.Equal(t, computersToBeAdded[0]["ip_address"], computer.IPAddress)
		assert.Equal(t, computersToBeAdded[0]["mac_address"], computer.MACAddress)
	})

	t.Run("Trying to get a non-existing computer returns 404", func(t *testing.T) {
		computerID := 999
		resp := getComputerByID(computerID)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get all computers returns all previously added computers", func(t *testing.T) {
		resp := getAllComputers()
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var computers handler.GetComputersResponse
		err = json.Unmarshal(body, &computers)
		require.NoError(t, err)

		assert.Len(t, computers.Computers, len(computersToBeAdded))
	})

	t.Run("Get all computers by employee abbriviation", func(t *testing.T) {
		employeeAbbreviation := "EMP"
		resp := getComputersByEmployee(employeeAbbreviation)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var computers handler.GetComputersResponse
		err = json.Unmarshal(body, &computers)
		require.NoError(t, err)

		assert.Len(t, computers.Computers, 2)
		assert.Equal(t, employeeAbbreviation, *computers.Computers[0].EmployeeAbbreviation)
		assert.Equal(t, employeeAbbreviation, *computers.Computers[1].EmployeeAbbreviation)
	})
}

func TestAddAndUpdateComputersIntegration(t *testing.T) {
	defer truncateTable()

	computersToBeAdded := []map[string]any{
		{
			"name":        "TestPC-01",
			"ip_address":  "192.168.1.100",
			"mac_address": "AA:BB:CC:DD:EE:F1",
		},
		{
			"name":                  "TestPC-02",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F2",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #1 for employee EMP",
		},
		{
			"name":                  "TestPC-03",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F3",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #2 for employee EMP",
		},
	}

	t.Run("Add all computers to the DB returns 201", func(t *testing.T) {
		for _, computer := range computersToBeAdded {
			resp, err := addComputer(computer)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)
		}
	})

	t.Run("Update a computer", func(t *testing.T) {
		// Get all computers
		resp := getAllComputers()
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var computers handler.GetComputersResponse
		err = json.Unmarshal(body, &computers)
		require.NoError(t, err)

		assert.Len(t, computers.Computers, len(computersToBeAdded))

		// Update the first computer
		updateRequest := map[string]any{
			"name":                  "UpdatedPC-01",
			"ip_address":            "192.168.1.200",
			"mac_address":           "AA:BB:CC:DD:EE:GG",
			"employee_abbreviation": "STR",
			"description":           "Updated description for TestPC-01",
		}

		resp, err = updateComputer(computers.Computers[0].ID, updateRequest)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Get updated computer by ID and validate the fields
		resp = getComputerByID(computers.Computers[0].ID)
		defer resp.Body.Close()

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var updatedComputer handler.GetComputerByIDResponse
		err = json.Unmarshal(body, &updatedComputer)
		require.NoError(t, err)

		assert.Equal(t, "UpdatedPC-01", updatedComputer.Name)
		assert.Equal(t, "192.168.1.200", updatedComputer.IPAddress)
		assert.Equal(t, "AA:BB:CC:DD:EE:GG", updatedComputer.MACAddress)
		assert.Equal(t, "STR", *updatedComputer.EmployeeAbbreviation)
		assert.Equal(t, "Updated description for TestPC-01", *updatedComputer.Description)
	})
}

func TestNotificationIntegration(t *testing.T) {
	defer truncateTable()

	computersToBeAdded := []map[string]any{
		{
			"name":        "TestPC-01",
			"ip_address":  "192.168.1.100",
			"mac_address": "AA:BB:CC:DD:EE:F1",
		},
		{
			"name":                  "TestPC-02",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F2",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #1 for employee EMP",
		},
		{
			"name":                  "TestPC-03",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F3",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #2 for employee EMP",
		},
	}

	t.Run("Add all computers to the DB returns 201", func(t *testing.T) {
		for _, computer := range computersToBeAdded {
			resp, err := addComputer(computer)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)
		}
	})

	t.Run("Adding a third computer to the same employee sends a notification", func(t *testing.T) {
		// Add a third computer for the same employee
		data := map[string]any{
			"name":                  "TestPC-04",
			"ip_address":            "192.168.1.100",
			"mac_address":           "AA:BB:CC:DD:EE:F4",
			"employee_abbreviation": "EMP",
			"description":           "Test computer #3 for employee EMP",
		}

		wg.Add(1)

		resp, err := addComputer(data)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Wait for the notification to be sent and verify the payload.
		wg.Wait()

		require.NotEmpty(t, notificationPayload, "Expected a notification to be sent")

		var sentMessage service.NotificationPayload

		err = json.Unmarshal(notificationPayload, &sentMessage)
		require.NoError(t, err)

		assert.Equal(t, "warning", sentMessage.Level)
		assert.Equal(t, "There are 3 or more computers assigned to the same employee.", sentMessage.Message)
		assert.Equal(t, "EMP", sentMessage.EmployeeAbbreviation)
	})

	t.Run("Delete one of the computers that has an employee set", func(t *testing.T) {
		// Get all computers of a specific employee
		employee := "EMP"

		resp := getComputersByEmployee(employee)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var computers handler.GetComputersResponse
		err = json.Unmarshal(body, &computers)
		require.NoError(t, err)
		require.Len(t, computers.Computers, 3)

		// Delete the first computer
		computerID := computers.Computers[0].ID
		resp = deleteComputer(computerID)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify that there is one computer less in the DB for the same employee
		resp = getComputersByEmployee(employee)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(body, &computers)
		require.NoError(t, err)

		require.Len(t, computers.Computers, 2)
	})

	t.Run("Update the computer that has no employee set and send a notification upon update", func(t *testing.T) {
		resp := getAllComputers()
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var computers handler.GetComputersResponse
		err = json.Unmarshal(body, &computers)
		require.NoError(t, err)

		require.Len(t, computers.Computers, 3)

		// Find the computer that has no employee set
		computersWithoutEmployee := make([]handler.GetComputerByIDResponse, 0)

		for _, computer := range computers.Computers {
			if computer.EmployeeAbbreviation == nil {
				computersWithoutEmployee = append(computersWithoutEmployee, computer)
			}
		}

		require.Len(t, computersWithoutEmployee, 1)

		// Updating a computer should trigger sending a notification
		updateReq := map[string]any{
			"description":           "Updated description for TestPC",
			"employee_abbreviation": "EMP",
		}

		wg.Add(1)

		resp, err = updateComputer(computersWithoutEmployee[0].ID, updateReq)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Wait for the notification to be sent and verify the payload.
		wg.Wait()

		require.NotEmpty(t, notificationPayload, "Expected a notification to be sent")

		var sentMessage service.NotificationPayload

		err = json.Unmarshal(notificationPayload, &sentMessage)
		require.NoError(t, err)

		assert.Equal(t, "warning", sentMessage.Level)
		assert.Equal(t, "There are 3 or more computers assigned to the same employee.", sentMessage.Message)
		assert.Equal(t, "EMP", sentMessage.EmployeeAbbreviation)
	})
}

// truncateTable clears the computers table and resets the identity column.
func truncateTable() {
	_, err := db.Exec("TRUNCATE TABLE computers RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("failed to truncate table: %v", err)
	}

	log.Println("Table truncated successfully")
}

func addComputer(data map[string]any) (*http.Response, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/computers", bytes.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	h.AddComputer(rec, req)

	return rec.Result(), nil
}

func getComputersByEmployee(employee string) *http.Response {
	targetURL := "/computers/employee/" + employee

	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"employee": employee,
	})

	rec := httptest.NewRecorder()
	h.GetComputersByEmployee(rec, req)

	return rec.Result()
}

func getAllComputers() *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/computers", nil)

	rec := httptest.NewRecorder()
	h.GetAllComputers(rec, req)

	return rec.Result()
}

func getComputerByID(computerID int) *http.Response {
	targetComputerID := strconv.Itoa(computerID)
	targetURL := "/computers/" + targetComputerID

	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"computerID": targetComputerID,
	})

	rec := httptest.NewRecorder()
	h.GetComputerByID(rec, req)

	return rec.Result()
}

func updateComputer(computerID int, updateReq map[string]any) (*http.Response, error) {
	jsonBody, err := json.Marshal(updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update request: %w", err)
	}

	targetComputerID := strconv.Itoa(computerID)
	targetURL := "/computers/" + targetComputerID
	req := httptest.NewRequest(http.MethodPut, targetURL, bytes.NewReader(jsonBody))
	req = mux.SetURLVars(req, map[string]string{
		"computerID": targetComputerID,
	})

	rec := httptest.NewRecorder()

	h.UpdateComputer(rec, req)

	return rec.Result(), nil
}

func deleteComputer(computerID int) *http.Response {
	targetComputerID := strconv.Itoa(computerID)
	targetURL := "/computers/" + targetComputerID

	deleteReq := httptest.NewRequest(http.MethodDelete, targetURL, nil)
	deleteReq = mux.SetURLVars(deleteReq, map[string]string{
		"computerID": targetComputerID,
	})

	rec := httptest.NewRecorder()
	h.DeleteComputer(rec, deleteReq)

	return rec.Result()
}
