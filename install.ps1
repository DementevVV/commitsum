$ErrorActionPreference = "Stop"

$Repo = "DementevVV/commitsum"
$App = "commitsum"

function Fail($msg) {
  Write-Error $msg
  exit 1
}

$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
  "AMD64" { $Arch = "amd64" }
  "ARM64" { $Arch = "arm64" }
  default { Fail "Unsupported architecture: $Arch" }
}

$Tag = (Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest").tag_name
if (-not $Tag) { Fail "Unable to determine latest release tag" }

$Asset = "$App-windows-$Arch.exe.zip"
$BaseUrl = "https://github.com/$Repo/releases/download/$Tag"

$TempDir = New-Item -ItemType Directory -Force -Path ([IO.Path]::Combine([IO.Path]::GetTempPath(), "commitsum-install"))
$Archive = Join-Path $TempDir $Asset
$Checksums = Join-Path $TempDir "checksums.txt"

Write-Host "Downloading $Asset for windows/$Arch ($Tag)..."
Invoke-WebRequest -Uri "$BaseUrl/$Asset" -OutFile $Archive
Invoke-WebRequest -Uri "$BaseUrl/checksums.txt" -OutFile $Checksums

$Expected = (Select-String -Path $Checksums -Pattern " $Asset$").Line.Split(" ")[0]
if (-not $Expected) { Fail "Checksum not found for $Asset" }

$Actual = (Get-FileHash -Algorithm SHA256 $Archive).Hash.ToLower()
if ($Expected.ToLower() -ne $Actual) { Fail "Checksum mismatch for $Asset" }

Expand-Archive -Path $Archive -DestinationPath $TempDir -Force
$Binary = Join-Path $TempDir "$App-windows-$Arch.exe"
if (-not (Test-Path $Binary)) { Fail "Binary not found in archive" }

$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\commitsum"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$InstallPath = Join-Path $InstallDir "$App.exe"
Move-Item -Force $Binary $InstallPath

Write-Host "Installed $App to $InstallPath"
Write-Host "Add to PATH:"
Write-Host "  `$env:Path += `";$InstallDir`""
