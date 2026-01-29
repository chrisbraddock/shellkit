# ============================================================
# Git Aliases
# ============================================================
# Consolidated git aliases for efficient git workflows
# ============================================================

# ------------------------------------------------------------
# Status & Information
# ------------------------------------------------------------

alias gs='git status'
alias gl='git log'
alias glg='git log --graph'
alias gls='git log --stat'
alias glp='git log --patch'
alias gsh='git show'

# ------------------------------------------------------------
# Branches
# ------------------------------------------------------------

alias gb="git --no-pager branch"
alias gbv='git branch -v'
alias gdlm='git-delete-local-merged'
alias gpromote='git-promote'

# ------------------------------------------------------------
# Checkout
# ------------------------------------------------------------

alias gco='git checkout'
alias gcob='git checkout -b'
alias gcom='git checkout main'
alias gcod='git checkout dev'

# ------------------------------------------------------------
# Staging & Commits
# ------------------------------------------------------------

alias ga='git add'
alias gaa='git add --all'
alias gc='git commit'
alias gcm='git commit -m'
alias gca='git commit --amend'
alias gall='git-all'
alias gamend='git-amend'
alias gcredit='git-credit'
alias gundo='git-undo'

# ------------------------------------------------------------
# Diffing
# ------------------------------------------------------------

alias gd='git diff'
alias gds='git diff --staged'
alias gdds="gdd --staged"
alias gdt='git difftool'

# Pretty diff function using dunk (if available)
gdd() {
    if command -v dunk &>/dev/null; then
        git diff "$@" | dunk
    else
        git diff "$@"
    fi
}

# ------------------------------------------------------------
# Fetch & Pull
# ------------------------------------------------------------

alias gf='git fetch'
alias gfa='git fetch --all'
alias gfap='git fetch --all --prune'
alias gpl='git pull'
alias gpod='git pull origin dev'
alias gprod='git pull --rebase origin dev'
alias gprom='git pull --rebase origin main'
alias gup='git-up'

# ------------------------------------------------------------
# Stash
# ------------------------------------------------------------

alias gst='git stash'
alias gstp='git stash pop'
alias gstl='git stash list'

# ------------------------------------------------------------
# Cleanup & Misc
# ------------------------------------------------------------

alias grm='git rm'
alias gmv='git mv'
alias gclean='git clean -fd'

# ------------------------------------------------------------
# Functions
# ------------------------------------------------------------

# Archive a branch with a tag before deleting
git-archive() {
    if [[ $# -eq 1 ]]; then
        git tag "archive/$1" "$1" && git branch -D "$1"
    else
        echo "Usage: git-archive <branch-name>"
    fi
}

# Normalize line endings according to .gitattributes
git-normalize() {
    echo "Normalizing line endings..."
    echo "Warning: Partially staged files will become fully staged."
    # Save list of currently staged files
    local staged_files
    staged_files=$(git diff --cached --name-only -z)
    # Unstage everything
    git reset HEAD --quiet 2>/dev/null || true
    # Renormalize entire repo
    git add --renormalize .
    if git diff --cached --quiet; then
        echo "No line ending changes needed."
    else
        git commit -m "Normalize line endings according to .gitattributes"
    fi
    # Re-stage original files if any
    if [[ -n "$staged_files" ]]; then
        echo -n "$staged_files" | xargs -0 git add --
        echo "Re-staged your original files."
    fi
}
