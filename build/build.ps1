[CmdletBinding()]
Param (
    [Parameter(Mandatory = $true)]
    [String] $PathToExecutable,
    [Parameter(Mandatory = $true)]
    [String] $Version,
    [Parameter(Mandatory = $false)]
    [ValidateSet("amd64", "386")]
    [String] $Arch = "amd64"
)
$ErrorActionPreference = "Stop"

# Get absolute path to executable before switching directories
$PathToExecutable = Resolve-Path $PathToExecutable
# Set working dir to this directory, reset previous on exit
Push-Location $PSScriptRoot
Trap {
    # Reset working dir on error
    Pop-Location
}

if ($PSVersionTable.PSVersion.Major -lt 5) {
    Write-Error "Powershell version 5 required"
    exit 1
}

$wc = New-Object System.Net.WebClient
function Get-FileIfNotExists {
    Param (
        $Url,
        $Destination
    )
    if (-not (Test-Path $Destination)) {
        Write-Verbose "Downloading $Url"
        $wc.DownloadFile($Url, $Destination)
    }
    else {
        Write-Verbose "${Destination} already exists. Skipping."
    }
}
$sourceDir = mkdir -Force Source
mkdir -Force Work, Output | Out-Null

Write-Verbose "Downloading WiX..."
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Get-FileIfNotExists "https://github.com/wixtoolset/wix3/releases/download/wix311rtm/wix311-binaries.zip" "$sourceDir\wix-binaries.zip"
mkdir -Force WiX | Out-Null

Expand-Archive -Path "${sourceDir}\wix-binaries.zip" -DestinationPath WiX -Force

Write-Verbose "Creating prom2lyrid-${Version}-${Arch}.msi"
$wixOpts = "-ext WixFirewallExtension -ext WixUtilExtension"

Invoke-Expression "WiX\heat.exe dir ..\..\web\build -nologo -cg AppFiles0 -gg -g1 -srd -sfrag -template fragment -dr APPDIR0 -var var.SourceDir0 -out AppFiles0.wxs"
Invoke-Expression "WiX\candle.exe $wixOpts -dSourceDir0=..\web\build AppFiles0.wxs LicenseAgreementDlg_HK.wxs WixUI_HK.wxs product.wxs"
Invoke-Expression "WiX\light.exe $wixOpts -sacl -spdb  -out ..\wix\prom2lyrid.msi AppFiles0.wixobj LicenseAgreementDlg_HK.wixobj WixUI_HK.wixobj product.wixobj"
