package configs

import (
	"os"
)

const (
	defaultAddress = "localhost:8080"
)

func GetAddress() string {
	host := os.Getenv("ADDRESS")
	if host == "" {
		return defaultAddress
	}
	return host
}
