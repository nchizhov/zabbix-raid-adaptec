#!/usr/bin/env python3

import json
import sys
import re
import os
import tempfile
from subprocess import PIPE, DEVNULL, Popen
from typing import Optional, List
from config import arcconf_path

"""
Description: Get Information From Adaptec RAID Controller.
Required external program: arcconf

Version: 1
Author: Chizhov Nikolay
E-Mail: nchizhov@inok.ru
"""

title_split = '----------------------------------------------------------------------'
subtitle_split = '--------------------------------------------------------'
logical_device_re = re.compile('^Logical device number (\d+)$')
physical_device_re = re.compile('^Device #(\d+)$')


def adapter_data(index: int) -> Optional[List[str]]:
    arcconf_process = Popen('{path} GETCONFIG {index} AL'.format(path=arcconf_path, index=index),
                            shell=True, stdout=PIPE, stderr=DEVNULL)
    arcconf_lines = arcconf_process.communicate()[0].decode('ascii').splitlines()
    if arcconf_process.returncode != 0:
        return None
    return arcconf_lines


def parse_data(data: List[str]) -> dict:
    controller_data = {}
    current_section = None
    current_subsection = None
    current_index = None
    for idx, line in enumerate(data):
        line = line.strip()
        if idx < 3:
            continue
        if line == title_split and data[idx - 2].strip() == title_split:
            current_section = data[idx - 1].strip()
            current_subsection = None
            current_index = None
            controller_data[current_section] = {}
            continue
        if line == subtitle_split and data[idx - 2].strip() == subtitle_split:
            current_subsection = data[idx - 1].strip()
            if current_index is None:
                controller_data[current_section][current_subsection] = {}
            else:
                controller_data[current_section][current_index][current_subsection] = {}
            continue
        logical_device = logical_device_re.search(line)
        if logical_device is not None:
            current_index = logical_device[1]
            current_subsection = None
            controller_data[current_section][current_index] = {}
            continue
        physical_device = physical_device_re.search(line)
        if physical_device is not None:
            current_index = physical_device[1]
            current_subsection = None
            controller_data[current_section][current_index] = {
                'Device Info': data[idx + 1].strip()
            }
            continue
        s_line = line.split(': ', 1)
        if len(s_line) == 2:
            field = s_line[0].strip()
            field_val = s_line[1].strip()
            if current_index is None:
                if current_subsection is None:
                    controller_data[current_section][field] = field_val
                else:
                    controller_data[current_section][current_subsection][field] = field_val
            else:
                if current_subsection is None:
                    controller_data[current_section][current_index][field] = field_val
                else:
                    controller_data[current_section][current_index][current_subsection][field] = field_val
    return controller_data


def save_info(index: int, data: dict) -> None:
    info_path = os.path.join(tempfile.gettempdir(), 'adaptec-{index}.json'.format(index=index))
    if os.path.exists(info_path):
        os.remove(info_path)
    with open(info_path, 'w') as f:
        json.dump(data, f, ensure_ascii=False)


if __name__ == '__main__':
    if len(sys.argv) != 2 or not sys.argv[1].isnumeric():
        sys.exit(1)
    index = int(sys.argv[1])
    data = adapter_data(index)
    if data is None:
        sys.exit(1)
    parsed_data = parse_data(data)
    save_info(index, parsed_data)
