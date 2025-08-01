package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uhuaha/computers-management/internal/mocks"
	"uhuaha/computers-management/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddComputerHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                 string
		requestBody          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "valid JSON without optional fields",
			requestBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF"
                }`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					AddComputer(model.Computer{
						Name:       "TestPC",
						IPAddress:  "192.168.0.1",
						MACAddress: "AA:BB:CC:DD:EE:FF",
					}).
					Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name: "valid JSON with all fields",
			requestBody: `{
                    "name": "DevPC",
                    "ip_address": "10.0.0.2",
                    "mac_address": "11:22:33:44:55:66",
                    "employee_abbreviation": "STR",
                    "description": "Development Machine"
                }`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					AddComputer(model.Computer{
						Name:                 "DevPC",
						IPAddress:            "10.0.0.2",
						MACAddress:           "11:22:33:44:55:66",
						EmployeeAbbreviation: toPointer("STR"),
						Description:          toPointer("Development Machine"),
					}).
					Return(2, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":2}`,
		},
		{
			name: "service layer returns error",
			requestBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF"
                }`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					AddComputer(model.Computer{
						Name:       "TestPC",
						IPAddress:  "192.168.0.1",
						MACAddress: "AA:BB:CC:DD:EE:FF",
					}).
					Return(0, fmt.Errorf("something went wrong in the service layer"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to add computer"}`,
		},
		{
			name: "invalid request: employee abbreviation is not 3 characters long",
			requestBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF",
                    "employee_abbreviation": "Stefan",
                    "description": "Development Machine"
                }`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// No service call expected because of prior validation error
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid employee abbreviation"}`,
		},
		{
			name: "invalid JSON request",
			requestBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF"
                `, // missing closing brace
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// No service call expected because JSON is invalid
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid request body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)
			tt.mockBehavior(mockComputerMgmtService)

			req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := New(mockComputerMgmtService)

			// Act
			handler.AddComputer(rec, req)

			// Assert
			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

			if tt.expectedResponseBody != "" {
				body, _ := io.ReadAll(res.Body)
				assert.JSONEq(t, tt.expectedResponseBody, string(body))
			}
		})
	}
}

func toPointer(s string) *string {
	return &s
}
