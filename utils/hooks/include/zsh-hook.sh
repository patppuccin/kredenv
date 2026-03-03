# Hook configuration for kredenv

# Initialize by adding the following to ~/.zshrc:
# eval "$(kredenv hook zsh)"

# Unload secrets tracked in KREDENV_LOADED_VARS
__kredenv_unload() {
    if [[ -n "${KREDENV_LOADED_VARS:-}" ]]; then
        for key in ${(s:,:)KREDENV_LOADED_VARS}; do
            unset "$key"
        done
        unset KREDENV_LOADED_VARS
        unset KREDENV_LOADED_COUNT
    fi
}

# Load secrets from inject output
__kredenv_load() {
    local keys=""
    while IFS= read -r line; do
        local key="${line%%=*}"
        local value="${line#*=}"
        export "$key=$value"
        keys="${keys:+$keys,}$key"
    done <<< "$1"
    export KREDENV_LOADED_VARS="$keys"
    local count=0
    if [[ -n "$keys" ]]; then
        count=$(tr ',' '\n' <<< "$keys" | wc -l | tr -d ' ')
    fi
    export KREDENV_LOADED_COUNT="$count"
}

# Hook to reload secrets on directory change
__kredenv_hook() {
    __kredenv_unload
    local secrets
    secrets="$(\command kredenv inject --format dotenv 2>/dev/null)"
    if [[ -n "$secrets" ]]; then
        __kredenv_load "$secrets"
    fi
}

# Initialize binary path and register chpwd hook (safe to source multiple times)
\builtin typeset -ga chpwd_functions
chpwd_functions=("${(@)chpwd_functions:#__kredenv_hook}")
chpwd_functions+=(__kredenv_hook)
export __KREDENV_BIN="$(\builtin whence -p kredenv)"

# Public interceptor (shadows the kredenv binary)
kredenv() {
    for arg in "$@"; do
        if [[ "$arg" == "--help" || "$arg" == "-h" ]]; then
            \command kredenv "$@"
            return
        fi
    done

    case "$1" in
        load)
            __kredenv_unload
            local ns_flags=()
            local i=2
            while [[ $i -le $# ]]; do
                local arg="${@[$i]}"
                if [[ "$arg" == "--namespace" || "$arg" == "-n" ]]; then
                    local next=$((i + 1))
                    ns_flags=("--namespace" "${@[$next]}")
                    break
                fi
                ((i++))
            done
            local secrets
            secrets="$(\command kredenv inject --format dotenv "${ns_flags[@]}" 2>/dev/null)"
            if [[ -n "$secrets" ]]; then
                __kredenv_load "$secrets"
            fi
            \command kredenv "$@"
            ;;
        unload)
            __kredenv_unload
            \command kredenv unload
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