# kmux
Start tmux session with `KUBECONFIG` fetched from different locations

## Usage
Run
`bin/km --gke --name cloud --location europe-west4`  
Then in session:
`export KUBECONFIG=$(tmux showenv | awk -F= '/KUBECONFIG/ {print$2}')`