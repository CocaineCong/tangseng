import json


def RespDefaultSuccess(data):
    return json.dumps({
        'code':200,
        'data':data,
        'msg':'ok',
        'error':'',
    })

def RespDefaultError(error):
    return json.dumps({
        'code':500,
        'data':'',
        'msg':'ok',
        'error':str(error),
    })
