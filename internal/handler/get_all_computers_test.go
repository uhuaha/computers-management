package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"uhuaha/computers-management/internal/mocks"
	"uhuaha/computers-management/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllComputersHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "success: return 200 with computers list",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					GetAllComputers().
					Return([]model.Computer{
						{ID: 1, Name: "PC1", IPAddress: "192.168.0.1", MACAddress: "AA:BB:CC:DD:EE:FF"},
						{ID: 2, Name: "PC2", IPAddress: "192.168.0.2", MACAddress: "11:22:33:44:55:66", EmployeeAbbreviation: toPointer("EMP"), Description: toPointer("Office PC")},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"computers":
			[
				{"id":1,"name":"PC1","ip_address":"192.168.0.1","mac_address":"AA:BB:CC:DD:EE:FF"},
				{"id":2,"name":"PC2","ip_address":"192.168.0.2","mac_address":"11:22:33:44:55:66","employee_abbreviation":"EMP","description":"Office PC"}]
			}`,
		},
		{
			name: "success: return 200 with empty list",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					GetAllComputers().
					Return([]model.Computer{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"computers":[]}`,
		},
		{
			name: "return 500 due to service error",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					GetAllComputers().
					Return(nil, fmt.Errorf("database unavailable"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockComputerMgmtService := mocks.NewMockComputerMgmtService(ctrl)
			tt.mockBehavior(mockComputerMgmtService)

			handler := New(mockComputerMgmtService)

			req := httptest.NewRequest(http.MethodGet, "/computers", nil)
			rec := httptest.NewRecorder()

			// Act
			handler.GetAllComputers(rec, req)

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
