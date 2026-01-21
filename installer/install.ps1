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


# --- MCP Configuration Logic ---
function Configure-MCPServer {
    param (
        [string]$ConfigPath
    )

    if (Test-Path $ConfigPath) {
        Write-Host "Found MCP config: $ConfigPath"
        try {
            $jsonContent = Get-Content $ConfigPath -Raw | ConvertFrom-Json -Depth 10 # Depth needed for nested objects

            # Ensure mcpServers object exists (PowerShell PSCustomObject weirdness handling)
            if (-not $jsonContent.PSObject.Properties.Match('mcpServers').Count) {
                $jsonContent | Add-Member -MemberType NoteProperty -Name "mcpServers" -Value @{}
            }

            # Create the ASDP server config object
            $asdpConfig = @{
                command = "$InstallDir\$AppName.exe"
                args    = @()
                env     = @{}
            }

            # Add or Update ASDP entry
            # Note: Manipulating nested PSCHustomObjects usually requires re-creating them or specific casting
            # A simpler way for this specific structure:
            if ($null -eq $jsonContent.mcpServers.asdp) {
                 $jsonContent.mcpServers | Add-Member -MemberType NoteProperty -Name "asdp" -Value $asdpConfig
            } else {
                 $jsonContent.mcpServers.asdp = $asdpConfig
            }
            
            $jsonContent | ConvertTo-Json -Depth 10 | Set-Content $ConfigPath
            Write-Host "Successfully updated ASDP in $ConfigPath" -ForegroundColor Green
        } catch {
            Write-Host "Error updating config $ConfigPath : $_" -ForegroundColor Red
        }
    }
}

Write-Host "Configuring IDE Integrations..."

# Common Windows Paths
$AppData = $env:APPDATA
$VSCodePath = "$AppData\Code\User\globalStorage\mcp-servers.json"
$ClaudePath = "$AppData\Claude\claude_desktop_config.json"
$CursorPath = "$AppData\Cursor\User\globalStorage\mcp-servers.json"

Configure-MCPServer -ConfigPath $VSCodePath
Configure-MCPServer -ConfigPath $ClaudePath
Configure-MCPServer -ConfigPath $CursorPath

Write-Host "ASDP installed successfully!" -ForegroundColor Green
Write-Host "Run 'asdp' to start."
