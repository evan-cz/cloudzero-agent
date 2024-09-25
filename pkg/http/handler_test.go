package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudzero/cloudzero-insights-controller/pkg/hook"
	"github.com/cloudzero/cloudzero-insights-controller/pkg/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"
	admission "k8s.io/api/admission/v1"
)

type MockHandler struct {
	hook.Handler
}

func NewMockHandler() hook.Handler {
	m := &MockHandler{}
	m.Handler.Create = m.Create()
	return m.Handler
}

func (m *MockHandler) Create() hook.AdmitFunc {
	return func(r *hook.Request) (*hook.Result, error) {
		return &hook.Result{Allowed: true}, nil
	}
}

func TestServe(t *testing.T) {
	// Setup
	handler := handler()
	mockHandler := NewMockHandler()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&storage.ResourceTags{})
	handler_func := handler.Serve(mockHandler, db)

	// Test cases
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid content type",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid body",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           "invalid body",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid request",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           admission.AdmissionReview{Request: &admission.AdmissionRequest{UID: "12345"}},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			mockRequest, _ := http.NewRequest(tt.method, "/validate/deployment", bytes.NewReader(jsonBody))
			mockRequest.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()
			handler_func(rr, mockRequest)
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
