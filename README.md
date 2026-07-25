# kmux
Personal Kubeconfig&tmuxinator manager tool  
Doesn't keep global kubeconfig prone to accidental modifications.  
Manages tmuxinator sessions per separate `KUBECONFIG`s.  

## Usage
Prepare `$XDG_CONFIG_HOME/kmux/kmrc.yaml` or `~/.config/kmux/kmrc.yaml` config file:

```
log:
  level: warn

environments: # leave empty to use defaults
  - /tmuxinator_config_1
  - /tmuxinator_config_2

tmuxinatorTemplate: |
  name: {{ .Name }}
  root: {{ .Root }}
  pre_window: export KUBECONFIG={{ .Kubeconfig }} && tmux setenv KUBECONFIG {{ .Kubeconfig }}
  windows:
    - main:
        layout: main-horizontal
        # Synchronize all panes of this window, can be enabled before or after the pane commands run.
        # 'before' represents legacy functionality and will be deprecated in a future release, in favour of 'after'
        # synchronize: after
        panes:
          - main: []
          - secondary:
            - kgpw -o wide

```

## Shell completions (zsh)
Generate the completion function and append it to your zshrc:
```
km completions zsh >> ~/.zshrc
```

Ensure `autoload -Uz compinit && compinit` runs earlier in your `.zshrc` (automatically done by oh-my-zsh/prezto).

## Start session
Start a new tmux session (configured in kmrc.yaml):
```
km start <session-name>
```

Run in background without attaching:
```
km start --bg <session-name>
```
