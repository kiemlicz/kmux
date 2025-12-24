package common

type Config struct {
	Log struct {
		Level string `koanf:"level"`
	} `koanf:"log"`

	TmuxinatorConfigPaths    []string `koanf:"environments"`       // tmuxinator config paths
	TmuxinatorConfigTemplate string   `koanf:"tmuxinatorTemplate"` // template for new tmuxinator configs
}

// todo instead of --start, --discover, --new flags, use positional args
type Operations struct {
	New      string `koanf:"new"`
	Discover string `koanf:"discover"`
	Start    string `koanf:"start"`

	Location   string `koanf:"location"`
	Kubeconfig string `koanf:"kubeconfig"`
	Root       string `koanf:"root"`

	Bg bool `koanf:"bg"` //spawn in background
}
