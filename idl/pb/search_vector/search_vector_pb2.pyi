from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SearchVectorRequest(_message.Message):
    __slots__ = ["query"]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    query: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, query: _Optional[_Iterable[str]] = ...) -> None: ...

class SearchVectorResponse(_message.Message):
    __slots__ = ["code", "doc_ids", "msg", "error"]
    CODE_FIELD_NUMBER: _ClassVar[int]
    DOC_IDS_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    code: int
    doc_ids: _containers.RepeatedScalarFieldContainer[str]
    msg: str
    error: str
    def __init__(self, code: _Optional[int] = ..., doc_ids: _Optional[_Iterable[str]] = ..., msg: _Optional[str] = ..., error: _Optional[str] = ...) -> None: ...
