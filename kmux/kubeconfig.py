import logging

log = logging.getLogger(__name__)


def find_name(kube_config: dict[str, any]) -> str:
    first_cluster_name = kube_config["clusters"][0]["name"]
    cc = kube_config['current-context']
    active_cluster_name = list(filter(lambda c: c['name'] == cc, kube_config['contexts']))
    return active_cluster_name[0]['context']['cluster'] if active_cluster_name else first_cluster_name
