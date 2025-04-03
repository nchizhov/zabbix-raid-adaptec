#!/usr/bin/env python3

import sys
from typing import Optional
from functions import get_info_data, get_adapter_info

"""
Description: Get Adaptec RAID Controller Info for Zabbix

Version: 1
Author: Chizhov Nikolay
E-Mail: nchizhov@inok.ru
"""

def get_info(data: dict, param: str) -> Optional[str | int]:
    if 'Controller information' not in data:
        return None
    return get_info_data(data['Controller information'], param)

if __name__ == '__main__':
    if len(sys.argv) != 3 or not sys.argv[1].isnumeric():
        sys.stderr.write('ZBX_NOTSUPPORTED\0Incorrect arguments')
        sys.exit(1)
    controller_info = get_adapter_info(int(sys.argv[1]))
    if controller_info is None:
        sys.stderr.write('ZBX_NOTSUPPORTED\0Controller info not found')
        sys.exit(1)
    param_data = get_info(controller_info, sys.argv[2])
    if param_data is None:
        sys.stderr.write('ZBX_NOTSUPPORTED\0Param not exists')
        sys.exit(1)
    print(param_data)
