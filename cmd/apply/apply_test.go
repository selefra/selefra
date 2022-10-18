package apply

import (
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestGetRules(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	modules, err := config.GetModulesByPath()
	if err != nil {
		t.Error(err)
	}
	rules := RunRulesWithoutModule()
	rulesTwo := CreateRulesByModule(modules)
	if len(*rules) == 0 && len(rulesTwo) == 0 {
		t.Error("rules is empty")
	}
}
