<#
    .SYNOPSIS
    Discovery Adaptec RAID Devices

    .DESCRIPTION
    Discovery Adaptec RAID Devices for Zabbix

    .PARAMETER index
    RAID Controller Device ID

    .PARAMETER type
    RAID Controller Device type:
      - pd - physical devices
      - ld - logical devices

    .OUTPUTS
    JSON string for Zabbix

    .NOTES
        Version: 1
        Author: Chizhov Nikolay
        E-Mail: nchizhov@inok.ru
#>

param(
    [Parameter(Mandatory=$true)]
    [int] $index,
    [Parameter(Mandatory=$true)]
    [string] $type
)

Set-Location $PSScriptRoot

. .\functions.ps1

$device_type = Get-Device-Type -type $type
if ($null -eq $device_type) {
    Write-Host "ZBX_NOTSUPPORTED`0Device type not found. Correct value: ld, pd"
    Exit 1
}

$controller_info = Get-Adapter-Info -index $index
if ($null -eq $controller_info) {
    Write-Host "ZBX_NOTSUPPORTED`0Controller device info not found"
    Exit 1
}

$discovery_info = @{'data' = @()}
if ($controller_info.ContainsKey($device_type['name'])) {
    Foreach ($device in $controller_info[$device_type['name']].GetEnumerator()) {
        if ($device.Value.ContainsKey($device_type['field'])) {
            $discovery_info['data'] += @{'{#DEVICEID}' = $device.Key; '{#DEVICENAME}' = $device.Value[$device_type['field']]}
        }
    }
}
$discovery_info_json = ConvertTo-Json $discovery_info -Depth 3
Write-Host $discovery_info_json
