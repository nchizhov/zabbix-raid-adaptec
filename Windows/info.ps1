<#
    .SYNOPSIS
    Get Adaptec RAID Controller Info

    .DESCRIPTION
    Get Adaptec RAID Controller Info for Zabbix

    .PARAMETER index
    RAID Controller Device ID

    .PARAMETER param
    RAID Controller Info required param. 
    Allow "." for subfields

    .OUTPUTS
    None - if any errors.
    String - required parameter value

    .EXAMPLE
    > .\info.ps1 1 "Controller Status"

    .EXAMPLE
    > .\info.ps1 1 "Controller Version Information.BIOS"

    .NOTES
        Version: 1
        Author: Chizhov Nikolay
        E-Mail: nchizhov@inok.ru
#>

param(
    [Parameter(Mandatory=$true)]
    [int] $index,
    [Parameter(Mandatory=$true)]
    [string] $param
)

Set-Location $PSScriptRoot

. .\functions.ps1

function Get-Info {
    [OutputType([string], [int], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [System.Collections.Hashtable] $data,
        [Parameter(Mandatory=$true)]
        [string] $param
    )
    if (-not ($data.ContainsKey('Controller information'))) {
        return $null
    }
    return Get-Info-Data -data $data['Controller information'] -param $param
}

$controller_info = Get-Adapter-Info -index $index
if ($null -eq $controller_info) {
    Write-Host "ZBX_NOTSUPPORTED`0Controller info not found"
    Exit 1
}

$param_data = Get-Info -data $controller_info -param $param
if ($null -eq $param_data) {
    Write-Host "ZBX_NOTSUPPORTED`0Param not exists"
    Exit 1
}
Write-Host $param_data