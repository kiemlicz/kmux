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

	glog "gopkg.in/op/go-logging.v1"
)

const (
	CONFIG_NAME = "kmrc.yaml"
)

var Log *logrus.Logger

func Setup(logLevel string) {
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

	lvl, _ := glog.LogLevel(logLevel)
	if lvl == glog.DEBUG {
		lvl = glog.INFO // map debug to info as yq-lib debug is too verbose
	}
	glog.SetLevel(lvl, "yq-lib")
}

func SetupConfig() (*Config, error) {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})
	parser := kyaml.Parser()

	files := []string{CONFIG_NAME}
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		files = append(files, filepath.Join(xdgConfigHome, "kmux", CONFIG_NAME))
	} else {
		files = append(files, filepath.Join(os.Getenv("HOME"), ".config", "kmux", CONFIG_NAME))
	}

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
