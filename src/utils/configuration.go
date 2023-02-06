package utils

import (
	"fmt"
	"os"
	"strings"
)

func GeneratePathForCss(path *string) (string, error) {
	if !isLocal(path) {
		return fmt.Sprintf("url('%s')", *path), nil
	}
	return ReadFile(path)
}

func GeneratePathForHTML(path *string) (string, error) {
	if !isLocal(path) {
		return *path, nil
	}
	return ReadFile(path)
}

func isLocal(path *string) bool {
	return !strings.HasPrefix(*path, "http")
}

func ReadFile(path *string) (string, error) {
	c, err := os.ReadFile(*path)
	if err != nil {
		return "", fmt.Errorf("impossible d'ouvrir le fichier : %s", err)
	}
	return string(c), nil
}
