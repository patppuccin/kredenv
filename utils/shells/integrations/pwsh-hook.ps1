# Hook configuration for kredenv

# Initialize by adding following command to $PROFILE:
# Invoke-Expression (& { (kredenv hook powershell | Out-String) })

# Unload Secrets tracked in KREDENV_LOADED_VARS
function global:_KredenvUnload {
    if ($null -ne $env:KREDENV_LOADED_VARS) {
        foreach ($Key in ($env:KREDENV_LOADED_VARS -split ",")) {
            Remove-Item -Path "Env:$Key" -ErrorAction SilentlyContinue
        }
        Remove-Item -Path "Env:KREDENV_LOADED_VARS" -ErrorAction SilentlyContinue
        Remove-Item -Path "Env:KREDENV_LOADED_COUNT" -ErrorAction SilentlyContinue
    }
}

# Load Secrets from inject output
function global:_KredenvLoad {
    param([string[]]$Secrets)
    $Keys = @()
    foreach ($Line in $Secrets) {
        $Stripped = $Line -replace '^export ', ''
        $Eq = $Stripped.IndexOf('=')
        if ($Eq -lt 0) { continue }

        $Key = $Stripped.Substring(0, $Eq)
        $Value = $Stripped.Substring($Eq + 1).Trim('"')

        Set-Item -Path "Env:$Key" -Value $Value
        $Keys += $Key
    }
    $env:KREDENV_LOADED_VARS = $Keys -join ","
    $env:KREDENV_LOADED_COUNT = $Keys.Length
}

# Hook to detect directory change
function global:_KredenvHook {
    $Result = (Get-Location).ProviderPath
    if ($Result -ne $global:_KredenvOldPwd) {
        $global:_KredenvOldPwd = $Result
        _KredenvUnload
        $Secrets = & $env:__KREDENV_BIN inject 2>$null
        if ($Secrets) { _KredenvLoad $Secrets }
    }
}

# Initialize hook — idempotent, safe to source multiple times
$global:_KredenvHooked = (Get-Variable _KredenvHooked -ErrorAction Ignore -ValueOnly)
if ($global:_KredenvHooked -ne 1) {
    $global:_KredenvHooked = 1
    $global:_KredenvOldPwd = (Get-Location).ProviderPath
    $global:_KredenvOldPrompt = $function:Prompt

    $env:__KREDENV_BIN = (Get-Command -CommandType Application kredenv).Source

    function global:Prompt {
        if ($null -ne $global:_KredenvOldPrompt) {
            & $global:_KredenvOldPrompt
        }
        $null = _KredenvHook
    }
}

# Public interceptor — shadows the kredenv binary
function global:kredenv {
    if ($args -contains "--help" -or $args -contains "-h") {
        & $env:__KREDENV_BIN @args
        return
    }
    switch ($args[0]) {
        "load" {
            _KredenvUnload
            $NSFlags = @()
            $NSFullIdx = [Array]::IndexOf($args, "--namespace")
            $NSShortIdx = [Array]::IndexOf($args, "-n")
            if ($NSFullIdx -ge 0 -and $NSFullIdx + 1 -lt $args.Length) {
                $NSFlags = @("--namespace", $args[$NSFullIdx + 1])
            }
            elseif ($NSShortIdx -ge 0 -and $NSShortIdx + 1 -lt $args.Length) {
                $NSFlags = @("--namespace", $args[$NSShortIdx + 1])
            }
            $Secrets = & $env:__KREDENV_BIN inject @NSFlags 2>$null
            if ($Secrets) { _KredenvLoad $Secrets }
            & $env:__KREDENV_BIN @args
        }
        "unload" {
            _KredenvUnload
            & $env:__KREDENV_BIN unload
        }
        "inject" {
            Write-Host "kredenv inject is for internal use only" -ForegroundColor Red
        }
        default {
            & $env:__KREDENV_BIN @args
        }
    }
}