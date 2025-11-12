#!/usr/bin/env python3
from typing import Optional, Any

from pydantic import BaseModel

class UserData(BaseModel):
    name: str
    link: str
    meta: Optional[dict]

class PushMessageCreate(BaseModel):
    id: int
    title: str
    body: str
    token: str
    user: UserData


class PushMessageResponse(BaseModel):
    id: int
    title: str
    body: str
    token: str
    status: str


class PaginationMeta(BaseModel):
    total: int
    limit: int
    page: int
    total_pages: int
    has_next: bool
    has_previous: bool


class RootResponse(BaseModel):

    success: bool
    data: Optional[Any]
    error: Optional[str]
    message: str
    meta: Optional[PaginationMeta]
