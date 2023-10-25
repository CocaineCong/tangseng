import sys
import etcd3

from ..utils.logs import LOGGER
from ..config.config import ETCD_HOST,ETCD_PORT


class Etcd_Client:

    def __init__(self, host=ETCD_HOST, port=ETCD_PORT) -> None:
        try:
            self.client = etcd3.client(host=host, port=int(port))
        except Exception as e:
            LOGGER.error(f"init etcd {e}")
            sys.exit(1)
    
    def get(self, key):
        try:
            return self.client.get(key)
        except Exception as e:
            LOGGER.error(f"failed to get key from etcd {e}")
            sys.exit(1)

    def set(self, key, value):
        try:
            return self.client.put(key, value)
        except Exception as e:
            LOGGER.error(f"failed to set key to etcd {e}")
            sys.exit(1)


etcd_client = Etcd_Client(ETCD_HOST, ETCD_PORT)
    