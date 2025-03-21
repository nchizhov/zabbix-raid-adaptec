#!/usr/bin/env python3

import sys
from functions import get_device_type, get_adapter_info, get_info_data

"""
Description: Get Adaptec RAID Controller Device Info for Zabbix

Version: 1
Author: Chizhov Nikolay
E-Mail: nchizhov@inok.ru
"""

if __name__ == '__main__':
    if len(sys.argv) != 5 or not sys.argv[1].isnumeric():
        sys.stderr.write('ZBX_NOTSUPPORTED\0Incorrect arguments')
        sys.exit(1)
    index = int(sys.argv[1])
    device = sys.argv[2]
    device_index = sys.argv[3]
    device_type = get_device_type(device)
    if device_type is None:
        sys.stderr.write('ZBX_NOTSUPPORTED\0Device type not found. Correct value: ld, pd')
        sys.exit(1)

    controller_info = get_adapter_info(index)
    if (controller_info is None
            or device_type['name'] not in controller_info or device_index not in controller_info[device_type['name']]):
        sys.stderr.write('ZBX_NOTSUPPORTED\0Controller device info not found')
        sys.exit(1)

    param_data = get_info_data(controller_info[device_type['name']][device_index], sys.argv[4])
    if param_data is None:
        sys.stderr.write('ZBX_NOTSUPPORTED\0Param not exists')
        sys.exit(1)
    print(param_data)
