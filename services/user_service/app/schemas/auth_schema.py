from pydantic import BaseModel, EmailStr,Field
from typing import Optional,Annotated




class RegisterRequest(BaseModel):
    name: str
    email: EmailStr
    password: Annotated[str, Field(min_length=6, max_length=72)]
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
