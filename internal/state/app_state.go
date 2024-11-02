package state

import (
	"fmt"
	"github.com/StarTerrarium/hisame/internal/config"
	"github.com/StarTerrarium/hisame/internal/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type AppState struct {
	mutex sync.RWMutex

	config    *config.UserConfig
	authToken string // TODO: Likely don't need to store this once the AniList client is created.
}

// Use a Singleton to manage the application state
var instance *AppState
var once sync.Once

// InitialiseAppState initialises the AppState with the provided configuration.
// It ensures that AppState is only initialised once.
func InitialiseAppState(cfg *config.UserConfig) *AppState {
	once.Do(func() {
		instance = &AppState{
			config: cfg,
		}

		// Set log level if it is configured in the user configuration
		if cfg.LogLevel != "" {
			level, err := logrus.ParseLevel(cfg.LogLevel)
			if err != nil {
				logrus.Warnf("Invalid log level '%s' in configuration; Continuing with level: ", logrus.GetLevel().String())
			} else {
				// Do not bypass log level env var when first loading user configuration
				utils.SetLogLevel(level, false)
			}
		}
	})
	return instance
}

// GetAppState returns the singleton instance of AppState.
// It panics if InitialiseAppState has not been called yet.
func GetAppState() *AppState {
	if instance == nil {
		logrus.Panic("Attempt to access AppState before initialising.  This should never happen and is a bug in the code.  Please open an issue.")
	}
	return instance
}

func (s *AppState) GetConfig() *config.UserConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.config
}

func getTokenFilePath() (string, error) {
	var dir string

	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		dir = filepath.Join(dir, "hisame")
	} else if runtime.GOOS == "darwin" {
		dir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "hisame")
	} else {
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
		dir = filepath.Join(dir, "hisame")
	}

	tokenFilePath := filepath.Join(dir, "token")
	return tokenFilePath, nil
}

// GetAuthToken gets the authentication token.
func (s *AppState) GetAuthToken() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.authToken
}

// SetAuthToken sets the authentication token.
func (s *AppState) SetAuthToken(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.authToken = token
}

// LoadAuthToken loads the token from disk and sets it in AppState.
func (s *AppState) LoadAuthToken() error {
	tokenFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Info("No token file found; starting without authentication")
			return nil
		}
		logrus.Errorf("Failed to read token file: %v", err)
		return err
	}

	s.authToken = string(data)
	logrus.Info("Authentication token loaded from disk")
	return nil
}

// SaveAuthToken saves the current authentication token to disk.
func (s *AppState) SaveAuthToken() error {
	tokenFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(tokenFilePath), 0700)
	if err != nil {
		logrus.Errorf("Failed to create token directory: %v", err)
		return err
	}

	// Write the token to the file
	err = os.WriteFile(tokenFilePath, []byte(s.authToken), 0600)
	if err != nil {
		logrus.Errorf("Failed to write token file: %v", err)
		return err
	}

	logrus.Info("Authentication token saved to disk")
	return nil
}

// ClearAuthToken clears the token from AppState and deletes it from disk.
func (s *AppState) ClearAuthToken() error {
	tokenFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.authToken = ""

	err = os.Remove(tokenFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("Failed to delete token file: %v", err)
			return err
		}
	}

	logrus.Info("Authentication token file deleted")
	return nil
}
