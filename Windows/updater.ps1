<#
    .SYNOPSIS
    Updater for Adaptec RAID Scripts

    .DESCRIPTION
    Updater for Adaptec RAID Scripts

    .NOTES
        Version: 1
        Author: Chizhov Nikolay
        E-Mail: nchizhov@inok.ru
#>

$update_url = "https://api.github.com/repos/nchizhov/zabbix-raid-adaptec/contents/Windows/{0}?ref=master"
$update_files = @("device-info.ps1", "discovery.ps1", "functions.ps1", "get-info.ps1", "info.ps1")

function Get-Server-File-Info {
    [OutputType([System.Management.Automation.PSCustomObject], $null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [string] $file_name
    )
    $update_file_url = [string]::Format($update_url, $file_name)
    try {
        $request = Invoke-WebRequest -UseBasicParsing -URI $update_file_url -TimeoutSec 10
    } catch {
        [void](Show-File-Error -file_name $file_name -reason $_)
        return $null
    }
    try {
        $data = $request.Content | ConvertFrom-Json
    } catch {
        [void](Show-File-Error -file_name $file_name -reason $_)
        return $null
    }
    return $data
}

function Calculate-File-Sha-Sum {
    [OutputType([string])]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [string] $data
    )
    $stringAsStream = [System.IO.MemoryStream]::new()
    $writer = [System.IO.StreamWriter]::new($stringAsStream)
    $writer.Write([string]::Format("blob {0}`0{1}", $data.Length, $data))
    $writer.Flush()
    $stringAsStream.Position = 0
    return (Get-FileHash -InputStream $stringAsStream -Algorithm SHA1).Hash.ToLower()
}

function Get-Data-From-Base64 {
    [OutputType([string])]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [string] $data
    )
    return [Text.Encoding]::Utf8.GetString([Convert]::FromBase64String($data))
}

function Rollback-Update {
    [OutputType($null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [System.Collections.Hashtable] $files
    )
    if ($files.Count -eq 0) {
        return
    }
    Write-Host "Updates failed! Rollback:"
    foreach ($file in $files.GetEnumerator()) {
        $is_rollback = $false
        while ($false -eq $is_rollback) {
            try {
                Set-Content -Path $file.Key -Value $file.Value -NoNewline -Encoding UTF8 -Force -ErrorAction Stop
            } catch {
                Start-Sleep -Seconds 1
                continue
            }
            $is_rollback = $true
        }
        $rollback_message = [string]::Format("{0} - rolled back")
        Write-Host $rollback_message
    }
}

function Show-File-Error {
    [OutputType($null)]
    [CmdletBinding(PositionalBinding=$false)]
    Param(
        [Parameter(Mandatory=$true)]
        [string] $file_name,
        [Parameter(Mandatory=$true)]
        $reason
    )
    $error_message = [string]::Format("Error update file {0}. Reason: {1}", $file_name, $reason)
    Write-Host $error_message -ForegroundColor Red
}

$rollback_files = @{}
foreach ($update_file in $update_files) {
    $data = Get-Server-File-Info -file $update_file
    if ($null -eq $data) {
        Rollback-Update -files $rollback_files
        Exit 1
    }
    try {
        $file_data = Get-Content -Raw -Path $update_file -ErrorAction Stop
    } catch {
        [void](Show-File-Error -file_name $update_file -reason $_)
        Rollback-Update -files $rollback_files
        Exit 1
    }
    $data_sha = Calculate-File-Sha-Sum -data $file_data
    if ($data_sha -eq $data.sha) {
        $message = [string]::Format("{0} - up to date", $update_file)
        Write-Host $message
        continue
    }
    $update_data = Get-Data-From-Base64 -data $data.content
    try {
        Set-Content -Path $update_file -Value $update_data -NoNewline -Encoding UTF8 -Force -ErrorAction Stop
    } catch {
        [void](Show-File-Error -file_name $update_file -reason $_)
        Rollback-Update -files $rollback_files
        Exit 1
    }
    $rollback_files[$update_file] = $file_data
    $updated_message = [string]::Format("{0} - updated", $update_file)
    Write-Host $updated_message
}
Write-Host "Update finished"
Write-Host "Don't forget to update zabbix template from: https://github.com/nchizhov/zabbix-raid-adaptec/blob/master/adaptec_template.xml"
