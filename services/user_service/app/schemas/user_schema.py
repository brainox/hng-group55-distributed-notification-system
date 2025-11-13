from pydantic import BaseModel, EmailStr
from typing import Optional
from uuid import UUID


class UserPreference(BaseModel):
    email: bool = True
    push: bool = True


class UserCreate(BaseModel):
    email: EmailStr
    password: str
    full_name: Optional[str] = None


class UserUpdate(BaseModel):
    full_name: Optional[str] = None
    email: Optional[EmailStr] = None
    push_token: Optional[str] = None
    preferences: Optional[UserPreference] = None


class UserOut(BaseModel):
    user_id: UUID
    email: EmailStr
    full_name: Optional[str] = None
    push_token: Optional[str] = None
    preferences: Optional[UserPreference] = None

    model_config = {
        "from_attributes": True
    }
