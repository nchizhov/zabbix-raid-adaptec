package support

import "os"

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
