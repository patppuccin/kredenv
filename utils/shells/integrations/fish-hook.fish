# Hook configuration for kredenv

# Initialize by adding the following to ~/.config/fish/config.fish:
# kredenv hook fish | source

# Unload secrets tracked in KREDENV_LOADED
function __kredenv_unload
    if set -q KREDENV_LOADED
        for key in (string split ',' $KREDENV_LOADED)
            set -e $key
        end
        set -e KREDENV_LOADED
    end
end

# Load secrets from inject output
function __kredenv_load
    set -l keys
    for line in $argv
        set -l stripped (string replace 'export ' '' $line)
        set -l parts (string split --max 1 '=' $stripped)
        set -l key $parts[1]
        set -l value (string trim --chars='"' $parts[2])
        set -gx $key $value
        set -a keys $key
    end
    set -gx KREDENV_LOADED (string join ',' $keys)
end

# Hook to detect directory change
function __kredenv_hook --on-variable PWD
    set -l secrets (command kredenv inject 2>/dev/null)
    if test -n "$secrets"
        __kredenv_unload
        __kredenv_load $secrets
    end
end

# Public interceptor — shadows the kredenv binary
function kredenv
    switch $argv[1]
        case load
            __kredenv_unload
            set -l secrets (command kredenv inject 2>/dev/null)
            if test -n "$secrets"
                __kredenv_load $secrets
            end
        case unload
            __kredenv_unload
        case inject
            echo "kredenv inject is for internal use only" >&2
            return 1
        case '*'
            command kredenv $argv
    end
end