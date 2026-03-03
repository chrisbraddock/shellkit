# ai-prompt.zsh — Alt-A to describe a command in English, get shell command back
# Requires: claude CLI

_ai_prompt_activate() {
    _AI_PROMPT_SAVED_BUFFER="$BUFFER"
    _AI_PROMPT_SAVED_CURSOR="$CURSOR"
    BUFFER=""
    CURSOR=0
    PREDISPLAY="ai> "
    _AI_PROMPT_ACTIVE=1
    zle -K ai-prompt
}

_ai_prompt_submit() {
    local query="$BUFFER"
    [[ -z "$query" ]] && { _ai_prompt_cancel; return; }

    PREDISPLAY="ai> thinking... "
    BUFFER=""
    zle -R

    local result
    result=$(claude --print --model haiku --max-tokens 200 \
        "You are a shell command generator. Respond with ONLY the command, no explanation, no markdown, no code fences. OS: $(uname -s). Shell: zsh. Query: $query" 2>/dev/null)

    _ai_prompt_cleanup
    if [[ -n "$result" ]]; then
        BUFFER="$result"
        CURSOR=${#BUFFER}
    else
        BUFFER="$_AI_PROMPT_SAVED_BUFFER"
        CURSOR="$_AI_PROMPT_SAVED_CURSOR"
        zle -M "ai: no response (is claude CLI working?)"
    fi
}

_ai_prompt_cancel() {
    _ai_prompt_cleanup
    BUFFER="$_AI_PROMPT_SAVED_BUFFER"
    CURSOR="$_AI_PROMPT_SAVED_CURSOR"
}

_ai_prompt_cleanup() {
    PREDISPLAY=""
    _AI_PROMPT_ACTIVE=0
    zle -K main
}

# Register widgets
zle -N _ai_prompt_activate
zle -N _ai_prompt_submit
zle -N _ai_prompt_cancel

# Custom keymap (inherits from main so normal editing works)
bindkey -N ai-prompt main
bindkey -M ai-prompt '^M' _ai_prompt_submit     # Enter
bindkey -M ai-prompt '^[' _ai_prompt_cancel      # Escape
bindkey -M ai-prompt '^C' _ai_prompt_cancel      # Ctrl-C

# Bind activation to Alt-A
bindkey '^[a' _ai_prompt_activate
