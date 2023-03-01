package utils

import (
	"fmt"
	"os"
	"strings"
)

const Config = ".config/local-searchengine/config.toml"
const Datas = ".local/share/local-searchengine/"

func GeneratePathForCss(path *string) string {
	if !isLocal(path) {
		return fmt.Sprintf("url('%s')", *path)
	}
	return fmt.Sprintf("url('/local%s')", *path)
}

func GeneratePathForHTML(path *string) string {
	if !isLocal(path) {
		return *path
	}
	return fmt.Sprintf("/local%s", *path)
}

func isLocal(path *string) bool {
	return !strings.HasPrefix(*path, "http")
}

func ReadFile(path *string) ([]byte, error) {
	c, err := os.ReadFile(*path)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier : %s", err)
	}
	return c, nil
}
