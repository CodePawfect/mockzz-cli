package internals

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envCache map[string]string
	once     sync.Once
)

// LoadEnv loads the environment variables from .env only once
func LoadEnv() {
	once.Do(func() {
		var err error
		envCache, err = godotenv.Read()
		if err != nil {
			log.Fatal("Could not load .env file")
		}
	})
}

// EnvString retrieves the value of the environment variable from the cache with a fallback
func EnvString(key, fallback string) string {
	LoadEnv()

	if val, ok := envCache[key]; ok {
		return val
	}
	return fallback
}
