package mode

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// Dev is used for development
	Dev = "dev"
	// Prod is used for production
	Prod = "prod"
	// TestDev is used for testing
	TestDev = "test_dev"
	// Default mode
)

var (
	defaultMode = Dev
	mode        = defaultMode
)

// parseMode parses the mode
func parseMode(s string) (string, error) {
	switch s := strings.ToLower(s); s {
	case Dev, Prod, TestDev:
		return s, nil
	case "development":
		return Dev, nil
	case "release":
		return Prod, nil
	case "test":
		return TestDev, nil
	case "test-dev":
		return TestDev, nil
	default:
		return "", ErrInvalidMode
	}
}

// Get returns the mode
func Get() string {
	return mode
}

// Set sets the mode
func Set(newMode string) {
	if parseMode, err := parseMode(newMode); err == nil {
		mode = parseMode
	} else {
		mode = defaultMode
	}

	updateGinMode()
}

// GetDefault returns the default mode
func GetDefault() string {
	return defaultMode
}

// SetDefaultMode sets the default mode
func SetDefaultMode(newDefaultMode string) {
	if parseMode, err := parseMode(newDefaultMode); err == nil {
		mode = parseMode
	} else {
		mode = defaultMode
	}
}

// IsDev returns true if the mode is Dev or TestDev
func IsDev() bool {
	mode := Get()

	return mode == Dev || mode == TestDev
}

// updateGinMode updates the gin mode
func updateGinMode() {
	var mode string

	switch Get() {
	case Dev:
		mode = gin.DebugMode
	case TestDev:
		mode = gin.TestMode
	case Prod:
		mode = gin.ReleaseMode
	default:
		panic("unknown mode")
	}

	gin.SetMode(mode)
}
