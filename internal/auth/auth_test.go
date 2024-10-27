package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAuth(t *testing.T) {
	a := NewAuth()
	if a == nil {
		t.Fatal("Expected NewAuth to return a non-nil Auth instance")
	}

	if a.LoginURL == nil {
		t.Fatal("Expected LoginURL to be initialised")
	}

	expectedURL := "https://anilist.co/api/v2/oauth/authorize?client_id=18776&response_type=token"
	if a.LoginURL.String() != expectedURL {
		t.Fatalf("Expected LoginURL to be %s, got %s", expectedURL, a.LoginURL.String())
	}

	if a.tokenChannel == nil {
		t.Fatal("Expected TokenChannel to be initialized")
	}

	if a.httpServer != nil {
		t.Fatal("Expected HTTPServer to be nil before starting the server")
	}
}

func TestStartCallbackServer(t *testing.T) {
	a1 := NewAuth()
	err := a1.StartCallbackServer()
	if err != nil {
		t.Fatalf("Expected StartCallbackServer to start without error, got %v", err)
	}

	// Attempt to start the server again on the same port to simulate an error
	a2 := NewAuth()
	err = a2.StartCallbackServer()
	if err == nil {
		t.Fatal("Expected StartCallbackServer to fail when port is already in use")
	}

	// Clean up
	a1.StopCallbackServer()
}

func TestHandleToken(t *testing.T) {
	a := NewAuth()

	// Create a test server using the handler
	handler := a.handleToken()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Prepare the token data
	tokenData := `{"token":"test_token"}`
	resp, err := http.Post(ts.URL, "application/json", strings.NewReader(tokenData))
	if err != nil {
		t.Fatalf("Expected POST request to succeed, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Read the response body
	var respData map[string]string
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if respData["status"] != "token stored" {
		t.Fatalf("Expected status 'token stored', got '%s'", respData["status"])
	}

	// Verify that the token was received
	select {
	case token := <-a.tokenChannel:
		if token != "test_token" {
			t.Fatalf("Expected token 'test_token', got '%s'", token)
		}
	default:
		t.Fatal("Expected token to be sent to TokenChannel")
	}
}

func TestHandleCallback(t *testing.T) {
	req := httptest.NewRequest("GET", "/callback", nil)
	w := httptest.NewRecorder()

	handleCallback(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/html" {
		t.Fatalf("Expected Content-Type 'text/html', got '%s'", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Simple check to see if the body contains expected content
	if !strings.Contains(string(body), "<title>Hisame Auth</title>") {
		t.Fatal("Response body does not contain expected HTML content")
	}
}

func TestWaitForToken(t *testing.T) {
	a := NewAuth()

	// Start a goroutine to send a token after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		a.tokenChannel <- "test_token"
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token, err := a.WaitForToken(ctx)
	if err != nil {
		t.Fatalf("Expected WaitForToken to return without error, got %v", err)
	}

	if token != "test_token" {
		t.Fatalf("Expected token 'test_token', got '%s'", token)
	}
}

func TestWaitForToken_EmptyToken(t *testing.T) {
	a := NewAuth()

	// Send an empty token
	go func() {
		time.Sleep(100 * time.Millisecond)
		a.tokenChannel <- ""
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token, err := a.WaitForToken(ctx)
	if err == nil {
		t.Fatal("Expected error when receiving empty token")
	}

	if token != "" {
		t.Fatalf("Expected empty token, got '%s'", token)
	}
}

func TestWaitForToken_ContextCancelled(t *testing.T) {
	a := NewAuth()

	// Create a context that we'll cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine to cancel the context after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	token, err := a.WaitForToken(ctx)
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context canceled error, got %v", err)
	}

	if token != "" {
		t.Fatalf("Expected empty token due to context cancellation, got '%s'", token)
	}
}

func TestStopCallbackServer(t *testing.T) {
	a := NewAuth()
	err := a.StartCallbackServer()
	if err != nil {
		t.Fatalf("Expected server to start without error, got %v", err)
	}

	// Stop the server
	a.StopCallbackServer()

	// Verify that the server is no longer accepting connections
	_, err = net.DialTimeout("tcp", ":"+callbackPort, 100*time.Millisecond)
	if err == nil {
		t.Fatal("Expected connection to fail after server is stopped")
	}
}
