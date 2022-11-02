package utils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/selefra/selefra/pkg/modules"
	"gopkg.in/yaml.v3"
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

func GetHomeModulesPath() (string, error) {
	path, _, err := Home()
	modulesPath := filepath.Join(path, "download/modules")
	_, err = os.Stat(modulesPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(modulesPath, 0755)
		if err != nil {
			return "", err
		}
	}
	return modulesPath, nil
}

func GetTempPath() (string, error) {
	path, _, err := Home()
	if err != nil {
		return "", err
	}
	ociPath := filepath.Join(path, "temp")
	_, err = os.Stat(ociPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(ociPath, 0755)
		if err != nil {
			return "", err
		}
	}
	return ociPath, nil
}

func GetCredentialsPath() (string, error) {
	path, _, err := Home()
	if err != nil {
		return "", err
	}
	cred := filepath.Join(path, "credentials.json")
	_, err = os.Stat(cred)
	if errors.Is(err, os.ErrNotExist) {
		os.WriteFile(cred, []byte("{}"), 0644)
	}
	return cred, nil
}

func SetCredentials(token string) error {
	credentials, err := GetCredentialsPath()
	if err != nil {
		return err
	}
	jsonbytes, err := os.ReadFile(credentials)
	if err != nil {
		return err
	}
	var jsonmap map[string]string
	err = json.Unmarshal(jsonbytes, &jsonmap)
	if err != nil {
		return err
	}
	jsonmap["token"] = token
	jsonbytes, err = json.Marshal(jsonmap)
	if err != nil {
		return err
	}
	err = os.Remove(credentials)
	if err != nil {
		return err
	}
	err = os.WriteFile(credentials, jsonbytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetCredentialsToken() (string, error) {
	path, err := GetCredentialsPath()
	if err != nil {
		return "", err
	}
	jsonbytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var jsonmap map[string]string
	err = json.Unmarshal(jsonbytes, &jsonmap)
	if err != nil {
		return "", err
	}
	return jsonmap["token"], nil
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

const ROW = "https://raw.githubusercontent.com/piexlmax/registry"

type ModulesMetadata struct {
	Name          string   `json:"name" yaml:"name"`
	LatestVersion string   `json:"latest-version" yaml:"latest-version"`
	LatestUpdate  string   `json:"latest-update" yaml:"latest-update"`
	Introduction  string   `json:"introduction" yaml:"introduction"`
	Versions      []string `json:"versions" yaml:"versions"`
}

type ModulesSupplement struct {
	PackageName string `json:"package-name" yaml:"package-name"`
	Source      string `json:"source" yaml:"source"`
	Checksums   string `json:"checksums" yaml:"checksums"`
}

func getModulesMetadata(ctx context.Context, modulesName string) (ModulesMetadata, error) {
	var metadata ModulesMetadata
	_url := ROW + "/main/module/" + modulesName + "/metadata.yaml"

	body, err := Request(ctx, "GET", _url, nil)
	if err != nil {
		return metadata, err
	}
	err = yaml.Unmarshal(body, &metadata)
	return metadata, err
}

func getModulesModulesSupplement(ctx context.Context, modulesName, version string) (ModulesSupplement, error) {
	var supplement ModulesSupplement
	_url := ROW + "/main/module/" + modulesName + "/" + version + "/supplement.yaml"

	body, err := Request(ctx, "GET", _url, nil)
	if err != nil {
		return supplement, err
	}
	err = yaml.Unmarshal(body, &supplement)
	return supplement, err
}

func ModulesUpdate(modulesName string) error {
	_, config, err := Home()
	if err != nil {
		return err
	}
	c, err := os.ReadFile(config)
	if err != nil {
		return err
	}
	var configMap = make(map[string]string)
	err = json.Unmarshal(c, &configMap)
	if err != nil {
		return err
	}
	metadata, err := getModulesMetadata(context.Background(), modulesName)
	if err != nil {
		return err
	}
	if configMap["modules"+"/"+modulesName] == metadata.LatestVersion {
		return nil
	} else {
		supplement, err := getModulesModulesSupplement(context.Background(), modulesName, metadata.LatestVersion)
		if err != nil {
			return err
		}
		modulesPath, err := GetHomeModulesPath()
		if err != nil {
			return err
		}
		url := supplement.Source + "/releases/download/" + metadata.LatestVersion + "/" + modulesName + ".zip"
		err = os.RemoveAll(filepath.Join(modulesPath, modulesName))
		if err != nil {
			return err
		}
		err = modules.DownloadModule(url, modulesPath)
		if err != nil {
			return err
		}
		configMap["modules"+"/"+modulesName] = metadata.LatestVersion
		c, err := json.Marshal(configMap)
		if err != nil {
			return err
		}
		err = os.Remove(config)
		if err != nil {
			return err
		}
		err = os.WriteFile(config, c, 0644)
	}
	return nil
}
