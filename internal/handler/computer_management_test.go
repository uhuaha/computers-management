package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uhuaha/computers-management/internal/mocks"
	"uhuaha/computers-management/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type addComputerMock struct {
	expectedModel model.Computer
	returnID      int
	returnError   error
}

type getComputersByEmployeeMock struct {
	employee        string
	returnComputers []model.Computer
	returnError     error
}

type servicesMocks struct {
	addComputer            addComputerMock
	getComputersByEmployee getComputersByEmployeeMock
}

func TestAddComputerHandler(t *testing.T) {
	t.Run("valid JSON without optional fields", func(t *testing.T) {
		jsonBody := `{
				"name": "TestPC",
				"ip_address": "192.168.0.1",
				"mac_address": "AA:BB:CC:DD:EE:FF"
			}`

		req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Set up mocks

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)

		servicesMocks := servicesMocks{
			addComputer: addComputerMock{
				expectedModel: model.Computer{
					Name:       "TestPC",
					IPAddress:  "192.168.0.1",
					MACAddress: "AA:BB:CC:DD:EE:FF",
				},
				returnID:    1,
				returnError: nil,
			},
		}

		mockComputerMgmtService.EXPECT().
			AddComputer(servicesMocks.addComputer.expectedModel).
			Return(servicesMocks.addComputer.returnID, servicesMocks.addComputer.returnError)

		mockNotifier := mocks.NewMockNotifier(ctrl)

		// Create system under test

		handler := New(mockComputerMgmtService, mockNotifier)
		handler.AddComputer(rec, req)

		result := rec.Result()
		defer result.Body.Close()

		// Define expections

		expectedStatusCode := http.StatusCreated
		expectedContentType := "application/json"
		expectedReturnID := servicesMocks.addComputer.returnID

		// Assertions

		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Equal(t, expectedContentType, result.Header.Get("Content-Type"))

		var resp AddComputerResponse
		err := json.NewDecoder(result.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, expectedReturnID, resp.ID)
	})

	t.Run("valid JSON with all fields", func(t *testing.T) {
		jsonBody := `{
				"name": "DevPC",
				"ip_address": "10.0.0.2",
				"mac_address": "11:22:33:44:55:66",
				"employee_abbreviation": "STR",
				"description": "Development Machine"
			}`

		req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Set up mocks

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)

		servicesMocks := servicesMocks{
			addComputer: addComputerMock{
				expectedModel: model.Computer{
					Name:                 "DevPC",
					IPAddress:            "10.0.0.2",
					MACAddress:           "11:22:33:44:55:66",
					EmployeeAbbreviation: toPointer("STR"),
					Description:          toPointer("Development Machine"),
				},
				returnID:    2,
				returnError: nil,
			},
			getComputersByEmployee: getComputersByEmployeeMock{
				employee:        "STR",
				returnComputers: []model.Computer{{}},
				returnError:     nil,
			},
		}

		mockComputerMgmtService.EXPECT().
			AddComputer(servicesMocks.addComputer.expectedModel).
			Return(servicesMocks.addComputer.returnID, servicesMocks.addComputer.returnError)

		mockComputerMgmtService.EXPECT().
			GetComputersByEmployee(servicesMocks.getComputersByEmployee.employee).
			Return(servicesMocks.getComputersByEmployee.returnComputers, servicesMocks.getComputersByEmployee.returnError)

		mockNotifier := mocks.NewMockNotifier(ctrl)

		// Create system under test

		handler := New(mockComputerMgmtService, mockNotifier)
		handler.AddComputer(rec, req)

		result := rec.Result()
		defer result.Body.Close()

		// Define expections

		expectedStatusCode := http.StatusCreated
		expectedContentType := "application/json"
		expectedReturnID := servicesMocks.addComputer.returnID

		// Assertions

		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Equal(t, expectedContentType, result.Header.Get("Content-Type"))

		var resp AddComputerResponse
		err := json.NewDecoder(result.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, expectedReturnID, resp.ID)
	})

	t.Run("handler.AddComputer() fails due to service layer returning an internal server error", func(t *testing.T) {
		jsonBody := `{
			"name": "TestPC",
			"ip_address": "192.168.0.1",
			"mac_address": "AA:BB:CC:DD:EE:FF"
		}`

		req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Set up mocks

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)

		servicesMocks := servicesMocks{
			addComputer: addComputerMock{
				expectedModel: model.Computer{
					Name:       "TestPC",
					IPAddress:  "192.168.0.1",
					MACAddress: "AA:BB:CC:DD:EE:FF",
				},
				returnID:    0,
				returnError: fmt.Errorf("something went wrong in the service layer"),
			},
		}

		mockComputerMgmtService.EXPECT().
			AddComputer(servicesMocks.addComputer.expectedModel).
			Return(servicesMocks.addComputer.returnID, servicesMocks.addComputer.returnError)

		mockNotifier := mocks.NewMockNotifier(ctrl)

		// Set up system under test

		handler := New(mockComputerMgmtService, mockNotifier)
		handler.AddComputer(rec, req)

		result := rec.Result()
		defer result.Body.Close()

		// Define expectations

		expectedStatusCode := http.StatusInternalServerError

		// Assertions

		assert.Equal(t, expectedStatusCode, result.StatusCode)
	})

	t.Run("handler.AddComputer() fails due to an invalid JSON request", func(t *testing.T) {
		jsonBody := `{
				"name": "TestPC",
				"ip_address": "192.168.0.1",
				"mac_address": "AA:BB:CC:DD:EE:FF"
			`

		req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Set up mocks

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)
		mockNotifier := mocks.NewMockNotifier(ctrl)

		// Set up system under test

		handler := New(mockComputerMgmtService, mockNotifier)
		handler.AddComputer(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		// Define expectations

		expectedStatusCode := http.StatusBadRequest

		// Assertions

		assert.Equal(t, expectedStatusCode, res.StatusCode)
	})
}

func toPointer(s string) *string {
	return &s
}
