# Hook configuration for kredenv

# Initialize by adding the following to ~/.bashrc:
# eval "$(kredenv hook bash)"

# Unload secrets tracked in KREDENV_LOADED
__kredenv_unload() {
    if [[ -n "${KREDENV_LOADED:-}" ]]; then
        for key in ${KREDENV_LOADED//,/ }; do
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
    local retval="$?"
    local pwd_tmp
    pwd_tmp="$(builtin pwd -P)"
    if [[ "${__kredenv_oldpwd}" != "${pwd_tmp}" ]]; then
        __kredenv_oldpwd="${pwd_tmp}"
        kredenv load
    fi
    return "${retval}"
}

# Initialize hook — idempotent, safe to source multiple times
__kredenv_oldpwd="$(builtin pwd -P)"
if [[ "${PROMPT_COMMAND:-}" != *'__kredenv_hook'* ]]; then
    PROMPT_COMMAND="__kredenv_hook;${PROMPT_COMMAND#;}"
fi

# Public interceptor — shadows the kredenv binary
kredenv() {
    case "$1" in
        load)
            __kredenv_unload
            local secrets
            secrets="$(command kredenv inject 2>/dev/null)"
            if [[ -n "$secrets" ]]; then
                __kredenv_load "$secrets"
            fi
            ;;
        unload)
            __kredenv_unload
            ;;
        inject)
            printf '%s\n' "kredenv inject is for internal use only" >&2
            return 1
            ;;
        *)
            command kredenv "$@"
            ;;
    esac
}