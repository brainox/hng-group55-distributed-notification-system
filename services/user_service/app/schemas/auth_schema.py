from pydantic import BaseModel, EmailStr
from typing import Optional


class RegisterRequest(BaseModel):
    name: str
    email: EmailStr
    password: str
    push_token: Optional[str] = None
    preferences: Optional["UserPreference"] = None  # reference to below


class LoginRequest(BaseModel):
    email: EmailStr
    password: str


class TokenResponse(BaseModel):
    access_token: str
    token_type: str = "bearer"


class UserPreference(BaseModel):
    email: bool = True
    push: bool = True


class UserOut(BaseModel):
    user_id: str
    name: str
    email: EmailStr
    push_token: Optional[str] = None
    preferences: Optional[UserPreference] = None

    class Config:
        orm_mode = True
