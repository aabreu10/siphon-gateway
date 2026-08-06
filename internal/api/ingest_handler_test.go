package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Since ingestHandler relies on the DB (WebhookRepo), RabbitMQ (Publisher), and SSE Hub, 
// a full unit test would require interface mocking. 
// For this SaaS pivot verification, we will test the basic routing and validation rejection.

func TestIngestHandler_InvalidEndpointID(t *testing.T) {
	// Create a dummy router
	r := chi.NewRouter()
	
	// Pass nil for dependencies as they shouldn't be reached due to early validation failure
	r.Post("/ingest/{endpoint_id}", ingestHandler(nil, nil, nil))

	// Invalid UUID format
	req, err := http.NewRequest("POST", "/ingest/invalid-uuid", bytes.NewBuffer([]byte(`{"test":true}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"invalid endpoint_id"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestIngestHandler_InvalidContentType(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/ingest/{endpoint_id}", ingestHandler(nil, nil, nil))

	validUUID := uuid.New().String()
	req, err := http.NewRequest("POST", "/ingest/"+validUUID, bytes.NewBuffer([]byte(`{"test":true}`)))
	if err != nil {
		t.Fatal(err)
	}
	// Intentionally omitting Content-Type: application/json

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnsupportedMediaType)
	}
}
