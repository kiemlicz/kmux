# kmux
Start tmux session with `KUBECONFIG` fetched from different locations

## Usage
Prepare `$XDG_CONFIG_HOME/kmrc`
```
gcp:
  project: "project-name"
kubeconfig_dir: "~/.kube/" # or some other dir
```
Run
`km --gke --name cloud --location europe-west4`  
