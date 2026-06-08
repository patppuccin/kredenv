# Hook configuration for kredenv

# Initialize by adding the following to ~/.bashrc:
# eval "$(kredenv hook bash)"

# Unload secrets tracked in __KREDENV_LOADED_VARS
__kredenv_unload() {
    if [[ -n "${__KREDENV_LOADED_VARS:-}" ]]; then
        for key in ${__KREDENV_LOADED_VARS//,/ }; do
            unset "$key"
        done
        unset __KREDENV_LOADED_NS
        unset __KREDENV_LOADED_VARS
        unset __KREDENV_LOADED_COUNT
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
        if [[ "$key" != __KREDENV_* ]]; then
            keys="${keys:+$keys,}$key"
        fi
    done <<< "$1"
    export __KREDENV_LOADED_VARS="$keys"
    local count=0
    if [[ -n "$keys" ]]; then
        count=$(tr ',' '\n' <<< "$keys" | wc -l | tr -d ' ')
    fi
    export __KREDENV_LOADED_COUNT="$count"
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

# Initialize binary path and PWD tracking (safe to source multiple times)
__kredenv_oldpwd="$(builtin pwd -P)"
if [[ "${PROMPT_COMMAND:-}" != *'__kredenv_hook'* ]]; then
    PROMPT_COMMAND="__kredenv_hook;${PROMPT_COMMAND#;}"
fi
export __KREDENV_BIN="$(builtin command -v kredenv)"

# Public interceptor (shadows the kredenv binary)
kredenv() {
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
            local i=2
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