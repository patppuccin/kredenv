# Hook configuration for kredenv

# Initialize by adding the following to ~/.zshrc:
# eval "$(kredenv hook zsh)"

# Unload secrets tracked in KREDENV_LOADED
__kredenv_unload() {
    if [[ -n "${KREDENV_LOADED:-}" ]]; then
        for key in ${(s:,:)KREDENV_LOADED}; do
            unset "$key"
        done
        unset KREDENV_LOADED
    fi
}

# Load secrets from inject output
__kredenv_load() {
    local keys=""
    while IFS= read -r line; do
        local key value
        line="${line#export }"
        key="${line%%=*}"
        value="${line#*=}"
        value="${value%\"}"
        value="${value#\"}"
        export "$key=$value"
        keys="${keys:+$keys,}$key"
    done <<< "$1"
    export KREDENV_LOADED="$keys"
}

# Hook to detect directory change
__kredenv_hook() {
    \command kredenv inject 2>/dev/null | __kredenv_load
}

# Initialize hook — idempotent, safe to source multiple times
\builtin typeset -ga chpwd_functions
chpwd_functions=("${(@)chpwd_functions:#__kredenv_hook}")
chpwd_functions+=(__kredenv_hook)

# Public interceptor — shadows the kredenv binary
kredenv() {
    case "$1" in
        load)
            __kredenv_unload
            local secrets
            secrets="$(\command kredenv inject 2>/dev/null)"
            if [[ -n "$secrets" ]]; then
                __kredenv_load "$secrets"
            fi
            ;;
        unload)
            __kredenv_unload
            ;;
        inject)
            \builtin printf '%s\n' "kredenv inject is for internal use only" >&2
            return 1
            ;;
        *)
            \command kredenv "$@"
            ;;
    esac
}