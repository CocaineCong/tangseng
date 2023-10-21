
def RespDefaultSuccess(data):
    return {
        'code':200,
        'data':data,
        'msg':'ok',
        'error':'',
    }

def RespDefaultError(error):
    return {
        'code':500,
        'data':'',
        'msg':'ok',
        'error':str(error),
    }
