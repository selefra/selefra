package registry

type ProviderMetadata struct {
	Name          string   `json:"name" yaml:"name"`
	LatestVersion string   `json:"latest-version" yaml:"latest-version"`
	LatestUpdate  string   `json:"latest-update" yaml:"latest-update"`
	Introduction  string   `json:"introduction" yaml:"introduction"`
	Versions      []string `json:"versions" yaml:"versions"`
}

type ProviderSupplement struct {
	PackageName string    `json:"package-name" yaml:"package-name"`
	Source      string    `json:"source" yaml:"source"`
	Checksums   Checksums `json:"checksums" yaml:"checksums"`
}

type Checksums struct {
	LinuxArm64   string `json:"linux-arm64" yaml:"linux-arm64"`
	LinuxAmd64   string `json:"linux-amd64" yaml:"linux-amd64"`
	WindowsArm64 string `json:"windows-arm64" yaml:"windows-arm64"`
	WindowsAmd64 string `json:"windows-amd64" yaml:"windows-amd64"`
	DarwinArm64  string `json:"darwin-arm64" yaml:"darwin-arm64"`
	DarwinAmd64  string `json:"darwin-amd64" yaml:"darwin-amd64"`
}

type FileMetadata struct {
	DownloadPath string `json:"download_path"`
	FileName     string `json:"file_name"`
	Suffix       string `json:"suffix"`
}

type Header struct {
	Key, Value string
}
