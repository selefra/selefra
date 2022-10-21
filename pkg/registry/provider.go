package registry

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/songzhibin97/gkit/ternary"
	yaml "gopkg.in/yaml.v3"

	"github.com/selefra/selefra/pkg/internal/getter"
)

const (
	row = "https://raw.githubusercontent.com/selefra/registry"
)

var (
	suffix = ternary.ReturnString(runtime.GOOS == "windows", ".exe", "")
)

type ProviderBinary struct {
	Provider
	Filepath string
}

type Provider struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
	Path    string `json:"path"`
}

func (p *Provider) String() string {
	return fmt.Sprintf("%s@%s", p.Name, p.Version)
}

type RegisterProvider interface {
	CheckUpdate(ctx context.Context, binary ProviderBinary) (ProviderBinary, error)
	Download(ctx context.Context, provider Provider, skipVerify bool) (ProviderBinary, error)
	DeleteProvider(binary ProviderBinary) error
}

type Providers struct {
	providers map[string]Provider
	sync.RWMutex
}

func (p *Providers) Set(provider Provider) {
	p.Lock()
	defer p.Unlock()
	if p.providers == nil {
		p.providers = make(map[string]Provider)
	}
	p.providers[provider.String()] = provider
}

func (p *Providers) get(name string) (Provider, bool) {
	if p.providers == nil {
		return Provider{}, false
	}
	provider, ok := p.providers[name]
	return provider, ok
}

func (p *Providers) GetMany(names ...string) []Provider {
	p.RLock()
	defer p.RUnlock()
	providers := make([]Provider, 0, len(names))
	for _, n := range names {
		p, ok := p.get(n)
		if !ok {
			continue
		}
		providers = append(providers, p)
	}
	return providers
}

func (p *Providers) Delete(name string) {
	p.Lock()
	defer p.Unlock()
	delete(p.providers, name)
}

func NewProviders() *Providers {
	return &Providers{
		providers: make(map[string]Provider),
	}
}

// ===========================================================

func request(ctx context.Context, method string, _url string, body []byte, headers ...Header) ([]byte, error) {
	client := &http.Client{}

	sBody := strings.NewReader(string(body))
	request, err := http.NewRequestWithContext(ctx, method, _url, sBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	for _, header := range headers {
		request.Header.Add(header.Key, header.Value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("code not equal 200")
	}
	rByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read body err :" + err.Error())
	}
	return rByte, err
}

// =====================================================================================================================

type provider struct {
	namespace string
}

func (p *provider) CheckUpdate(ctx context.Context, binary ProviderBinary) (ProviderBinary, error) {
	metadata, err := p.getProviderMetadata(ctx, &binary.Provider)
	if err != nil {
		return ProviderBinary{}, err
	}
	if binary.Provider.Version == metadata.LatestVersion {
		return ProviderBinary{}, nil
	}
	if binary.Provider.Version != "" && binary.Provider.Version != "latest" {
		_ = p.deleteProviderBinary(binary)
	}

	binary.Provider.Version = metadata.LatestVersion
	return p.download(ctx, binary.Provider, true)
}

func (p *provider) DeleteProvider(binary ProviderBinary) error {
	return p.deleteProviderBinary(binary)
}

func (p *provider) deleteProviderBinary(binary ProviderBinary) error {
	if _, err := os.Stat(binary.Filepath); err != nil {
		return err
	}

	return os.RemoveAll(filepath.Dir(binary.Filepath))
}

func (p *provider) getSupplement(ctx context.Context, provider *Provider) (ProviderSupplement, error) {
	var supplement ProviderSupplement
	_url := row + "/main/provider/" + provider.Name + "/" + provider.Version + "/" + "supplement.yaml"

	body, err := request(ctx, "GET", _url, nil)
	if err != nil {
		return supplement, err
	}

	err = yaml.Unmarshal(body, &supplement)
	return supplement, err
	//downloadUrl := supplement.Supplement.Source + "/releases/download/" + provider.Version + "/" + provider.Name + "_" + runtime.GOOS + "_" + runtime.GOARCH + ".tar.gz"
}

func (p *provider) fillVersion(ctx context.Context, provider *Provider, skipVerify bool) error {
	if provider.Version != "" && provider.Version != "latest" && skipVerify {
		return nil
	}

	metadata, err := p.getProviderMetadata(ctx, provider)
	if err != nil {
		return err
	}
	if provider.Version != "" && provider.Version != "latest" {
		for _, version := range metadata.Versions {
			if provider.Version != version {
				continue
			}
			return nil
		}
		return errors.New("version not found")
	}
	provider.Version = metadata.LatestVersion
	return nil
}

func (p *provider) getProviderMetadata(ctx context.Context, provider *Provider) (ProviderMetadata, error) {
	var metadata ProviderMetadata
	_url := row + "/main/provider/" + provider.Name + "/metadata.yaml"

	body, err := request(ctx, "GET", _url, nil)
	if err != nil {
		return metadata, err
	}

	err = yaml.Unmarshal(body, &metadata)
	return metadata, err
}

func (p *provider) verify(supplement ProviderSupplement, skipVerify bool) bool {
	if skipVerify {
		return true
	}
	// TODO Not realized
	return true
}

func (p *provider) Download(ctx context.Context, provider Provider, skipVerify bool) (ProviderBinary, error) {

	var pp ProviderBinary

	if provider.Version == "" || provider.Version == "latest" {
		err := p.fillVersion(ctx, &provider, skipVerify)
		if err != nil {
			return pp, err
		}
	}
	return p.download(ctx, provider, skipVerify)
}

func (p *provider) download(ctx context.Context, provider Provider, skipVerify bool) (ProviderBinary, error) {
	pp := ProviderBinary{Provider: provider}
	if provider.Path != "" {
		pp.Filepath = provider.Path
		return pp, nil
	}
	abs, err := filepath.Abs(p.namespace)
	if err != nil {
		return pp, err
	}
	supplement, err := p.getSupplement(ctx, &provider)
	if err != nil {
		return pp, err
	}

	if !p.verify(supplement, skipVerify) {
		return pp, errors.New("verify failed")
	}

	pp.Filepath = abs + "/download/providers/" + provider.Name + "_" + provider.Version + "/"
	fileName := supplement.PackageName + "_" + strings.Replace(provider.Version, "v", "", 1) + "_" + runtime.GOOS + "_" + runtime.GOARCH
	_, err = os.Stat(filepath.Join(pp.Filepath, supplement.PackageName, suffix))
	if err == nil {
		pp.Filepath = filepath.Join(pp.Filepath, supplement.PackageName, suffix)
		return pp, nil
	}

	getUrl := supplement.Source + "/releases/download/" + provider.Version + "/" + fileName + ".tar.gz"

	if !skipVerify {
		checksum, err := selectChecksums(supplement)
		if err != nil {
			return pp, err
		}
		getUrl = getUrl + "?checksum=sha256:" + checksum
	}

	err = getter.Get(ctx, pp.Filepath, getUrl)
	if err != nil {
		return ProviderBinary{}, err
	}
	pp.Filepath = filepath.Join(pp.Filepath, supplement.PackageName, suffix)
	return pp, nil
}

func NewProviderRegistry(namespace string) RegisterProvider {
	return &provider{namespace: namespace}
}

func selectChecksums(supplement ProviderSupplement) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return supplement.Checksums.DarwinAmd64, nil
		case "arm64":
			return supplement.Checksums.DarwinArm64, nil
		default:
			return "", errors.New("unsupported arch")
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return supplement.Checksums.WindowsAmd64, nil
		case "arm64":
			return supplement.Checksums.WindowsArm64, nil
		default:
			return "", errors.New("unsupported arch")
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return supplement.Checksums.LinuxAmd64, nil
		case "arm64":
			return supplement.Checksums.LinuxArm64, nil
		default:
			return "", errors.New("unsupported arch")
		}
	default:
		return "", errors.New("unsupported os")
	}
}
