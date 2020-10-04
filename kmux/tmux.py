import libtmux


def ensure_server() -> libtmux.Server:
    '''
    Either create new or return existing server
    '''
    return libtmux.Server()


def spawn_session(name: str, kubeconfig_location: str, server: libtmux.Server):
    if server.has_session(name):
        return
    else:
        session = server.new_session(name)
        session.set_environment("KUBECONFIG", kubeconfig_location)
