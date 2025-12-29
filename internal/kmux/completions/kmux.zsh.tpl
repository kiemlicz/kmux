_km() {
    local commands projects
    commands=({{.AllCommands}})
    projects=({{.Projects}})

    if (( CURRENT == 2 )); then
      _arguments '1:commands:($commands)'
    elif (( CURRENT == 3 )); then
      case $words[2] in
        {{.RunnableCommands}})
          _arguments '*:projects:_path_files -W "($projects)" -g "*(:r)"'
          ;;
        completions)
          _arguments '*:shell:(zsh)'
          ;;
      esac
    fi
}
compdef _km km
