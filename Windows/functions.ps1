<#
    .SYNOPSIS
    Additional functions

    .DESCRIPTION
    Additional functions

    .NOTES
        Version: 1
        Author: Chizhov Nikolay
        E-Mail: nchizhov@inok.ru
#>


$devices = @{'ld' = @{'name' = 'Logical device information';
                      'field' = 'Logical device name'}; 
             'pd' = @{'name' = 'Physical Device information';
                      'field' = 'Reported Location'}}

function Get-Device-Type {
    [OutputType([System.Collections.Hashtable], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [string] $type
    )
    if ($devices.ContainsKey($type)) {
        return $devices[$type]
    }
    return $null
}

function ConvertPSObjectToHashtable {
    Param (
        [Parameter(ValueFromPipeline)]
        $InputObject
    )

    process {
        if ($null -eq $InputObject) { return $null }
        if ($InputObject -is [System.Collections.IEnumerable] -and $InputObject -isnot [string]) {
            $collection = @(
                foreach ($object in $InputObject) { ConvertPSObjectToHashtable $object }
            )
            Write-Output -NoEnumerate $collection
        } elseif ($InputObject -is [psobject]) {
            $hash = @{}
            foreach ($property in $InputObject.PSObject.Properties) {
                $hash[$property.Name] = (ConvertPSObjectToHashtable $property.Value).PSObject.BaseObject
            }
            $hash
        } else {
            $InputObject
        }
    }
}

function Get-Adapter-Info {
    [OutputType([System.Collections.Hashtable], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [int] $index
    )
    $info_path = [System.IO.Path]::Combine($env:TEMP, [string]::Format('adaptec-{0}.json', $index))
    if (-not (Test-Path $info_path -PathType Leaf)) {
        return $null
    }
    $retries = 0
    while ($retries -le 4) {
        $retries++
        try {
            $data = (Get-Content -Raw -Path $info_path | ConvertFrom-Json | ConvertPSObjectToHashtable)
        } catch {
            $data = $null
            Start-Sleep -Seconds 1
        }
        if ($null -eq $data) {
            break
        }
    }
    return $data
}

function Get-Info-Data {
    [OutputType([string], [int], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [System.Collections.Hashtable] $data,
        [Parameter(Mandatory=$true)]
        [string] $param
    )
    $s_param = $param.Split('.')
    if ($s_param.Count -eq 1) {
        if (Check-Info-Data -info $data -param $s_param[0]) {
            return $data[$s_param[0]]
        }
        return $null
    }
    $s_param_length = $s_param.Count
    $main_param = ''
    $last_idx = $null
    for ($idx = 0; $idx -lt $s_param_length; $idx++) {
        $main_param += $s_param[$idx]
        if ($data.ContainsKey($main_param)) {
            $last_idx = $idx
            break
        }
        $main_param += '.'
        if ($data.ContainsKey($main_param)) {
            $last_idx = $idx
        }
    }
    if ($null -eq $last_idx) {
        return $null
    }
    $last_idx = ($last_idx + 1)
    if ($last_idx -eq $s_param_length) {
        if (Check-Info-Data -info $data -param $main_param) {
            return $data[$param]
        }
        return $null
    }
    if ($data[$main_param] -is [System.Collections.Hashtable]) {
        $second_param = $s_param[$last_idx..$s_param_length] -join '.'
        if (Check-Info-Data -info $data[$main_param] -param $second_param) {
            return $data[$main_param][$second_param]
        }
    }
    return $null
}

function Check-Info-Data {
    [OutputType([bool])]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [System.Collections.Hashtable] $info,
        [Parameter(Mandatory=$true)]
        [string] $param
    )
    return $info.ContainsKey($param) -and (($info[$param] -is [int]) -or ($info[$param] -is [float]) -or ($info[$param] -is [string]))
}