#!/usr/bin/env python3

import requests
import hashlib
import base64
from time import sleep
from sys import exit
from typing import Optional
from json import loads, JSONDecodeError

"""
Description: Updater for Adaptec RAID Scripts

Version: 1
Author: Chizhov Nikolay
E-Mail: nchizhov@inok.ru
"""

update_url = "https://api.github.com/repos/nchizhov/zabbix-raid-adaptec/contents/Linux/{file}?ref=master"
update_files = ("device-info.py", "discovery.py", "functions.py", "get-info.py", "info.py")

def get_server_file_info(file_name: str) -> Optional[dict]:
    update_file_url = update_url.format(file=file_name)
    try:
        request = requests.get(update_file_url, timeout=10)
        request.raise_for_status()
    except (requests.exceptions.HTTPError, requests.exceptions.ReadTimeout, requests.exceptions.ConnectionError,
            requests.exceptions.RequestException) as e:
        show_file_error(file_name, str(e))
        return None
    try:
        data = loads(request.content)
    except JSONDecodeError as e:
        show_file_error(file_name, str(e))
        return None
    return data

def calculate_file_sha_sum(data: str) -> str:
    sha_data = "blob {len}\0{data}".format(len=len(data), data=data)
    return hashlib.sha1(sha_data.encode()).hexdigest().lower()

def get_data_from_base64(data: str) -> str:
    return base64.b64decode(data).decode()

def rollback_update(files: dict) -> None:
    if len(files) == 0:
        return None
    print("Updates failed! Rollback:")
    for file, data in files.items():
        is_rollback = False
        while not is_rollback:
            try:
                with open(file, 'w', newline='\n') as f:
                    f.write(data)
            except:
                sleep(1)
                continue
            is_rollback = True
        print("{file} - rolled back".format(file=file))

def show_file_error(file_name: str, reason: str) -> None:
    print("Error update file {file}. Reason {reason}".format(file=file_name, reason=reason))


rollback_files = {}
for update_file in update_files:
    data = get_server_file_info(update_file)
    if data is None:
        rollback_update(rollback_files)
        exit(1)
    try:
        with open(update_file, 'r') as f:
            file_data = f.read()
    except Exception as e:
        show_file_error(update_file, str(e))
        rollback_update(rollback_files)
        exit(1)
    data_sha = calculate_file_sha_sum(file_data)
    if data_sha == data['sha']:
        print("{file} - up to date".format(file=update_file))
        continue
    update_data = get_data_from_base64(data['content'])
    try:
        with open(update_file, 'w', newline='\n') as f:
            f.write(update_data)
    except Exception as e:
        show_file_error(update_file, str(e))
        rollback_update(rollback_files)
        exit(1)
    rollback_files[update_file] = file_data
    print("{file}- updated".format(file=update_file))
print("Update finished")
print("Don't forget to update zabbix template from: https://github.com/nchizhov/zabbix-raid-adaptec/blob/master/adaptec_template.xml")