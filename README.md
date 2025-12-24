# kmux
Personal Kubeconfig manager tool  
Doesn't keep global kubeconfig prone to accidental modifications

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
