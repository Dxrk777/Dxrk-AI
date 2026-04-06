set -gx EDITOR nvim
set -gx PATH $HOME/.local/bin $PATH
alias ls "eza --icons"
alias ll "eza -la --icons"
alias cat "bat"
alias lg "lazygit"
alias vim "nvim"
zoxide init fish | source
starship init fish | source
