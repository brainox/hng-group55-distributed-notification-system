from pydantic import BaseModel, EmailStr
from typing import Optional
import uuid

class UserCreate(BaseModel):
    email: EmailStr
    password: str
    full_name: Optional[str] = None

class UserOut(BaseModel):
    user_id: uuid.UUID
    email: EmailStr
    full_name: Optional[str] = None

    class Config:
        orm_mode = True
