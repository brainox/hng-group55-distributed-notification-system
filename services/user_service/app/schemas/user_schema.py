from pydantic import BaseModel, EmailStr
from typing import Optional
from uuid import UUID


class UserCreate(BaseModel):
    email: EmailStr
    password: str
    full_name: Optional[str] = None


class UserOut(BaseModel):
    user_id: UUID
    email: EmailStr
    full_name: Optional[str] = None

    model_config = {
        "from_attributes": True
    }
