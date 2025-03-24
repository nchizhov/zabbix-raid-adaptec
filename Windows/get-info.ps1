<#
    .SYNOPSIS
    Get Information From Adaptec RAID Controller

    .DESCRIPTION
    Get Information From Adaptec RAID Controller.
    Required external program: arcconf

    .PARAMETER index
    RAID Controller Device ID

    .OUTPUTS
    None. Save RAID Controller information to tmp-folder in JSON-format

    .NOTES
        Version: 1
        Author: Chizhov Nikolay
        E-Mail: nchizhov@inok.ru
#>

param(
    [Parameter(Mandatory=$true)]
    [int] $index
)

$arcconf_path = 'C:\Scripts\RAID\arcconf.exe'
$title_split = '----------------------------------------------------------------------'
$subtitle_split = '--------------------------------------------------------'
$logical_device_re = "^Logical device number (\d+)$"
$physical_device_re = "^Device #(\d+)$"

function Adapter-Data {
    [OutputType([string[]], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [int] $index
    )
    $pinfo = New-Object System.Diagnostics.ProcessStartInfo
    $pinfo.FileName = $arcconf_path
    $pinfo.RedirectStandardOutput = $true
    $pinfo.UseShellExecute = $false
    $pinfo.Arguments = "GETCONFIG",$index,"AL"
    $arcconf_process = New-Object System.Diagnostics.Process
    $arcconf_process.StartInfo = $pinfo
    [void]$arcconf_process.Start()
    $lines = $arcconf_process.StandardOutput.ReadToEnd().Split([Environment]::NewLine,[StringSplitOptions]::RemoveEmptyEntries)
    $arcconf_process.WaitForExit()
    if ($arcconf_process.ExitCode -ne 0) {
        return $null
    }
    return $lines
}

function Parse-Data {
    [OutputType([System.Collections.Hashtable], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [System.Object[]] $data
    )
    $controller_data = @{}
    $current_section = $null
    $current_subsection = $null
    $current_index = $null
    for ($idx = 0; $idx -lt $data.Count; $idx++) {
        if ($idx -lt 3) {
            continue
        }
        $line = $data[$idx].Trim()
        if (($line -eq $title_split) -and ($data[$idx-2].Trim() -eq $title_split)) {
            $current_section = $data[$idx-1].Trim()
            $current_subsection = $null
            $current_index = $null
            $controller_data[$current_section] = @{}
            continue
        }
        if (($line -eq $subtitle_split) -and ($data[$idx-2].Trim() -eq $subtitle_split)) {
            $current_subsection = $data[$idx-1].Trim()
            if ($null -eq $current_index) {
                $controller_data[$current_section][$current_subsection] = @{}
            } else {
                $controller_data[$current_section][$current_index][$current_subsection] = @{}
            }
            continue
        }
        if ($line -match $logical_device_re) {
            $current_index = $Matches.1
            $current_subsection = $null
            $controller_data[$current_section][$current_index] = @{}
            continue
        }
        if ($line -match $physical_device_re) {
            $current_index = $Matches.1
            $current_subsection = $null
            $controller_data[$current_section][$current_index] = @{'Device Info' = $data[$idx+1].Trim()}
            continue
        }
        $s_line = $line -split ": ",2
        if ($s_line.Count -eq 2) {
            $field = $s_line[0].Trim()
            $field_val = $s_line[1].Trim()
            if ($null -eq $current_index) {
                if ($null -eq $current_subsection) {
                    $controller_data[$current_section][$field] = $field_val
                } else {
                    $controller_data[$current_section][$current_subsection][$field] = $field_val
                }
            } else {
                if ($null -eq $current_subsection) {
                    $controller_data[$current_section][$current_index][$field] = $field_val
                } else {
                    $controller_data[$current_section][$current_index][$current_subsection][$field] = $field_val
                }
            }
        }
    }
    return $controller_data
}

function Save-Info {
    [OutputType([System.Void], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [int] $index,
        [Parameter(Mandatory=$true)]
        [System.Collections.Hashtable] $data
    )
    $info_path = [System.IO.Path]::Combine($env:TEMP, [string]::Format('adaptec-{0}.json', $index))
    if (Test-Path $info_path -PathType Leaf) {
        [void](Remove-Item $info_path -Force)
    }
    $data | ConvertTo-Json -Depth 10 | Out-File $info_path
}

$data = Adapter-Data -index $index
if ($null -eq $data) {
    Exit 1
}
$parsed_data = Parse-Data -data $data
Save-Info -index $index -data $parsed_data
