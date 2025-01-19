package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type mockAdmitHandler struct{}

func (m *mockAdmitHandler) v1(ar v1.AdmissionReview) *v1.AdmissionResponse {
	return &v1.AdmissionResponse{
		Allowed: true,
	}
}

type mockHeaders struct {
	key   string
	value string
}

func TestServe(t *testing.T) {
	tests := []struct {
		name           string
		object         map[string]interface{}
		headers        mockHeaders
		wantHTTPStatus int
	}{
		{
			name: "happy-path",
			object: map[string]interface{}{
				"kind":       "AdmissionReview",
				"apiVersion": "admission.k8s.io/v1",
				"request": map[string]interface{}{
					"uid": "123",
				},
			},
			headers: mockHeaders{
				key:   "Content-Type",
				value: "application/json",
			},
			wantHTTPStatus: 200,
		},
		{
			name: "bad-headers",
			object: map[string]interface{}{
				"kind":       "AdmissionReview",
				"apiVersion": "admission.k8s.io/v1",
				"request": map[string]interface{}{
					"uid": "123",
				},
			},
			headers: mockHeaders{
				key:   "Content-Type",
				value: "application/not-json",
			},
			wantHTTPStatus: 415,
		},
		{
			name: "bad-kind",
			object: map[string]interface{}{
				"kind":       "FailingReview",
				"apiVersion": "admission.k8s.io/v1",
				"request": map[string]interface{}{
					"uid": "123",
				},
			},
			headers: mockHeaders{
				key:   "Content-Type",
				value: "application/json",
			},
			wantHTTPStatus: 400,
		},
		{
			name: "unsupported-api-version",
			object: map[string]interface{}{
				"kind":       "AdmissionReview",
				"apiVersion": "admission.k8s.io/v1beta1",
				"request": map[string]interface{}{
					"uid": "123",
				},
			},
			headers: mockHeaders{
				key:   "Content-Type",
				value: "application/json",
			},
			wantHTTPStatus: 400,
		},
	}

	// os pkg must be imported to properly register logger init function for testing, it seems
	fmt.Println("ENVIRONMENT value is:", os.Getenv("ENVIRONMENT"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cr := &unstructured.Unstructured{
				Object: tt.object,
			}

			data, err := json.Marshal(cr)
			if err != nil {
				t.Fatal(err)
			}

			body := bytes.NewBuffer(data)
			req, err := http.NewRequest("POST", "/mutate", body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Header.Set(tt.headers.key, tt.headers.value)
				serve(w, r, newAdmitHandler((&mockAdmitHandler{}).v1))
			})

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantHTTPStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantHTTPStatus)
			}
		})
	}
}

func TestServeIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestServeVersion(t *testing.T) {
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeVersion)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "v0.0.1"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
