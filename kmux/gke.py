import getpass

import google.auth
import google.auth.transport.requests
from google.cloud import container_v1
from google.cloud.container_v1.types import GetClusterRequest

from kmux.kubeconfig import KubeConfigBase


class GKEUserKubeConfig(KubeConfigBase):
    def __init__(self):
        creds, project = google.auth.default()
        creds.refresh(google.auth.transport.requests.Request())

        self.project = project
        self.client = container_v1.ClusterManagerClient(credentials=creds)
        self.users = [{
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

    def set_cluster(self, name: str, location: str):
        cluster = self.client.get_cluster(
            GetClusterRequest(name="projects/{}/locations/{}/clusters/{}".format(self.project, location, name)))
        self.clusters = [{
            'cluster': {
                'certificate-authority-data': cluster.master_auth.cluster_ca_certificate,
                'server': "https://{}".format(cluster.endpoint)
            },
            'name': "gke-{}-{}".format(cluster.name, cluster.location)
        }]
