class KubeConfig:
    def __init__(self):
        self.apiVersion = 'v1'
        self.kind = 'Config'
        self.clusters = []
        self.users = []
        self.contexts = []


def _kubeconfig_contexts(clusters_entry, users_entry):
    return [{'context': {'cluster': c['name'], 'namespace': "default", 'user': u['name']}, 'name': c['name']} for c, u
            in zip(clusters_entry, users_entry)]


def generate_kubeconfig(users, clusters):
    k = {
        'apiVersion': 'v1',
        'kind': 'Config',
    }
    k['clusters'] = clusters
    k['users'] = users
    k['contexts'] = _kubeconfig_contexts(clusters, users)
    k['current-context'] = k['contexts'][0]['name']
    return k
