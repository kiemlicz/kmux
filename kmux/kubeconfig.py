import os
from typing import Tuple, List

import yaml


class KubeConfigBase:
    def __init__(self, clusters=[], users=[], contexts=[]):
        self.clusters = clusters
        self.users = users
        self.contexts = contexts

    def generate(self):
        k = {
            'apiVersion': 'v1',
            'kind': 'Config',
            'clusters': self.clusters,
            'contexts': self.contexts,
            'current-context': self.contexts[0]['name'],
            'users': self.users
        }
        return k

    def set_contexts(self, name_namespace: List[Tuple]):
        if len(self.clusters) != 1 or len(self.users) != 1:
            raise KubeConfigException(
                "Only one users+cluster is supported per KUBECONFIG: clusters ({}) or users ({})".format(
                    len(self.clusters), len(self.users)))
        u = self.users[0]
        c = self.clusters[0]

        self.contexts = [{
            'context': {
                'cluster': c['name'], 'namespace': ns, 'user': u['name']
            },
            'name': n
        } for n, ns in name_namespace]


class KubeConfigExistsException(RuntimeError):
    def __init__(self, message):
        self.message = message


class KubeConfigException(RuntimeError):
    def __init__(self, message):
        self.message = message


def save_kube_config(location: str, kube_config: KubeConfigBase, overwrite: bool = False):
    if os.path.exists(location) and not overwrite:
        raise KubeConfigExistsException("KUBECONFIG: {} already exists".format(location))
    with open(location, 'w') as f:
        yaml.dump(kube_config.generate(), f)
