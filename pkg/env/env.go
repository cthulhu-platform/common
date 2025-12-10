package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var v *viper.Viper
var initialized bool

// Init initializes Viper with support for .env files and environment variables
// It supports both direct environment variable names and APP_ prefixed names
// Example: GITHUB_CLIENT_ID or APP_GITHUB_CLIENT_ID
func Init(envFilePaths ...string) error {
	if initialized {
		return nil
	}

	v = viper.New()

	// Enable automatic environment variable reading
	// This allows Viper to read from os.Getenv()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set config type to env for .env files
	v.SetConfigType("env")

	// Try to load .env files from provided paths or default locations
	if len(envFilePaths) > 0 {
		for _, path := range envFilePaths {
			v.SetConfigFile(path)
			if err := v.ReadInConfig(); err == nil {
				// Successfully loaded a config file
				break
			}
		}
	} else {
		// Try default locations
		v.SetConfigName(".env")
		v.AddConfigPath(".")
		v.AddConfigPath("./")
		v.ReadInConfig() // Ignore error if .env doesn't exist
	}

	initialized = true
	return nil
}

// GetEnv retrieves an environment variable with support for:
// 1. Direct environment variable (e.g., GITHUB_CLIENT_ID)
// 2. APP_ prefixed environment variable (e.g., APP_GITHUB_CLIENT_ID)
// 3. Value from .env file (supports both GITHUB_CLIENT_ID and APP_GITHUB_CLIENT_ID)
// 4. Default value if none found
//
// Priority order:
// - APP_ prefixed env var (highest priority)
// - Direct env var
// - APP_ prefixed from .env file
// - Direct key from .env file
// - Default value
func GetEnv(key, def string) string {
	if !initialized {
		// Auto-initialize if not already done
		Init()
	}

	appKey := fmt.Sprintf("APP_%s", key)

	// Priority 1: APP_ prefixed environment variable (highest priority)
	if val := os.Getenv(appKey); val != "" {
		return val
	}

	// Priority 2: Direct environment variable
	if val := os.Getenv(key); val != "" {
		return val
	}

	// Priority 3: APP_ prefixed from .env file
	if v.IsSet(appKey) {
		if val := v.GetString(appKey); val != "" {
			return val
		}
	}

	// Priority 4: Direct key from .env file
	if v.IsSet(key) {
		if val := v.GetString(key); val != "" {
			return val
		}
	}

	// Return default
	return def
}

// GetViper returns the underlying Viper instance for advanced usage
func GetViper() *viper.Viper {
	if !initialized {
		Init()
	}
	return v
}
