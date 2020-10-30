package adapter

import (
	"io/ioutil"
	"os"
)

type FileOutputProvider interface {
	WriteFile(content string, filename string) (int, error)
	ReadFile(filename string) (string, error)
	RemoveFile(filename string) error
}

type OsFileOutputProvider struct {
}

func (*OsFileOutputProvider) WriteFile(content string, filename string) (int, error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	f.Sync()
	size, err := f.WriteString(content)
	return size, err
}

func (*OsFileOutputProvider) ReadFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (*OsFileOutputProvider) RemoveFile(filename string) error {
	return os.Remove(filename)
}
