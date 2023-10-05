from utils.encode import word2vec


def word2vec_test():
    res = word2vec("哈哈哈哈")
    print(res)


if __name__ == '__main__':
    word2vec_test()