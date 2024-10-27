package state

import (
	"github.com/StarTerrarium/hisame/internal/config"
	"github.com/StarTerrarium/hisame/internal/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

type AppState struct {
	mutex sync.RWMutex

	config *config.UserConfig
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
