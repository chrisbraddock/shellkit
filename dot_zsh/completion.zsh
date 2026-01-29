# ============================================================
# Zsh Completion Configuration
# ============================================================

# Case-insensitive matching for lowercase
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Z}'

# Pasting with tabs doesn't perform completion
zstyle ':completion:*' insert-tab pending

# Completion menu with selection
zstyle ':completion:*' menu select

# Group matches by category
zstyle ':completion:*' group-name ''

# Colors in completion menu
zstyle ':completion:*:default' list-colors ${(s.:.)LS_COLORS}

# Process completion
zstyle ':completion:*:*:kill:*:processes' list-colors '=(#b) #([0-9]#)*=0=01;31'
zstyle ':completion:*:kill:*' command 'ps -u $USER -o pid,%cpu,tty,cputime,cmd'

# Don't complete uninteresting users
zstyle ':completion:*:*:*:users' ignored-patterns \
    adm amanda apache avahi beaglidx bin cacti canna clamav daemon \
    dbus distcache dovecot fax ftp games gdm gkrellmd gopher \
    hacluster haldaemon halt hsqldb ident junkbust ldap lp mail \
    mailman mailnull mldonkey mysql nagios named netdump news nfsnobody \
    nobody nscd ntp nut nx openvpn operator pcap postfix postgres \
    privoxy pulse pvm quagga radvd rpc rpcuser rpm shutdown squid \
    sshd sync uucp vcsa xfs '_*'

# SSH/SCP/rsync completion
zstyle ':completion:*:(scp|rsync):*' tag-order 'hosts:-host:host hosts:-domain:domain hosts:-ipaddr:ip\ address *'
zstyle ':completion:*:ssh:*' tag-order 'hosts:-host:host hosts:-domain:domain hosts:-ipaddr:ip\ address *'
