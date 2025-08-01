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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetComputersByEmployeeHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                 string
		urlParam             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "valid request returns computers",
			urlParam: "ABC",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				computers := []model.Computer{
					{ID: 1, Name: "PC1", IPAddress: "192.168.0.1", MACAddress: "AA:BB:CC:DD:EE:FF"},
					{ID: 2, Name: "PC2", IPAddress: "192.168.0.2", MACAddress: "11:22:33:44:55:66"},
				}
				m.EXPECT().GetComputersByEmployee("ABC").Return(computers, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"computers":[
				{"id":1,"name":"PC1","ip_address":"192.168.0.1","mac_address":"AA:BB:CC:DD:EE:FF"},
				{"id":2,"name":"PC2","ip_address":"192.168.0.2","mac_address":"11:22:33:44:55:66"}
			]}`,
		},
		{
			name:     "invalid employee abbreviation length",
			urlParam: "AB",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// no call expected
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid URL parameter 'employee'"}`,
		},
		{
			name:     "service returns error",
			urlParam: "XYZ",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().GetComputersByEmployee("XYZ").Return(nil, fmt.Errorf("db failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to get computers by employee"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/computers/employee/"+tt.urlParam, nil)
			req = mux.SetURLVars(req, map[string]string{"employee": tt.urlParam})
			rec := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockComputerMgmtService(ctrl)
			tt.mockBehavior(mockService)

			handler := New(mockService)
			handler.GetComputersByEmployee(rec, req)

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
