# import libtmux

NAME = 'name'
ROOT_DIR = 'root'
PRE_WINDOW = 'pre_window'


# def ensure_server() -> libtmux.Server:
# '''
# Either create new or return existing server
# '''
#    return libtmux.Server()


# def spawn_session(name: str, kubeconfig_location: str, server: libtmux.Server):
#     if server.has_session(name):
#         return
#     else:
#         session = server.new_session(name)
#         session.set_environment("KUBECONFIG", kubeconfig_location)
#         # the new_session will create default window and pane which will not contain KUBECONFIG, add manually
#         session.attached_window.attached_pane.send_keys("export KUBECONFIG={}".format(kubeconfig_location))

def _merge(source, destination):
    for k, v in source.items():
        if isinstance(v, dict):
            n = destination.setdefault(k, {})
            _merge(v, n)
        # consider list append
        else:
            destination[k] = v
    return destination


def generate_tmuxinator_template(name: str, root_dir: str, kubeconfig_location: str, template: dict[str, any]) -> dict[str, any]:
    c = {
        NAME: name,
        ROOT_DIR: root_dir,
        PRE_WINDOW: f"export KUBECONFIG={kubeconfig_location} && tmux setenv KUBECONFIG {kubeconfig_location}"
    }

    return _merge(template, c)  # add template
