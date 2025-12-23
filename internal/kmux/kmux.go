package kmux

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/kiemlicz/kmux/internal/common"
)

func NewEnvironment(config *common.Config) error {
	name := config.New
	root := config.Root
	kubeconfig := config.Kubeconfig
	location := config.Location
	tpl := config.TmuxinatorConfigTemplate

	if _, exists := config.TmuxinatorConfigs[name]; exists {
		common.Log.Warnf("Environment '%s' already exists. Choose a different name or remove existing.", name)
		return fmt.Errorf("environment '%s' already exists", name)
	}

	data := map[string]string{
		"Name":       name,
		"Root":       root,
		"Kubeconfig": kubeconfig,
	}
	t := template.Must(template.New("config").Parse(tpl))
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

	// Update in-memory config
	config.TmuxinatorConfigs[name] = common.EnvConfig{
		Kubeconfig:       kubeconfig,
		TmuxinatorConfig: location,
	}
	return nil
}

func StartEnvironment(config common.Config) error {
	name := config.Start
	envConfig, exists := config.TmuxinatorConfigs[name]
	if !exists {
		common.Log.Errorf("Environment '%s' does not exist.", name)
		return fmt.Errorf("environment '%s' does not exist", name)
	}
	cmd := exec.Command("tmuxinator", "start", name)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", common.TmuxinatorConfig, envConfig.TmuxinatorConfig))

	// Detach the process from the parent
	// TODO research better this and detach properly
	//cmd.Stdin = nil
	//cmd.Stdout = nil
	//cmd.Stderr = nil
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	Setpgid: true,
	//}

	err := cmd.Start()
	if err != nil {
		common.Log.Errorf("Failed to start environment '%s': %v", name, err)
		return fmt.Errorf("failed to start environment '%s': %v", name, err)
	}

	common.Log.Infof("Started environment '%s'", name)
	return nil
}

func DiscoverEnvironment(config common.Config) error {
	name := config.Discover

	common.Log.Infof("Updated environment '%s'", name)
	return nil
}
