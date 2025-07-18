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
	tests := []struct {
		name               string
		jsonBody           string
		expectedModel      model.Computer
		mockReturnID       int
		mockReturnError    error
		expectedStatusCode int
		expectResponse     bool
	}{
		{
			name: "valid JSON without optional fields",
			jsonBody: `{
				"name": "TestPC",
				"ip_address": "192.168.0.1",
				"mac_address": "AA:BB:CC:DD:EE:FF"
			}`,
			expectedModel: model.Computer{
				Name:       "TestPC",
				IPAddress:  "192.168.0.1",
				MACAddress: "AA:BB:CC:DD:EE:FF",
			},
			mockReturnID:       1,
			mockReturnError:    nil,
			expectedStatusCode: http.StatusCreated,
			expectResponse:     true,
		},
		{
			name: "valid JSON with all fields",
			jsonBody: `{
				"name": "DevPC",
				"ip_address": "10.0.0.2",
				"mac_address": "11:22:33:44:55:66",
				"employee_abbreviation": "STR",
				"description": "Development Machine"
			}`,
			expectedModel: model.Computer{
				Name:                 "DevPC",
				IPAddress:            "10.0.0.2",
				MACAddress:           "11:22:33:44:55:66",
				EmployeeAbbreviation: toPointer("STR"),
				Description:          toPointer("Development Machine"),
			},
			mockReturnID:       2,
			mockReturnError:    nil,
			expectedStatusCode: http.StatusCreated,
			expectResponse:     true,
		},
		{
			name: "handler.AddComputer() fails due to service layer returning an internal server error",
			jsonBody: `{
				"name": "TestPC",
				"ip_address": "192.168.0.1",
				"mac_address": "AA:BB:CC:DD:EE:FF"
			}`,
			expectedModel: model.Computer{
				Name:       "TestPC",
				IPAddress:  "192.168.0.1",
				MACAddress: "AA:BB:CC:DD:EE:FF",
			},
			mockReturnID:       0,
			mockReturnError:    fmt.Errorf("something went wrong in the service layer"),
			expectedStatusCode: http.StatusInternalServerError,
			expectResponse:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(tt.jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)
			mockComputerMgmtService.EXPECT().
				AddComputer(tt.expectedModel).
				Return(tt.mockReturnID, tt.mockReturnError)

			handler := New(mockComputerMgmtService)
			handler.AddComputer(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

			if tt.expectResponse {
				assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
				var resp AddComputerResponse
				err := json.NewDecoder(res.Body).Decode(&resp)
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturnID, resp.ID)
			}
		})
	}

	t.Run("handler.AddComputer() fails due to an invalid JSON request", func(t *testing.T) {
		jsonBody := `{
				"name": "TestPC",
				"ip_address": "192.168.0.1",
				"mac_address": "AA:BB:CC:DD:EE:FF"
			`
		expectedStatusCode := http.StatusBadRequest

		req := httptest.NewRequest(http.MethodPost, "/computers", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)

		handler := New(mockComputerMgmtService)
		handler.AddComputer(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, expectedStatusCode, res.StatusCode)
	})
}

func toPointer(s string) *string {
	return &s
}
