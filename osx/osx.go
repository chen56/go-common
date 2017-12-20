package osx

import "os"

func GetEnvOr(key string, default_ string) string {
	result := os.Getenv(key)
	if result==""{
		return default_
	}
	return result
}

