package apply

import (
	"context"
	"github.com/google/uuid"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/ui/client"
	"testing"
)

func TestGetRules(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	modules, err := config.GetModulesByPath()
	if err != nil {
		t.Error(err)
		return
	}
	rules := RunRulesWithoutModule()
	rulesTwo := CreateRulesByModule(modules)
	if len(*rules) == 0 && len(rulesTwo) == 0 {
		t.Error("rules is empty")
	}

	var useRules []config.Rule
	if len(rulesTwo) != 0 {
		useRules = rulesTwo
	}
	if len(*rules) != 0 {
		useRules = *rules
	}
	ctx := context.Background()
	uid, _ := uuid.NewUUID()
	s := config.SelefraConfig{}
	err = s.GetConfig()
	c, e := client.CreateClientFromConfig(ctx, &s.Selefra, uid)
	if err != nil {
		t.Error(e)
	}
	err = RunRules(ctx, c, "", useRules)
	if err != nil {
		t.Error(err)
	}
}

func TestApply(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
