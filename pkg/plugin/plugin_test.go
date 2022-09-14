package plugin

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"testing"
)

func Test_newManagedPlugin(t *testing.T) {
	plug, err := newManagedPlugin("/Users/songzhibin/go/src/aqgs/foo-provider-test/foo-provider-test", "test", "1.0.0", "test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(plug.Provider().SetProviderConfiguration(context.Background(), &shard.SetProviderConfigurationRequest{}))
}
