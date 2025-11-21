package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_WithDefaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	config := NewConfig()

	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "5432", config.Port)
	assert.Equal(t, "postgres", config.User)
	assert.Equal(t, "password", config.Password)
	assert.Equal(t, "roadmap", config.DBName)
	assert.Equal(t, "disable", config.SSLMode)
}

func TestNewConfig_WithEnvVars(t *testing.T) {
	os.Setenv("DB_HOST", "test-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "test-user")
	os.Setenv("DB_PASSWORD", "test-password")
	os.Setenv("DB_NAME", "test-db")
	os.Setenv("DB_SSLMODE", "require")

	config := NewConfig()

	assert.Equal(t, "test-host", config.Host)
	assert.Equal(t, "5433", config.Port)
	assert.Equal(t, "test-user", config.User)
	assert.Equal(t, "test-password", config.Password)
	assert.Equal(t, "test-db", config.DBName)
	assert.Equal(t, "require", config.SSLMode)

	// Cleanup
	os.Clearenv()
}

func TestConfig_DSN(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "roadmap",
		SSLMode:  "disable",
	}

	dsn := config.DSN()
	expected := "host=localhost port=5432 user=postgres password=password dbname=roadmap sslmode=disable"

	assert.Equal(t, expected, dsn)
}

func TestConfig_DSNForMigrate(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "roadmap",
		SSLMode:  "disable",
	}

	dsn := config.DSNForMigrate()
	expected := "postgres:password@localhost:5432/roadmap?sslmode=disable"

	assert.Equal(t, expected, dsn)
}

func TestGetEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_KEY", "test-value")

	value := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test-value", value)

	os.Unsetenv("TEST_KEY")
}

func TestGetEnv_WithDefault(t *testing.T) {
	os.Unsetenv("TEST_KEY")

	value := getEnv("TEST_KEY", "default-value")
	assert.Equal(t, "default-value", value)
}
