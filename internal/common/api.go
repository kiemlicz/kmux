package common

type Config struct {
	Log struct {
		Level string `koanf:"level"`
	} `koanf:"log"`

	TmuxinatorConfigs        map[string]EnvConfig `koanf:"environments"`        // env name to tmuxinator config path
	TmuxinatorConfigTemplate string               `koanf:"tmuxinator_template"` // template for new tmuxinator configs

	KmuxConfigFile string `koanf:"kmux_config_file"` // when multiple configs found and not set in file, this will be set to top-most one, required to persist changes of configs

	//CLI options
	New      string `koanf:"new"`
	Discover string `koanf:"discover"`
	Start    string `koanf:"start"`

	Location   string `koanf:"location"`
	Kubeconfig string `koanf:"kubeconfig"`
	Root       string `koanf:"root"`
}

type EnvConfig struct {
	Kubeconfig       string `koanf:"kubeconfig"`
	TmuxinatorConfig string `koanf:"tmuxinator_config"`
}
