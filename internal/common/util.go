package common

import (
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
)

const (
	ConfigName = "kmrc.yaml"
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

func SetupConfig() (*Config, error) {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.String("new", "", "Create new [name] configuration")
	f.String("discover", "", "Update kubeconfig for [name] with all namespaces")
	f.String("start", "", "Start tmux for [name] configuration")
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	k := koanf.NewWithConf(koanf.Conf{
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
			if err := k.Load(kfile.Provider(file), parser); err != nil {
				log.Fatalf("error loading config: %v", err)
			}
		}
	}
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	var config Config
	err := k.Unmarshal("", &config)
	if err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return &config, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
