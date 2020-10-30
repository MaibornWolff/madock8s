package adapter

import "errors"

type FileOutputProviderSuccessMock struct {
}

type FileOutputProviderErrorMock struct {
}

func (*FileOutputProviderSuccessMock) WriteFile(content string, filename string) (int, error) {
	return 100, nil
}

func (*FileOutputProviderSuccessMock) ReadFile(filename string) (string, error) {
	return "Content", nil
}

func (*FileOutputProviderSuccessMock) RemoveFile(filename string) error {
	return nil
}

func (*FileOutputProviderErrorMock) WriteFile(content string, filename string) (int, error) {
	return 0, errors.New("cannot write file")
}

func (*FileOutputProviderErrorMock) ReadFile(filename string) (string, error) {
	return "", errors.New("cannot read file")
}

func (*FileOutputProviderErrorMock) RemoveFile(filename string) error {
	return errors.New("file not found")
}
