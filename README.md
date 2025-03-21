## Zabbix Monitoring for Adaptec RAID Controllers

### Contains:
- Scripts for Linux (Python 3) and Windows (PowerShell)
- Template for Zabbix 6.XX, 7.XX

### Requirements:
- Python 3.6 and higher for Linux
- PowerShell 3.0 and higher for Windows
- Zabbix agent - any version or Zabbix agent 2 - version 6
- Zabbix version 6.XX and higher

### Install support scripts
#### Windows:
1. Copy ```.ps1``` files from ```Windows``` folder of repository to ```C:\Scripts\RAID```
2. Change in file ```get-info.ps1``` variable ```$arcconf_path``` - path to arcconf.exe for your Adaptec RAID Controller 
3. Import to Windows Scheduler repository file ```Windows\Scheduler.xml```:
   1. Open created task in Windows Scheduler
   2. Open tab ```Actions```
   3. Change action - edit argument and set index of Adaptec Controller in Windows
   4. Save 
- For Zabbix Agent:
  1. Open Zabbix Agent config file
  2. Change ```Timeout``` option to 30
  3. Insert to the end of config file:
     ```
     UserParameter=raid.adaptec.info[*],powershell -NoProfile -File "C:\Scripts\RAID\info.ps1" "$1" "$2"
     UserParameter=raid.adaptec.discovery[*],powershell -NoProfile -File "C:\Scripts\RAID\discovery.ps1" "$1" "$2"
     UserParameter=raid.adaptec.ld[*],powershell -NoProfile -File "C:\Scripts\RAID\device-info.ps1" "$1" ld "$2" "$3"
     UserParameter=raid.adaptec.pd[*],powershell -NoProfile -File "C:\Scripts\RAID\device-info.ps1" "$1" pd "$2" "$3"
     ``` 
- For Zabbix Agent 2:
  1. Copy repository file ```Windows\adaptec.conf``` to Zabbix Agent 2: ```zabbix_agent2.d\plugins.d```
  2. Open Zabbix Agent config file
  3. Change ```Timeout``` option to 30

### Linux:
1. Copy ```*.py``` files from ```Linux``` folder if repository to ```/opt/adaptec```
2. Change in file ```get-info.py``` variable ```arcconf_path``` - path to arcconf for you Adaptec RAID Controller
3. Set ```*.py``` files as execution:
   ```bash
   chmod +x /opt/adaptec/*.py
   ```
4. Add schedule to cron:
   ```
   */2 * * * * /opt/adaptec/get-info.py 1
   ```
   , where 1 - is index of Adaptec Controller in Linux
- For Zabbix Agent:
  1. Open Zabbix Agent config file
  2. Insert to the end of config file:
     ```
     UserParameter=raid.adaptec.info[*],/opt/adaptec/info.py "$1" "$2"
     UserParameter=raid.adaptec.discovery[*],/opt/adaptec/discovery.py "$1" "$2"
     UserParameter=raid.adaptec.ld[*],/opt/adaptec/device-info.py "$1" ld "$2" "$3"
     UserParameter=raid.adaptec.pd[*],/opt/adaptec/device-info.py "$1" pd "$2" "$3"
     ```
- For Zabbix Agent 2:
  1. Copy repository file ```Linux/adaptec.conf``` to Zabbix Agent 2: ```zabbix_agent2.d/plugins.d```

## Zabbix Server
1. Import template ```adaptec_template.xml``` from repository to Zabbix Server Templates
2. Add imported template to host
3. Edit Macros ```{$ADAPTEC.CONTROLLER.INDEX}``` to set current index of Adaptec RAID Controller in Host 