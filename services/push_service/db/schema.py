#!/usr/bin/env python3
from pydantic import BaseModel


class PushMessageCreate(BaseModel):
    title: str
    body: str
    token: str


class PushMessageResponse(BaseModel):
    id: int
    title: str
    body: str
    token: str
    status: str

    class Config:
        orm_mode = True
