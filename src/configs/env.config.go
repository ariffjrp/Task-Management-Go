package configs

import "os"

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "Error get Key error"
	}
	return value
}
