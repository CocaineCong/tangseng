from sentence_transformers import SentenceTransformer

from config.config import TRANSFORMER_MODEL_NAME


# word文本转向量格式
def word2vec(word):
    model = SentenceTransformer(TRANSFORMER_MODEL_NAME)
    return model.encode(word)
