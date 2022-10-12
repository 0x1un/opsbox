package main

import (
	cp "github.com/otiai10/copy"
	"os"
)

func MoveFile(sourcePath, destPath string) error {
	if err := cp.Copy(sourcePath, destPath); err != nil {
		return err
	}
	return os.RemoveAll(sourcePath)
}
