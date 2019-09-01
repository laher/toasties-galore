package tpi

import "os"

func Env(name string, default_ string) string {
	val := os.Getenv(name)
	if val == "" {
		return default_
	}
	return val
}
