package kmux

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/kiemlicz/kmux/internal/common"
)

//go:embed completions/kmux.zsh.tpl
var zshCompletionsTemplate string

func CompletionsZsh(c *common.Config) (string, error) {
	tmuxinatorConfigPaths := c.TmuxinatorConfigPaths
	projectFiles := strings.Join(tmuxinatorConfigPaths, " ")
	allCommands := strings.Join(common.AllCommands, " ")
	runnableCommands := strings.Join(common.AllCommands, "|")
	runnableCommands = strings.ReplaceAll(strings.ReplaceAll(runnableCommands, common.OptionCompletions+"|", ""), "|"+common.OptionCompletions, "")

	data := map[string]string{
		"AllCommands":      allCommands,
		"RunnableCommands": runnableCommands,
		"Projects":         projectFiles,
	}
	t := template.Must(template.New("completions").Parse(zshCompletionsTemplate))
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		common.Log.Errorf("Error executing template: %v", err)
		return "", err
	}
	return buf.String(), nil
}
