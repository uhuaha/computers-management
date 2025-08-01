package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"uhuaha/computers-management/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDeleteComputerHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                 string
		urlParam             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "valid request deletes computer",
			urlParam: "1",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().DeleteComputer(1).Return(nil)
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
		},
		{
			name:     "invalid computerID parameter",
			urlParam: "abc",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// no call expected
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid URL parameter 'computerID'"}`,
		},
		{
			name:     "service returns error",
			urlParam: "2",
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				m.EXPECT().DeleteComputer(2).Return(fmt.Errorf("db error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to delete computer"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/computers/"+tt.urlParam, nil)
			req = mux.SetURLVars(req, map[string]string{"computerID": tt.urlParam})
			rec := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockComputerMgmtService(ctrl)
			tt.mockBehavior(mockService)

			handler := New(mockService)
			handler.DeleteComputer(rec, req)

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
