// internal/state/app_state_test.go
package state

import (
	"github.com/sirupsen/logrus"
	"sync"
	"testing"

	"github.com/StarTerrarium/hisame/internal/config"
)

// Test that InitialiseAppState initializes the singleton correctly
func TestInitialiseAppState(t *testing.T) {
	instance = nil
	once = sync.Once{}

	cfg := &config.UserConfig{
		LogLevel: "trace",
	}

	originalLogLevel := logrus.GetLevel()
	defer logrus.SetLevel(originalLogLevel)

	appState := InitialiseAppState(cfg)

	if appState == nil {
		t.Fatal("Expected AppState instance, got nil")
	}

	if appState.GetConfig() != cfg {
		t.Fatal("Expected configuration to be set in AppState")
	}

	// Ensure that GetAppState returns the same instance
	retrievedAppState := GetAppState()
	if appState != retrievedAppState {
		t.Fatal("GetAppState did not return the same instance as InitialiseAppState")
	}

	if logrus.GetLevel() != logrus.TraceLevel {
		t.Fatal("InitialiseAppState did not return the correct LogLevel")
	}

}

func TestGetAppStateWithoutInitialization(t *testing.T) {
	instance = nil
	once = sync.Once{}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic when calling GetAppState without initialization, but no panic occurred")
		}
	}()

	// This should cause a panic
	_ = GetAppState()
}
