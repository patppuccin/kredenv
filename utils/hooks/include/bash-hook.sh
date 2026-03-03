# Hook configuration for kredenv

# Initialize by adding the following to ~/.bashrc:
# eval "$(kredenv hook bash)"

# Unload secrets tracked in KREDENV_LOADED_VARS
__kredenv_unload() {
    if [[ -n "${KREDENV_LOADED_VARS:-}" ]]; then
        for key in ${KREDENV_LOADED_VARS//,/ }; do
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
        local key value
        key="${line%%=*}"
        value="${line#*=}"
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

# Hook to detect directory change
__kredenv_hook() {
    local retval="$?"
    local pwd_tmp
    pwd_tmp="$(builtin pwd -P)"
    if [[ "${__kredenv_oldpwd}" != "${pwd_tmp}" ]]; then
        __kredenv_oldpwd="${pwd_tmp}"
        __kredenv_unload
        local secrets
        secrets="$(command kredenv inject --format dotenv 2>/dev/null)"
        if [[ -n "$secrets" ]]; then
            __kredenv_load "$secrets"
        fi
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
    # passthrough help flags
    for arg in "$@"; do
        if [[ "$arg" == "--help" || "$arg" == "-h" ]]; then
            command kredenv "$@"
            return
        fi
    done

    case "$1" in
        load)
            __kredenv_unload
            local ns_flags=()
            local i=1
            while [[ $i -le $# ]]; do
                local arg="${!i}"
                if [[ "$arg" == "--namespace" || "$arg" == "-n" ]]; then
                    local next=$((i + 1))
                    ns_flags=("--namespace" "${!next}")
                    break
                fi
                ((i++))
            done
            local secrets
            secrets="$(command kredenv inject --format dotenv "${ns_flags[@]}" 2>/dev/null)"
            if [[ -n "$secrets" ]]; then
                __kredenv_load "$secrets"
            fi
            command kredenv "$@"
            ;;
        unload)
            __kredenv_unload
            command kredenv unload
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