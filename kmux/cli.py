import argparse
import os
from os.path import expanduser, expandvars

import yaml

CONTEXTS = "contexts"
ACTIVE_CTX = "active"
TMUXINATOR_TEMPLATE = "tmuxinator_template"
TMUXINATOR_DIR = "tdir"
KUBECONFIG_DIR = "kdir"


# the config is required as it hosts templates for tmux sessions
def _kmrc_default_location():
    if 'XDG_CONFIG_HOME' in os.environ:
        return os.path.join(expandvars("$XDG_CONFIG_HOME"), "kmrc")
    return os.path.join(expanduser("~"), ".config", "kmrc")


class RequireOneOfTwo(argparse.Action):
    def __call__(self, parser, namespace, values, option_string=None):
        if getattr(namespace, 'start', None) or getattr(namespace, 'create', None):
            setattr(namespace, self.dest, values)
        else:
            parser.error("One of --start or --create is required")


class KmuxConfig:

    def __init__(self, c: dict[str, any], ctx: str = None, name: str = None):
        self.ctx_override = ctx
        self.name_override = name
        self.active_context = c[ACTIVE_CTX]
        self.all_contexts = c[CONTEXTS]
        self.tmuxinator_templates = c[TMUXINATOR_TEMPLATE]

    def get_context_name(self) -> str:
        if self.ctx_override:
            return self.ctx_override
        return self.active_context

    def get_context(self):
        cc = self.get_context_name()
        return list(filter(lambda c: c['name'] == cc, self.all_contexts))[0]

    def get_tdir(self) -> str:
        return self.get_context()[TMUXINATOR_DIR]

    def get_kdir(self) -> str:
        return self.get_context()[KUBECONFIG_DIR]


parser = argparse.ArgumentParser(description='Prepare tmux session for given KUBECONFIG')
parser.add_argument('--name', help="Override cluster name, otherwise cluster name from Kubeconfig is taken", required=False)
parser.add_argument('--start', help="Start cluster tmux", required=False, type=str, action=RequireOneOfTwo)
parser.add_argument('--create', help="Create cluster setup, run to update", required=False, type=str, action=RequireOneOfTwo)
parser.add_argument('--rdir', help="Tmuxinator root dir when creating", required=False, type=str)
parser.add_argument('--ctx', help="Use and change active context", type=str, required=False)
parser.add_argument('--config', help="Configuration file", default=_kmrc_default_location())
parser.add_argument('--log', help="log level (TRACE, DEBUG, INFO, WARN, ERROR)", required=False, default="INFO")
args = parser.parse_args()

with open(args.config, 'r') as f:
    config = KmuxConfig(yaml.safe_load(f), args.ctx, args.name)


def save_kmux_config(c: KmuxConfig):
    if c.name_override or c.ctx_override:
        with open(args.config, 'w') as f:
            d = {CONTEXTS: c.all_contexts, ACTIVE_CTX: c.active_context, TMUXINATOR_TEMPLATE: c.tmuxinator_templates}
            yaml.dump(d, f)
