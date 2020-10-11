# kmux
Start tmux session with `KUBECONFIG` fetched from different locations

## Usage
Prepare `$XDG_CONFIG_HOME/kmrc`
```
gcp:
  project: "project-name"
kubeconfig_dir: "~/.kube/" # or some other dir
```
### Run
GKE:  
`km --gke --name cloud --location europe-west4`  
Google Drive:  
`km --gdrive /provide/path/kubeconfig --name localname`