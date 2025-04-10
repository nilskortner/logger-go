package filesupport

import (
	"fmt"
	"os"
)

func DeleteIfExists(path string) bool {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		// File does not exist
		return false
	} else if err != nil {
		// Some other error occurred
		fmt.Println("Error deleting file:", err)
		return false
	}
	// File was successfully deleted
	return true
}
