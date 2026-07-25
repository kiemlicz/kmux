_km() {
  local -a commands
  commands=({{.AllCommands}})
  if (( CURRENT == 2 )); then
    _describe 'command' commands
    return
  fi
  case ${words[2]} in
    new)
      _arguments \
        '--location=[tmuxinator config directory]:location:(({{.Projects}}))' \
        '--root=[project root directory]:root:_path_files -/' \
        '--kubeconfig=[kubeconfig file]:kubeconfig:_files' \
        '2:name:_path_files -W "({{.Projects}})" -g "*.(yml|yaml)(:r)"'
      ;;
    discover)
      _arguments '2::name (optional, defaults to current TMUX/KUBECONFIG env):_path_files -W "({{.Projects}})" -g "*.(yml|yaml)(:r)"'
      ;;
    start)
      _arguments \
        '--bg[spawn session in background, do not attach]' \
        '2:name:_path_files -W "({{.Projects}})" -g "*.(yml|yaml)(:r)"'
      ;;
    completions)
      _arguments '2:shell:(zsh)'
      ;;
  esac
}
compdef _km km
