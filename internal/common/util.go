package common

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	kyaml "github.com/knadh/koanf/parsers/yaml"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const (
	ConfigName       = "kmrc.yaml"
	TmuxinatorConfig = "TMUXINATOR_CONFIG"
	Tmuxinator       = "tmuxinator"
)

var Log *logrus.Logger

func SetupLog(logLevel string) {
	Log = logrus.New()
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		Log.Warnf("Invalid Log level in config: %s. Using 'info'.", logLevel)
		level = logrus.InfoLevel
	}

	Log.SetLevel(level)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func SetupConfig() (*Config, *Operations, error) {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.String("new", "", "Create new [name] configuration")
	f.String("discover", "", "Update kubeconfig for [name] with all namespaces")
	f.String("start", "", "Start tmux for [name] configuration")

	f.String("location", "", "TMUXINATOR_CONFIG - tmuxinator config directory")
	f.String("root", "", "Tmuxinator's config root directory (dir auto changed to)")
	f.String("kubeconfig", "", "KUBECONFIG for new environment")
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	// Get positional arguments
	args := f.Args()
	if len(args) < 1 {
		f.Usage()
		log.Fatal("error: must provide command (new, discover, or start)")
	}
	if len(args) < 2 {
		f.Usage()
		log.Fatal("error: must provide name argument")
	}

	// Parse Operations from CLI flags before setting positional arguments
	var ops Operations
	command := args[0]
	name := args[1]
	opsKoanf := koanf.New(".")
	if err := opsKoanf.Load(posflag.Provider(f, ".", opsKoanf), nil); err != nil {
		log.Fatalf("error loading operations from flags: %v", err)
	}
	if err := opsKoanf.Unmarshal("", &ops); err != nil {
		log.Fatalf("error unmarshalling operations: %v", err)
	}
	// Create operations based on positional arguments
	switch command {
	case "new":
		ops.New = name
	case "discover":
		ops.Discover = name
	case "start":
		ops.Start = name
	default:
		log.Fatalf("error: unknown command '%s'. Must be one of: new, discover, start", command)
	}

	fileConfigKoanf := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})
	parser := kyaml.Parser()

	var files []string
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		files = append(files, filepath.Join(xdgConfigHome, "kmux", ConfigName))
	} else {
		files = append(files, filepath.Join(os.Getenv("HOME"), ".config", "kmux", ConfigName))
	}
	files = append(files, filepath.Join(".local", ConfigName)) // for local dev

	for _, file := range files {
		if fileExists(file) {
			if err := fileConfigKoanf.Load(kfile.Provider(file), parser); err != nil {
				log.Fatalf("error loading config: %v", err)
			}
		}
	}
	if err := fileConfigKoanf.Load(posflag.Provider(f, ".", fileConfigKoanf), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	var config Config
	err := fileConfigKoanf.Unmarshal("", &config)
	if err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	setupDefaults(&config, &ops)

	return &config, &ops, nil
}

func setupDefaults(c *Config, o *Operations) {
	if c.TmuxinatorConfigPaths == nil {
		c.TmuxinatorConfigPaths = []string{defaultTmuxinatorConfigs()}
	}
	if o.Root == "" {
		o.Root = "~/"
	}
	if o.Location == "" {
		o.Location = c.TmuxinatorConfigPaths[0]
	}
	if o.Kubeconfig == "" {
		o.Kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
}

func DumpYamlToFile(buf bytes.Buffer, dir string, filename string) error {
	// Parse as YAML to ensure valid format
	var tmuxinatorFile map[string]any
	err := yaml.Unmarshal(buf.Bytes(), &tmuxinatorFile)
	if err != nil {
		Log.Errorf("Error unmarshalling template to YAML after templating: %v", err)
		return err
	}
	// Write YAML to file
	file, err := os.Create(filepath.Join(dir, filename+".yml"))
	if err != nil {
		Log.Errorf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	err = encoder.Encode(tmuxinatorFile)
	if err != nil {
		Log.Errorf("Error encoding YAML to file: %v", err)
		return err
	}
	return nil
}

// defaultTmuxinatorConfigs stops on first found
// https://github.com/tmuxinator/tmuxinator?tab=readme-ov-file#project-configuration-location
func defaultTmuxinatorConfigs() string {
	tmuxinatorConfigEnv := os.Getenv("TMUXINATOR_CONFIG")
	xdgConfigHomeEnv := os.Getenv("XDG_CONFIG_HOME")

	if tmuxinatorConfigEnv != "" {
		return tmuxinatorConfigEnv
	}
	if xdgConfigHomeEnv != "" {
		return filepath.Join(xdgConfigHomeEnv, "tmuxinator")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".config", "tmuxinator") // $HOME/.tmuxinator seems deprecated
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
