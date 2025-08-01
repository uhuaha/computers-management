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

func TestAddComputerHandler(t *testing.T) {
	type args struct {
		jsonBody string
	}
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                string
		args                args
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedContentType string
		expectedID          int
	}{
		{
			name: "valid JSON without optional fields",
			args: args{
				jsonBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF"
                }`,
			},
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					AddComputer(model.Computer{
						Name:       "TestPC",
						IPAddress:  "192.168.0.1",
						MACAddress: "AA:BB:CC:DD:EE:FF",
					}).
					Return(1, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: "application/json",
			expectedID:          1,
		},
		{
			name: "valid JSON with all fields",
			args: args{
				jsonBody: `{
                    "name": "DevPC",
                    "ip_address": "10.0.0.2",
                    "mac_address": "11:22:33:44:55:66",
                    "employee_abbreviation": "STR",
                    "description": "Development Machine"
                }`,
			},
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
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: "application/json",
			expectedID:          2,
		},
		{
			name: "service layer returns error",
			args: args{
				jsonBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF"
                }`,
			},
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					AddComputer(model.Computer{
						Name:       "TestPC",
						IPAddress:  "192.168.0.1",
						MACAddress: "AA:BB:CC:DD:EE:FF",
					}).
					Return(0, fmt.Errorf("something went wrong in the service layer"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "invalid request: employee abbreviation is not 3 characters long",
			args: args{
				jsonBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF",
                    "employee_abbreviation": "Stefan",
                    "description": "Development Machine"
                }`,
			},
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// No service call expected because of prior validation error
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid JSON request",
			args: args{
				jsonBody: `{
                    "name": "TestPC",
                    "ip_address": "192.168.0.1",
                    "mac_address": "AA:BB:CC:DD:EE:FF"
                `, // missing closing brace
			},
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// No service call expected because JSON is invalid
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)
			tt.mockBehavior(mockComputerMgmtService)

			req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(tt.args.jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := New(mockComputerMgmtService)

			// Act
			handler.AddComputer(rec, req)

			// Assert
			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, res.Header.Get("Content-Type"))
			}

			if tt.expectedID != 0 {
				var resp AddComputerResponse
				err := json.NewDecoder(res.Body).Decode(&resp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, resp.ID)
			}
		})
	}
}

func toPointer(s string) *string {
	return &s
}
