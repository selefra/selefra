package tools

import (
	"context"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"testing"
)

func getProviderAndConfig() (registry.ProviderBinary, *config.SelefraConfig, error) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	ctx := context.Background()
	var cof = new(config.SelefraConfig)
	err := cof.GetConfig()
	if err != nil {
		return registry.ProviderBinary{}, nil, err
	}
	pr := registry.Provider{
		Name:    "aws",
		Version: "latest",
		Source:  "",
	}
	namespace, _, err := utils.Home()
	if err != nil {
		return registry.ProviderBinary{}, nil, err
	}
	provider := registry.NewProviderRegistry(namespace)
	p, err := provider.Download(ctx, pr, true)
	return p, cof, err
}

func TestGetProviders(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	var s = new(config.SelefraConfig)
	err := s.GetConfig()
	if err != nil {
		t.Error(err)
	}
	provider, err := GetProviders(s, "aws")
	t.Log(provider)
	if err != nil {
		t.Error(err)
	}
	if len(provider) == 0 {
		t.Error("Provider is empty")
	}
}

func TestSetProviders(t *testing.T) {
	p, cof, err := getProviderAndConfig()
	s := `
      ##  Optional, Repeated. Add an accounts block for every account you want to assume-role into and fetch data from.
      #accounts:
      #    #     Optional. User identification
      #  - account_name: <UNIQUE ACCOUNT IDENTIFIER>
      #    #    Optional. Named profile in config or credential file from where Selefra should grab credentials
      #    shared_config_profile: < PROFILE_NAME >
      #    #    Optional. Location of shared configuration files
      #    shared_config_files:
      #      - <FILE_PATH>
      #    #   Optional. Location of shared credentials files
      #    shared_credentials_files:
      #      - <FILE_PATH>
      #    #    Optional. Role ARN we want to assume when accessing this account
      #    role_arn: < YOUR_ROLE_ARN >
      #    #    Optional. Named role session to grab specific operation under the assumed role
      #    role_session_name: <SESSION_NAME>
      #    #    Optional. Any outside of the org account id that has additional control
      #    external_id: <ID>
      #    #    Optional. Designated region of servers
      #    default_region: <REGION_CODE>
      #    #    Optional. by default assumes all regions
      #    regions:
      #      - us-east-1
      #      - us-west-2
      ##    The maximum number of times that a request will be retried for failures. Defaults to 10 retry attempts.
      #max_attempts: 10
      ##    The maximum back off delay between attempts. The backoff delays exponentially with a jitter based on the number of attempts. Defaults to 30 seconds.
      #max_backoff: 30
`
	err = SetProviders(s, p, cof)
	if err != nil {
		t.Error(err)
	}
}

func TestSetSelefraProvider(t *testing.T) {
	p, cof, err := getProviderAndConfig()
	if err != nil {
		t.Error(err)
	}
	err = SetSelefraProvider(p, cof)
	if err != nil {
		t.Error(err)
	}
}
