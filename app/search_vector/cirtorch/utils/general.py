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

import os
import hashlib

def get_root():
    return os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))


def get_data_root():
    return os.path.join(get_root(), 'data')


def htime(c):
    c = round(c)
    
    days = c // 86400
    hours = c // 3600 % 24
    minutes = c // 60 % 60
    seconds = c % 60

    if days > 0:
        return '{:d}d {:d}h {:d}m {:d}s'.format(days, hours, minutes, seconds)
    if hours > 0:
        return '{:d}h {:d}m {:d}s'.format(hours, minutes, seconds)
    if minutes > 0:
        return '{:d}m {:d}s'.format(minutes, seconds)
    return '{:d}s'.format(seconds)


def sha256_hash(filename, block_size=65536, length=8):
    sha256 = hashlib.sha256()
    with open(filename, 'rb') as f:
        for block in iter(lambda: f.read(block_size), b''):
            sha256.update(block)
    return sha256.hexdigest()[:length-1]
