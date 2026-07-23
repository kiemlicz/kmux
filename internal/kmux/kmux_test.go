package kmux

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kiemlicz/kmux/internal/common"
)

func TestKubeconfigIsLazilyResolved(t *testing.T) {
	common.SetupLog("debug")
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "myenv.yaml")
	if err := os.WriteFile(cfgPath, []byte("KUBECONFIG=/path/to/kubeconfig\n"), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	km := NewKmux(common.Config{
		TmuxinatorConfigPaths: []string{dir},
	})

	env, exists := km.environments["myenv"]
	if !exists {
		t.Fatalf("expected environment 'myenv' to be discovered")
	}

	// Field must not be resolved yet - NewKmux/addEnvironment must not have read the file's contents.
	if env.kubeconfigVal != "" {
		t.Fatalf("expected kubeconfig to be unresolved before first access, got %q", env.kubeconfigVal)
	}

	// First access triggers resolution.
	got := env.Kubeconfig()
	if got != "/path/to/kubeconfig" {
		t.Fatalf("expected kubeconfig '/path/to/kubeconfig', got %q", got)
	}

	// Mutate the underlying file and confirm the resolved value is memoized (not re-read).
	if err := os.WriteFile(cfgPath, []byte("KUBECONFIG=/changed\n"), 0644); err != nil {
		t.Fatalf("failed to rewrite test config: %v", err)
	}
	if got := env.Kubeconfig(); got != "/path/to/kubeconfig" {
		t.Fatalf("expected memoized kubeconfig '/path/to/kubeconfig', got %q", got)
	}

	// Memoization must persist through the map (pointer semantics), not just on a copied value.
	if got := km.environments["myenv"].Kubeconfig(); got != "/path/to/kubeconfig" {
		t.Fatalf("expected memoized kubeconfig via map lookup, got %q", got)
	}
}

func TestExplicitKubeconfigSkipsExtraction(t *testing.T) {
	common.SetupLog("debug")
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "myenv.yaml")
	// File has no KUBECONFIG line at all - if extraction ran, it would return "".
	if err := os.WriteFile(cfgPath, []byte("root: /somewhere\n"), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	environments := make(map[string]*KmuxEnvironment)
	addEnvironment(environments, cfgPath, "/explicit/kubeconfig")

	env := environments["myenv"]
	if got := env.Kubeconfig(); got != "/explicit/kubeconfig" {
		t.Fatalf("expected explicit kubeconfig '/explicit/kubeconfig', got %q", got)
	}
}
