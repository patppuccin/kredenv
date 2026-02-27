# Hook configuration for kredenv

# Initialize by saving to your autoload directory:
# kredenv hook nushell | save -f ($nu.default-config-dir | path join "autoload" "kredenv.nu")

# Unload secrets tracked in KREDENV_LOADED_VARS
def --env _kredenv_unload [] {
     if ($env.KREDENV_LOADED_VARS? | is-not-empty) {
        for key in ($env.KREDENV_LOADED | split row ",") {
            hide-env --ignore-errors $key
        }
        hide-env --ignore-errors KREDENV_LOADED_VARS
        hide-env --ignore-errors KREDENV_LOADED_COUNT
    }
}

# Load secrets from inject output
def --env _kredenv_load [secrets: list<string>] {
    let pairs = $secrets | reduce --fold {} {|line, acc|
        let stripped = $line | str replace "export " ""
        let parts    = $stripped | split row "=" | collect
        let key      = $parts | first
        let value    = $parts | last | str trim -c '"'
        $acc | insert $key $value
    }
    load-env $pairs
    $env.KREDENV_LOADED_VARS = ($pairs | columns | str join ",")
    $env.KREDENV_LOADED_COUNT = ($pairs | columns | length | into string)
}

# Initialize PWD hook — idempotent, safe to source multiple times
export-env {
    $env.__KREDENV_BIN = (which ^kredenv | where type == "external" | get path.0)

    $env.config = (
        $env.config?
        | default {}
        | upsert hooks { default {} }
        | upsert hooks.env_change { default {} }
        | upsert hooks.env_change.PWD { default [] }
    )

    let already_hooked = (
        $env.config.hooks.env_change.PWD | any { try { get __kredenv_hook } catch { false } }
    )

    if not $already_hooked {
        $env.config.hooks.env_change.PWD = ($env.config.hooks.env_change.PWD | append {
            __kredenv_hook: true,
            code: { | _, _dir |
                _kredenv_unload
                let secrets = (nu -c $"($env.__KREDENV_BIN) inject" | lines)
                if ($secrets | is-not-empty) { _kredenv_load $secrets }
            }
        })
    }
}

# Public interceptor — shadows the kredenv binary
def --env --wrapped kredenv [...args: string] {
    if ($args | is-empty) {
        nu -c $"($env.__KREDENV_BIN)"
        return
    }

    if ("--help" in $args) or ("-h" in $args) {
        let args_str = ($args | str join " ")
        nu -c $"($env.__KREDENV_BIN) ($args_str)"
        return
    }

    match ($args | first) {
        "load" => {
            _kredenv_unload
            let ns_flag = if ("--namespace" in $args) {
                let idx = ($args | index-of "--namespace")
                $" --namespace ($args | get ($idx + 1))"
            } else if ("-n" in $args) {
                let idx = ($args | index-of "-n")
                $" --namespace ($args | get ($idx + 1))"
            } else {
                ""
            }
            let secrets = (nu -c $"($env.__KREDENV_BIN) inject($ns_flag)" | lines)
            if ($secrets | is-not-empty) { _kredenv_load $secrets }
            let args_str = ($args | str join " ")
            nu -c $"($env.__KREDENV_BIN) ($args_str)"
        }
        "unload" => {
            _kredenv_unload
            nu -c $"($env.__KREDENV_BIN) unload"
        }
        "inject" => {
            print -e $"(ansi red)kredenv inject is for internal use only(ansi reset)"
        }
        _ => {
            let args_str = ($args | str join " ")
            nu -c $"($env.__KREDENV_BIN) ($args_str)"
        }
    }
}