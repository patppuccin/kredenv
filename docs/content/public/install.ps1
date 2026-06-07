# KREDENV Installation Convenience Script for Windows
# Usage: powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://raw.githubusercontent.com/patppuccin/kredenv/main/scripts/install.ps1 | iex"

if ($PSVersionTable.PSVersion.Major -lt 5 -or
    ($PSVersionTable.PSVersion.Major -eq 5 -and $PSVersionTable.PSVersion.Minor -lt 1)) {
    Write-Host "kredenv requires PowerShell 5.1 or later" -ForegroundColor Red
    Write-Host "Current version: $($PSVersionTable.PSVersion)" -ForegroundColor Red
    exit 1
}

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$Repo = "patppuccin/kredenv"
$BinName = "kredenv.exe"
$InstallDir = "$env:USERPROFILE\.local\bin"

function Write-Log {
    param(
        [ValidateSet("INF", "WRN", "ERR")]$Level = "INF",
        [Parameter(Mandatory)][string]$Message
    )
    $colors = @{ INF = "Blue"; WRN = "Yellow"; ERR = "Red" }
    Write-Host "$Level " -ForegroundColor $colors[$Level] -NoNewline
    Write-Host $Message
}

function Get-LatestVersion {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    return $Release.tag_name
}

function Get-Arch {
    if (-not [Environment]::Is64BitOperatingSystem) {
        Write-Log ERR "kredenv requires a 64-bit operating system"
        exit 1
    }

    switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        default {
            Write-Log ERR "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"
            exit 1
        }
    }
}

function Install-Kredenv {
    $Version = Get-LatestVersion
    $CleanVersion = $Version.TrimStart('v')
    $Arch = Get-Arch

    $ArchiveName = "kredenv-windows-$Arch.zip"
    $ChecksumFile = "kredenv_$CleanVersion`_checksums.txt"

    $DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$ArchiveName"
    $ChecksumUrl = "https://github.com/$Repo/releases/download/$Version/$ChecksumFile"

    $TempDir = Join-Path $env:TEMP "kredenv-install"
    $ArchivePath = Join-Path $TempDir $ArchiveName
    $ChecksumPath = Join-Path $TempDir $ChecksumFile

    Write-Log INF "Installing kredenv $Version for $Arch"
    Write-Log INF "Downloading from $DownloadUrl"
    Write-Log INF "Installing to $InstallDir"

    $ExistingBin = Join-Path $InstallDir $BinName
    if (Test-Path $ExistingBin) {
        $ExistingVersion = & $ExistingBin --version 2>$null
        Write-Log WRN "Found existing $ExistingVersion (will be upgraded to $Version)"
    }

    New-Item -ItemType Directory -Force -Path $TempDir | Out-Null

    Write-Log INF "Downloading artifacts (archive, checksums)"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchivePath -UseBasicParsing
    Invoke-WebRequest -Uri $ChecksumUrl -OutFile $ChecksumPath -UseBasicParsing

    $ChecksumLine = Select-String -Path $ChecksumPath -Pattern $ArchiveName |
    Select-Object -ExpandProperty Line

    if (-not $ChecksumLine) {
        Write-Log ERR "Checksum entry not found for $ArchiveName"
        exit 1
    }

    $ExpectedHash = ($ChecksumLine -split '\s+')[0].ToUpper()
    $ActualHash = (Get-FileHash -Path $ArchivePath -Algorithm SHA256).Hash.ToUpper()

    if ($ActualHash -ne $ExpectedHash) {
        Write-Log ERR "Checksum verification failed"
        Write-Log ERR "Expected: $ExpectedHash"
        Write-Log ERR "Actual:   $ActualHash"
        exit 1
    }

    Write-Log INF "Checksum verified"

    Expand-Archive -Path $ArchivePath -DestinationPath $TempDir -Force

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

    $ExtractedBin = Join-Path $TempDir $BinName
    if (-not (Test-Path $ExtractedBin)) {
        Write-Log ERR "Extracted binary not found at expected location"
        exit 1
    }

    Move-Item -Path $ExtractedBin -Destination (Join-Path $InstallDir $BinName) -Force

    Write-Log INF "Installed kredenv $Version to $InstallDir\$BinName"

    Remove-Item -Path $TempDir -Recurse -Force

    $CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")

    if ([string]::IsNullOrEmpty($CurrentPath)) {
        [Environment]::SetEnvironmentVariable("PATH", $InstallDir, "User")
        Write-Log INF "Added $InstallDir to PATH (requires shell restart)"
    }
    else {
        $PathEntries = $CurrentPath -split ';'
        if ($PathEntries -notcontains $InstallDir) {
            $NewPath = "$CurrentPath;$InstallDir"
            [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
            Write-Log INF "Added $InstallDir to PATH (requires shell restart)"
        }
    }

    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Green
    Write-Host "  1. Restart your terminal"
    Write-Host "  2. Run: 'kredenv hook --help' for help with configuring the shell of your choice"
}

try {
    Write-Host ""
    Write-Host "KREDENV Windows Installer" -ForegroundColor Green
    Write-Host "Source Repo: https://github.com/$Repo" -ForegroundColor Gray
    Write-Host ""
    Install-Kredenv
}
catch {
    Write-Host ""
    Write-Log ERR "An error occurred when installing kredenv"
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}