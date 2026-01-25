# ASDP Installer for Windows
$ErrorActionPreference = "Stop"

$AppName = "asdp"
$Repo = "Josepavese/asdp"
$InstallDir = "$env:USERPROFILE\.asdp\bin"
$BinaryName = "asdp-windows-amd64.exe.zst"

# Authentication for private repos
$Headers = @{}
if ($env:GITHUB_TOKEN) {
    $Headers["Authorization"] = "token $env:GITHUB_TOKEN"
}

Write-Host "Starting ASDP Installer..." -ForegroundColor Green

# 1. Structure
if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-Host "Created $InstallDir"
}

# 2. Prerequisites
$PrereqMissing = $false
if (Get-Command "ctags" -ErrorAction SilentlyContinue) {
    Write-Host "Prerequisite 'ctags' found." -ForegroundColor Green
} else {
    Write-Host "Warning: 'ctags' not found." -ForegroundColor Yellow
    Write-Host "Please install Universal Ctags (e.g., 'winget install UniversalCtags')"
}

if (Get-Command "zstd" -ErrorAction SilentlyContinue) {
    Write-Host "Prerequisite 'zstd' found." -ForegroundColor Green
} else {
    Write-Host "Error: 'zstd' not found. Required for decompressing the binary." -ForegroundColor Red
    Write-Host "Please install zstd (e.g., 'winget install facebook.zstd')"
    $PrereqMissing = $true
}

if ($PrereqMissing) { exit 1 }

# 3. Download
$DownloadUrl = "https://github.com/$Repo/releases/latest/download/$BinaryName"
$CoreUrl = "https://github.com/$Repo/releases/latest/download/asdp-core.zip"

$OutputFile = Join-Path $InstallDir "$AppName.exe"
$CoreFile = Join-Path $env:TEMP "asdp-core.zip"

Write-Host "Downloading Binary from $DownloadUrl..."
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile "$OutputFile.zst" -Headers $Headers
    Write-Host "Binary download successful." -ForegroundColor Green
    Write-Host "Decompressing binary..."
    & zstd -d --rm "$OutputFile.zst" -o "$OutputFile"
} catch {
    Write-Host "Binary download failed: $_" -ForegroundColor Red
    if ($_.Exception.Response.StatusCode -eq 404) {
        Write-Host "Note: If the repository is private, ensure GITHUB_TOKEN is set." -ForegroundColor Yellow
        Write-Host "Tip: Try: `$env:GITHUB_TOKEN='your_token'; .\installer\install.ps1`" -ForegroundColor Green
    }
    exit 1
}

Write-Host "Downloading Core Assets from $CoreUrl..."
try {
    Invoke-WebRequest -Uri $CoreUrl -OutFile $CoreFile -Headers $Headers
    Write-Host "Extracting core assets..."
    # Expand-Archive requires destination. We extract to parent of bin (.asdp)
    $BaseDir = "$env:USERPROFILE\.asdp"
    $CoreDir = "$BaseDir\core"
    if (Test-Path $CoreDir) {
        Write-Host "Cleaning up old core assets..."
        Remove-Item -Path $CoreDir -Recurse -Force
    }
    Expand-Archive -Path $CoreFile -DestinationPath $BaseDir -Force
    Remove-Item $CoreFile
    Write-Host "Core assets installed." -ForegroundColor Green
} catch {
    Write-Host "Core download failed (optional): $_" -ForegroundColor Yellow
    Write-Host "If this is a private repo, ensure GITHUB_TOKEN is set."
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
Configure-MCPServer -ConfigPath "$env:USERPROFILE\.gemini\antigravity\mcp_config.json"

Write-Host "ASDP installed successfully!" -ForegroundColor Green
Write-Host "Run 'asdp' to start."

# 5. Interactive Project Initialization
Write-Host ""
$choice = Read-Host "Do you want to initialize ASDP in the current directory? (y/N)"
if ($choice -eq "y") {
    $AgentDir = ".\.agent"
    Write-Host "Initializing ASDP in $((Get-Location).Path)..."
    if (!(Test-Path $AgentDir)) {
        New-Item -ItemType Directory -Path $AgentDir -Force | Out-Null
    }

    # Copy from global core/agent if it exists
    $SrcAgent = "$env:USERPROFILE\.asdp\core\agent"
    if (Test-Path $SrcAgent) {
        Copy-Item -Path "$SrcAgent\*" -Destination $AgentDir -Recursive -Force
        Write-Host "Project initialized successfully in $AgentDir" -ForegroundColor Green
    } else {
        Write-Host "Warning: Global agent templates not found at $SrcAgent" -ForegroundColor Yellow
    }
}
