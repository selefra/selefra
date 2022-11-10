package fetch

//
//func TestFetchOnline(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//		return
//	}
//	*global.WORKSPACE = "../../tests/workspace/online"
//	global.SERVER = "dev-api.selefra.io"
//	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
//	ctx := context.Background()
//	var cof = &config.SelefraConfig{}
//	err := cof.GetConfig()
//	for _, p := range cof.Selefra.Providers {
//		err = Fetch(ctx, cof, p)
//		if err != nil {
//			t.Error(err)
//		}
//	}
//}
