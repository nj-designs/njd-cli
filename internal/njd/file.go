package njd

import (
	"fmt"
	"os"
)

// IsFileReadable checks if supplied file exists and is not a directory
func IsFileReadable(absFilePath string) error {
	info, err := os.Stat(absFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", absFilePath)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", absFilePath)
	}
	return nil
}
