from pydantic import BaseModel, EmailStr, Field
from typing import Optional, Annotated
from uuid import UUID


class UserPreference(BaseModel):
    email: bool = True
    push: bool = True


class RegisterRequest(BaseModel):
    name: str
    email: EmailStr
    password: Annotated[str, Field(min_length=6, max_length=72)]
    push_token: Optional[str] = None
    preferences: Optional[UserPreference] = None


class LoginRequest(BaseModel):
    email: EmailStr
    password: str


class TokenResponse(BaseModel):
    access_token: str
    token_type: str = "bearer"


class UserOut(BaseModel):
    user_id: UUID
    name: str = Field(alias="full_name")
    email: EmailStr
    push_token: Optional[str] = None
    preferences: Optional[UserPreference] = None

    model_config = {
        "from_attributes": True,
        "populate_by_name": True
    }
