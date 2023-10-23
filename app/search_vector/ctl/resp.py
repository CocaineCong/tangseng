import json


def RespDefaultSuccess(data):
    return json.dumps({
        'code':200,
        'doc_ids':data,
        'msg':'ok',
        'error':'',
    })

def RespDefaultError(error):
    return json.dumps({
        'code':500,
        'doc_ids':'',
        'msg':'ok',
        'error':str(error),
    })
