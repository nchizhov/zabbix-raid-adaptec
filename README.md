## Zabbix Monitoring for Adaptec RAID Controllers

### Contains:
- Utility for get adapter info, data for Zabbix, self-update 
- Template for Zabbix 7.XX

### Requirements:
- Zabbix agent - any version or Zabbix agent 2 - version 6
- Zabbix version 7.XX and higher

### Install support scripts
#### Windows:
1. Download latest release asset https://github.com/nchizhov/zabbix-raid-adaptec/releases/latest to folder ```C:\Scripts\RAID```:
   - For ```x64```-arch postfix of filename is ```windows-amd64.exe```
   - For ```x86```-arch postfix of filename is ```windows-386.exe```
2. Rename downloaded file to ```adapter.exe```  
3. Download ```config.json``` file from https://raw.githubusercontent.com/nchizhov/zabbix-raid-adaptec/refs/heads/master/config.json and save it to ```C:\Scripts\RAID```
4. Change in file ```config.json``` value of ```arcconf_path``` to path to arcconf.exe for your Adaptec RAID Controller 
5. Import to Windows Scheduler repository file ```Windows\Scheduler.xml```:
   1. Open created task in Windows Scheduler
   2. Open tab ```Actions```
   3. Change action - edit argument and set index of Adaptec Controller in Windows
   4. Save 
- For Zabbix Agent:
  1. Open Zabbix Agent config file
  2. Change ```Timeout``` option to 30
  3. Insert to the end of config file:
     ```
     UserParameter=raid.adaptec.info[*],"C:\Scripts\RAID\adapter.exe" -info -index $1 -field "$2"
     UserParameter=raid.adaptec.discovery[*],"C:\Scripts\RAID\adapter.exe" -discovery -index $1 -$2
     UserParameter=raid.adaptec.ld[*],"C:\Scripts\RAID\adapter.exe" -info -index $1 -ld -drive-index "$2" -field "$3"
     UserParameter=raid.adaptec.pd[*],"C:\Scripts\RAID\adapter.exe" -info -index $1 -pd -drive-index "$2" -field "$3"
     ``` 
- For Zabbix Agent 2:
  1. Copy repository file ```Windows\adaptec.conf``` to Zabbix Agent 2: ```zabbix_agent2.d\plugins.d```
  2. Open Zabbix Agent config file
  3. Change ```Timeout``` option to 30

### Linux:
1. Download latest release asset https://github.com/nchizhov/zabbix-raid-adaptec/releases/latest to folder ```/opt/adaptec```:
    - For ```x64```-arch postfix of filename is ```linux-amd64```
    - For ```x86```-arch postfix of filename is ```linux-386```
2. Rename downloaded file to ```adapter```
3. Set ```adapter``` file execution flag:
   ```bash
   chmod +x /opt/adaptec/adapter 
   ```
4. Download ```config.json``` file from https://raw.githubusercontent.com/nchizhov/zabbix-raid-adaptec/refs/heads/master/config.json and save it to ```/opt/adaptec```
5. Change in file ```config.json``` value of ```arcconf_path``` to path to arcconf.exe for your Adaptec RAID Controller
6. Add schedule to cron:
   ```
   */2 * * * * /opt/adaptec/adapter -get-info -index 1
   ```
   , where 1 - is index of Adaptec Controller in Linux
- For Zabbix Agent:
  1. Open Zabbix Agent config file
  2. Insert to the end of config file:
     ```
     UserParameter=raid.adaptec.info[*],/opt/adaptec/adapter -info -index $1 -field $2
     UserParameter=raid.adaptec.discovery[*],/opt/adaptec/adapter -discovery -index $1 -$2
     UserParameter=raid.adaptec.ld[*],/opt/adaptec/adapter -info -index $1 -ld -drive-index "$2" -field "$3"
     UserParameter=raid.adaptec.pd[*],/opt/adaptec/adapter -info -index $1 -pd -drive-index "$2" -field "$3"
     ```
- For Zabbix Agent 2:
  1. Copy repository file ```Linux/adaptec.conf``` to Zabbix Agent 2: ```zabbix_agent2.d/plugins.d```

## Updates
### Windows
1. Open PowerShell console
2. Change current dir to ```C:\Scripts\RAID```
3. Run ```adapter.exe -update```
### Linux
1. Open console
2. Change current dir to ```/opt/adaptec```
3. Run ```adapter -update```

## Zabbix Server
1. Import template ```adaptec_template.xml``` from repository to Zabbix Server Templates
2. Add imported template to host
3. Edit Macros ```{$ADAPTEC.CONTROLLER.INDEX}``` to set current index of Adaptec RAID Controller in Host 