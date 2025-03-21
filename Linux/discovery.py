#!/usr/bin/env python3

import sys
import json
from functions import get_device_type, get_adapter_info

"""
Discovery Adaptec RAID Devices for Zabbix

Version: 1
Author: Chizhov Nikolay
E-Mail: nchizhov@inok.ru
"""

if __name__ == '__main__':
    if len(sys.argv) != 3 or not sys.argv[1].isnumeric():
        sys.stderr.write('ZBX_NOTSUPPORTED\0Incorrect arguments')
        sys.exit(1)
    index = int(sys.argv[1])
    device = sys.argv[2]
    device_type = get_device_type(device)
    if device_type is None:
        sys.stderr.write('ZBX_NOTSUPPORTED\0Device type not found. Correct value: ld, pd')
        sys.exit(1)

    controller_info = get_adapter_info(index)
    if controller_info is None:
        sys.stderr.write('ZBX_NOTSUPPORTED\0Controller device info not found')
        sys.exit(1)

    discovery_info = {'data': []}
    if device_type['name'] in controller_info:
        for device_k, device_v in controller_info[device_type['name']].items():
            if device_type['field'] in device_v:
                discovery_info['data'].append({
                    '{#DEVICEID}': device_k,
                    '{#DEVICENAME}': device_v[device_type['field']]
                })
    print(json.dumps(discovery_info))
