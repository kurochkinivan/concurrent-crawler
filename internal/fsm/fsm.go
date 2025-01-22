package fsm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func InitializeDirectory(hostname string) error {
	err := os.RemoveAll(hostname)
	if err != nil {
		return fmt.Errorf("failed to remove directory, err: %w", err)
	}

	err = os.Mkdir(hostname, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory, err: %w", err)
	}

	return nil
}

func SaveHTMLFile(path string, content io.Reader) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory, err: %w", err)
	}
	
	fileName := filepath.Join(path, "index.html")
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file, err: %w", err)
	}

	_, err = io.Copy(f, content)
	if err != nil {
		return fmt.Errorf("failed to copy html, err: %w", err)
	}

	return nil
}
