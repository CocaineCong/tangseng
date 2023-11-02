"""
default the struct to response
"""
import json


def resp_default_success(data):
    """
    default the struct to response success request
    
    params: data
    return: json string of resp data
    """
    resp = json.dumps({
        'code': 200,
        'doc_ids': data,
        'msg': 'ok',
        'error': '',
    })
    return resp


def resp_default_error(error):
    """
    default the struct to response error request
    
    params: error
    return: json string of resp data
    """
    return json.dumps({
        'code': 500,
        'doc_ids': '',
        'msg': 'ok',
        'error': str(error),
    })
