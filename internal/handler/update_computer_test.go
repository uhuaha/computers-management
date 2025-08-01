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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateComputerHandler(t *testing.T) {
	type mockBehavior func(m *mocks.MockComputerMgmtService)

	tests := []struct {
		name                 string
		urlParam             string
		requestBody          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "success: valid request updates computer",
			urlParam: "1",
			requestBody: `{
				"name": "UpdatedPC",
				"ip_address": "10.0.0.10",
				"mac_address": "AA:BB:CC:DD:EE:00"
			}`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				expected := model.Computer{
					Name:       "UpdatedPC",
					IPAddress:  "10.0.0.10",
					MACAddress: "AA:BB:CC:DD:EE:00",
				}
				m.EXPECT().UpdateComputer(1, expected).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:        "invalid computerID in URL",
			urlParam:    "abc",
			requestBody: `{}`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// no call expected
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid URL parameter"}`,
		},
		{
			name:        "invalid JSON body",
			urlParam:    "2",
			requestBody: `{ "name": "Broken", `,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// no call expected
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid request body"}`,
		},
		{
			name:     "invalid employee abbreviation length",
			urlParam: "3",
			requestBody: `{
				"name": "PCWithBadAbbrev",
				"employee_abbreviation": "TOOLONG"
			}`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				// no call expected
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid employee abbreviation"}`,
		},
		{
			name:     "service layer returns error",
			urlParam: "4",
			requestBody: `{
				"name": "ErrorPC",
				"ip_address": "10.0.0.11",
				"mac_address": "AA:BB:CC:DD:EE:11"
			}`,
			mockBehavior: func(m *mocks.MockComputerMgmtService) {
				expected := model.Computer{
					Name:       "ErrorPC",
					IPAddress:  "10.0.0.11",
					MACAddress: "AA:BB:CC:DD:EE:11",
				}
				m.EXPECT().UpdateComputer(4, expected).Return(fmt.Errorf("db failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to update computer"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/computers/"+tt.urlParam, strings.NewReader(tt.requestBody))
			req = mux.SetURLVars(req, map[string]string{"computerID": tt.urlParam})
			rec := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockComputerMgmtService(ctrl)
			tt.mockBehavior(mockService)

			handler := New(mockService)
			handler.UpdateComputer(rec, req)

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
