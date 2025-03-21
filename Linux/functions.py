from typing import Optional
from os import path
import json
import tempfile

devices = {
    'ld': {
        'name': 'Logical device information',
        'field': 'Logical device name'
    },
    'pd': {
        'name': 'Physical Device information',
        'field': 'Reported Location'
    }
}

def get_device_type(device: str) -> Optional[dict]:
    return devices[device] if device in devices else None

def get_adapter_info(index: int) -> Optional[dict]:
    info_path = path.join(tempfile.gettempdir(), 'adaptec-{index}.json'.format(index=index))
    if not path.exists(info_path):
        return None
    try:
        with open(info_path, 'r') as f:
            data = json.load(f)
        return data
    except json.JSONDecodeError:
        return None

def get_info_data(data: dict, param: str) -> Optional[str | int]:
    s_param = param.split('.')
    if len(s_param) == 1:
        if check_info_data(data, s_param[0]):
            return data[s_param[0]]
        return None
    s_param_length = len(s_param)
    main_param = ''
    last_idx = None
    for idx, t_param in enumerate(s_param):
        main_param = "".join((main_param, t_param))
        if main_param in data:
            last_idx = idx
            break
        main_param = "".join((main_param, '.'))
        if main_param in data:
            last_idx = idx
    if last_idx is None:
        return None
    last_idx += 1
    if last_idx == s_param_length:
        return data[main_param] if check_info_data(data, main_param) else None
    if isinstance(data[main_param], dict):
        second_param = ".".join(s_param[last_idx:])
        return data[main_param][second_param] if check_info_data(data[main_param], second_param) else None
    return None

def check_info_data(info, param: str) -> bool:
    return param in info and isinstance(info[param], (str, int, float))