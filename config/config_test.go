package config

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestConfig_GetDSN(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	s := SelefraConfig{}
	err := s.GetConfig()
	if err != nil {
		t.Error(err)
	}
	dsn := s.Selefra.GetDSN()
	if len(dsn) == 0 {
		t.Error("dsn is empty")
	}
}

func TestGetAllConfig(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	fileMap, err := GetAllConfig(*global.WORKSPACE, nil)
	if err != nil {
		t.Error(err)
	}
	if len(fileMap) == 0 {
		t.Error("fileMap is empty")
	}
}

func TestIsSelefra(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	err := IsSelefra()
	if err != nil {
		t.Error(err)
	}
}

func TestGetModulesByPath(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	modules, err := GetModulesByPath()
	if err != nil {
		t.Error(err)
	}
	if len(modules) == 0 {
		t.Error("modules is empty")
	}
}

func TestGetConfigPath(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	path, err := GetConfigPath()
	if err != nil {
		t.Error(err)
	}
	if len(path) == 0 {
		t.Error("path is empty")
	}
}

func TestGetClientStr(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	clientStr, err := GetClientStr()
	if err != nil {
		t.Error(err)
	}
	if len(clientStr) == 0 {
		t.Error("clientStr is empty")
	}
}

func TestGetModulesStr(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	modulesStr, err := GetModulesStr()
	if err != nil {
		t.Error(err)
	}
	if len(modulesStr) == 0 {
		t.Error("modulesStr is empty")
	}
}

func TestGetRules(t *testing.T) {
	*global.WORKSPACE = "../tests/workspace/offline"
	rules, err := GetRules()
	if err != nil {
		t.Error(err)
	}
	for i := range rules.Rules {
		if len(rules.Rules[i].Name) == 0 {
			t.Error("rule name is empty")
		}
	}
}
