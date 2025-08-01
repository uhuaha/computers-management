package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	errs "uhuaha/computers-management/internal/errors"
	"uhuaha/computers-management/internal/mocks"
	"uhuaha/computers-management/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetComputerByIDHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                 string
		urlParam             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "valid request returns 200",
			urlParam: "1",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					GetComputer(1).
					Return(model.Computer{
						ID:         1,
						Name:       "TestPC",
						IPAddress:  "192.168.1.10",
						MACAddress: "AA:BB:CC:DD:EE:FF",
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1, "name":"TestPC", "ip_address":"192.168.1.10", "mac_address":"AA:BB:CC:DD:EE:FF"}`,
		},
		{
			name:     "invalid computerID in URL returns 400",
			urlParam: "abc",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// Service should not be called because parsing fails
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid URL parameter"}`,
		},
		{
			name:     "computer not found returns 404",
			urlParam: "42",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					GetComputer(42).
					Return(model.Computer{}, &errs.NotFoundError{Msg: "computer not found"})
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"computer not found"}`,
		},
		{
			name:     "internal service error returns 500",
			urlParam: "5",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().
					GetComputer(5).
					Return(model.Computer{}, fmt.Errorf("db connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to get computer by ID"}`,
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

			req := httptest.NewRequest(http.MethodGet, "/computers/"+tt.urlParam, nil)
			req = mux.SetURLVars(req, map[string]string{
				"computerID": tt.urlParam,
			})
			rec := httptest.NewRecorder()

			// Act
			handler.GetComputerByID(rec, req)

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
