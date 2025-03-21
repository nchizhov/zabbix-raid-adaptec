<#
    .SYNOPSIS
    Get Adaptec RAID Controller Device Info

    .DESCRIPTION
    Get Adaptec RAID Controller Device Info for Zabbix

    .PARAMETER index
    RAID Controller Device ID

    .PARAMETER type
    RAID Controller Device type:
      - pd - physical devices
      - ld - logical devices

    .PARAMETER device_index
    Physical/Logical RAID Controller Device Index

    .PARAMETER param
    Physical/Logical RAID Controller Device required param. 
    Allow "." for subfields

    .OUTPUTS
    None - if any errors.
    String - required parameter value

    .EXAMPLE
    > .\info-device.ps1 1 ld 0 "Size"

    .EXAMPLE
    > .\info-device.ps1 1 ld 0 "Logical device segment information.Segment 0"

    .NOTES
        Version: 1
        Author: Chizhov Nikolay
        E-Mail: nchizhov@inok.ru
#>


param(
    [Parameter(Mandatory=$true)]
    [int] $index,
    [Parameter(Mandatory=$true)]
    [string] $type,
    [Parameter(Mandatory=$true)]
    [string] $device_index,
    [Parameter(Mandatory=$true)]
    [string] $param
)

Set-Location $PSScriptRoot

. .\functions.ps1

$device_type = Get-Device-Type -type $type
if ($null -eq $device_type) {
    Write-Host "ZBX_NOTSUPPORTED`0Device type not found. Correct value: ld, pd"
    Exit 1
}

$controller_info = Get-Adapter-Info -index $index
if (($null -eq $controller_info) -or (-not ($controller_info.ContainsKey($device_type['name'])) -or (-not ($controller_info[$device_type['name']].ContainsKey($device_index))))) {
    Write-Host "ZBX_NOTSUPPORTED`0Controller device info not found"
    Exit 1
}

$param_data = Get-Info-Data -data $controller_info[$device_type['name']][$device_index] -param $param
if ($null -eq $param_data) {
    Write-Host "ZBX_NOTSUPPORTED`0Param not exists"
    Exit 1
}
Write-Host $param_data
