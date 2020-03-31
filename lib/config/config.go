package config

import (
	"log"
	"os"
	"strconv"
)

// It will panic if such config value don't exists
func GetStr(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panic(`Environmental variable [` + key + `] don't exists`)
	}
	return value
}

// It will panic if such config value don't exists
// or the value converts to integer unsuccessfully
func GetInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Panic(`Environmental variable [` + key + `] don't exists`)
	}
	output, err := strconv.Atoi(value)
	if err != nil {
		log.Panic(`Environmental variable [` + key + `] is not an integer`)
	}

	return output
}

// It will panic if such config value don't exists
func GetBytes(key string) []byte {
	value := os.Getenv(key)
	if value == "" {
		log.Panic(`Environmental variable [` + key + `] don't exists`)
	}
	return []byte(value)
}
