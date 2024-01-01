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
import numpy as np

from app.search_vector.utils.logs import LOGGER


def whitenapply(X, m, P, dimensions=None):
    
    if not dimensions:
        dimensions = P.shape[0]

    X = np.dot(P[:dimensions, :], X-m)
    X = X / (np.linalg.norm(X, ord=2, axis=0, keepdims=True) + 1e-6)

    return X

def pcawhitenlearn(X):

    N = X.shape[1]

    # Learning PCA w/o annotations
    m = X.mean(axis=1, keepdims=True)
    Xc = X - m
    Xcov = np.dot(Xc, Xc.T)
    Xcov = (Xcov + Xcov.T) / (2*N)
    eigval, eigvec = np.linalg.eig(Xcov)
    order = eigval.argsort()[::-1]
    eigval = eigval[order]
    eigvec = eigvec[:, order]

    P = np.dot(np.linalg.inv(np.sqrt(np.diag(eigval))), eigvec.T)
    
    return m, P

def whitenlearn(X, qidxs, pidxs):

    # Learning Lw w annotations
    m = X[:, qidxs].mean(axis=1, keepdims=True)
    df = X[:, qidxs] - X[:, pidxs]
    S = np.dot(df, df.T) / df.shape[1]
    P = np.linalg.inv(cholesky(S))
    df = np.dot(P, X-m)
    D = np.dot(df, df.T)
    eigval, eigvec = np.linalg.eig(D)
    order = eigval.argsort()[::-1]
    eigval = eigval[order]
    eigvec = eigvec[:, order]

    P = np.dot(eigvec.T, P)

    return m, P


def cholesky(S):
    # Cholesky decomposition
    # with adding a small value on the diagonal
    # until matrix is positive definite
    alpha = 0
    while 1:
        try:
            L = np.linalg.cholesky(S + alpha*np.eye(*S.shape))
            return L
        except Exception as e:
            LOGGER.error(f"cholesky error {e}")
            if alpha == 0:
                alpha = 1e-10
            else:
                alpha *= 10
            print(">>>> {}::cholesky: Matrix is not positive definite, adding {:.0e} on the diagonal"
                .format(os.path.basename(__file__), alpha))
