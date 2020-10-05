import getpass

from google.cloud import container_v1
from google.cloud.container_v1.types import ListClustersRequest, GetClusterRequest
from google.oauth2.service_account import Credentials
from google.cloud.container_v1.types.cluster_service import Cluster


def list_clusters(project: str, credentials: Credentials = None):
    gclient = container_v1.ClusterManagerClient(credentials=credentials)
    c = gclient.list_clusters(ListClustersRequest(parent="projects/{}/locations/-".format(project)))
    return c.clusters


def get_cluster(project: str, name: str, location: str, credentials: Credentials = None):
    gclient = container_v1.ClusterManagerClient(credentials=credentials)
    return gclient.get_cluster(
        GetClusterRequest(name="projects/{}/locations/{}/clusters/{}".format(project, location, name)))


def kubeconfig_clusters(cluster: Cluster):
    return [{
        'cluster': {
            'certificate-authority-data': cluster.master_auth.cluster_ca_certificate,
            'server': "https://{}".format(cluster.endpoint)
        },
        'name': "gke-{}-{}".format(cluster.name, cluster.location)
    }]


def kubeconfig_users(creds):
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
