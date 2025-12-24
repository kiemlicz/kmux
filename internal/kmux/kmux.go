package kmux

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/template"

	"github.com/kiemlicz/kmux/internal/common"
)

type Kmux struct {
	environments       map[string]string // map of environment name to its Tmuxinator config path
	tmuxinatorTemplate string
}

func NewKmux(c common.Config) *Kmux {
	tmuxinatorConfigPaths := c.TmuxinatorConfigPaths
	tpl := c.TmuxinatorConfigTemplate
	environments := make(map[string]string)

	for _, path := range tmuxinatorConfigPaths {
		common.Log.Debugf("Loading environments from: %s", path)
		files, err := os.ReadDir(path)
		if err != nil {
			common.Log.Errorf("Failed to read directory %s: %v", path, err)
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				name := file.Name()
				if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
					common.Log.Debugf("Found YAML file: %s", name)
					basename := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
					environments[basename] = path
				}
			}
		}
	}

	return &Kmux{
		environments:       environments,
		tmuxinatorTemplate: tpl,
	}
}

func (km *Kmux) NewEnvironment(ops *common.Operations) error {
	name := ops.New
	root := ops.Root
	kubeconfig := ops.Kubeconfig
	location := ops.Location

	if _, exists := km.environments[name]; exists {
		common.Log.Warnf("Environment '%s' already exists. Choose a different name or remove existing.", name)
		return fmt.Errorf("environment '%s' already exists", name)
	}

	data := map[string]string{
		"Name":       name,
		"Root":       root,
		"Kubeconfig": kubeconfig,
	}
	t := template.Must(template.New("config").Parse(km.tmuxinatorTemplate))
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		common.Log.Errorf("Error executing template: %v", err)
		return err
	}

	err = common.DumpYamlToFile(buf, location, name)
	if err != nil {
		return err
	}
	// Not creating KUBECONFIG as must be either populated using external tool or manually

	return nil
}

func (km *Kmux) StartEnvironment(ops common.Operations) error {
	name := ops.Start
	tmuxinatorConfig, exists := km.environments[name]
	var err error

	if !exists {
		common.Log.Errorf("Environment '%s' does not exist.", name)
		return fmt.Errorf("environment '%s' does not exist", name)
	}

	if ops.Bg {
		err = km.spawnTmuxinatorBg(name, tmuxinatorConfig)
	} else {
		err = km.spawnTmuxinatorFg(name, tmuxinatorConfig)
	}
	if err != nil {
		common.Log.Errorf("Failed to start environment '%s': %v", name, err)
		return fmt.Errorf("failed to start environment '%s': %v", name, err)
	}
	common.Log.Infof("Started environment '%s'", name)
	return nil
}

func (km *Kmux) spawnTmuxinatorFg(name, tmuxinatorConfig string) error {
	// Find the full path to tmuxinator
	tmuxinatorPath, err := exec.LookPath(common.Tmuxinator)
	if err != nil {
		common.Log.Errorf("Failed to find tmuxinator: %v", err)
		return fmt.Errorf("failed to find tmuxinator: %v", err)
	}

	// Prepare arguments
	args := []string{common.Tmuxinator, "start", name}
	env := envAddTmuxinatorConfig(tmuxinatorConfig)

	// Replace current process with tmuxinator, mind that when `go run` the top-level go will remain
	return syscall.Exec(tmuxinatorPath, args, env)
}

// spawnTmuxinatorBg starts a tmuxinator session in the background
// doesn't immediately attach to it
func (km *Kmux) spawnTmuxinatorBg(name, tmuxinatorConfig string) error {
	cmd := exec.Command(common.Tmuxinator, "start", name)
	cmd.Env = envAddTmuxinatorConfig(tmuxinatorConfig)

	// Detach the process from the parent
	// TODO research better this and detach properly
	//cmd.Stdin = nil
	//cmd.Stdout = nil
	//cmd.Stderr = nil
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	Setpgid: true,
	//}
	return cmd.Start()
}

func (km *Kmux) DiscoverEnvironment(ops common.Operations) error {
	name := ops.Discover

	common.Log.Infof("Updated environment '%s'", name)
	return nil
}

func envAddTmuxinatorConfig(tmuxinatorConfig string) []string {
	return append(os.Environ(), fmt.Sprintf("%s=%s", common.TmuxinatorConfig, tmuxinatorConfig))
}
