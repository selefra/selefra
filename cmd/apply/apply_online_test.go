package apply

//
//func TestGetRulesOnline(t *testing.T) {
//	global.SERVER = "dev-api.selefra.io"
//	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
//	*global.WORKSPACE = "../../tests/workspace/online"
//	modules, err := config.GetModulesByPath()
//	if err != nil {
//		t.Error(err)
//	}
//	rules := RunRulesWithoutModule()
//	rulesTwo := CreateRulesByModule(modules)
//	if len(*rules) == 0 && len(rulesTwo) == 0 {
//		t.Error("rules is empty")
//	}
//
//	var useRules []config.Rule
//	if len(rulesTwo) != 0 {
//		useRules = rulesTwo
//	}
//	if len(*rules) != 0 {
//		useRules = *rules
//	}
//	ctx := context.Background()
//	uid, _ := uuid.NewUUID()
//	s := config.SelefraConfig{}
//	err = s.GetConfig()
//	c, e := client.CreateClientFromConfig(ctx, &s.Selefra, uid)
//	if err != nil {
//		t.Error(e)
//	}
//	err = RunRules(ctx, c, "", useRules)
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func TestApplyOnLine(t *testing.T) {
//	global.SERVER = "dev-api.selefra.io"
//	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
//	*global.WORKSPACE = "../../tests/workspace/online"
//	err := Apply(context.Background())
//	if err != nil {
//		t.Error(err)
//	}
//}
