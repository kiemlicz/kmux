package kmux

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/kiemlicz/kmux/internal/common"
)

// extractCaseBlock returns the body of a `<label>) ... ;;` case block from a
// rendered completion script, so assertions can be scoped to one subcommand
// instead of matching anywhere in the whole script.
func extractCaseBlock(t *testing.T, script, label string) string {
	t.Helper()
	marker := label + ")"
	start := strings.Index(script, marker)
	if start == -1 {
		t.Fatalf("case block %q not found in generated script:\n%s", label, script)
	}
	rest := script[start+len(marker):]
	end := strings.Index(rest, ";;")
	if end == -1 {
		t.Fatalf("case block %q has no terminating ';;' in generated script:\n%s", label, script)
	}
	return rest[:end]
}

func TestCompletionsZshNewFlags(t *testing.T) {
	common.SetupLog("debug")
	cfg := &common.Config{TmuxinatorConfigPaths: []string{"/tmuxinator_config_1", "/tmuxinator_config_2"}}

	out, err := CompletionsZsh(cfg)
	if err != nil {
		t.Fatalf("CompletionsZsh returned error: %v", err)
	}

	newBlock := extractCaseBlock(t, out, "new")
	for _, want := range []string{
		"--location=[tmuxinator config directory]:location:((/tmuxinator_config_1 /tmuxinator_config_2))",
		"--root=[project root directory]:root:_path_files -/",
		"--kubeconfig=[kubeconfig file]:kubeconfig:_files",
		`2:name:_path_files -W "(/tmuxinator_config_1 /tmuxinator_config_2)" -g "*.(yml|yaml)(:r)"`,
	} {
		if !strings.Contains(newBlock, want) {
			t.Errorf("expected 'new' case block to contain %q, got:\n%s", want, newBlock)
		}
	}
}

func TestCompletionsZshDiscoverOptionalName(t *testing.T) {
	common.SetupLog("debug")
	cfg := &common.Config{TmuxinatorConfigPaths: []string{"/tmuxinator_config_1"}}

	out, err := CompletionsZsh(cfg)
	if err != nil {
		t.Fatalf("CompletionsZsh returned error: %v", err)
	}

	discoverBlock := extractCaseBlock(t, out, "discover")
	if strings.Contains(discoverBlock, "--") {
		t.Errorf("expected 'discover' case block to have no flag specs, got:\n%s", discoverBlock)
	}
	if !strings.Contains(discoverBlock, "2::name") {
		t.Errorf("expected 'discover' case block to use optional positional '2::name', got:\n%s", discoverBlock)
	}
	if !strings.Contains(discoverBlock, `-g "*.(yml|yaml)(:r)"`) {
		t.Errorf("expected 'discover' case block to use extension-restricted glob, got:\n%s", discoverBlock)
	}
}

func TestCompletionsZshStartBgFlag(t *testing.T) {
	common.SetupLog("debug")
	cfg := &common.Config{TmuxinatorConfigPaths: []string{"/tmuxinator_config_1"}}

	out, err := CompletionsZsh(cfg)
	if err != nil {
		t.Fatalf("CompletionsZsh returned error: %v", err)
	}

	startBlock := extractCaseBlock(t, out, "start")
	if !strings.Contains(startBlock, "--bg[spawn session in background, do not attach]") {
		t.Errorf("expected 'start' case block to contain --bg flag spec, got:\n%s", startBlock)
	}
	if strings.Contains(startBlock, "--location") || strings.Contains(startBlock, "--kubeconfig") {
		t.Errorf("expected 'start' case block to not contain --location/--kubeconfig, got:\n%s", startBlock)
	}
	if !strings.Contains(startBlock, "2:name:") || strings.Contains(startBlock, "2::name") {
		t.Errorf("expected 'start' case block to use mandatory positional '2:name', got:\n%s", startBlock)
	}
}

func TestCompletionsZshSyntaxIsValid(t *testing.T) {
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh not available")
	}
	common.SetupLog("debug")
	cfg := &common.Config{TmuxinatorConfigPaths: []string{"/tmuxinator_config_1"}}

	out, err := CompletionsZsh(cfg)
	if err != nil {
		t.Fatalf("CompletionsZsh returned error: %v", err)
	}

	cmd := exec.Command("zsh", "-n")
	cmd.Stdin = strings.NewReader(out)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated completion script failed zsh -n syntax check: %v\n%s", err, output)
	}
}
