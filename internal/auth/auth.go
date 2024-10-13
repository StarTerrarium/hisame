package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	callbackPort = "19331"
	callbackPath = "/callback"
	tokenPath    = "/token"
	clientID     = "18776"
)

type Auth struct {
	LoginURL     *url.URL
	tokenChannel chan string
	httpServer   *http.Server
}

func NewAuth() *Auth {
	return &Auth{
		LoginURL:     generateAuthURL(),
		tokenChannel: make(chan string, 1),
		httpServer:   nil,
	}
}

// StartCallbackServer starts the HTTP server listening for the callback from AniList.
func (auth *Auth) StartCallbackServer() error {
	logrus.Info("Starting auth callback server.")

	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, handleCallback)
	mux.HandleFunc(tokenPath, auth.handleToken())

	// Create auth listener early so we can report an error if we can't secure the port.
	listener, err := net.Listen("tcp", ":"+callbackPort)
	if err != nil {
		logrus.Errorf("Could not listen on port %s: %v", callbackPort, err)
		return err
	}

	auth.httpServer = &http.Server{
		Handler: mux,
	}

	go func() {
		if err := auth.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("Server error: %v", err)
		}
	}()

	return nil
}

// WaitForToken sits and waits for a token to be received on the channel.  This is a way to block and wait
// for a token.  Also accepts a context as an arg so we can stop waiting if the user cancels the login flow.
func (auth *Auth) WaitForToken(ctx context.Context) (string, error) {
	logrus.Debug("Waiting for token to arrive on /token endpoint")
	// Ensure the callback server is stopped after we finish waiting
	defer auth.StopCallbackServer()

	// Wait for the token to be received
	select {
	case <-ctx.Done():
		logrus.Debug("WaitForToken exiting because context is done")
		return "", ctx.Err()
	case token, ok := <-auth.tokenChannel:
		if !ok || token == "" {
			logrus.Warn("Failed to receive token")
			return "", errors.New("failed to receive token")
		}
		logrus.Info("Received token")
		return token, nil
	}
}

func (auth *Auth) StopCallbackServer() {
	if auth.httpServer == nil {
		logrus.Warn("Call to StopCallbackServer when server was not started")
		return
	}
	logrus.Debug("Stopping callback server..")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := auth.httpServer.Shutdown(ctx); err != nil {
		logrus.Error("Server shutdown failed: ", err)
	}
	logrus.Debug("Callback server shutdown successfully")
}

func generateAuthURL() *url.URL {
	loginURL, err := url.Parse(fmt.Sprintf("https://anilist.co/api/v2/oauth/authorize?client_id=%s&response_type=token", clientID))
	if err != nil {
		// For simplicity simply kill the application for now.
		logrus.Panicf("Failed to generate auth url: %v", err)
		panic("Failed to generate auth url.  Exiting application.")
	}
	return loginURL
}

func (auth *Auth) handleToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("Received post to %s endpoint", tokenPath)
		var data struct {
			Token string `json:"token"`
		}

		// Parse the token from the POST request body
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		logrus.Debugf("Token decoded: %s", data.Token)

		// Send the token to the channel
		auth.tokenChannel <- data.Token

		// Send auth success response back
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "token stored"})
	}
}

// handleCallback handles the callback from AniList after auth is successful.
// As we are using the implicit grant, the token is passed along as a URL fragment.  This is why we are returning
// some javascript in the page to have the browser extract that token, and forward it to our /token POST endpoint
// so that our application can finally receive the token.
func handleCallback(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Hisame Auth</title>
        <script>
            window.onload = function() {
                const fragment = window.location.hash.substring(1);
                const params = new URLSearchParams(fragment);
                const token = params.get("access_token");

                if (token) {
                    fetch("/token", {
                        method: "POST",
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ token: token })
                    }).then(response => response.json())
                    .then(data => {
                        document.body.innerHTML = "<h1>Token fetched successfully.  You can close this window.</h1>";
                    }).catch((error) => {
                        document.body.innerHTML = "<h1>Error retrieving token: " + error + "</h1>";
                    });
                } else {
                    document.body.innerHTML = "<h1>No token found in the URL fragment</h1>";
                }
            };
        </script>
    </head>
    <body>
        <h1>Processing OAuth Token...</h1>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, htmlContent)
	if err != nil {
		logrus.Errorf("Error handling callback: %v", err)
	}
}
