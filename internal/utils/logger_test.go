package utils

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func Test_getLogLevelFromEnv(t *testing.T) {
	originalEnv := os.Getenv(logLevelEnvVar)
	defer os.Setenv(logLevelEnvVar, originalEnv)

	testCases := []struct {
		name          string
		envValue      string
		expectedLevel logrus.Level
		expectedErr   error
	}{
		{"EnvVarNotSet", "", logrus.Level(0), errEnvVarNotSet},
		{"InvalidLevel", "invalid", logrus.Level(0), errInvalidLogLevel},
		{"ValidLevelTrace", "trace", logrus.TraceLevel, nil},
		{"ValidLevelDebug", "debug", logrus.DebugLevel, nil},
		{"ValidLevelInfo", "info", logrus.InfoLevel, nil},
		{"ValidLevelWarn", "warn", logrus.WarnLevel, nil},
		{"ValidLevelError", "error", logrus.ErrorLevel, nil},
		{"ValidLevelFatal", "fatal", logrus.FatalLevel, nil},
		{"ValidLevelPanic", "panic", logrus.PanicLevel, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(logLevelEnvVar, tc.envValue)
			level, err := getLogLevelFromEnv()

			if err != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected error '%v', got '%v'", tc.expectedErr, err)
				}
			} else {
				if level != tc.expectedLevel {
					t.Errorf("Expected level '%v', got '%v'", tc.expectedLevel, level)
				}
			}
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	originalEnv := os.Getenv(logLevelEnvVar)
	defer os.Setenv(logLevelEnvVar, originalEnv)

	testCases := []struct {
		name          string
		envValue      string
		bypassEnvVar  bool
		setLevel      logrus.Level
		expectedLevel logrus.Level
	}{
		{"EnvVarNotSet_BypassFalse", "", false, logrus.DebugLevel, logrus.DebugLevel},
		{"EnvVarNotSet_BypassTrue", "", true, logrus.DebugLevel, logrus.DebugLevel},
		{"EnvVarValid_BypassFalse", "info", false, logrus.DebugLevel, logrus.InfoLevel},
		{"EnvVarValid_BypassTrue", "info", true, logrus.DebugLevel, logrus.DebugLevel},
		{"EnvVarInvalid_BypassFalse", "invalid", false, logrus.DebugLevel, logrus.DebugLevel},
		{"EnvVarInvalid_BypassTrue", "invalid", true, logrus.DebugLevel, logrus.DebugLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(logLevelEnvVar, tc.envValue)
			logrus.SetLevel(logrus.InfoLevel) // Reset to a known state

			SetLogLevel(tc.setLevel, tc.bypassEnvVar)

			currentLevel := logrus.GetLevel()
			if currentLevel != tc.expectedLevel {
				t.Errorf("Expected log level '%s', got '%s'", tc.expectedLevel, currentLevel)
			}
		})
	}
}
