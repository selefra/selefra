package utils

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"
)

func Home() (string, string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", "", err
	}
	registryPath := filepath.Join(home, ".selefra")
	_, err = os.Stat(registryPath)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(registryPath, 0755)
		if err != nil {
			return "", "", err
		}
	}

	providerPath := filepath.Join(home, ".selefra", ".path")

	_, err = os.Stat(providerPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(providerPath, 0755)
		if err != nil {
			return "", "", err
		}
	}

	config := filepath.Join(home, ".selefra", ".path", "config.json")

	_, err = os.Stat(config)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(config, []byte("{}"), 0644)
		if err != nil {
			return "", "", err
		}
	}
	return registryPath, config, nil
}

func CreateSource(path, version string) string {
	return "selefra/" + path + "@" + version
}

func GetNameBySource(source string) string {
	path := filepath.Base(source)
	arr := strings.Split(path, "@")
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func GetPathBySource(source string) string {
	_, config, err := Home()
	if err != nil {
		return ""
	}
	c, err := os.ReadFile(config)
	if err != nil {
		return ""
	}
	var configMap = make(map[string]string)
	err = json.Unmarshal(c, &configMap)
	if err != nil {
		return ""
	}
	return configMap[source]
}
