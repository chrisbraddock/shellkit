# pet-select.zsh — Ctrl-S to search and insert snippets via pet
# Requires: pet (brew install pet)

_pet_select() {
    BUFFER=$(pet search --query "$LBUFFER")
    CURSOR=$#BUFFER
    zle redisplay
}
zle -N _pet_select

# Disable terminal flow control so Ctrl-S is available
stty -ixon 2>/dev/null

bindkey '^s' _pet_select
