package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogLevel = logrus.InfoLevel
	logLevelEnvVar  = "HISAME_LOG_LEVEL"
)

var (
	errEnvVarNotSet    = errors.New("environment variable not set")
	errInvalidLogLevel = errors.New("invalid log level")
)

// InitLogger sets up the global logger with a level and file.
// It returns a cleanup function to be called when the application exits.
func InitLogger() func() {
	level := defaultLogLevel

	envLevel, err := getLogLevelFromEnv()
	if err != nil {
		switch {
		case errors.Is(err, errInvalidLogLevel):
			logrus.Warnf("Invalid log level '%s' in environment variable %s; using default level '%s'.",
				os.Getenv(logLevelEnvVar), logLevelEnvVar, defaultLogLevel)
		case errors.Is(err, errEnvVarNotSet):
			// Environment variable not set; proceed with default level
		default:
			logrus.Warnf("Error retrieving log level from environment: %v", err)
		}
	} else {
		level = envLevel
	}

	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout) // Default output

	var logFile *os.File

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logrus.Warnf("Error getting cache directory; file logging will be disabled: %v", err)
	} else {
		logPath := filepath.Join(cacheDir, "hisame", "log", "hisame.log")
		if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
			logrus.Warnf("Error creating log directory; file logging will be disabled: %v", err)
		} else {
			logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				logrus.Warnf("Error opening log file; file logging will be disabled: %v", err)
			} else {
				logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
				logrus.Debugf("Logging to file %s", logPath)
			}
		}
	}

	logrus.Infof("===== Welcome to Hisame (Log Level: %s) =====", level)

	// Return a cleanup function
	return func() {
		logrus.Info("Hisame is shutting down")
		if logFile != nil {
			if err := logFile.Close(); err != nil {
				logrus.Errorf("Error closing log file: %v", err)
			}
		}
	}
}

// SetLogLevel updates the log level of logrus.
// If bypassEnvVar is true, it sets the log level regardless of the environment variable.
func SetLogLevel(level logrus.Level, bypassEnvVar bool) {
	if !bypassEnvVar {
		if _, err := getLogLevelFromEnv(); err == nil {
			// Environment variable is set; do not override
			logrus.Infof("Log level not changed due to %s being set. Current level: %s", logLevelEnvVar, logrus.GetLevel())
			return
		}
	}
	logrus.Infof("Setting log level to %s", level)
	logrus.SetLevel(level)
}

// getLogLevelFromEnv parses the log level from the environment variable if set.
// Returns errEnvVarNotSet if the environment variable is not set,
// or errInvalidLogLevel if the value cannot be parsed into a log level.
func getLogLevelFromEnv() (logrus.Level, error) {
	envLogLevel := os.Getenv(logLevelEnvVar)
	if envLogLevel == "" {
		return logrus.Level(0), errEnvVarNotSet
	}

	parsedLevel, err := logrus.ParseLevel(envLogLevel)
	if err != nil {
		return logrus.Level(0), fmt.Errorf("%w: '%s'", errInvalidLogLevel, envLogLevel)
	}
	return parsedLevel, nil
}
