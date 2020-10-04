import getpass

import google.auth
import google.auth.transport.requests
from google.cloud import container_v1
from google.cloud.container_v1.types import ListClustersRequest, GetClusterRequest
from google.cloud.container_v1.types.cluster_service import Cluster
from google.oauth2.service_account import Credentials


def list_clusters(project: str, credentials: Credentials = None):
    gclient = container_v1.ClusterManagerClient(credentials=credentials)
    c = gclient.list_clusters(ListClustersRequest(parent="projects/{}/locations/-".format(project)))
    return c.clusters


def get_cluster(project: str, name: str, location: str, credentials: Credentials = None):
    gclient = container_v1.ClusterManagerClient(credentials=credentials)
    return gclient.get_cluster(
        GetClusterRequest(name="projects/{}/locations/{}/clusters/{}".format(project, location, name)))


def _kubeconfig_users(creds):
    return [{
        'name': getpass.getuser(),
        'user': {
            'auth-provider': {
                'config': {
                    'access-token': creds.token,
                    'cmd-args': "config config-helper --format=json",
                    'cmd-path': '/usr/lib/google-cloud-sdk/bin/gcloud',
                    'expiry': creds.expiry,
                    'expiry-key': '{.credential.token_expiry}',
                    'token-key': '{.credential.access_token}'
                },
                'name': 'gcp'
            }
        }
    }]


def _kubeconfig_clusters(cluster: Cluster):
    return [{
        'cluster': {
            'certificate-authority-data': cluster.master_auth.cluster_ca_certificate,
            'server': "https://{}".format(cluster.endpoint)
        },
        'name': "gke-{}-{}".format(cluster.name, cluster.location)
    }]


def _kubeconfig_contexts(clusters_entry, users_entry):
    return [{'context': {'cluster': c['name'], 'namespace': "default", 'user':u['name']}, 'name': c['name']} for c,u in zip(clusters_entry, users_entry)]


def kubeconfig(name: str, location: str):
    k = {
        'apiVersion': 'v1',
        'kind': 'Config',
    }
    creds, project = google.auth.default()
    cluster = get_cluster(project, name, location)
    auth_req = google.auth.transport.requests.Request()
    creds.refresh(auth_req)
    users = _kubeconfig_users(creds)
    clusters = _kubeconfig_clusters(cluster)
    k['clusters'] = clusters
    k['users'] = users
    k['contexts'] = _kubeconfig_contexts(clusters, users)
    k['current-context'] = k['contexts'][0]['name']
    return k
