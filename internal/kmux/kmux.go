package kmux

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"text/template"

	"github.com/kiemlicz/kmux/internal/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	KubeconfigRegex = regexp.MustCompile(`KUBECONFIG=([^\s]+)`)
)

type Kmux struct {
	environments       map[string]KmuxEnvironment // map of environment name to its Tmuxinator config path
	tmuxinatorTemplate string
}
type KmuxEnvironment struct {
	name     string
	fullpath string
}

func NewKmux(c common.Config) *Kmux {
	tmuxinatorConfigPaths := c.TmuxinatorConfigPaths
	tpl := c.TmuxinatorConfigTemplate
	environments := make(map[string]KmuxEnvironment)

	for _, path := range tmuxinatorConfigPaths {
		common.Log.Debugf("Loading environments from: %s", path)
		files, err := os.ReadDir(path)
		if err != nil {
			common.Log.Errorf("Failed to read directory %s: %v", path, err)
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				filename := file.Name()
				if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
					common.Log.Debugf("Found YAML file: %s", filename)
					basename := strings.TrimSuffix(strings.TrimSuffix(filename, ".yaml"), ".yml")
					environments[basename] = KmuxEnvironment{
						fullpath: filepath.Join(path, filename),
						name:     basename,
					}
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
	name := ops.OperationArgs
	root := ops.Root
	kubeconfig := ops.Kubeconfig
	location := ops.Location

	if _, exists := km.environments[name]; exists {
		common.Log.Warnf("Environment '%s' already exists. Choose a different name or remove existing.", name)
		return fmt.Errorf("environment '%s' already exists", name)
	}

	common.Log.Infof("Creating environment '%s', TMUXINATOR_CONFIG=%s, KUBECONFIG=%s, root=%s", name, location, kubeconfig, root)

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
	name := ops.OperationArgs
	kmuxEnv, exists := km.environments[name]
	var err error

	if !exists {
		common.Log.Errorf("Environment '%s' does not exist.", name)
		return fmt.Errorf("environment '%s' does not exist", name)
	}

	common.Log.Infof("Starting environment '%s'", name)

	tmuxinatorConfig := filepath.Dir(kmuxEnv.fullpath)
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
	name := ops.OperationArgs
	kmuxEnv, exists := km.environments[name]
	if !exists {
		return fmt.Errorf("environment '%s' does not exist", name)
	}

	fullpath := kmuxEnv.fullpath
	common.Log.Infof("Discovering environment '%s'", fullpath)
	content, err := os.ReadFile(fullpath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	matches := KubeconfigRegex.FindStringSubmatch(string(content))
	var kubeconfig string
	if len(matches) > 1 {
		kubeconfig = matches[1]
		common.Log.Debugf("Found KUBECONFIG: %s", kubeconfig)
	} else {
		return fmt.Errorf("KUBECONFIG not found in environment file")
	}
	// Load kubeconfig
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}
	kctx, err := kubeconfigCtx(config)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig context: %v", err)
	}
	namespaces, err := listNamespaces(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %v", err)
	}

	// Generate new contexts for each namespace
	newContexts := make(map[string]*api.Context)
	for _, ns := range namespaces {
		contextName := ns + "-" + name
		newContexts[contextName] = &api.Context{
			Cluster:   kctx.Cluster,
			Namespace: ns,
			AuthInfo:  kctx.AuthInfo,
		}
	}
	// Replace contexts in kubeconfig
	config.Contexts = newContexts
	// Save updated kubeconfig
	err = clientcmd.WriteToFile(*config, kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to save kubeconfig: %v", err)
	}

	common.Log.Infof("Updated environment '%s'", name)
	return nil
}

func listNamespaces(kubeconfig string) ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaceNames []string
	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}

	return namespaceNames, nil
}

// kubeconfigCtx extracts the active context from the kubeconfig file
func kubeconfigCtx(config *api.Config) (*api.Context, error) {
	if config.CurrentContext == "" {
		return nil, fmt.Errorf("no current context set in kubeconfig")
	}

	kcontext, exists := config.Contexts[config.CurrentContext]
	if !exists {
		return nil, fmt.Errorf("current context '%s' not found in kubeconfig", config.CurrentContext)
	}

	if kcontext.AuthInfo == "" {
		return nil, fmt.Errorf("no user set in current context")
	}
	if kcontext.Cluster == "" {
		return nil, fmt.Errorf("no cluster set in current context")
	}

	return kcontext, nil
}

func envAddTmuxinatorConfig(tmuxinatorConfig string) []string {
	return append(os.Environ(), fmt.Sprintf("%s=%s", common.TmuxinatorConfig, tmuxinatorConfig))
}
