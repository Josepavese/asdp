# ASDP Installer for Windows
$ErrorActionPreference = "Stop"

$AppName = "asdp"
$Repo = "Josepavese/asdp"
$InstallDir = "$env:USERPROFILE\.asdp\bin"
$BinaryName = "asdp-windows-amd64.exe"

Write-Host "Starting ASDP Installer..." -ForegroundColor Green

# 1. Structure
if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-Host "Created $InstallDir"
}

# 2. Prerequisites
if (Get-Command "ctags" -ErrorAction SilentlyContinue) {
    Write-Host "Prerequisite 'ctags' found." -ForegroundColor Green
} else {
    Write-Host "Warning: 'ctags' not found." -ForegroundColor Yellow
    Write-Host "Please install Universal Ctags via Winget or Chocolatey:"
    Write-Host "winget install UniversalCtags"
}

# 3. Download
$DownloadUrl = "https://github.com/$Repo/releases/latest/download/$BinaryName"
$CoreUrl = "https://github.com/$Repo/releases/latest/download/asdp-core.zip"

$OutputFile = Join-Path $InstallDir "$AppName.exe"
$CoreFile = Join-Path $env:TEMP "asdp-core.zip"

Write-Host "Downloading Binary from $DownloadUrl..."
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $OutputFile
    Write-Host "Binary download successful." -ForegroundColor Green
} catch {
    Write-Host "Binary download failed: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Downloading Core Assets from $CoreUrl..."
try {
    Invoke-WebRequest -Uri $CoreUrl -OutFile $CoreFile
    Write-Host "Extracting core assets..."
    # Expand-Archive requires destination. We extract to parent of bin (.asdp)
    $BaseDir = "$env:USERPROFILE\.asdp"
    Expand-Archive -Path $CoreFile -DestinationPath $BaseDir -Force
    Remove-Item $CoreFile
    Write-Host "Core assets installed." -ForegroundColor Green
} catch {
    Write-Host "Core download failed (optional): $_" -ForegroundColor Yellow
}

# 4. Path Setup
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to User Path..."
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host "Path updated. Please restart PowerShell." -ForegroundColor Green
} else {
    Write-Host "Path already configured."
}

Write-Host "ASDP installed successfully!" -ForegroundColor Green
