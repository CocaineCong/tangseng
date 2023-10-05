from sentence_transformers import SentenceTransformer


# word文本转向量格式
def word2vec(word):
    model = SentenceTransformer('uer/sbert-base-chinese-nli')
    return model.encode(word)
