# Hook configuration for kredenv

# Initialize by adding following command to $PROFILE:
# Invoke-Expression (& { (kredenv hook powershell | Out-String) })

# Unload secrets tracked in KREDENV_LOADED
function global:_kredenv_unload {
    if ($null -ne $env:KREDENV_LOADED) {
        foreach ($key in ($env:KREDENV_LOADED -split ",")) {
            Remove-Item -Path "Env:$key" -ErrorAction SilentlyContinue
        }
        Remove-Item -Path "Env:KREDENV_LOADED" -ErrorAction SilentlyContinue
    }
}

# Load secrets from inject output
function global:_kredenv_load {
    param([string[]]$secrets)
    $keys = @()
    foreach ($line in $secrets) {
        $stripped = $line -replace '^export ', ''
        $eq = $stripped.IndexOf('=')
        if ($eq -lt 0) { continue }

        $key = $stripped.Substring(0, $eq)
        $value = $stripped.Substring($eq + 1).Trim('"')

        Set-Item -Path "Env:$key" -Value $value
        $keys += $key
    }
    $env:KREDENV_LOADED = $keys -join ","
}

# Hook to detect directory change
function global:_kredenv_hook {
    $result = (Get-Location).ProviderPath
    if ($result -ne $global:_kredenv_oldpwd) {
        $global:_kredenv_oldpwd = $result
        _kredenv_unload
        $secrets = & $env:__KREDENV_BIN inject 2>$null
        if ($secrets) { _kredenv_load $secrets }
    }
}

# Initialize hook — idempotent, safe to source multiple times
$global:_kredenv_hooked = (Get-Variable _kredenv_hooked -ErrorAction Ignore -ValueOnly)
if ($global:_kredenv_hooked -ne 1) {
    $global:_kredenv_hooked = 1
    $global:_kredenv_oldpwd = (Get-Location).ProviderPath
    $global:_kredenv_prompt_old = $function:prompt

    $env:__KREDENV_BIN = (Get-Command -CommandType Application kredenv).Source

    function global:prompt {
        if ($null -ne $global:_kredenv_prompt_old) {
            & $global:_kredenv_prompt_old
        }
        $null = _kredenv_hook
    }
}

# Public interceptor — shadows the kredenv binary
function global:kredenv {
    switch ($args[0]) {
        "load" {
            _kredenv_unload
            $secrets = & $env:__KREDENV_BIN inject 2>$null
            if ($secrets) { _kredenv_load $secrets }
            & $env:__KREDENV_BIN load
        }
        "unload" {
            _kredenv_unload
            & $env:__KREDENV_BIN unload
        }
        "inject" {
            Write-Host "kredenv inject is for internal use only" -ForegroundColor Red
            exit 1
        }
        default {
            & $env:__KREDENV_BIN @args
        }
    }
}