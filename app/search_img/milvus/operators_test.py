from milvus.milvus import Milvus
from milvus.operators import do_upload


def do_upload_test():
    client = Milvus()
    ids = do_upload("test", 1, "title_test",
                    "test something like test", client)
    print(ids)


if __name__ == '__main__':
    do_upload_test()
