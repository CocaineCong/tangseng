# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

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
        'msg': 'failed',
        'error': str(error),
    })
