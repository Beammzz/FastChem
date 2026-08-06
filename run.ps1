#Requires -Version 5.1
# Local dev: build the frontend static export, then run the backend that serves it.
# Windows counterpart of run.sh. Ctrl+C stops the backend.

$ErrorActionPreference = 'Stop'

$RootDir = $PSScriptRoot
$LogDir = Join-Path $RootDir 'logs'
New-Item -ItemType Directory -Force -Path $LogDir | Out-Null

# Prefer pnpm if available, same as run.sh
if (Get-Command pnpm -ErrorAction SilentlyContinue) {
    $PkgCmd = 'pnpm'
} elseif (Get-Command npm -ErrorAction SilentlyContinue) {
    $PkgCmd = 'npm'
} else {
    Write-Error 'No package manager found (pnpm or npm).'
}

# Transcript instead of piping to Tee-Object: in Windows PowerShell 5.1, piping a
# native command with 2>&1 wraps every stderr line in a NativeCommandError.
Start-Transcript -Path (Join-Path $LogDir 'dev.log') | Out-Null

try {
    Write-Host '=== FastChem: Building & starting ==='

    Write-Host 'Building frontend...'
    Push-Location (Join-Path $RootDir 'frontend')
    try {
        & $PkgCmd run build
        if ($LASTEXITCODE -ne 0) { Write-Error "Frontend build failed (exit $LASTEXITCODE)." }
    } finally {
        Pop-Location
    }
    Write-Host 'Frontend build done -> frontend/out/'

    Write-Host ''
    Write-Host '===> Open http://localhost:8080 to use FastChem <==='
    Write-Host 'Ctrl+C to stop.'
    Write-Host ''

    $env:FRONTEND_DIR = Join-Path $RootDir 'frontend\out'
    Push-Location (Join-Path $RootDir 'backend')
    try {
        & go run ./cmd/server
    } finally {
        Pop-Location
        Write-Host "Stopped. Log is $LogDir\dev.log"
    }
} finally {
    Stop-Transcript | Out-Null
}
