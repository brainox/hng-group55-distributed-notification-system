from pydantic import BaseModel
import uuid


class PreferenceCreate(BaseModel):
    email: bool = True
    push: bool = True


class PreferenceOut(BaseModel):
    preference_id: uuid.UUID
    email: bool
    push: bool

    class Config:
        orm_mode = True
