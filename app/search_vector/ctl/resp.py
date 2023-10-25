import json


def RespDefaultSuccess(data):
    resp = json.dumps({
        'code':200,
        'doc_ids':data,
        'msg':'ok',
        'error':'',
    })
    print(resp)
    return resp

def RespDefaultError(error):
    return json.dumps({
        'code':500,
        'doc_ids':'',
        'msg':'ok',
        'error':str(error),
    })
