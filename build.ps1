# Windows build script for switcher

param(
    [switch]$Install
)

Write-Host "Building switcher for Windows..." -ForegroundColor Green

# Build the executable
go build -o switcher.exe .

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Build successful!" -ForegroundColor Green

if ($Install) {
    Write-Host "Installing switcher..." -ForegroundColor Green

    # Create installation directory
    $installPath = "$env:LOCALAPPDATA\Programs\switcher"
    New-Item -ItemType Directory -Force -Path $installPath | Out-Null

    # Copy executable
    Copy-Item -Path "switcher.exe" -Destination $installPath -Force

    # Add to PATH if not already present
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$installPath*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$installPath", "User")
        Write-Host "Added to PATH: $installPath" -ForegroundColor Green
        Write-Host "Please restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
    }

    Write-Host "Installation complete!" -ForegroundColor Green
    Write-Host "Installed to: $installPath\switcher.exe" -ForegroundColor Cyan
}
