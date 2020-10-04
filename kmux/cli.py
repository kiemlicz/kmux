import argparse
import yaml
import os
from os.path import expanduser


def _merge(source, destination):
    for k,v in source.items():
        if isinstance(v, dict):
            n = destination.setdefault(k, {})
            _merge(v, n)
        # consider list append
        else:
            destination[k] = v
    return destination


parser = argparse.ArgumentParser(description='Spawn TMUX session with KUBECONFIG fetched automatically')
parser.add_argument('--gke', help="Fech KUBECONFIG from GKE", required=False, action='store_true')
parser.add_argument('--eks', help="Fech KUBECONFIG from EKS", required=False, action='store_true')
parser.add_argument('--name', help="Cluster name", required=False)
parser.add_argument('--location', help="Cluster location", required=False)
parser.add_argument('--project', help="Cloud project name (e.g. GCP project)", required=False)
parser.add_argument('--config', help="Configuration file", default=os.path.join(expanduser("~"), ".kmrc"))
args = parser.parse_args()

with open(args.config, 'r') as f:
    config = yaml.safe_load(f)


def gcp_get_or_default(key: str):
    return vars(args)[key] if vars(args)[key] is not None else config['gcp'][key]
