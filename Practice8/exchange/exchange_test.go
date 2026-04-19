package exchange

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//helper
func newTestServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}


func TestGetRate_Success(t *testing.T) {
	srv := newTestServer(http.StatusOK, `{"base":"USD","target":"EUR","rate":0.92}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	rate, err := svc.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if rate != 0.92 {
		t.Errorf("expected rate 0.92, got %f", rate)
	}
}

//API business error 

func TestGetRate_APIBusinessError_404(t *testing.T) {
	srv := newTestServer(http.StatusNotFound, `{"error":"invalid currency pair"}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "XYZ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "api error: invalid currency pair" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetRate_APIBusinessError_400(t *testing.T) {
	srv := newTestServer(http.StatusBadRequest, `{"error":"invalid currency pair"}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("INVALID", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "api error: invalid currency pair" {
		t.Errorf("unexpected error message: %v", err)
	}
}


func TestGetRate_MalformedJSON(t *testing.T) {
	srv := newTestServer(http.StatusOK, `Internal Server Error`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
	if !strings.HasPrefix(err.Error(), "decode error:") {
		t.Errorf("expected 'decode error:' prefix, got: %v", err)
	}
}

func TestGetRate_TruncatedJSON(t *testing.T) {
	srv := newTestServer(http.StatusOK, `{"base":"USD","rate":0.`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
	if !strings.HasPrefix(err.Error(), "decode error:") {
		t.Errorf("expected 'decode error:' prefix, got: %v", err)
	}
}



func TestGetRate_Timeout(t *testing.T) {
	// Server that hangs for 200ms
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rate":1.0}`))
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	svc.Client.Timeout = 50 * time.Millisecond // much shorter than server delay

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("expected 'network error' in message, got: %v", err)
	}
}


func TestGetRate_ServerPanic_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error on 500, got nil")
	}
	if err.Error() != "api error: internal server error" {
		t.Errorf("unexpected error: %v", err)
	}
}


func TestGetRate_EmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// write nothing
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error on empty body, got nil")
	}
	if !strings.HasPrefix(err.Error(), "decode error:") {
		t.Errorf("expected 'decode error:' prefix, got: %v", err)
	}
}
