package system

import "os"

func GetProperty(name string) string {
	return os.Getenv(name)
}
