# Hook configuration for kredenv

# Initialize by adding the following to ~/.config/fish/config.fish:
# kredenv hook fish | source

# Unload secrets tracked in __KREDENV_LOADED_VARS
function __kredenv_unload
    if set -q __KREDENV_LOADED_VARS
        for key in (string split ',' $__KREDENV_LOADED_VARS)
            if test -n "$key"
                set -e $key
            end
        end
        set -e __KREDENV_LOADED_NS
        set -e __KREDENV_LOADED_VARS
        set -e __KREDENV_LOADED_COUNT
    end
end

# Load secrets from inject output
function __kredenv_load
    set -l keys
    for line in $argv
        set -l parts (string split --max 1 '=' $line)
        set -l key $parts[1]
        set -l value $parts[2]
        set -gx $key $value
        if not string match -q '__KREDENV_*' $key
            set -a keys $key
        end
    end
    set -gx __KREDENV_LOADED_VARS (string join ',' $keys)
    set -gx __KREDENV_LOADED_COUNT (count $keys)
end

# Hook to detect directory change
function __kredenv_hook --on-variable PWD
    __kredenv_unload
    set -l secrets (command kredenv inject --format dotenv 2>/dev/null)
    if test -n "$secrets"
        __kredenv_load $secrets
    end
end

# Initialize binary path (safe to source multiple times)
if not set -q __KREDENV_BIN
    set -gx __KREDENV_BIN (command -v kredenv)
end

# Public interceptor (shadows the kredenv binary)
function kredenv
    if contains -- --help $argv; or contains -- -h $argv
        command kredenv $argv
        return
    end

    switch $argv[1]
        case load
            __kredenv_unload
            set -l ns_flags
            set -l i 2
            while test $i -le (count $argv)
                if test $argv[$i] = --namespace -o $argv[$i] = -n
                    set next (math $i + 1)
                    set ns_flags --namespace $argv[$next]
                    break
                end
                set i (math $i + 1)
            end
            set -l secrets (command kredenv inject --format dotenv $ns_flags 2>/dev/null)
            if test -n "$secrets"
                __kredenv_load $secrets
            end
            command kredenv $argv
        case unload
            __kredenv_unload
            command kredenv unload
        case inject
            printf '%s\n' "kredenv inject is for internal use only" >&2
            return 1
        case '*'
            command kredenv $argv
    end
end