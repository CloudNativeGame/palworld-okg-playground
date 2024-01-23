package cmd

import (
	"os"
)
import "path"

const (
	defaultEnvFile = ".pal.env"
)

type EnvFile struct {
	path string
}

func NewEnvFile() *EnvFile {
	p, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return &EnvFile{
		path: path.Join(p, defaultEnvFile),
	}
}

func (e *EnvFile) Exists() bool {
	if _, err := os.Stat(e.path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (e *EnvFile) Read() (string, error) {
	return "", nil
}

func (e *EnvFile) Write() error {
	return nil
}
