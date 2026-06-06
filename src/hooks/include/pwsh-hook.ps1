# Hook configuration for kredenv

# Initialize by adding following command to $PROFILE:
# Invoke-Expression (& { (kredenv hook powershell | Out-String) })

# Unload secrets tracked in KREDENV_LOADED_VARS
function global:_KredenvUnload {
    if ($null -ne $env:KREDENV_LOADED_VARS -and $env:KREDENV_LOADED_VARS -ne "") {
        foreach ($Key in ($env:KREDENV_LOADED_VARS -split ",")) {
            if ($Key -ne "") {
                Remove-Item -Path "Env:$Key" -ErrorAction SilentlyContinue
            }
        }
        Remove-Item -Path "Env:KREDENV_LOADED_VARS" -ErrorAction SilentlyContinue
        Remove-Item -Path "Env:KREDENV_LOADED_COUNT" -ErrorAction SilentlyContinue
    }
}

# Load secrets from inject output
function global:_KredenvLoad {
    param([hashtable]$Secrets)
    foreach ($Entry in $Secrets.GetEnumerator()) {
        Set-Item -Path "Env:$($Entry.Key)" -Value $Entry.Value
    }
    $env:KREDENV_LOADED_VARS = ($Secrets.Keys -join ",")
    $env:KREDENV_LOADED_COUNT = $Secrets.Count
}

# Hook to detect directory change
function global:_KredenvHook {
    $Result = (Get-Location).ProviderPath
    if ($Result -ne $global:_KredenvOldPwd) {
        $global:_KredenvOldPwd = $Result
        _KredenvUnload
        $Raw = & $env:__KREDENV_BIN inject --format json 2>$null
        if ($Raw) {
            $Secrets = $Raw | ConvertFrom-Json -AsHashtable
            if ($Secrets.Count -gt 0) { _KredenvLoad $Secrets }
        }
    }
}

# Initialization hook (safe to source multiple times)
$global:_KredenvHooked = (Get-Variable _KredenvHooked -ErrorAction Ignore -ValueOnly)
if ($global:_KredenvHooked -ne 1) {
    $global:_KredenvHooked = 1
    $global:_KredenvOldPwd = (Get-Location).ProviderPath
    $global:_KredenvOldPrompt = $function:Prompt

    $env:__KREDENV_BIN = (Get-Command -CommandType Application kredenv | Select-Object -First 1).Source

    function global:Prompt {
        $null = _KredenvHook
        if ($null -ne $global:_KredenvOldPrompt) {
            & $global:_KredenvOldPrompt
        }
    }
}

# Public interceptor (shadows the kredenv binary)
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
            $Raw = & $env:__KREDENV_BIN inject --format json @NSFlags 2>$null
            if ($Raw) {
                $Secrets = $Raw | ConvertFrom-Json -AsHashtable
                if ($Secrets.Count -gt 0) { _KredenvLoad $Secrets }
            }
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